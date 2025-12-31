package config

import (
	"fmt"
	"salehelp/internal/db"
	"salehelp/pkg/auth"

	"github.com/spf13/viper"
)

type (
	ServiceConfig struct {
		Listen         string
		Environment    string
		TelegramToken  string
		GormLogLevel   int
		ClientSecret   string
		AllowedOrigins []string
		Auth           auth.Config
		DB             db.Config
		//	GCS            gcs.Config
	}

	Config struct {
		Environment string
		Service     *ServiceConfig
	}
)

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
