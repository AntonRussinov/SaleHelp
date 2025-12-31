package monitoring

import (
	"context"
	"fmt"
)

func (s *Service) GetLatestStateByShopID(ctx context.Context, shopID int) (map[int]*LatestProductState, error) {
	productIDs, err := s.productsRepo.ListIDsByShopID(ctx, shopID)
	if err != nil {
		return nil, fmt.Errorf("get latest state by shop: list product ids: %w", err)
	}

	stateByProductID, err := s.GetLatestStateByProductIDs(ctx, productIDs)
	if err != nil {
		return nil, fmt.Errorf("get latest state by shop: get latest state by product ids: %w", err)
	}

	return stateByProductID, nil
}
