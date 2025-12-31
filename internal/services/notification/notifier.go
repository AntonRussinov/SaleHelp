package notifications

// import (
// 	"context"
// 	"fmt"

// 	"salehelp/pkg/notify"
// )

// type Router interface {
// 	GetRecipient(ctx context.Context, shopID int) (notify.Recipient, error)
// 	EnabledChannels(ctx context.Context, shopID int) ([]notify.Channel, error)
// }

// type Notifier struct {
// 	router  Router
// 	senders map[notify.Channel]notify.Sender
// }

// func NewNotifier(router Router, senders []notify.Sender) *Notifier {
// 	m := make(map[notify.Channel]notify.Sender, len(senders))
// 	for _, s := range senders {
// 		m[s.Channel()] = s
// 	}
// 	return &Notifier{
// 		router:  router,
// 		senders: m,
// 	}
// }

// // MVP: fail-fast. Если один канал упал — считаем доставку неуспешной.
// // Позже можно сделать best-effort и таблицу delivery по каналам.
// func (n *Notifier) SendToEnabledChannels(ctx context.Context, shopID int, msg notify.Message) error {
// 	rcpt, err := n.router.GetRecipient(ctx, shopID)
// 	if err != nil {
// 		return err
// 	}

// 	channels, err := n.router.EnabledChannels(ctx, shopID)
// 	if err != nil {
// 		return err
// 	}
// 	if len(channels) == 0 {
// 		return fmt.Errorf("no enabled channels for shop_id=%d", shopID)
// 	}

// 	for _, ch := range channels {
// 		s, ok := n.senders[ch]
// 		if !ok {
// 			// канал включён, но sender не зарегистрирован — это конфиг/код ошибка
// 			return fmt.Errorf("sender not registered for channel=%s", ch)
// 		}
// 		if err := s.Send(ctx, rcpt, msg); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
