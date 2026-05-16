package domain

import "time"

// ToneOfVoice — допустимые значения тона ответов AI.
const (
	ToneStandard  = "standard"
	ToneFriendly   = "friendly"
	ToneOfficial   = "official"
	ToneEmpathic   = "empathic"
	ToneMinimalist = "minimalist"
)

// WBFeedback — отзыв покупателя из WB API (Вопросы и отзывы).
// Используется для уведомлений о новых отзывах в Telegram и для AI-ответов.
// OrderStatusBuyout — покупатель выкупил товар.
const OrderStatusBuyout = "buyout"

// OrderStatusReturned — товар возвращён (отказ от товара).
const OrderStatusReturned = "returned"

type WBFeedback struct {
	ID               int
	ShopID           int
	WBFeedbackID     string    // id из WB API (уникален в рамках WB)
	NmID             int64
	ProductName      string
	Text             string
	Pros             string
	Cons             string
	ProductValuation int
	UserName         string
	CreatedAt        time.Time
	NotifiedAt       *time.Time
	IsReplied        bool       // true, если сгенерирован и сохранён AI-ответ
	RepliedAt        *time.Time // время проставления ответа
	OrderStatus      string     // "buyout" — выкуп, "returned" — товар возвращён (из WB API)
	Extra            []byte     // JSON при необходимости
}

// ShopAISettings — настройки AI-ответов для магазина.
type ShopAISettings struct {
	ShopID      int
	ToneOfVoice string // standard | friendly | official | empathic | minimalist
	BrandName   string
	CustomRules string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// AIFeedbackReplyStatus — статус записи генерации ответа.
const (
	AIReplyStatusGenerated = "generated"
	AIReplyStatusError     = "error"
	AIReplyStatusSkipped   = "skipped"
	AIReplyStatusNoQuota   = "no_quota"
)

// AIFeedbackReply — лог попытки/результата генерации AI-ответа на отзыв.
type AIFeedbackReply struct {
	ID                   int
	FeedbackID           int
	ShopID               int
	AnswerText           string
	ModelUsed            string
	Status               string // generated | error | skipped | no_quota
	ErrorMessage         string
	InputTokens          int
	OutputTokens         int
	TokensUsed           int
	AITier               string // paid | free_data_sharing
	EstimatedCostKopeks  int
	WBSentAt             *time.Time
	CreatedAt            time.Time
}

// AIFeedbackPendingSend — сгенерированный ответ, ещё не отправленный в WB.
type AIFeedbackPendingSend struct {
	ReplyID      int
	ShopID       int
	AnswerText   string
	WBFeedbackID string // WB-идентификатор отзыва (из wb_feedbacks.wb_feedback_id)
}

// ShopAIReplyQuota — квота/баланс ответов по магазину за месяц.
type ShopAIReplyQuota struct {
	ID             int
	ShopID         int
	Month          time.Time // первый день месяца
	IncludedLimit  int       // лимит бесплатных по подписке
	IncludedUsed   int
	PaidBalance    int
	PaidUsed       int
	UpdatedAt      time.Time
}

// ShopAIReplyPurchaseStatus — статус покупки пакета.
const (
	AIPurchaseStatusPending  = "pending"
	AIPurchaseStatusPaid     = "paid"
	AIPurchaseStatusCanceled = "canceled"
)

// ShopAIReplyPurchase — покупка пакета ответов.
type ShopAIReplyPurchase struct {
	ID           int
	ShopID       int
	PackageName  string
	AnswersCount int
	PriceRub     int
	PaymentRef   string
	Status       string // pending | paid | canceled
	CreatedAt    time.Time
	PaidAt       *time.Time
}

// LedgerSource — источник списания.
const (
	LedgerSourceIncluded = "included"
	LedgerSourcePaid     = "paid"
)

// ShopAIReplyLedger — запись аудита списания ответа.
type ShopAIReplyLedger struct {
	ID         int
	ShopID     int
	FeedbackID *int
	Month      time.Time
	Source     string // included | paid
	Delta      int    // обычно -1
	Reason     string
	CreatedAt  time.Time
}
