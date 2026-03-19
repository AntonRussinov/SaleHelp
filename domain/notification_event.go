package domain

// Notification event keys (used for per-shop notification preferences).
// For alert-triggered events we reuse alert types (e.g. "price_drop").
const (
	NotificationEventFBSOrders          = "fbs_orders"
	NotificationEventDBWOrders          = "dbw_orders"
	NotificationEventDBSOrders          = "dbs_orders"
	NotificationEventSalesAndRefusals   = "sales_and_refusals"
	NotificationEventWBFeedbacks        = "wb_feedbacks"
	NotificationEventWBQuestions        = "wb_questions"
	NotificationEventStockReplenish     = "stock_replenish"
	NotificationEventSubscriptionStatus = "subscription_status"
	NotificationEventFirstSync          = "first_sync"
)
