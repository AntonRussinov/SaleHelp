package product

// func (service *Service) GetAllProducts(ctx context.Context, merchantID string) ([]*domain.Product, error) {
// 	const fn = "ProductsService.GetProducts"

// 	// add FindProductsByID in repo
// 	// add search in portal API
// 	// think about how to call each of them in what cases.

// 	api := connector.New()

// 	products, err := api.GetAllProducts(ctx, "30345212")

// 	for _, data := range products.Data {

// 		// service.ImageRepo.Save(ctx, data.Images )

// 	}

// 	product, err := service.repo.FindAllByMerchID(ctx, merchantID)
// 	if err != nil {
// 		if errors.Is(err, storage.ErrNotFound) {
// 			return nil, notFoundError(merchantID, err)
// 		}

// 		return nil, service.handleError(err, fn, "failed to get product from repository")
// 	}

// 	return product, nil
// }
