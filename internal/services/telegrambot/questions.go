package telegrambot

import (
	"context"
	"fmt"
	"salehelp/domain"
	"strconv"
	"strings"
	"time"
)

const (
	StateWaitQuestionAnswer = "wait_question_answer"
	payloadQuestionID       = "question_id"
	payloadShopID           = "shop_id_q"
)

// handleAnswerQuestion вызывается при нажатии кнопки «💬 Ответить» под уведомлением о вопросе.
// payload format: "<shopID>:<wbQuestionID>"
func (s *Service) handleAnswerQuestion(ctx context.Context, chatID int64, payload string) {
	parts := strings.SplitN(payload, ":", 2)
	if len(parts) != 2 {
		_ = s.tg.SendMessage(ctx, chatID, "Неверный формат запроса.")
		return
	}
	shopID, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil || shopID <= 0 {
		_ = s.tg.SendMessage(ctx, chatID, "Неверный идентификатор магазина.")
		return
	}
	wbQuestionID := strings.TrimSpace(parts[1])

	// Проверяем, что этот чат связан с магазином.
	shops, err := s.shopsRepo.ListByTelegramChat(ctx, chatID)
	if err != nil {
		s.log.Error("handleAnswerQuestion: list shops by chat", "err", err)
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка проверки доступа.")
		return
	}
	hasAccess := false
	for _, sh := range shops {
		if sh.ID == shopID {
			hasAccess = true
			break
		}
	}
	if !hasAccess {
		_ = s.tg.SendMessage(ctx, chatID, "У вас нет доступа к этому магазину.")
		return
	}

	// Пытаемся найти текст вопроса для контекста.
	questionContext := ""
	if s.wbQuestionsRepo != nil {
		q, errQ := s.wbQuestionsRepo.GetByWBQuestionID(ctx, shopID, wbQuestionID)
		if errQ == nil && q != nil {
			productName := q.ProductName
			if productName == "" {
				productName = fmt.Sprintf("nmId %d", q.NmID)
			}
			questionContext = fmt.Sprintf("\nТовар: %s\nВопрос: \"%s\"\n", productName, q.Text)
		}
	}

	// Сохраняем состояние ожидания ответа.
	p := payloadMap{
		payloadShopID:     strconv.Itoa(shopID),
		payloadQuestionID: wbQuestionID,
	}
	st := &domain.TelegramChatState{
		ChatID:  chatID,
		State:   StateWaitQuestionAnswer,
		Payload: mustJSONAny(p),
	}
	if err := s.statesRepo.Upsert(ctx, st, time.Now()); err != nil {
		s.log.Error("handleAnswerQuestion: upsert state", "err", err)
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка сохранения состояния. Попробуйте ещё раз.")
		return
	}

	prompt := "✏️ Введите ваш ответ на вопрос (до 1000 символов)."
	if questionContext != "" {
		prompt = questionContext + "\n" + prompt
	}
	prompt += "\n\nДля отмены напишите /cancel."
	_ = s.tg.SendMessage(ctx, chatID, prompt)
}

// stepQuestionAnswer обрабатывает введённый пользователем текст ответа на вопрос.
func (s *Service) stepQuestionAnswer(ctx context.Context, chatID int64, text string) {
	text = strings.TrimSpace(text)
	if len([]rune(text)) < 2 {
		_ = s.tg.SendMessage(ctx, chatID, "Ответ слишком короткий. Введите текст ответа или /cancel для отмены.")
		return
	}
	if len([]rune(text)) > 1000 {
		_ = s.tg.SendMessage(ctx, chatID, "Ответ слишком длинный (максимум 1000 символов). Сократите текст.")
		return
	}

	st, err := s.statesRepo.GetByChatID(ctx, chatID)
	if err != nil || st == nil {
		_ = s.tg.SendMessage(ctx, chatID, "Состояние потеряно. Нажмите кнопку «Ответить» в уведомлении снова.")
		return
	}

	p, err := parsePayload(st.Payload)
	if err != nil {
		_ = s.statesRepo.Clear(ctx, chatID)
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка состояния. Нажмите кнопку «Ответить» снова.")
		return
	}

	shopID, err := strconv.Atoi(p[payloadShopID])
	if err != nil || shopID <= 0 {
		_ = s.statesRepo.Clear(ctx, chatID)
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка: не удалось определить магазин. Нажмите «Ответить» снова.")
		return
	}
	wbQuestionID := p[payloadQuestionID]

	shop, err := s.shopsRepo.GetByID(ctx, shopID)
	if err != nil || shop == nil {
		_ = s.statesRepo.Clear(ctx, chatID)
		_ = s.tg.SendMessage(ctx, chatID, "Магазин не найден.")
		return
	}

	if err := s.wbClient.AnswerQuestion(ctx, shop.APIToken, wbQuestionID, text); err != nil {
		s.log.Error("stepQuestionAnswer: AnswerQuestion WB API failed", "shop_id", shopID, "question_id", wbQuestionID, "err", err)
		_ = s.tg.SendMessage(ctx, chatID, "❌ Не удалось отправить ответ в WB. Попробуйте через личный кабинет WB.")
		_ = s.statesRepo.Clear(ctx, chatID)
		return
	}

	_ = s.statesRepo.Clear(ctx, chatID)
	_ = s.tg.SendMessage(ctx, chatID, "✅ Ответ на вопрос отправлен!")
}
