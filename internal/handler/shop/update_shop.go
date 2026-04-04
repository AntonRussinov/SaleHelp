package shop

import (
	"net/http"
	"strings"
	"time"

	"salehelp/api"

	"github.com/labstack/echo/v4"
)

// UpdateShop обновляет настройки магазина (timezone и др.).
//
// @Summary     Обновить магазин
// @Description Обновление настроек магазина: часовой пояс отчётов WB. Требует JWT.
// @Tags        shops
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       shopId path int true "ID магазина"
// @Param       body   body api.UpdateShopRequest true "Поля для обновления"
// @Success     200 {object} api.ShopResponse
// @Failure     400 {object} api.HTTPError
// @Failure     401 {object} api.HTTPError
// @Failure     403 {object} api.HTTPError
// @Router      /v1/shops/{shopId} [patch]
func (h *Handler) UpdateShop(c echo.Context) error {
	shopID, err := h.ensureShopAccessByID(c)
	if err != nil {
		return err
	}

	var req api.UpdateShopRequest
	if err = c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}

	shop, err := h.shopsRepo.GetByID(c.Request().Context(), shopID)
	if err != nil || shop == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "shop not found"})
	}

	if req.WBReportTimezone != nil {
		tz := strings.TrimSpace(*req.WBReportTimezone)
		if tz != "" {
			if _, tzErr := time.LoadLocation(tz); tzErr != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid timezone: " + tz})
			}
		}
		shop.WBReportTimezone = tz
	}

	if err = h.shopsRepo.Update(c.Request().Context(), shop); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, shopToResponse(shop))
}
