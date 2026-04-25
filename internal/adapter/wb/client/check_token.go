package client

import (
	"context"
	"fmt"
	"salehelp/internal/adapter/wb/dto"
)

func (c *Client) CheckToken(ctx context.Context, token string) error {
	req := dto.GetCardsRequest{
		Settings: dto.Settings{
			Cursor: dto.Cursor{
				Limit: 1,
			},
		},
	}

	_, _, err := c.GetAllCards(ctx, token, req, nil)
	if err != nil {
		return fmt.Errorf("wb token check failed: %w", err)
	}

	return nil
}
