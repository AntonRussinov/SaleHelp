package ai_settings

import (
	"context"
	"fmt"

	"salehelp/domain"
)

// LoadBundle загружает всё нужное для генерации промпта одним вызовом:
// настройки, правила бренда и шаблон нужной категории.
// Возвращает nil если для магазина нет настроек (AI не настроен).
func (r *Repository) LoadBundle(ctx context.Context, shopID int, category string) (*domain.AISettingsBundle, error) {
	settings, err := r.GetByShopID(ctx, shopID)
	if err != nil {
		return nil, fmt.Errorf("ai_settings load bundle: %w", err)
	}
	if settings == nil {
		return nil, nil
	}

	rules, err := r.ListRules(ctx, settings.ID)
	if err != nil {
		return nil, fmt.Errorf("ai_settings load bundle rules: %w", err)
	}

	tmpl, err := r.GetTemplate(ctx, settings.ID, category)
	if err != nil {
		return nil, fmt.Errorf("ai_settings load bundle template: %w", err)
	}

	return &domain.AISettingsBundle{
		Settings: settings,
		Rules:    rules,
		Template: tmpl,
	}, nil
}
