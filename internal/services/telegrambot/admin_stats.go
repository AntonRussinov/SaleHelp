package telegrambot

import (
	"context"
	"fmt"
	"salehelp/internal/storage/pggorm/model"
	"strings"
	"time"
)

func (s *Service) handleAdminStats(ctx context.Context, chatID int64) {
	if s.adminChatID != 0 && chatID != s.adminChatID {
		_ = s.tg.SendMessage(ctx, chatID, "Неизвестная команда. Напиши /start")
		return
	}

	now := time.Now()

	// --- helpers ---
	formatMinor := func(amount int64) string {
		sign := ""
		if amount < 0 {
			sign = "-"
			amount = -amount
		}
		return fmt.Sprintf("%s%d.%02d", sign, amount/100, amount%100)
	}

	formatMoneyRows := func(rows []model.MoneyByCurrency) string {
		if len(rows) == 0 {
			return "- none\n"
		}
		var b strings.Builder
		for _, r := range rows {
			b.WriteString(fmt.Sprintf("- %s: %s (%d shops)\n", r.CurrencyCode, formatMinor(r.AmountMinor), r.Count))
		}
		return b.String()
	}

	// --- shops ---
	total, err := s.shopsRepo.CountAll(ctx)
	if err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка получения статистики (shops total).")
		return
	}

	active, err := s.shopsRepo.CountActive(ctx)
	if err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка получения статистики (shops active).")
		return
	}

	unauth401, err := s.shopsRepo.CountByLastSyncStatusCode(ctx, 401)
	if err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка получения статистики (401).")
		return
	}

	// --- subscriptions counters ---
	activeSubs, graceSubs, expiredSubs, err := s.subscriptionRepo.CountByState(ctx, now)
	if err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка получения статистики (subscriptions state).")
		return
	}

	// --- money in flight ---
	activeMoney, err := s.subscriptionRepo.SumRevenueActiveByCurrency(ctx, now)
	if err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка получения статистики (active revenue).")
		return
	}

	graceMoney, err := s.subscriptionRepo.SumRevenueGraceByCurrency(ctx, now)
	if err != nil {
		_ = s.tg.SendMessage(ctx, chatID, "Ошибка получения статистики (grace revenue).")
		return
	}

	// --- render ---
	var out strings.Builder
	out.WriteString("🛠 Admin stats\n\n")

	out.WriteString("🏪 Shops:\n")
	out.WriteString(fmt.Sprintf("- total: %d\n- active: %d\n\n", total, active))

	out.WriteString("💳 Subscriptions:\n")
	out.WriteString(fmt.Sprintf("- active: %d\n- grace: %d\n- expired: %d\n\n", activeSubs, graceSubs, expiredSubs))

	out.WriteString("💰 Money in flight:\n")
	out.WriteString("Active:\n")
	out.WriteString(formatMoneyRows(activeMoney))
	out.WriteString("\nGrace:\n")
	out.WriteString(formatMoneyRows(graceMoney))
	out.WriteString("\n")

	out.WriteString("⚠️ Errors:\n")
	out.WriteString(fmt.Sprintf("- 401 unauthorized shops: %d\n", unauth401))

	_ = s.tg.SendMessage(ctx, chatID, out.String())
}
