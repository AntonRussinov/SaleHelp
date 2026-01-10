package monitoring

import (
	"context"
	"fmt"
	"time"
)

func (s *Service) CheckWBToken(ctx context.Context, token string) error {
	if token == "" {
		return fmt.Errorf("wb token is empty")
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := s.wbClient.CheckToken(ctx, token); err != nil {
		return err
	}

	return nil
}
