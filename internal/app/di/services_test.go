package di_test

import (
	"salehelp/internal/app/di"
	"salehelp/internal/config"
	"salehelp/pkg/logger/handlers/slogdiscard"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitServices(t *testing.T) {
	t.Parallel()

	cfg := &config.ServiceConfig{
		// Минимальный конфиг чтобы telegram sender не паниковал
		TelegramToken: "1234567890:fake_token_for_test",
	}
	prov := di.NewProvider(slogdiscard.NewDiscardLogger(), cfg)
	prov.InitServices()

	// Обязательные сервисы — всегда должны быть инициализированы
	assert.NotNil(t, prov.Services.AlertService)
	assert.NotNil(t, prov.Services.Subscription)
	assert.NotNil(t, prov.Services.NotificationSettings)
	assert.NotNil(t, prov.Services.ProductList)
	assert.NotNil(t, prov.Services.Monitoring)
	assert.NotNil(t, prov.Services.Notifier)
	assert.NotNil(t, prov.Services.AlertEngine)
	assert.NotNil(t, prov.Services.AlertWorker)
	assert.NotNil(t, prov.Services.TelegramBot)
	assert.NotNil(t, prov.Services.SubscriptionMonitor)
	assert.NotNil(t, prov.Services.StockReplenish)
	assert.NotNil(t, prov.Services.PriceRestore)
	assert.NotNil(t, prov.Services.ABCAnalysis)

	// Опциональные сервисы (AIReply, AuthGoogle) — зависят от внешних ключей,
	// не проверяем их инициализацию без реального конфига.
}
