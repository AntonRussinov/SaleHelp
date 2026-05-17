package account

import (
	"math"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type SubscriptionStatus struct {
	Tier         string    `json:"tier"`
	Status       string    `json:"status"`
	ExpiresAt    time.Time `json:"expires_at"`
	DaysRemaining int      `json:"days_remaining"`
	MaxShops     int       `json:"max_shops"`
	ExtraShops   int       `json:"extra_shops"`
	AIQuota      int       `json:"ai_quota"`
	AIUsed       int       `json:"ai_used"`
	AIResetsAt   time.Time `json:"ai_resets_at"`
}

type AIBalance struct {
	SubscriptionRemaining int `json:"subscription_remaining"`
	PackagesRemaining     int `json:"packages_remaining"`
	TotalAvailable        int `json:"total_available"`
}

type AccountStatusResponse struct {
	Subscription *SubscriptionStatus `json:"subscription"`
	AIBalance    AIBalance           `json:"ai_balance"`
	ShopsCount   int                 `json:"shops_count"`
	ShopsLimit   int                 `json:"shops_limit"`
}

// GetStatus returns account status: subscription, AI balance, shop counts.
//
// @Summary     Статус аккаунта
// @Description Текущая подписка, AI-баланс, количество и лимит магазинов
// @Tags        account
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} account.AccountStatusResponse
// @Failure     401 {object} api.HTTPError
// @Failure     500 {object} api.HTTPError
// @Router      /v1/account/status [get]
func (h *Handler) GetStatus(c echo.Context) error {
	ctx := c.Request().Context()
	userID := userIDFromCtx(c)
	if userID == 0 {
		return echo.ErrUnauthorized
	}

	sub, err := h.subs.GetActiveByUserID(ctx, userID)
	if err != nil {
		h.log.Error("get active subscription", "err", err, "user_id", userID)
		return echo.ErrInternalServerError
	}

	shopsCount, err := h.shops.CountByUser(ctx, userID)
	if err != nil {
		h.log.Error("count shops by user", "err", err, "user_id", userID)
		return echo.ErrInternalServerError
	}

	packagesRemaining, err := h.packages.TotalRemaining(ctx, userID)
	if err != nil {
		h.log.Error("ai packages total remaining", "err", err, "user_id", userID)
		return echo.ErrInternalServerError
	}

	resp := AccountStatusResponse{
		ShopsCount: shopsCount,
	}

	if sub != nil {
		daysRemaining := int(math.Ceil(time.Until(sub.ExpiresAt).Hours() / 24))
		if daysRemaining < 0 {
			daysRemaining = 0
		}
		resp.Subscription = &SubscriptionStatus{
			Tier:          sub.Tier,
			Status:        sub.Status,
			ExpiresAt:     sub.ExpiresAt,
			DaysRemaining: daysRemaining,
			MaxShops:      sub.MaxShops,
			ExtraShops:    sub.ExtraShops,
			AIQuota:       sub.AIQuota,
			AIUsed:        sub.AIUsed,
			AIResetsAt:    sub.AIResetsAt,
		}
		resp.AIBalance = AIBalance{
			SubscriptionRemaining: sub.AIRemaining(),
			PackagesRemaining:     packagesRemaining,
			TotalAvailable:        sub.AIRemaining() + packagesRemaining,
		}
		resp.ShopsLimit = sub.MaxShops + sub.ExtraShops
	} else {
		resp.AIBalance = AIBalance{
			PackagesRemaining: packagesRemaining,
			TotalAvailable:    packagesRemaining,
		}
	}

	return c.JSON(http.StatusOK, resp)
}
