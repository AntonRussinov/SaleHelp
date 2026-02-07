package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"salehelp/internal/adapter/wb/dto"
)

// GetFeedbacks возвращает список отзывов (GET /api/v1/feedbacks).
// isAnswered: false — необработанные, true — обработанные.
// Лимит WB: 3 запроса/сек для категории «Вопросы и отзывы».
func (c *Client) GetFeedbacks(ctx context.Context, apiKey string, isAnswered bool, take, skip int, order string, dateFrom, dateTo *int64) (*dto.GetFeedbacksResponse, error) {
	u, err := url.Parse(c.feedbacksURL + "/api/v1/feedbacks")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("isAnswered", strconv.FormatBool(isAnswered))
	q.Set("take", strconv.Itoa(take))
	q.Set("skip", strconv.Itoa(skip))
	if order != "" {
		q.Set("order", order)
	}
	if dateFrom != nil {
		q.Set("dateFrom", strconv.FormatInt(*dateFrom, 10))
	}
	if dateTo != nil {
		q.Set("dateTo", strconv.FormatInt(*dateTo, 10))
	}
	u.RawQuery = q.Encode()

	var result dto.GetFeedbacksResponse
	if err := c.doJSON(ctx, http.MethodGet, u.String(), apiKey, RateFeedbacks, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetQuestions возвращает список вопросов (GET /api/v1/questions).
// isAnswered: false — неотвеченные, true — отвеченные.
// Лимит WB: 3 запроса/сек для категории «Вопросы и отзывы».
func (c *Client) GetQuestions(ctx context.Context, apiKey string, isAnswered bool, take, skip int, order string, dateFrom, dateTo *int64) (*dto.GetQuestionsResponse, error) {
	u, err := url.Parse(c.feedbacksURL + "/api/v1/questions")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("isAnswered", strconv.FormatBool(isAnswered))
	q.Set("take", strconv.Itoa(take))
	q.Set("skip", strconv.Itoa(skip))
	if order != "" {
		q.Set("order", order)
	}
	if dateFrom != nil {
		q.Set("dateFrom", strconv.FormatInt(*dateFrom, 10))
	}
	if dateTo != nil {
		q.Set("dateTo", strconv.FormatInt(*dateTo, 10))
	}
	u.RawQuery = q.Encode()

	var result dto.GetQuestionsResponse
	if err := c.doJSON(ctx, http.MethodGet, u.String(), apiKey, RateFeedbacks, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AnswerFeedback отправляет ответ на отзыв (POST /api/v1/feedbacks/answer).
// feedbackID — ID отзыва, text — текст ответа (2–5000 символов). Ответ при успехе: 204 No Content.
// Документация: https://dev.wildberries.ru/docs/openapi/user-communication#tag/Otzyvy/paths/~1api~1v1~1feedbacks~1answer/post
func (c *Client) AnswerFeedback(ctx context.Context, apiKey string, feedbackID string, text string) error {
	req := dto.AnswerFeedbackRequest{ID: feedbackID, Text: text}
	return c.doJSON(ctx, http.MethodPost, c.feedbacksURL+"/api/v1/feedbacks/answer", apiKey, RateFeedbacks, req, nil)
}

// PatchQuestion выполняет действие над вопросом (PATCH /api/v1/questions).
// body — один из: QuestionViewRequest (просмотр), QuestionAnswerRequest (ответ/редактирование) и т.д.
// Документация: https://dev.wildberries.ru/docs/openapi/user-communication#tag/Voprosy/paths/~1api~1v1~1questions/patch
func (c *Client) PatchQuestion(ctx context.Context, apiKey string, body any) error {
	return c.doJSON(ctx, http.MethodPatch, c.feedbacksURL+"/api/v1/questions", apiKey, RateFeedbacks, body, nil)
}

// MarkQuestionViewed отмечает вопрос как просмотренный (PATCH /api/v1/questions, body: id + wasViewed).
func (c *Client) MarkQuestionViewed(ctx context.Context, apiKey string, questionID string, viewed bool) error {
	return c.PatchQuestion(ctx, apiKey, dto.QuestionViewRequest{ID: questionID, WasViewed: viewed})
}

// AnswerQuestion отправляет ответ на вопрос или редактирует существующий ответ (PATCH /api/v1/questions).
// Тело запроса: id + answer.text. Если API вернёт 400, возможно WB ожидает другой формат — см. документацию.
func (c *Client) AnswerQuestion(ctx context.Context, apiKey string, questionID string, answerText string) error {
	return c.PatchQuestion(ctx, apiKey, dto.QuestionAnswerRequest{
		ID:     questionID,
		Answer: dto.WBQuestionAnswerIn{Text: answerText},
	})
}
