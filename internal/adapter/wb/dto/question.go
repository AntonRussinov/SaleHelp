package dto

// GetQuestionsResponse — ответ GET https://feedbacks-api.wildberries.ru/api/v1/questions.
// Документация: https://dev.wildberries.ru/docs/openapi/user-communication#tag/Voprosy/paths/~1api~1v1~1questions/get
type GetQuestionsResponse struct {
	Data  GetQuestionsData `json:"data"`
	Error bool             `json:"error"`
}

type GetQuestionsData struct {
	CountUnanswered int          `json:"countUnanswered"`
	CountArchive    int          `json:"countArchive"`
	Questions       []WBQuestion `json:"questions"`
}

type WBQuestion struct {
	ID             string            `json:"id"`
	Text           string            `json:"text"`
	CreatedDate    string            `json:"createdDate"`
	State          string            `json:"state"`
	Answer         *WBQuestionAnswer `json:"answer"`
	ProductDetails WBProductDetails `json:"productDetails"`
	WasViewed      bool             `json:"wasViewed"`
	IsWarned       bool             `json:"isWarned"`
}

type WBQuestionAnswer struct {
	Text       string `json:"text"`
	Editable   bool   `json:"editable"`
	CreateDate string `json:"createDate"`
}

// Тела запросов PATCH /api/v1/questions (работа с вопросами).
// Документация: https://dev.wildberries.ru/docs/openapi/user-communication#tag/Voprosy/paths/~1api~1v1~1questions/patch
// В зависимости от тела: просмотр (wasViewed), отклонить вопрос, ответить или отредактировать ответ.

// QuestionViewRequest — отметить вопрос как просмотренный.
type QuestionViewRequest struct {
	ID        string `json:"id"`
	WasViewed bool   `json:"wasViewed"`
}

// QuestionAnswerRequest — ответить на вопрос или отредактировать ответ.
// Формат по аналогии с ответом в API (answer.text); если WB вернёт 400, возможно другое имя поля.
type QuestionAnswerRequest struct {
	ID     string             `json:"id"`
	Answer WBQuestionAnswerIn `json:"answer"`
}

// WBQuestionAnswerIn — вложенный объект ответа в запросе (только текст).
type WBQuestionAnswerIn struct {
	Text string `json:"text"`
}
