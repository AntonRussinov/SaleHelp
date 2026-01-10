package subscriptionaddon

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
)

func (r *Repository) ListBySubscriptionID(
	ctx context.Context,
	subscriptionID int,
) ([]domain.SubscriptionAddon, error) {

	var models []model.SubscriptionAddon

	err := r.DB.WithContext(ctx).
		Where("subscription_id = ?", subscriptionID).
		Order("id ASC").
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf(
			"list subscription addons by subscription_id=%d: %w",
			subscriptionID,
			err,
		)
	}

	res := make([]domain.SubscriptionAddon, 0, len(models))
	for _, m := range models {
		res = append(res, *m.ToDomain())
	}

	return res, nil
}
