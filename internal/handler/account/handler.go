package account

import (
	"context"
	"log/slog"
	"salehelp/domain"

	auth_google "salehelp/internal/handler/auth_google"

	"github.com/labstack/echo/v4"
)

type SubscriptionsRepo interface {
	GetActiveByUserID(ctx context.Context, userID int) (*domain.Subscription, error)
	ExpireActiveByUserID(ctx context.Context, userID int) error
}

type ShopsRepo interface {
	CountByUser(ctx context.Context, userID int) (int, error)
}

type AIPackagesRepo interface {
	TotalRemaining(ctx context.Context, userID int) (int, error)
	ZeroOutByUserID(ctx context.Context, userID int) error
}

type UsersRepo interface {
	DeleteByID(ctx context.Context, id int) error
}

type Handler struct {
	log      *slog.Logger
	subs     SubscriptionsRepo
	shops    ShopsRepo
	packages AIPackagesRepo
	users    UsersRepo
}

func NewHandler(
	log *slog.Logger,
	subs SubscriptionsRepo,
	shops ShopsRepo,
	packages AIPackagesRepo,
	users UsersRepo,
) *Handler {
	return &Handler{
		log:      log.With("handler", "account"),
		subs:     subs,
		shops:    shops,
		packages: packages,
		users:    users,
	}
}

func userIDFromCtx(c echo.Context) int {
	return auth_google.UserIDFromContext(c.Request().Context())
}
