package price_restore

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/adapter/wb/dto"
	"strconv"
	"strings"
	"time"
)

// wbAlreadySetError returns true when WB responds that prices are already at the requested values.
// Treat this as success: the goal (prices at target level) is already achieved.
func wbAlreadySetError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "already set")
}

// Restore выполняет восстановление цен по пакету: загружает данные из предыдущего снапшота и шлёт в WB.
func (s *Service) Restore(ctx context.Context, batchID int) RestoreResult {
	batch, err := s.batchRepo.GetByID(ctx, batchID)
	if err != nil {
		return RestoreResult{Error: err}
	}
	if batch == nil {
		return RestoreResult{Error: fmt.Errorf("batch id=%d not found", batchID)}
	}
	if batch.RestoredAt != nil || batch.RestoreStatus == domain.RestoreStatusRestored {
		return RestoreResult{AlreadyRestored: true, RestoredAt: batch.RestoredAt}
	}

	won, err := s.batchRepo.TrySetRestoring(ctx, batchID)
	if err != nil {
		return RestoreResult{Error: err}
	}
	if !won {
		return RestoreResult{AlreadyRestoring: true}
	}

	// After acquiring the lock every error path must call SetRestoreFailed.
	failAndReturn := func(restoredCount int, wrapErr error) RestoreResult {
		_ = s.batchRepo.SetRestoreFailed(ctx, batchID)
		return RestoreResult{Error: wrapErr, RestoredCount: restoredCount}
	}

	triggeredIDs, err := s.batchRepo.GetTriggeredIDsByBatchID(ctx, batchID)
	if err != nil {
		return failAndReturn(0, err)
	}
	if len(triggeredIDs) == 0 {
		return failAndReturn(0, fmt.Errorf("batch has no triggered items"))
	}

	triggeredMap, err := s.triggered.GetByIDs(ctx, triggeredIDs)
	if err != nil {
		return failAndReturn(0, err)
	}
	var shopID int
	for _, tr := range triggeredMap {
		if tr != nil {
			shopID = tr.ShopID
			break
		}
	}
	if shopID == 0 {
		return failAndReturn(0, fmt.Errorf("shop_id not found for batch"))
	}

	shop, err := s.shopRepo.GetByID(ctx, shopID)
	if err != nil || shop == nil {
		return failAndReturn(0, fmt.Errorf("shop id=%d: %w", shopID, err))
	}
	apiKey := strings.TrimSpace(shop.APIToken)

	var taskGoods []dto.UploadTaskGood
	var taskSizeGoods []dto.UploadTaskSizeGood

	productIDs := make([]int, 0, len(triggeredMap))
	for _, tr := range triggeredMap {
		if tr == nil {
			continue
		}
		productIDs = append(productIDs, tr.ProductID)
	}
	productsByID, err := s.productRepo.GetByIDs(ctx, productIDs)
	if err != nil {
		return failAndReturn(0, err)
	}
	sizesByProduct, err := s.sizeRepo.ListByProductIDs(ctx, productIDs)
	if err != nil {
		return failAndReturn(0, err)
	}
	sizeByID := make(map[int]*domain.ProductSize)
	for _, sz := range sizesByProduct {
		if sz != nil {
			sizeByID[sz.ID] = sz
		}
	}

	for _, tr := range triggeredMap {
		if tr == nil || tr.AlertType != "price_drop" {
			continue
		}
		prevSnap, err := s.snapRepo.GetPreviousBySnapshotID(ctx, tr.SnapshotID)
		if err != nil {
			return failAndReturn(0, fmt.Errorf("previous snapshot for triggered %d: %w", tr.ID, err))
		}
		if prevSnap == nil {
			s.log.Warn("no previous snapshot for triggered", "triggered_id", tr.ID, "snapshot_id", tr.SnapshotID)
			continue
		}

		product := productsByID[tr.ProductID]
		if product == nil {
			s.log.Warn("product not found", "product_id", tr.ProductID)
			continue
		}
		nmID, err := strconv.ParseInt(product.ExternalID, 10, 64)
		if err != nil {
			s.log.Warn("invalid nmID", "product_id", tr.ProductID, "external_id", product.ExternalID)
			continue
		}

		editableSize := prevSnap.EditableSizePrice != nil && *prevSnap.EditableSizePrice
		if !editableSize {
			price := 0.0
			if prevSnap.BasePrice != nil {
				price = *prevSnap.BasePrice
			}
			discount := 0
			if prevSnap.Discount != nil {
				discount = int(*prevSnap.Discount)
			}
			taskGoods = append(taskGoods, dto.UploadTaskGood{NmID: nmID, Price: price, Discount: discount})
			continue
		}

		sizeSnaps, err := s.sizeSnapRepo.ListBySnapshotID(ctx, prevSnap.ID)
		if err != nil {
			return failAndReturn(0, fmt.Errorf("size snapshots for snapshot %d: %w", prevSnap.ID, err))
		}
		for _, ss := range sizeSnaps {
			if ss == nil {
				continue
			}
			sz := sizeByID[ss.ProductSizeID]
			if sz == nil {
				continue
			}
			sizeID, err := strconv.ParseInt(sz.ExternalID, 10, 64)
			if err != nil {
				continue
			}
			price := 0.0
			if ss.BasePrice != nil {
				price = *ss.BasePrice
			} else if ss.FinalPrice != nil {
				price = *ss.FinalPrice
			}
			taskSizeGoods = append(taskSizeGoods, dto.UploadTaskSizeGood{NmID: nmID, SizeID: sizeID, Price: price})
		}
	}

	restoredCount := 0
	if len(taskGoods) > 0 {
		for i := 0; i < len(taskGoods); i += maxGoodsPerRequest {
			end := i + maxGoodsPerRequest
			if end > len(taskGoods) {
				end = len(taskGoods)
			}
			chunk := taskGoods[i:end]
			_, err := s.wbClient.UploadPrices(ctx, apiKey, dto.UploadTaskRequest{Data: chunk})
			if err != nil && !wbAlreadySetError(err) {
				return failAndReturn(restoredCount, fmt.Errorf("wb upload prices: %w", err))
			}
			restoredCount += len(chunk)
		}
	}
	if len(taskSizeGoods) > 0 {
		for i := 0; i < len(taskSizeGoods); i += maxGoodsPerRequest {
			end := i + maxGoodsPerRequest
			if end > len(taskSizeGoods) {
				end = len(taskSizeGoods)
			}
			chunk := taskSizeGoods[i:end]
			_, err := s.wbClient.UploadPricesBySize(ctx, apiKey, dto.UploadTaskSizeRequest{Data: chunk})
			if err != nil && !wbAlreadySetError(err) {
				return failAndReturn(restoredCount, fmt.Errorf("wb upload prices by size: %w", err))
			}
			restoredCount += len(chunk)
		}
	}

	now := time.Now()
	if err := s.batchRepo.SetRestoredAt(ctx, batchID, now); err != nil {
		return failAndReturn(restoredCount, err)
	}
	return RestoreResult{ShopID: shopID, RestoredCount: restoredCount, RestoredAt: &now}
}
