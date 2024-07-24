package config

import (
	"fmt"
	"salehelp/internal/app"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string
	Service     app.Config
}

func Read(file string) (*Config, error) {
	const op = "config.Read"

	vip := viper.New()
	vip.SetConfigFile(file)

	if err := vip.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var cfg Config
	if err := vip.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &cfg, nil
}
