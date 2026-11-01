package ai_settings

import (
	"context"
	"fmt"
	"time"

	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm"
)

func (r *Repository) ListRules(ctx context.Context, settingsID int) ([]*domain.AIRule, error) {
	var ms []model.AIRule
	err := r.DB.WithContext(ctx).
		Where("ai_settings_id = ?", settingsID).
		Order("order_index ASC, id ASC").
		Find(&ms).Error
	if err != nil {
		return nil, fmt.Errorf("ai_rules list: %w", err)
	}
	result := make([]*domain.AIRule, len(ms))
	for i := range ms {
		result[i] = ms[i].ToDomain()
	}
	return result, nil
}

func (r *Repository) CreateRule(ctx context.Context, d *domain.AIRule) error {
	m := &model.AIRule{
		AISettingsID: d.AISettingsID,
		RuleType:     d.RuleType,
		RuleText:     d.RuleText,
		OrderIndex:   d.OrderIndex,
		CreatedAt:    time.Now(),
	}
	if err := r.DB.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("ai_rules create: %w", err)
	}
	d.ID = m.ID
	d.CreatedAt = m.CreatedAt
	return nil
}

func (r *Repository) UpdateRule(ctx context.Context, d *domain.AIRule) error {
	err := r.DB.WithContext(ctx).
		Model(&model.AIRule{}).
		Where("id = ? AND ai_settings_id = ?", d.ID, d.AISettingsID).
		Updates(map[string]any{
			"rule_type":   d.RuleType,
			"rule_text":   d.RuleText,
			"order_index": d.OrderIndex,
		}).Error
	if err != nil {
		return fmt.Errorf("ai_rules update: %w", err)
	}
	return nil
}

func (r *Repository) DeleteRule(ctx context.Context, settingsID, ruleID int) error {
	result := r.DB.WithContext(ctx).
		Delete(&model.AIRule{}, "id = ? AND ai_settings_id = ?", ruleID, settingsID)
	if result.Error != nil {
		return fmt.Errorf("ai_rules delete: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
