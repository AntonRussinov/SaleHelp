package product

import (
	"context"
	"strconv"

	"salehelp/domain"
)

// ListRepo — минимальный контракт для выдачи списка товаров по магазину (реализует pggorm product.Repository).
type ListRepo interface {
	ListByShop(ctx context.Context, shopID int, limit, offset int) ([]*domain.Product, error)
}

// ListService реализует выдачу списка товаров по ID магазина (для handler/product).
type ListService struct {
	repo ListRepo
}

// NewListService создаёт сервис списка товаров по магазину.
func NewListService(repo ListRepo) *ListService {
	return &ListService{repo: repo}
}

const defaultListLimit = 100

// GetAllProducts возвращает товары по идентификатору магазина (merchantID = shop ID).
// Реализует интерфейс handler/product.Service.
func (s *ListService) GetAllProducts(ctx context.Context, merchantID string) ([]*domain.Product, error) {
	shopID, err := strconv.Atoi(merchantID)
	if err != nil || shopID <= 0 {
		return nil, domain.NewValidationErr("id must be a positive integer (shop id)")
	}
	return s.repo.ListByShop(ctx, shopID, defaultListLimit, 0)
}
