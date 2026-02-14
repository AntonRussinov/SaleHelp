package config

import (
	"fmt"
	"salehelp/internal/db"
	"salehelp/internal/scheduler/jobs"
	"salehelp/pkg/auth"

	"github.com/spf13/viper"
)

type (
	ServiceConfig struct {
		Listen              string
		Environment         string
		TelegramToken       string
		TelegramSecret      string
		TelegramBotUsername string
		TelegramAdminID     int64
		GormLogLevel        int
		ClientSecret        string
		AllowedOrigins      []string
		Auth                auth.Config
		DB                  db.Config
		InstanceID          string
		Scheduler           jobs.SyncActiveShopsConfig
		// WBReportTimezone — часовой пояс по умолчанию для отчётов WB (sales, orders). Любой IANA: Europe/Moscow, Asia/Almaty и т.д.
		// У магазина может быть свой пояс (shop.wb_report_timezone); если не задан — используется этот. Пусто = Europe/Moscow.
		WBReportTimezone string `mapstructure:"wb_report_timezone"`
		// AI — настройки AI-ответов на отзывы (OpenAI).
		AI AIConfig `mapstructure:"ai"`
		//	GCS            gcs.Config
	}
	AIConfig struct {
		// APIKey — ключ OpenAI API. Можно задать в конфиге или через env OPENAI_API_KEY (env имеет приоритет при пустом значении в конфиге).
		APIKey string `mapstructure:"api_key"`
		// ModelForPositive — модель для позитивных отзывов (рейтинг 5). По умолчанию gpt-4o-mini.
		ModelForPositive string `mapstructure:"model_positive"`
		// ModelForOther — модель для отзывов 1–4 звёзд. По умолчанию gpt-4o.
		ModelForOther string `mapstructure:"model_other"`
		// IncludedLimit — бесплатных AI-ответов в месяц по умолчанию (если не задано по подписке).
		IncludedLimit int `mapstructure:"included_limit"`
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
