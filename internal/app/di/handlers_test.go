package di_test

import (
	"salehelp/internal/app/di"
	"salehelp/internal/config"
	"salehelp/pkg/logger/handlers/slogdiscard"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitHandlers(t *testing.T) {
	t.Parallel()

	prov := di.NewProvider(slogdiscard.NewDiscardLogger(), &config.ServiceConfig{})
	err := prov.InitHandlers()
	require.NoError(t, err)

	// Обязательные хэндлеры — всегда должны быть инициализированы
	assert.NotNil(t, prov.Handlers.Shop)
	assert.NotNil(t, prov.Handlers.TelegramBot)
	assert.NotNil(t, prov.Handlers.Product)
	assert.NotNil(t, prov.Handlers.Subscription)
	assert.NotNil(t, prov.Handlers.Alert)
	assert.NotNil(t, prov.Handlers.ShopAISettings)

	// Опциональные хэндлеры (AuthGoogle, NotificationSettings) зависят от внешних ключей,
	// не проверяем их без реального конфига.
}
