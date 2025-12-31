package mapper

import (
	"fmt"
	"strconv"
	"strings"

	"salehelp/domain"
	"salehelp/internal/adapter/wb/dto"
)

// productURLBase – например "https://www.wildberries.ru"
func CardToProduct(card dto.Card, shopID int, productURLBase string) domain.Product {
	urlBase := strings.TrimRight(productURLBase, "/")

	productURL := ""
	if card.NmID != 0 && urlBase != "" {
		productURL = fmt.Sprintf("%s/catalog/%d/detail.aspx", urlBase, card.NmID)
	}

	imageURL := ""
	if len(card.Photos) > 0 {
		imageURL = card.Photos[0].Big
	}

	return domain.Product{
		ShopID:      shopID,
		ExternalID:  strconv.Itoa(card.NmID), // можно заменить на card.NmUUID, если захочешь
		Title:       card.Title,
		URL:         productURL,
		ImageURL:    imageURL,
		Brand:       card.Brand,
		Category:    card.SubjectName,
		Description: card.Description,
	}
}
