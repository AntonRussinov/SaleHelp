package alert

import (
	"encoding/json"
	"salehelp/domain"
)

func parseStockLowCfg(payload []byte) (stockLowCfg, error) {
	// дефолт
	cfg := stockLowCfg{Threshold: 5}

	// payload может быть nil/empty/"null"
	if len(payload) == 0 {
		return cfg, nil
	}

	// если в БД лежит literal null, json.Unmarshal([]byte("null"), &cfg) не упадёт,
	// но cfg останется дефолтным — это ок.
	if err := json.Unmarshal(payload, &cfg); err != nil {
		return stockLowCfg{}, err
	}

	if cfg.Threshold <= 0 {
		cfg.Threshold = 1
	}
	return cfg, nil
}

func buildEffectiveAlertsByProduct(
	alerts []*domain.Alert,
	alertType string,
) (global *domain.Alert, perProduct map[int]*domain.Alert) {

	perProduct = make(map[int]*domain.Alert)

	for _, a := range alerts {
		if a == nil {
			continue
		}
		if a.Type != alertType {
			continue
		}

		// status фильтруем тут, чтобы RunShop был чище
		if a.Status != "active" && a.Status != "new" {
			continue
		}

		if a.ProductID == nil {
			// один global на тип
			// если вдруг их несколько — можно выбрать последний/первый, но лучше валидировать на уровне API
			global = a
			continue
		}

		perProduct[*a.ProductID] = a
	}

	return global, perProduct
}

func derefInt(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}
