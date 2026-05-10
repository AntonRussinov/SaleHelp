package config

import (
	"fmt"
	"os"
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
		// WBXClientSecret — сервисный секрет WB для заголовка X-Client-Secret. Через env: WB_X_CLIENT_SECRET. Пока не получен — оставить пустым, запросы уйдут без заголовка.
		WBXClientSecret string `mapstructure:"wb_x_client_secret"`
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
		// ABCAnalysis — расписание и параметры ABC-анализа товаров.
		ABCAnalysis ABCAnalysisConfig `mapstructure:"abc_analysis"`
		// GoogleAuth — для входа через Google (фронт).
		GoogleAuth GoogleAuthConfig `mapstructure:"google_auth"`
		// JWTSecret — секрет для подписи JWT (выдача после успешного Google-логина). Через env: JWT_SECRET.
		JWTSecret string `mapstructure:"jwt_secret"`
		//	GCS            gcs.Config
	}
	GoogleAuthConfig struct {
		// ClientID — OAuth2 Client ID из Google Cloud Console (audience для проверки id_token).
		ClientID string `mapstructure:"client_id"`
	}
	ABCAnalysisConfig struct {
		WeeklySchedule    string `mapstructure:"weekly_schedule"`    // cron, напр. "0 23 * * 0" (воскресенье 23:00)
		MonthlySchedule   string `mapstructure:"monthly_schedule"`   // cron, напр. "0 1 1 * *" (1-го числа в 01:00)
		MinSalesThreshold int    `mapstructure:"min_sales_threshold"` // минимум продаж за период для отправки отчёта (по умолчанию 10)
		DefaultPeriodDays int    `mapstructure:"default_period_days"` // период по умолчанию для расчёта (30)
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

	// Env vars override config file values (for Cloud Run / Secret Manager).
	if v := os.Getenv("DATABASE_URL"); v != "" {
		cfg.Service.DB.ConnString = v
	}
	if v := os.Getenv("TELEGRAM_BOT_TOKEN"); v != "" {
		cfg.Service.TelegramToken = v
	}
	if v := os.Getenv("OPENAI_API_KEY"); v != "" {
		cfg.Service.AI.APIKey = v
	}
	if v := os.Getenv("JWT_SECRET"); v != "" {
		cfg.Service.JWTSecret = v
	}
	if v := os.Getenv("WB_X_CLIENT_SECRET"); v != "" {
		cfg.Service.WBXClientSecret = v
	}

	return &cfg, nil
}
