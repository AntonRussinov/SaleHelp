package alert

import (
	"encoding/json"
	"fmt"
)

type priceDropCfg struct {
	Percent float64 `json:"percent"` // 7 = 7%
	Amount  int     `json:"amount"`  // в валюте (тенге и т.п.)
	Mode    string  `json:"mode"`    // "any" (по умолчанию) | "all" (опционально)
}

const (
	minPriceDropPercent = 0.5
	minPriceDropAmount  = 10
)

func parsePriceDropCfg(payload []byte) (priceDropCfg, error) {
	cfg := priceDropCfg{Mode: "any"}

	if len(payload) == 0 {
		return priceDropCfg{}, fmt.Errorf("empty payload")
	}

	if err := json.Unmarshal(payload, &cfg); err != nil {
		return priceDropCfg{}, err
	}

	if cfg.Mode == "" {
		cfg.Mode = "any"
	}
	if cfg.Mode != "any" && cfg.Mode != "all" {
		return priceDropCfg{}, fmt.Errorf("invalid mode: %s", cfg.Mode)
	}

	// нормализация
	if cfg.Percent < 0 {
		cfg.Percent = 0
	}
	if cfg.Amount < 0 {
		cfg.Amount = 0
	}

	// должен быть хотя бы один порог
	if cfg.Percent == 0 && cfg.Amount == 0 {
		return priceDropCfg{}, fmt.Errorf("percent and amount are both 0")
	}

	// защита от микро-колебаний
	if cfg.Percent > 0 && cfg.Percent < minPriceDropPercent {
		cfg.Percent = minPriceDropPercent
	}
	if cfg.Amount > 0 && cfg.Amount < minPriceDropAmount {
		cfg.Amount = minPriceDropAmount
	}

	return cfg, nil
}

func derefFloat64(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}
