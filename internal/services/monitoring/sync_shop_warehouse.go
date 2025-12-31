package monitoring

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/adapter/wb/dto"
)

func (s *Service) syncShopWarehouses(ctx context.Context, shop *domain.Shop) error {
	// сначала глобальные склады WB (их можно считать одинаковыми для всех магазинов)
	// if err := s.SyncWBWarehouses(ctx, shop); err != nil {
	// 	return err
	// }

	// потом склады именно этого продавца
	if err := s.syncSellerWarehouses(ctx, shop); err != nil {
		return err
	}

	return nil
}

func (s *Service) SyncWBWarehouses(ctx context.Context, shop *domain.Shop) error {
	whDTOs, err := s.wbClient.GetWBWarehouses(ctx, shop.APIToken)
	if err != nil {
		return fmt.Errorf("sync wb warehouses: %w", err)
	}
	if len(whDTOs) == 0 {
		return nil
	}

	for _, dtoWh := range whDTOs {
		wh := dtoWBWarehouseToDomain(dtoWh)

		if err := s.wbWarehousesRepo.Upsert(ctx, &wh); err != nil {
			return fmt.Errorf("sync wb warehouses: upsert wb_id=%d: %w", wh.WBID, err)
		}
	}

	return nil
}

// маппер DTO → domain
func dtoWBWarehouseToDomain(d dto.WBWarehouse) domain.WBWarehouse {
	// маленький хелпер для указателей
	ptrF := func(v float64) *float64 { return &v }
	ptrI := func(v int) *int { return &v }

	var lon, lat *float64
	var cargoType, deliveryType *int

	if d.Longitude != 0 {
		lon = ptrF(d.Longitude)
	}
	if d.Latitude != 0 {
		lat = ptrF(d.Latitude)
	}
	if d.CargoType != 0 {
		cargoType = ptrI(d.CargoType)
	}
	if d.DeliveryType != 0 {
		deliveryType = ptrI(d.DeliveryType)
	}

	return domain.WBWarehouse{
		WBID:            d.ID,
		Name:            d.Name,
		Address:         &d.Address,
		City:            &d.City,
		Longitude:       lon,
		Latitude:        lat,
		CargoType:       cargoType,
		DeliveryType:    deliveryType,
		FederalDistrict: &d.FederalDistrict,
		Selected:        &d.Selected,
	}
}

func (s *Service) syncSellerWarehouses(ctx context.Context, shop *domain.Shop) error {
	dtoWhs, err := s.wbClient.GetSellerWarehouses(ctx, shop.APIToken)
	if err != nil {
		return fmt.Errorf("sync seller warehouses: %w", err)
	}

	warehouses := make([]domain.SellerWarehouse, 0, len(dtoWhs))

	for _, d := range dtoWhs {
		warehouses = append(warehouses, dtoSellerWarehouseToDomain(shop.ID, d))
	}

	if err := s.sellerWarehousesRepo.ReplaceForShop(ctx, int64(shop.ID), warehouses); err != nil {
		return fmt.Errorf("sync seller warehouses: replace for shop %d: %w", shop.ID, err)
	}

	return nil
}

func dtoSellerWarehouseToDomain(shopID int, d dto.SellerWarehouse) domain.SellerWarehouse {
	ptrI := func(v int) *int { return &v }

	var cargoType, deliveryType *int
	if d.CargoType != 0 {
		cargoType = ptrI(d.CargoType)
	}
	if d.DeliveryType != 0 {
		deliveryType = ptrI(d.DeliveryType)
	}

	return domain.SellerWarehouse{
		ShopID:       shopID,
		WBOfficeID:   d.OfficeID,
		WBID:         d.ID,
		Name:         d.Name,
		CargoType:    cargoType,
		DeliveryType: deliveryType,
		IsDeleting:   d.IsDeleting,
		IsProcessing: d.IsProcessing,
	}
}
