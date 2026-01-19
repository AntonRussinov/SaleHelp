package notifications

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/pkg/notify"
	"time"
)

type Dispatcher struct {
	deliveriesRepo   DeliveriesRepo
	triggeredRepo    TriggeredRepo
	productsRepo     ProductsRepo
	productSizesRepo ProductSizesRepo
	endpointsRepo    EndpointsRepo

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
	ProductID  int
	AlertType  string
}

func NewDispatcher(
	deliveriesRepo DeliveriesRepo,
	triggeredRepo TriggeredRepo,
	productsRepo ProductsRepo,
	productSizesRepo ProductSizesRepo,
	endpointsRepo EndpointsRepo,
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
			ProductID:  tr.ProductID,
			AlertType:  tr.AlertType,
		}
		groups[k] = append(groups[k], del)
	}

	// --- send each group as one message
	for k, dels := range groups {
		ep, err := d.endpointsRepo.GetByID(ctx, k.EndpointID)
		if err != nil || ep == nil || !ep.Enabled {
			_ = d.deliveriesRepo.MarkFailedMany(ctx, deliveryIDs(dels), "endpoint not found or disabled")
			continue
		}

		// sender
		ch := notify.Channel(ep.Channel)
		sender, ok := d.registry.Get(ch)
		if !ok {
			_ = d.deliveriesRepo.MarkFailedMany(ctx, deliveryIDs(dels), "sender not registered for channel="+ep.Channel)
			continue
		}

		// triggered list
		trList := make([]*domain.AlertTriggered, 0, len(dels))
		for _, del := range dels {
			if tr := triggeredByID[del.TriggeredID]; tr != nil {
				trList = append(trList, tr)
			}
		}
		if len(trList) == 0 {
			_ = d.deliveriesRepo.MarkFailedMany(ctx, deliveryIDs(dels), "triggered events not found")
			continue
		}

		// product
		p := productsByID[k.ProductID]

		// message
		msg := FormatAggregatedMessageRU(k.AlertType, p, trList, sizesByID)

		// recipient
		rcpt, err := recipientFromEndpoint(ep, trList[0].ShopID)
		if err != nil {
			_ = d.deliveriesRepo.MarkFailedMany(ctx, deliveryIDs(dels), err.Error())
			continue
		}

		// send
		if err := sender.Send(ctx, rcpt, msg); err != nil {
			_ = d.deliveriesRepo.MarkFailedMany(ctx, deliveryIDs(dels), err.Error())
			continue
		}

		// mark sent
		if err := d.deliveriesRepo.MarkSentMany(ctx, deliveryIDs(dels), now); err != nil {
			return err
		}
	}

	return nil
}

func deliveryIDs(dels []*domain.AlertTriggeredDelivery) []int {
	ids := make([]int, 0, len(dels))
	for _, d := range dels {
		ids = append(ids, d.ID)
	}
	return ids
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
