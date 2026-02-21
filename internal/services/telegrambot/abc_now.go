package telegrambot

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"salehelp/domain"
	"salehelp/internal/services/abc_analysis/notify"
	"salehelp/internal/services/subscription"
)


func (s *Service) handleABCNow(ctx context.Context, chatID int64, args []string) {
	u, err := s.ensureUserByChatID(ctx, chatID)
	if err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка пользователя.")
		return
	}

	ownerShops, _ := s.shopsRepo.ListByUser(ctx, u.ID, 0, 0)
	memberShops, _ := s.shopsRepo.ListByTelegramChat(ctx, chatID)
	shopsByID := make(map[int]*domain.Shop)
	for _, sh := range ownerShops {
		shopsByID[sh.ID] = sh
	}
	for _, sh := range memberShops {
		if _, ok := shopsByID[sh.ID]; !ok {
			shopsByID[sh.ID] = sh
		}
	}

	if len(shopsByID) == 0 {
		_ = s.tg.SendMessage(ctx, chatID, "У тебя пока нет магазинов. Подключи магазин через /addshop.")
		return
	}

	// /abc_now без id — подсказка
	if len(args) < 2 || args[1] == "" {
		var b strings.Builder
		b.WriteString("📊 ABC-анализ за последние 30 дней\n\n")
		b.WriteString("Напиши /abc_now <id>, где id — номер магазина:\n\n")
		for _, sh := range shopsByID {
			b.WriteString(fmt.Sprintf("• %s — /abc_now %d\n", sh.Name, sh.ID))
		}
		_ = s.tg.SendMessage(ctx, chatID, b.String())
		return
	}

	shopID, err := strconv.Atoi(strings.TrimSpace(args[1]))
	if err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Неверный номер магазина. Пример: /abc_now 1")
		return
	}

	shop, ok := shopsByID[shopID]
	if !ok {
		_ = s.tg.SendMessage(ctx, chatID, "Нет доступа к этому магазину. Используй id из списка /myshops.")
		return
	}

	lvl, _, _ := s.subEntitlements.GetShopAccessLevel(ctx, shopID)
	if lvl != subscription.AccessActive {
		_ = s.tg.SendMessage(ctx, chatID, "ABC-отчёт доступен при активной подписке. Продли подписку для этого магазина.")
		return
	}

	_ = s.tg.SendMessage(ctx, chatID, "Считаю ABC-анализ за 30 дней…")

	snap, err := s.abcSvc.SaveSnapshot(ctx, shopID, 30)
	if err != nil {
		s.log.Error("abc_now: save snapshot failed", "shop_id", shopID, "err", err)
		_ = s.tg.SendMessage(ctx, chatID, "Не удалось посчитать отчёт. Попробуй позже.")
		return
	}

	two, _ := s.abcSvc.GetTwoLatestSnapshots(ctx, shopID, 30)
	var prev *domain.ABCSnapshot
	if len(two) >= 2 {
		prev = two[1]
	}

	products, _ := s.productsRepo.ListByShop(ctx, shopID, 0, 0)
	productByNm := make(map[int64]*domain.Product)
	for _, p := range products {
		var nmID int64
		if n, err := parseNmID(p.ExternalID); err == nil {
			nmID = n
			productByNm[nmID] = p
		}
	}

	msg := notify.FormatMonthlyReportRU(snap, prev, productByNm)
	header := fmt.Sprintf("📊 ABC-анализ «%s» (за последние 30 дней)\n\n", shop.Name)
	_ = s.tg.SendMessage(ctx, chatID, header+msg)
}

func parseNmID(s string) (int64, error) {
	return strconv.ParseInt(strings.TrimSpace(s), 10, 64)
}
