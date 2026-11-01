package notification_settings

import (
	"context"
	"fmt"
	"salehelp/internal/notificationevents"
	"time"
)

type Repo interface {
	MapByShopID(ctx context.Context, shopID int) (map[string]bool, error)
	Upsert(ctx context.Context, shopID int, eventKey string, enabled bool, now time.Time) error
}

// Service provides effective (default + overrides) notification settings per shop.
type Service struct {
	repo Repo
	now  func() time.Time
}

func New(repo Repo) *Service {
	return &Service{repo: repo, now: time.Now}
}

func (s *Service) SetNow(fn func() time.Time) { s.now = fn }

func (s *Service) EffectiveEnabled(ctx context.Context, shopID int, eventKey string) (bool, error) {
	if shopID <= 0 {
		return false, fmt.Errorf("notification_settings: invalid shop_id=%d", shopID)
	}
	if eventKey == "" {
		return false, fmt.Errorf("notification_settings: empty event_key")
	}

	// Unknown keys are enabled by default (forward compatible).
	def, known := notificationevents.Find(eventKey)
	defaultEnabled := true
	if known {
		defaultEnabled = def.DefaultEnabled
	}

	m, err := s.repo.MapByShopID(ctx, shopID)
	if err != nil {
		return false, err
	}
	if v, ok := m[eventKey]; ok {
		return v, nil
	}
	return defaultEnabled, nil
}

func (s *Service) EffectiveMap(ctx context.Context, shopID int) (map[string]bool, error) {
	if shopID <= 0 {
		return nil, fmt.Errorf("notification_settings: invalid shop_id=%d", shopID)
	}

	// Start with defaults for known keys.
	out := make(map[string]bool, len(notificationevents.Catalog)+8)
	for _, e := range notificationevents.All() {
		out[e.Key] = e.DefaultEnabled
	}

	// Apply overrides from DB.
	m, err := s.repo.MapByShopID(ctx, shopID)
	if err != nil {
		return nil, err
	}
	for k, v := range m {
		if k != "" {
			out[k] = v
		}
	}

	return out, nil
}

func (s *Service) Set(ctx context.Context, shopID int, eventKey string, enabled bool) error {
	if shopID <= 0 {
		return fmt.Errorf("notification_settings: invalid shop_id=%d", shopID)
	}
	if eventKey == "" {
		return fmt.Errorf("notification_settings: empty event_key")
	}
	return s.repo.Upsert(ctx, shopID, eventKey, enabled, s.now())
}
