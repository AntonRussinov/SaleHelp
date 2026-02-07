package dto

// WBProductDetails — данные о товаре в отзыве/вопросе (общее для feedbacks и questions API).
type WBProductDetails struct {
	ImtID            int64  `json:"imtId"`
	NmID             int64  `json:"nmId"`
	ProductName      string `json:"productName"`
	SupplierArticle string `json:"supplierArticle"`
	SupplierName     string `json:"supplierName"`
	BrandName        string `json:"brandName"`
	Size             string `json:"size"`
}

// GetFeedbacksResponse — ответ GET https://feedbacks-api.wildberries.ru/api/v1/feedbacks.
// Документация: https://dev.wildberries.ru/docs/openapi/user-communication#tag/Otzyvy/paths/~1api~1v1~1feedbacks/get
type GetFeedbacksResponse struct {
	Data  GetFeedbacksData `json:"data"`
	Error bool             `json:"error"`
}

type GetFeedbacksData struct {
	CountUnanswered int          `json:"countUnanswered"`
	CountArchive    int          `json:"countArchive"`
	Feedbacks       []WBFeedback `json:"feedbacks"`
}

type WBFeedback struct {
	ID                 string              `json:"id"`
	Text               string              `json:"text"`
	Pros               string              `json:"pros"`
	Cons               string              `json:"cons"`
	ProductValuation   int                 `json:"productValuation"`
	CreatedDate        string              `json:"createdDate"`
	Answer             *WBFeedbackAnswer   `json:"answer"`
	State              string              `json:"state"`
	ProductDetails     WBProductDetails   `json:"productDetails"`
	WasViewed          bool                `json:"wasViewed"`
	UserName string `json:"userName"`
}

type WBFeedbackAnswer struct {
	Text     string `json:"text"`
	State    string `json:"state"`
	Editable bool   `json:"editable"`
}

// AnswerFeedbackRequest — тело запроса POST /api/v1/feedbacks/answer (ответить на отзыв).
// Документация: https://dev.wildberries.ru/docs/openapi/user-communication#tag/Otzyvy/paths/~1api~1v1~1feedbacks~1answer/post
type AnswerFeedbackRequest struct {
	ID   string `json:"id"`   // ID отзыва
	Text string `json:"text"` // Текст ответа, 2–5000 символов
}
