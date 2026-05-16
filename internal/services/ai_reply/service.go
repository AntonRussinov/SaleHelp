package ai_reply

import (
	"context"
	"errors"
	"log/slog"
	"strings"
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

type AISettingsRepo interface {
	GetByShopID(ctx context.Context, shopID int) (*domain.ShopAISettings, error)
}

type WBFeedbackRepo interface {
	ListUnrepliedByShop(ctx context.Context, shopID int, limit int) ([]*domain.WBFeedback, error)
	MarkReplied(ctx context.Context, ids []int, at time.Time) error
}

type AIFeedbackReplyRepo interface {
	Create(ctx context.Context, d *domain.AIFeedbackReply) error
	MarkWBSent(ctx context.Context, replyID int, at time.Time) error
	ListPendingWBSend(ctx context.Context, limit int) ([]domain.AIFeedbackPendingSend, error)
}

// QuotaRepo — legacy per-shop quota (используется если нет активной подписки).
type QuotaRepo interface {
	GetForMonth(ctx context.Context, shopID int, month time.Time) (*domain.ShopAIReplyQuota, error)
	EnsureForMonth(ctx context.Context, shopID int, month time.Time, includedLimit int) error
	DebitOne(ctx context.Context, shopID int, month time.Time, feedbackID int) (source string, err error)
}

type LedgerRepo interface {
	Create(ctx context.Context, d *domain.ShopAIReplyLedger) error
}

type WBFeedbackSender interface {
	AnswerFeedback(ctx context.Context, apiKey string, feedbackID string, text string) error
}

type ShopByIDRepo interface {
	GetByID(ctx context.Context, id int) (*domain.Shop, error)
}

// SubscriptionsRepo — квота из подписки (новая схема, Phase 2).
type SubscriptionsRepo interface {
	GetActiveByUserID(ctx context.Context, userID int) (*domain.Subscription, error)
	IncrementAIUsed(ctx context.Context, userID int) (bool, error)
}

type Service struct {
	settingsRepo   AISettingsRepo
	feedbackRepo   WBFeedbackRepo
	replyRepo      AIFeedbackReplyRepo
	quotaRepo      QuotaRepo
	ledgerRepo     LedgerRepo
	shopRepo       ShopByIDRepo
	subsRepo       SubscriptionsRepo // nil если подписки не используются
	wbFeedbackSend WBFeedbackSender
	aiClient       *client.Client
	modelConfig    ai.ModelConfig
	includedLimit  int
	log            *slog.Logger
}

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

// SetSubscriptionsRepo подключает квоту из подписки (вызывается из DI после инициализации).
func (s *Service) SetSubscriptionsRepo(r SubscriptionsRepo) { s.subsRepo = r }

// subQuotaStatus — результат проверки квоты из подписки.
type subQuotaStatus int

const (
	subQuotaNoSub     subQuotaStatus = iota // нет активной подписки → использовать legacy
	subQuotaConsumed                        // квота подписки списана → продолжать
	subQuotaExhausted                       // квота подписки исчерпана → остановить
)

// trySubscriptionQuota проверяет и атомарно списывает квоту из подписки.
func (s *Service) trySubscriptionQuota(ctx context.Context, shopID int) (subQuotaStatus, error) {
	if s.subsRepo == nil {
		return subQuotaNoSub, nil
	}
	shop, err := s.shopRepo.GetByID(ctx, shopID)
	if err != nil || shop == nil {
		return subQuotaNoSub, err
	}
	sub, err := s.subsRepo.GetActiveByUserID(ctx, shop.UserID)
	if err != nil {
		return subQuotaNoSub, err
	}
	if sub == nil || sub.AIQuota == 0 {
		// Нет подписки или тариф без AI-квоты → legacy
		return subQuotaNoSub, nil
	}
	ok, err := s.subsRepo.IncrementAIUsed(ctx, shop.UserID)
	if err != nil {
		return subQuotaNoSub, err
	}
	if !ok {
		return subQuotaExhausted, nil
	}
	return subQuotaConsumed, nil
}

// ProcessOne генерирует ответ на один отзыв.
func (s *Service) ProcessOne(ctx context.Context, f *domain.WBFeedback) error {
	settings, err := s.settingsRepo.GetByShopID(ctx, f.ShopID)
	if err != nil {
		return err
	}
	if settings == nil || !settings.IsActive {
		s.saveReplyLog(ctx, f.ID, f.ShopID, "", "", domain.AIReplyStatusSkipped, "AI disabled", nil)
		return nil
	}

	// Квота: сначала проверяем подписку (новая схема), если нет — legacy per-shop quota.
	subStatus, err := s.trySubscriptionQuota(ctx, f.ShopID)
	if err != nil {
		return err
	}
	if subStatus == subQuotaExhausted {
		s.saveReplyLog(ctx, f.ID, f.ShopID, "", "", domain.AIReplyStatusNoQuota, "", nil)
		return nil
	}

	if subStatus == subQuotaNoSub {
		// Legacy: per-shop monthly quota
		month := billingMonth()
		if err := s.quotaRepo.EnsureForMonth(ctx, f.ShopID, month, s.includedLimit); err != nil {
			return err
		}
		quota, err := s.quotaRepo.GetForMonth(ctx, f.ShopID, month)
		if err != nil {
			return err
		}
		if quota == nil || (quota.IncludedUsed >= quota.IncludedLimit && quota.PaidBalance <= 0) {
			s.saveReplyLog(ctx, f.ID, f.ShopID, "", "", domain.AIReplyStatusNoQuota, "", nil)
			return nil
		}
	}

	toneInstruction := tone.Instruction(settings.ToneOfVoice)
	brandName := settings.BrandName
	if brandName == "" {
		brandName = "магазина"
	}
	in := &ai.FeedbackReplyInput{
		BrandName:            brandName,
		ToneInstruction:      toneInstruction,
		ProductName:          sanitize.Text(f.ProductName),
		Rating:               f.ProductValuation,
		Text:                 sanitize.Text(f.Text),
		Pros:                 sanitize.Text(f.Pros),
		Cons:                 sanitize.Text(f.Cons),
		CustomRules:          splitLines(settings.CustomRules),
		OrderStatus:          f.OrderStatus,
		CategoryCustomPrompt: "",
	}
	promptStr, err := prompt.FormatFeedbackReply(in)
	if err != nil {
		s.saveReplyLog(ctx, f.ID, f.ShopID, "", "", domain.AIReplyStatusError, err.Error(), nil)
		return err
	}
	result, err := s.aiClient.GenerateFeedbackReply(ctx, promptStr, f.ProductValuation)
	if err != nil {
		s.saveReplyLog(ctx, f.ID, f.ShopID, "", s.modelConfig.SelectModel(f.ProductValuation), domain.AIReplyStatusError, err.Error(), nil)
		return err
	}
	answerText := truncate500(sanitize.Text(result.Text))

	// Debit legacy quota (только если не использовали подписку).
	if subStatus == subQuotaNoSub {
		month := billingMonth()
		source, err := s.quotaRepo.DebitOne(ctx, f.ShopID, month, f.ID)
		if err != nil {
			if errors.Is(err, ai.ErrNoQuota) {
				s.saveReplyLog(ctx, f.ID, f.ShopID, answerText, result.ModelUsed, domain.AIReplyStatusNoQuota, "", result)
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
	}

	s.saveReplyLog(ctx, f.ID, f.ShopID, answerText, result.ModelUsed, domain.AIReplyStatusGenerated, "", result)

	if err := s.feedbackRepo.MarkReplied(ctx, []int{f.ID}, time.Now()); err != nil {
		s.log.Error("ai_reply: mark replied failed", "feedback_id", f.ID, "err", err)
	}
	return nil
}

func billingMonth() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func ptr(i int) *int { return &i }

func splitLines(s string) []string {
	var result []string
	for _, line := range strings.Split(s, "\n") {
		if t := strings.TrimSpace(line); t != "" {
			result = append(result, t)
		}
	}
	return result
}

func truncate500(s string) string {
	const max = 500
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max])
}

func (s *Service) saveReplyLog(ctx context.Context, feedbackID, shopID int, answerText, modelUsed, status, errMsg string, result *ai.FeedbackReplyResult) {
	if modelUsed == "" {
		modelUsed = "n/a"
	}
	d := &domain.AIFeedbackReply{
		FeedbackID:   feedbackID,
		ShopID:       shopID,
		AnswerText:   answerText,
		ModelUsed:    modelUsed,
		Status:       status,
		ErrorMessage: errMsg,
		AITier:       "paid",
	}
	if result != nil {
		d.InputTokens = result.InputTokens
		d.OutputTokens = result.OutputTokens
		d.TokensUsed = result.TokensUsed
		d.EstimatedCostKopeks = ai.EstimateCostKopeks(result.ModelUsed, result.InputTokens, result.OutputTokens)
	}
	_ = s.replyRepo.Create(ctx, d)
}
