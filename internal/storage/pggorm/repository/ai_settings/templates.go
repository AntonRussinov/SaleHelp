package ai_settings

import (
	"context"
	"fmt"
	"time"

	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
)

func (r *Repository) ListTemplates(ctx context.Context, settingsID int) ([]*domain.AITemplate, error) {
	var ms []model.AITemplate
	err := r.DB.WithContext(ctx).
		Where("ai_settings_id = ?", settingsID).
		Order("category ASC").
		Find(&ms).Error
	if err != nil {
		return nil, fmt.Errorf("ai_templates list: %w", err)
	}
	result := make([]*domain.AITemplate, len(ms))
	for i := range ms {
		result[i] = ms[i].ToDomain()
	}
	return result, nil
}

func (r *Repository) UpsertTemplate(ctx context.Context, d *domain.AITemplate) error {
	now := time.Now()
	m := &model.AITemplate{
		AISettingsID:      d.AISettingsID,
		Category:          d.Category,
		CustomPrompt:      d.CustomPrompt,
		ToneVoiceOverride: d.ToneVoiceOverride,
		Enabled:           d.Enabled,
		UpdatedAt:         now,
	}
	err := r.DB.WithContext(ctx).
		Where("ai_settings_id = ? AND category = ?", d.AISettingsID, d.Category).
		Assign(m).
		FirstOrCreate(m).Error
	if err != nil {
		return fmt.Errorf("ai_templates upsert: %w", err)
	}
	d.ID = m.ID
	d.CreatedAt = m.CreatedAt
	d.UpdatedAt = m.UpdatedAt
	return nil
}

// GetTemplate возвращает шаблон для конкретной категории или nil если не задан.
func (r *Repository) GetTemplate(ctx context.Context, settingsID int, category string) (*domain.AITemplate, error) {
	var ms []model.AITemplate
	err := r.DB.WithContext(ctx).
		Where("ai_settings_id = ? AND category = ?", settingsID, category).
		Limit(1).
		Find(&ms).Error
	if err != nil {
		return nil, fmt.Errorf("ai_templates get: %w", err)
	}
	if len(ms) == 0 {
		return nil, nil
	}
	return ms[0].ToDomain(), nil
}
