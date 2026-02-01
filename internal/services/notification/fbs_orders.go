package notifications

import "salehelp/domain"

// FormatFBSOrdersMessageRU формирует текст сообщения о новых FBS-заказах для отправки в Telegram.
// productTitleByNmID — название товара по nm_id (из карточки WB); если нет — используется артикул или "nm ID".
func FormatFBSOrdersMessageRU(orders []*domain.FBSOrder, productTitleByNmID map[int64]string) string {
	return FormatAssemblyOrdersMessageRU("FBS", fbsOrdersToRows(orders), productTitleByNmID)
}

func fbsOrdersToRows(orders []*domain.FBSOrder) []AssemblyOrderRow {
	rows := make([]AssemblyOrderRow, 0, len(orders))
	for _, o := range orders {
		rows = append(rows, AssemblyOrderRow{
			NmID:      o.NmID,
			Article:   o.Article,
			WBOrderID: o.WBOrderID,
			CreatedAt: o.CreatedAt,
		})
	}
	return rows
}
