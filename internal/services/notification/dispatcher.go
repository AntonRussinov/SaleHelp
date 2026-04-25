package notifications

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/services/notification_settings"
	"salehelp/pkg/notify"
	"strconv"
	"time"
)

const callbackRestorePricesPrefix = "restore_prices:"

type PriceRestoreBatchRepo interface {
	Create(ctx context.Context, batch *domain.PriceRestoreBatch, triggeredIDs []int) (int, error)
}

type Dispatcher struct {
	deliveriesRepo   DeliveriesRepo
	triggeredRepo    TriggeredRepo
	productsRepo     ProductsRepo
	productSizesRepo ProductSizesRepo
	endpointsRepo    EndpointsRepo
	restoreBatchRepo PriceRestoreBatchRepo
	settingsSvc      *notification_settings.Service

	registry *SenderRegistry

	workerID     string
	limit        int
	leaseSeconds int
	maxAttempts  int
}

type SenderRegistry struct {
	senders map[notify.Channel]notify.Sender
}

type groupKey struct {
	EndpointID int
	ShopID     int
}

func NewDispatcher(
	deliveriesRepo DeliveriesRepo,
	triggeredRepo TriggeredRepo,
	productsRepo ProductsRepo,
	productSizesRepo ProductSizesRepo,
	endpointsRepo EndpointsRepo,
	restoreBatchRepo PriceRestoreBatchRepo,
	settingsSvc *notification_settings.Service,
	registry *SenderRegistry,
	workerID string,
) *Dispatcher {

	if workerID == "" {
		workerID = "worker-unknown"
	}

	return &Dispatcher{
		deliveriesRepo:   deliveriesRepo,
		triggeredRepo:    triggeredRepo,
		productsRepo:     productsRepo,
		productSizesRepo: productSizesRepo,
		endpointsRepo:    endpointsRepo,
		restoreBatchRepo: restoreBatchRepo,
		settingsSvc:      settingsSvc,
		registry:         registry,
		workerID:         workerID,
		limit:            200,
		leaseSeconds:     300,
		maxAttempts:      10,
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

	// --- collect triggered IDs
	trigIDs := make([]int, 0, len(ds))
	for _, del := range ds {
		if del.Attempts >= d.maxAttempts {
			_ = d.deliveriesRepo.MarkFailed(ctx, del.ID, fmt.Sprintf("max attempts reached: %d", del.Attempts))
			continue
		}
		trigIDs = append(trigIDs, del.TriggeredID)
	}

	triggeredByID, err := d.triggeredRepo.GetByIDs(ctx, trigIDs)
	if err != nil {
		return err
	}

	// --- collect productIDs + sizeIDs
	productSet := map[int]struct{}{}
	sizeSet := map[int]struct{}{}

	for _, del := range ds {
		tr := triggeredByID[del.TriggeredID]
		if tr == nil {
			continue
		}
		productSet[tr.ProductID] = struct{}{}
		if tr.ProductSizeID != nil {
			sizeSet[*tr.ProductSizeID] = struct{}{}
		}
	}

	productIDs := make([]int, 0, len(productSet))
	for id := range productSet {
		productIDs = append(productIDs, id)
	}
	sizeIDs := make([]int, 0, len(sizeSet))
	for id := range sizeSet {
		sizeIDs = append(sizeIDs, id)
	}

	productsByID, err := d.productsRepo.GetByIDs(ctx, productIDs)
	if err != nil {
		return err
	}

	sizesByID := map[int]*domain.ProductSize{}
	if len(sizeIDs) > 0 {
		sizesByID, err = d.productSizesRepo.GetByIDs(ctx, sizeIDs)
		if err != nil {
			return err
		}
	}

	// sizesCountByProduct: нужно, чтобы скрывать "фантомный" размер только когда он единственный
	sizesCountByProduct := map[int]int{}

	if len(productIDs) > 0 {
		allSizes, err := d.productSizesRepo.ListByProductIDs(ctx, productIDs)
		if err != nil {
			return err
		}
		for _, sz := range allSizes {
			if sz == nil {
				continue
			}
			sizesCountByProduct[sz.ProductID]++
		}
	}

	// --- group deliveries
	groups := make(map[groupKey][]*domain.AlertTriggeredDelivery, 64)

	for _, del := range ds {
		tr := triggeredByID[del.TriggeredID]
		if tr == nil {
			// triggered missing => fail this delivery
			_ = d.deliveriesRepo.MarkFailed(ctx, del.ID, "triggered event not found")
			continue
		}
		if tr.AlertType == "" {
			_ = d.deliveriesRepo.MarkFailed(ctx, del.ID, "missing alert_type in triggered")
			continue
		}

		k := groupKey{
			EndpointID: del.EndpointID,
			ShopID:     tr.ShopID,
		}
		groups[k] = append(groups[k], del)

	}

	// --- send each group as one message
	// --- send each group as one (or multiple chunked) messages
	for k, dels := range groups {
		ep, err := d.endpointsRepo.GetByID(ctx, k.EndpointID)
		if err != nil || ep == nil || !ep.Enabled {
			_ = d.deliveriesRepo.MarkFailedMany(ctx, deliveryIDs(dels), "endpoint not found or disabled")
			continue
		}

		ch := notify.Channel(ep.Channel)
		sender, ok := d.registry.Get(ch)
		if !ok {
			_ = d.deliveriesRepo.MarkFailedMany(ctx, deliveryIDs(dels), "sender not registered for channel="+ep.Channel)
			continue
		}

		// triggered list
		trList := make([]*domain.AlertTriggered, 0, len(dels))
		enabledDels := make([]*domain.AlertTriggeredDelivery, 0, len(dels))
		disabledIDs := make([]int, 0, 8)

		var enabledMap map[string]bool
		if d.settingsSvc != nil {
			shopID := trListShopID(triggeredByID, dels)
			if shopID > 0 {
				m, errM := d.settingsSvc.EffectiveMap(ctx, shopID)
				if errM == nil {
					enabledMap = m
				}
			}
		}

		for _, del := range dels {
			if tr := triggeredByID[del.TriggeredID]; tr != nil {
				enabled := true
				if enabledMap != nil {
					if v, ok := enabledMap[tr.AlertType]; ok {
						enabled = v
					}
				}
				if !enabled {
					disabledIDs = append(disabledIDs, del.ID)
					continue
				}
				trList = append(trList, tr)
				enabledDels = append(enabledDels, del)
			}
		}

		// If everything is disabled by settings, just mark as sent (skipped).
		if len(trList) == 0 && len(disabledIDs) > 0 {
			_ = d.deliveriesRepo.MarkSentMany(ctx, disabledIDs, now)
			continue
		}
		if len(trList) == 0 {
			_ = d.deliveriesRepo.MarkFailedMany(ctx, deliveryIDs(dels), "triggered events not found")
			continue
		}

		// recipient
		rcpt, err := recipientFromEndpoint(ep, trList[0].ShopID)
		if err != nil {
			_ = d.deliveriesRepo.MarkFailedMany(ctx, deliveryIDs(enabledDels), err.Error())
			continue
		}

		// build digest messages (chunked)
		msgs := FormatShopDigestMessagesRU(trList, productsByID, sizesByID, sizesCountByProduct)

		// Для Telegram при наличии price_drop: один батч на все такие триггеры группы,
		// кнопка «Вернуть цены назад» — только в последнем сообщении (при разбиении отчёта на несколько чанков).
		if ch == notify.ChannelTelegram && len(msgs) > 0 {
			var priceDropTriggeredIDs []int
			for _, tr := range trList {
				if tr != nil && tr.AlertType == "price_drop" {
					priceDropTriggeredIDs = append(priceDropTriggeredIDs, tr.ID)
				}
			}
			if len(priceDropTriggeredIDs) > 0 && d.restoreBatchRepo != nil {
				batch := &domain.PriceRestoreBatch{CreatedAt: now, ShopID: k.ShopID}
				batchID, errCreate := d.restoreBatchRepo.Create(ctx, batch, priceDropTriggeredIDs)
				if errCreate == nil {
					lastIdx := len(msgs) - 1
					msgs[lastIdx].InlineKeyboard = inlineKeyboardRestorePrices(batchID)
				}
			}
		}

		// send all parts; only mark sent if ALL succeed
		sentOK := true
		for _, msg := range msgs {
			if err := sender.Send(ctx, rcpt, msg); err != nil {
				_ = d.deliveriesRepo.MarkFailedMany(ctx, deliveryIDs(enabledDels), err.Error())
				if len(disabledIDs) > 0 {
					_ = d.deliveriesRepo.MarkSentMany(ctx, disabledIDs, now)
				}
				sentOK = false
				break
			}
		}
		if !sentOK {
			continue
		}

		// mark sent (enabled + disabled)
		if len(disabledIDs) > 0 {
			_ = d.deliveriesRepo.MarkSentMany(ctx, disabledIDs, now)
		}
		if err := d.deliveriesRepo.MarkSentMany(ctx, deliveryIDs(enabledDels), now); err != nil {
			return err
		}
	}

	return nil
}

func trListShopID(triggeredByID map[int]*domain.AlertTriggered, dels []*domain.AlertTriggeredDelivery) int {
	for _, del := range dels {
		if del == nil {
			continue
		}
		if tr := triggeredByID[del.TriggeredID]; tr != nil {
			return tr.ShopID
		}
	}
	return 0
}

func deliveryIDs(dels []*domain.AlertTriggeredDelivery) []int {
	ids := make([]int, 0, len(dels))
	for _, d := range dels {
		ids = append(ids, d.ID)
	}
	return ids
}

func inlineKeyboardRestorePrices(batchID int) map[string]any {
	return map[string]any{
		"inline_keyboard": [][]map[string]string{
			{{"text": "Вернуть цены назад", "callback_data": callbackRestorePricesPrefix + strconv.Itoa(batchID)}},
		},
	}
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
