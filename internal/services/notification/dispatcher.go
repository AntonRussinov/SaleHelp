package notifications

import (
	"context"
	"fmt"
	"salehelp/pkg/notify"
	"strconv"
	"time"
)

type Dispatcher struct {
	deliveriesRepo DeliveriesRepo
	triggeredRepo  TriggeredRepo
	endpointsRepo  EndpointsRepo

	registry *SenderRegistry

	workerID     string
	limit        int
	leaseSeconds int
	maxAttempts  int
}

type SenderRegistry struct {
	senders map[notify.Channel]notify.Sender
}

func NewDispatcher(
	deliveriesRepo DeliveriesRepo,
	triggeredRepo TriggeredRepo,
	endpointsRepo EndpointsRepo,
	registry *SenderRegistry,
	workerID string,
) *Dispatcher {

	if workerID == "" {
		workerID = "worker-unknown"
	}

	return &Dispatcher{
		deliveriesRepo: deliveriesRepo,
		triggeredRepo:  triggeredRepo,
		endpointsRepo:  endpointsRepo,
		registry:       registry,
		workerID:       workerID,
		limit:          200,
		leaseSeconds:   300,
		maxAttempts:    10,
	}
}

func (d *Dispatcher) RunOnce(ctx context.Context) error {
	ds, err := d.deliveriesRepo.LeaseUnsent(ctx, d.workerID, d.limit, d.leaseSeconds)
	if err != nil {
		return err
	}
	if len(ds) == 0 {
		return nil
	}

	now := time.Now()

	for _, del := range ds {
		if del.Attempts >= d.maxAttempts {
			_ = d.deliveriesRepo.MarkFailed(ctx, del.ID, fmt.Sprintf("max attempts reached: %d", del.Attempts))
			continue
		}

		tr, err := d.triggeredRepo.GetByID(ctx, del.TriggeredID)
		if err != nil || tr == nil {
			_ = d.deliveriesRepo.MarkFailed(ctx, del.ID, "triggered event not found")
			continue
		}

		ep, err := d.endpointsRepo.GetByID(ctx, del.EndpointID)
		if err != nil || ep == nil || !ep.Enabled {
			_ = d.deliveriesRepo.MarkFailed(ctx, del.ID, "endpoint not found or disabled")
			continue
		}

		msg := FormatMessage(tr)

		ch := notify.Channel(ep.Channel)
		s, ok := d.registry.Get(ch)
		if !ok {
			_ = d.deliveriesRepo.MarkFailed(ctx, del.ID, "sender not registered for channel="+ep.Channel)
			continue
		}

		rcpt := notify.Recipient{ShopID: tr.ShopID}

		switch ch {
		case notify.ChannelTelegram:
			chatID, parseErr := strconv.ParseInt(ep.Address, 10, 64)
			if parseErr != nil {
				_ = d.deliveriesRepo.MarkFailed(ctx, del.ID, "invalid telegram chat_id")
				continue
			}
			rcpt.TelegramChatID = &chatID

		case notify.ChannelEmail:
			email := ep.Address
			rcpt.Email = &email

		case notify.ChannelWhatsApp:
			phone := ep.Address
			rcpt.WhatsAppPhone = &phone
		}

		if err := s.Send(ctx, rcpt, msg); err != nil {
			_ = d.deliveriesRepo.MarkFailed(ctx, del.ID, err.Error())
			continue
		}

		if err := d.deliveriesRepo.MarkSent(ctx, del.ID, now); err != nil {
			return err
		}
	}

	return nil
}

func NewSenderRegistry(senders []notify.Sender) *SenderRegistry {
	m := make(map[notify.Channel]notify.Sender, len(senders))
	for _, s := range senders {
		m[s.Channel()] = s
	}
	return &SenderRegistry{senders: m}
}

func (r *SenderRegistry) Get(ch notify.Channel) (notify.Sender, bool) {
	s, ok := r.senders[ch]
	return s, ok
}
