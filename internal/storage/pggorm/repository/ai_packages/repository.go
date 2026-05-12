package ai_packages

import (
	"context"
	"fmt"
	"log/slog"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm"
	"salehelp/internal/storage/pggorm/model"
	"salehelp/pkg/logger/sl"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	pggorm.Repository
}

func NewRepository(log *slog.Logger, db *gorm.DB) *Repository {
	repo := &Repository{}
	repo.Log = log.With(sl.Module("AIPackagesRepository"))
	repo.DB = db
	return repo
}

// ConsumeOne atomically decrements remaining by 1 from the oldest package with remaining > 0.
// Returns (packageID, true, nil) on success, (0, false, nil) if no packages available.
func (r *Repository) ConsumeOne(ctx context.Context, userID int) (int, bool, error) {
	var pkgID int
	err := r.DB.WithContext(ctx).Raw(`
		UPDATE ai_packages SET remaining = remaining - 1
		WHERE id = (
			SELECT id FROM ai_packages
			WHERE user_id = ? AND remaining > 0
			ORDER BY purchased_at ASC
			LIMIT 1
			FOR UPDATE SKIP LOCKED
		)
		RETURNING id`, userID).Scan(&pkgID).Error
	if err != nil {
		return 0, false, fmt.Errorf("consume ai package: %w", err)
	}
	return pkgID, pkgID > 0, nil
}

// Create inserts a new AI package.
func (r *Repository) Create(ctx context.Context, p *domain.AIPackage) error {
	m := &model.AIPackage{
		UserID:      p.UserID,
		Size:        p.Size,
		Remaining:   p.Size,
		PriceMinor:  p.PriceMinor,
		PurchasedAt: time.Now(),
	}
	if err := r.DB.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("ai package create: %w", err)
	}
	p.ID = m.ID
	p.PurchasedAt = m.PurchasedAt
	return nil
}

// TotalRemaining returns the sum of remaining replies across all packages for a user.
func (r *Repository) TotalRemaining(ctx context.Context, userID int) (int, error) {
	var total int64
	err := r.DB.WithContext(ctx).Raw(
		"SELECT COALESCE(SUM(remaining), 0) FROM ai_packages WHERE user_id = ? AND remaining > 0", userID,
	).Scan(&total).Error
	if err != nil {
		return 0, fmt.Errorf("ai package total remaining: %w", err)
	}
	return int(total), nil
}
