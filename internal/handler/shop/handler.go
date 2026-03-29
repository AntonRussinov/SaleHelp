package shop

import (
	"context"
	"net/http"
	"strconv"
	"time"

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

	ShopsRepo interface {
		GetByID(ctx context.Context, id int) (*domain.Shop, error)
		ListByUser(ctx context.Context, userID, limit, offset int) ([]*domain.Shop, error)
		Create(ctx context.Context, s *domain.Shop) error
	}

	ProductRepo interface {
		ListByShop(ctx context.Context, shopID, limit, offset int) ([]*domain.Product, error)
		CountByShop(ctx context.Context, shopID int) (int, error)
	}

	SubscriptionRepo interface {
		GetActiveByShopID(ctx context.Context, shopID int, now time.Time) (*domain.Subscription, error)
	}

	WBTokenChecker interface {
		CheckWBToken(ctx context.Context, token string) error
	}

	Handler struct {
		service         Service
		shopsRepo       ShopsRepo
		productRepo     ProductRepo
		subscriptionRepo SubscriptionRepo
		wbCheck         WBTokenChecker
	}
)

func NewHandler(service Service, shopsRepo ShopsRepo, productRepo ProductRepo, wbCheck WBTokenChecker, subscriptionRepo SubscriptionRepo) *Handler {
	return &Handler{
		service:         service,
		shopsRepo:       shopsRepo,
		productRepo:     productRepo,
		subscriptionRepo: subscriptionRepo,
		wbCheck:         wbCheck,
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
