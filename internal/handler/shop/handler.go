package shop

import (
	"context"
	"net/http"
	"strconv"

	"salehelp/domain"
	auth_google_handler "salehelp/internal/handler/auth_google"

	"github.com/labstack/echo/v4"
)

type (
	Presenter interface {
		// shop(shops *domain.shops) *view.shopView
	}

	Service interface {
		SyncShop(ctx context.Context, shopID int) error
	}

	ShopCreator interface {
		CreateShop(ctx context.Context, shop *domain.Shop) error
	}

	ShopsRepo interface {
		GetByID(ctx context.Context, id int) (*domain.Shop, error)
		GetByExternalID(ctx context.Context, externalID string) (*domain.Shop, error)
		ListByUser(ctx context.Context, userID, limit, offset int) ([]*domain.Shop, error)
		Update(ctx context.Context, s *domain.Shop) error
		CountByUser(ctx context.Context, userID int) (int, error)
	}

	ProductRepo interface {
		ListByShop(ctx context.Context, shopID, limit, offset int) ([]*domain.Product, error)
		CountByShop(ctx context.Context, shopID int) (int, error)
	}

	SubscriptionsRepo interface {
		GetActiveByUserID(ctx context.Context, userID int) (*domain.Subscription, error)
	}

	Handler struct {
		service     Service
		shopCreator ShopCreator
		shopsRepo   ShopsRepo
		productRepo ProductRepo
		subsRepo    SubscriptionsRepo
	}
)

func NewHandler(service Service, shopCreator ShopCreator, shopsRepo ShopsRepo, productRepo ProductRepo, subsRepo SubscriptionsRepo) *Handler {
	return &Handler{
		service:     service,
		shopCreator: shopCreator,
		shopsRepo:   shopsRepo,
		productRepo: productRepo,
		subsRepo:    subsRepo,
	}
}

func (hand *Handler) getIDFromParam(c echo.Context) (int, error) {
	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil || id <= 0 {
		return 0, echo.ErrBadRequest
	}
	return id, nil
}

// ensureShopAccessByID проверяет JWT и что пользователь — владелец магазина (параметр пути "id"). Для роутов /shop/:id/...
func (hand *Handler) ensureShopAccessByID(c echo.Context) (int, error) {
	userID := auth_google_handler.UserIDFromContext(c.Request().Context())
	if userID == 0 {
		return 0, c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	shopID, err := hand.getIDFromParam(c)
	if err != nil {
		return 0, err
	}
	shop, err := hand.shopsRepo.GetByID(c.Request().Context(), shopID)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if shop == nil || shop.UserID != userID {
		return 0, c.JSON(http.StatusForbidden, map[string]string{"error": "access denied"})
	}
	return shopID, nil
}
