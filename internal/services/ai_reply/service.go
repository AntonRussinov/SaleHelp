package ai_reply

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"salehelp/domain"
	"salehelp/pkg/ai"
	"salehelp/pkg/ai/client"
	"salehelp/pkg/ai/prompt"
	"salehelp/pkg/ai/sanitize"
	"salehelp/pkg/ai/tone"
)

const (
	defaultIncludedLimit = 50
	ledgerReason         = "ai_feedback_reply"
)

// AISettingsRepo доступ к настройкам AI магазина.
type AISettingsRepo interface {
	GetByShopID(ctx context.Context, shopID int) (*domain.ShopAISettings, error)
}

// WBFeedbackRepo отзывы без ответа и пометка ответа.
type WBFeedbackRepo interface {
	ListUnrepliedByShop(ctx context.Context, shopID int, limit int) ([]*domain.WBFeedback, error)
	MarkReplied(ctx context.Context, ids []int, at time.Time) error
}

// AIFeedbackReplyRepo сохранение логов генерации.
type AIFeedbackReplyRepo interface {
	Create(ctx context.Context, d *domain.AIFeedbackReply) error
}

// QuotaRepo квоты и списание.
type QuotaRepo interface {
	GetForMonth(ctx context.Context, shopID int, month time.Time) (*domain.ShopAIReplyQuota, error)
	EnsureForMonth(ctx context.Context, shopID int, month time.Time, includedLimit int) error
	DebitOne(ctx context.Context, shopID int, month time.Time, feedbackID int) (source string, err error)
}

// LedgerRepo аудит списаний.
type LedgerRepo interface {
	Create(ctx context.Context, d *domain.ShopAIReplyLedger) error
}

// WBFeedbackSender отправка ответа на отзыв в WB (POST /api/v1/feedbacks/answer).
type WBFeedbackSender interface {
	AnswerFeedback(ctx context.Context, apiKey string, feedbackID string, text string) error
}

// ShopByIDRepo получение магазина по ID (нужен APIToken для вызова WB).
type ShopByIDRepo interface {
	GetByID(ctx context.Context, id int) (*domain.Shop, error)
}

// Service сервис генерации AI-ответов на отзывы.
type Service struct {
	settingsRepo   AISettingsRepo
	feedbackRepo   WBFeedbackRepo
	replyRepo      AIFeedbackReplyRepo
	quotaRepo      QuotaRepo
	ledgerRepo     LedgerRepo
	shopRepo       ShopByIDRepo
	wbFeedbackSend WBFeedbackSender
	aiClient       *client.Client
	modelConfig    ai.ModelConfig
	includedLimit  int // бесплатных ответов в месяц по умолчанию (из подписки)
	log            *slog.Logger
}

// New создаёт сервис.
func New(
	settingsRepo AISettingsRepo,
	feedbackRepo WBFeedbackRepo,
	replyRepo AIFeedbackReplyRepo,
	quotaRepo QuotaRepo,
	ledgerRepo LedgerRepo,
	shopRepo ShopByIDRepo,
	wbFeedbackSend WBFeedbackSender,
	aiClient *client.Client,
	modelConfig ai.ModelConfig,
	includedLimit int,
	log *slog.Logger,
) *Service {
	if includedLimit <= 0 {
		includedLimit = defaultIncludedLimit
	}
	return &Service{
		settingsRepo:   settingsRepo,
		feedbackRepo:   feedbackRepo,
		replyRepo:      replyRepo,
		quotaRepo:      quotaRepo,
		ledgerRepo:     ledgerRepo,
		shopRepo:       shopRepo,
		wbFeedbackSend: wbFeedbackSend,
		aiClient:       aiClient,
		modelConfig:    modelConfig,
		includedLimit:  includedLimit,
		log:            log,
	}
}

// ProcessOne генерирует ответ на один отзыв: проверяет настройки, квоту, санитизирует, вызывает LLM, сохраняет ответ и списание.
func (s *Service) ProcessOne(ctx context.Context, f *domain.WBFeedback) error {
	settings, err := s.settingsRepo.GetByShopID(ctx, f.ShopID)
	if err != nil {
		return err
	}
	if settings == nil || !settings.IsActive {
		s.saveReplyLog(ctx, f.ID, f.ShopID, "", "", domain.AIReplyStatusSkipped, "AI disabled", 0)
		return nil
	}
	month := time.Now().UTC()
	month = time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC)
	if err := s.quotaRepo.EnsureForMonth(ctx, f.ShopID, month, s.includedLimit); err != nil {
		return err
	}
	quota, err := s.quotaRepo.GetForMonth(ctx, f.ShopID, month)
	if err != nil {
		return err
	}
	hasQuota := (quota != nil && (quota.IncludedUsed < quota.IncludedLimit || quota.PaidBalance > 0))
	if !hasQuota {
		s.saveReplyLog(ctx, f.ID, f.ShopID, "", "", domain.AIReplyStatusNoQuota, "", 0)
		return nil
	}
	toneInstruction := tone.Instruction(settings.ToneOfVoice)
	brandName := settings.BrandName
	if brandName == "" {
		brandName = "магазина"
	}
	in := &ai.FeedbackReplyInput{
		BrandName:       brandName,
		ToneInstruction: toneInstruction,
		ProductName:     sanitize.Text(f.ProductName),
		Rating:          f.ProductValuation,
		Text:            sanitize.Text(f.Text),
		Pros:            sanitize.Text(f.Pros),
		Cons:            sanitize.Text(f.Cons),
		CustomRules:     settings.CustomRules,
		OrderStatus:     f.OrderStatus, // "returned" — товар уже возвращён, не предлагать возврат в ответе
	}
	promptStr, err := prompt.FormatFeedbackReply(in)
	if err != nil {
		s.saveReplyLog(ctx, f.ID, f.ShopID, "", "", domain.AIReplyStatusError, err.Error(), 0)
		return err
	}
	result, err := s.aiClient.GenerateFeedbackReply(ctx, promptStr, f.ProductValuation)
	if err != nil {
		s.saveReplyLog(ctx, f.ID, f.ShopID, "", s.modelConfig.ModelForOther, domain.AIReplyStatusError, err.Error(), 0)
		return err
	}
	answerText := truncate500(sanitize.Text(result.Text))
	source, err := s.quotaRepo.DebitOne(ctx, f.ShopID, month, f.ID)
	if err != nil {
		if errors.Is(err, ai.ErrNoQuota) {
			s.saveReplyLog(ctx, f.ID, f.ShopID, answerText, result.ModelUsed, domain.AIReplyStatusNoQuota, "", result.TokensUsed)
			return nil
		}
		return err
	}
	if err := s.ledgerRepo.Create(ctx, &domain.ShopAIReplyLedger{
		ShopID:     f.ShopID,
		FeedbackID: ptr(f.ID),
		Month:      month,
		Source:     source,
		Delta:      -1,
		Reason:     ledgerReason,
	}); err != nil {
		s.log.Warn("ai_reply: ledger create failed", "feedback_id", f.ID, "err", err)
	}
	s.saveReplyLog(ctx, f.ID, f.ShopID, answerText, result.ModelUsed, domain.AIReplyStatusGenerated, "", result.TokensUsed)

	// Помечаем как обработанный ДО отправки в WB: квота уже списана, повторная обработка недопустима.
	// Если WB-отправка упадёт — отзыв всё равно не будет обработан снова и не сожжёт квоту.
	now := time.Now()
	if err := s.feedbackRepo.MarkReplied(ctx, []int{f.ID}, now); err != nil {
		s.log.Error("ai_reply: mark replied failed", "feedback_id", f.ID, "err", err)
	}

	// Отправляем ответ в WB, если есть текст (WB требует 2–5000 символов).
	if len([]rune(answerText)) >= 2 && s.wbFeedbackSend != nil && s.shopRepo != nil {
		shop, err := s.shopRepo.GetByID(ctx, f.ShopID)
		if err != nil {
			s.log.Warn("ai_reply: get shop for WB send failed", "feedback_id", f.ID, "shop_id", f.ShopID, "err", err)
		} else if shop != nil && shop.APIToken != "" {
			if err := s.wbFeedbackSend.AnswerFeedback(ctx, shop.APIToken, f.WBFeedbackID, answerText); err != nil {
				s.log.Error("ai_reply: send reply to WB failed", "feedback_id", f.ID, "wb_feedback_id", f.WBFeedbackID, "err", err)
			}
		}
	}

	return nil
}

func ptr(i int) *int { return &i }

func truncate500(s string) string {
	const max = 500
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max])
}

func (s *Service) saveReplyLog(ctx context.Context, feedbackID, shopID int, answerText, modelUsed, status, errMsg string, tokensUsed int) {
	if modelUsed == "" {
		modelUsed = "n/a"
	}
	_ = s.replyRepo.Create(ctx, &domain.AIFeedbackReply{
		FeedbackID:   feedbackID,
		ShopID:       shopID,
		AnswerText:   answerText,
		ModelUsed:    modelUsed,
		Status:       status,
		ErrorMessage: errMsg,
		TokensUsed:   tokensUsed,
	})
}
