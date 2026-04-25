package notifications

import (
	"salehelp/domain"
	"time"
)

// FormatFBSOrdersMessageRU формирует текст сообщения о новых FBS-заказах для отправки в Telegram.
// productTitleByNmID — название товара по nm_id; stockWarningByNmID — nmID у которых заказы исчерпали остаток.
func FormatFBSOrdersMessageRU(orders []*domain.FBSOrder, productTitleByNmID map[int64]string, stockWarningByNmID map[int64]bool, loc *time.Location) string {
	return FormatAssemblyOrdersMessageRU("FBS", fbsOrdersToRows(orders, stockWarningByNmID), productTitleByNmID, loc)
}

func fbsOrdersToRows(orders []*domain.FBSOrder, stockWarningByNmID map[int64]bool) []AssemblyOrderRow {
	rows := make([]AssemblyOrderRow, 0, len(orders))
	for _, o := range orders {
		rows = append(rows, AssemblyOrderRow{
			NmID:         o.NmID,
			Article:      o.Article,
			WBOrderID:    o.WBOrderID,
			CreatedAt:    o.CreatedAt,
			StockWarning: stockWarningByNmID[o.NmID],
		})
	}
	return rows
}
