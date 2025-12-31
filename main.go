package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	wbclient "salehelp/internal/adapter/wb/client"
	"salehelp/internal/adapter/wb/dto"
)

func main() {
	apiKey := "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJleGFtcGxlIn0.ZXhhbXBsZQ"
	if apiKey == "" {
		log.Fatal("WB_API_KEY is not set")
	}

	httpClient := &http.Client{
		Timeout: 15 * time.Second,
	}

	client := wbclient.New(httpClient)

	// минимальный запрос — можно потом усложнять
	req := dto.GetCardsRequest{
		Settings: dto.Settings{
			Cursor: dto.Cursor{
				Limit: 10,
			},
			Filter: dto.Filter{
				TextSearch:            "",
				AllowedCategoriesOnly: false,
				TagIDs:                []int{},
				ObjectIDs:             []int{},
				Brands:                []string{},
				ImtID:                 0,
				WithPhoto:             -1,
			},
		},
	}

	ctx := context.Background()
	resp, err := client.GetCards(ctx, apiKey, req)
	if err != nil {
		log.Fatalf("GetCards error: %v", err)
	}

	fmt.Printf("Total: %d\n", resp.Cursor.Total)
	fmt.Printf("Cards returned: %d\n", len(resp.Cards))

	for i, card := range resp.Cards {
		fmt.Printf("\n--- Card #%d ---\n", i+1)
		fmt.Printf("nmID: %d\n", card.NmID)
		fmt.Printf("vendorCode: %s\n", card.VendorCode)
		fmt.Printf("brand: %s\n", card.Brand)
		fmt.Printf("title: %s\n", card.Title)
		fmt.Printf("createdAt: %s\n", card.CreatedAt)
		if len(card.Sizes) > 0 {
			fmt.Printf("sizes[0].chrtID: %d\n", card.Sizes[0].ChrtID)
		}
		if i >= 4 { // чтобы не спамить консоль
			break
		}
	}
}
