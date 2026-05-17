package account

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type deleteAccountRequest struct {
	Confirmation string `json:"confirmation" validate:"required"`
}

// DeleteAccount permanently deletes the user account and all related data.
//
// @Summary     Удаление аккаунта
// @Description Необратимо удаляет аккаунт, магазины, AI-настройки и все данные пользователя
// @Tags        account
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body body account.deleteAccountRequest true "Подтверждение (значение должно быть 'DELETE')"
// @Success     204
// @Failure     400 {object} api.HTTPError
// @Failure     401 {object} api.HTTPError
// @Failure     500 {object} api.HTTPError
// @Router      /v1/account [delete]
func (h *Handler) DeleteAccount(c echo.Context) error {
	userID := userIDFromCtx(c)
	if userID == 0 {
		return echo.ErrUnauthorized
	}

	var req deleteAccountRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if req.Confirmation != "DELETE" {
		return echo.NewHTTPError(http.StatusBadRequest, "confirmation must be 'DELETE'")
	}

	ctx := c.Request().Context()

	if err := h.subs.ExpireActiveByUserID(ctx, userID); err != nil {
		h.log.Error("expire subscriptions on account delete", "err", err, "user_id", userID)
		return echo.ErrInternalServerError
	}

	if err := h.packages.ZeroOutByUserID(ctx, userID); err != nil {
		h.log.Error("zero out ai packages on account delete", "err", err, "user_id", userID)
		return echo.ErrInternalServerError
	}

	if err := h.users.DeleteByID(ctx, userID); err != nil {
		h.log.Error("delete user account", "err", err, "user_id", userID)
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusNoContent)
}
