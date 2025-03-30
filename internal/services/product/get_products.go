package product

import (
	"context"
	"errors"
	"salehelp/domain"
	"salehelp/internal/storage"
)

func (service *Service) GetAllProducts(ctx context.Context, merchantID string) ([]*domain.Product, error) {
	const fn = "ProductsService.GetProducts"

	// add FindProductsByID in repo
	// add search in portal API
	// think about how to call each of them in what cases.

	products, err := service.repo.FindAllByMerchID(ctx, merchantID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, notFoundError(merchantID, err)
		}

		return nil, service.handleError(err, fn, "failed to get product from repository")
	}

	return products, nil
}
