package subscription

import (
	"context"
	"fmt"
	"time"
)

// ExpireActiveByUserID marks all active/trial subscriptions for a user as expired.
func (r *Repository) ExpireActiveByUserID(ctx context.Context, userID int) error {
	return r.DB.WithContext(ctx).
		Exec("UPDATE subscriptions SET status='expired', updated_at=now() WHERE user_id=? AND status IN ('active','trial')", userID).
		Error
}

// IncrementAIUsed atomically increments ai_used by 1 if quota is not exhausted.
// Returns true if the increment succeeded (quota was available).
func (r *Repository) IncrementAIUsed(ctx context.Context, userID int) (bool, error) {
	res := r.DB.WithContext(ctx).
		Exec(`UPDATE subscriptions SET ai_used = ai_used + 1, updated_at = now()
		      WHERE user_id = ? AND status IN ('active','trial') AND expires_at > ? AND ai_used < ai_quota`,
			userID, time.Now())
	if res.Error != nil {
		return false, fmt.Errorf("increment ai_used: %w", res.Error)
	}
	return res.RowsAffected > 0, nil
}

// ResetMonthlyQuota resets ai_used=0 for all active subscriptions and advances ai_resets_at.
func (r *Repository) ResetMonthlyQuota(ctx context.Context, now time.Time) error {
	nextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
	return r.DB.WithContext(ctx).
		Exec("UPDATE subscriptions SET ai_used=0, ai_resets_at=?, updated_at=now() WHERE status='active'", nextMonth).
		Error
}
