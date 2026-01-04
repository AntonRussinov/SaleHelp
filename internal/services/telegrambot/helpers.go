package telegrambot

import (
	"encoding/json"
	"salehelp/domain"
)

func encodeInts(xs []int) string {
	b, _ := json.Marshal(xs)
	return string(b)
}

func decodeInts(s string) ([]int, error) {
	var xs []int
	if s == "" {
		return xs, nil
	}
	err := json.Unmarshal([]byte(s), &xs)
	return xs, err
}

func extractShopIDs(shops []*domain.Shop) []int {
	ids := make([]int, 0, len(shops))
	for _, sh := range shops {
		ids = append(ids, sh.ID)
	}
	return ids
}
