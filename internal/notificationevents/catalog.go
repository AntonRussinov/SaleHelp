package notificationevents

import "salehelp/domain"

type Event struct {
	Key            string
	Title          string
	Description    string
	Category       string
	DefaultEnabled bool
}

// Catalog is the single source of truth for user-togglable notification events.
// Add new events here to make them appear in both API and Telegram settings.
var Catalog = []Event{
	// Alerts (digest via dispatcher)
	{
		Key:            domain.AlertTypePriceDrop,
		Title:          "Падение цены",
		Description:    "Алерт, когда цена товара снизилась относительно предыдущего снимка.",
		Category:       "alerts",
		DefaultEnabled: true,
	},
	{
		Key:            domain.AlertTypeStockLow,
		Title:          "Остатки ниже порога",
		Description:    "Алерт, когда остатки товара становятся меньше заданного порога.",
		Category:       "alerts",
		DefaultEnabled: true,
	},
	{
		Key:            domain.AlertTypeStockZero,
		Title:          "Остатки закончились",
		Description:    "Алерт, когда остатки товара становятся равны нулю.",
		Category:       "alerts",
		DefaultEnabled: true,
	},

	// Orders / sales
	{
		Key:            domain.NotificationEventFBSOrders,
		Title:          "Новые заказы FBS",
		Description:    "Уведомления о новых сборочных заданиях FBS.",
		Category:       "orders",
		DefaultEnabled: true,
	},
	{
		Key:            domain.NotificationEventDBWOrders,
		Title:          "Новые заказы DBW",
		Description:    "Уведомления о новых сборочных заданиях (DBW).",
		Category:       "orders",
		DefaultEnabled: true,
	},
	{
		Key:            domain.NotificationEventDBSOrders,
		Title:          "Новые заказы DBS",
		Description:    "Уведомления о новых сборочных заданиях (DBS).",
		Category:       "orders",
		DefaultEnabled: true,
	},
	{
		Key:            domain.NotificationEventSalesAndRefusals,
		Title:          "Выкупы/возвраты/отказы",
		Description:    "Уведомления о новых выкупах, возвратах и отказах (не забрали).",
		Category:       "orders",
		DefaultEnabled: true,
	},

	// Feedback
	{
		Key:            domain.NotificationEventWBFeedbacks,
		Title:          "Новые отзывы",
		Description:    "Уведомления о новых отзывах покупателей.",
		Category:       "feedback",
		DefaultEnabled: true,
	},
	{
		Key:            domain.NotificationEventWBQuestions,
		Title:          "Новые вопросы",
		Description:    "Уведомления о новых вопросах покупателей.",
		Category:       "feedback",
		DefaultEnabled: true,
	},

	// Reports / system
	{
		Key:            domain.NotificationEventStockReplenish,
		Title:          "Отчёт: пополнение остатков",
		Description:    "Периодический отчёт по товарам, которые скоро закончатся.",
		Category:       "reports",
		DefaultEnabled: true,
	},
	{
		Key:            domain.NotificationEventSubscriptionStatus,
		Title:          "Подписка: напоминания",
		Description:    "Напоминания о скором окончании подписки и её истечении.",
		Category:       "system",
		DefaultEnabled: true,
	},
	{
		Key:            domain.NotificationEventFirstSync,
		Title:          "Системные: первый синк",
		Description:    "Служебные уведомления о первом синке и первичной загрузке данных.",
		Category:       "system",
		DefaultEnabled: true,
	},
}

func All() []Event {
	out := make([]Event, len(Catalog))
	copy(out, Catalog)
	return out
}

func Find(key string) (Event, bool) {
	for _, e := range Catalog {
		if e.Key == key {
			return e, true
		}
	}
	return Event{}, false
}
