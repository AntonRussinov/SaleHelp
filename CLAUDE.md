# salehelp

Go-бэкенд SaaS для продавцов на маркетплейсах (Wildberries).
Интеграции: AI-ответы на отзывы (OpenAI), Telegram-уведомления, Google OAuth, GCS.

## Команды

```bash
make run          # запуск локально (configs/config.yml)
make test         # все тесты
make lint         # golangci-lint (enable-all, очень строгий)
make fix          # автоисправление lint-ошибок
make swagger      # регенерация docs/ из cmd/swagger/main.go
make migrate-up   # применить миграции (goose)
make compose-db   # только PostgreSQL в Docker
```

## Архитектура

Слоёная архитектура с ручным DI через `Provider` (`internal/app/di/`):

```
Handler → Service → Repository → DB (GORM + PostgreSQL/pgx)
```

- Middleware парсят query-параметры (`cursor`, `pagination`, `filters`, `order`) и кладут в контекст
- Presenters форматируют ответы
- `Provider.Init()` — точка сборки всех зависимостей

## Ключевые пакеты

| Путь | Назначение |
|------|-----------|
| `internal/app/di/` | DI Provider, инициализация всех слоёв |
| `internal/handler/` | Echo-хэндлеры |
| `internal/services/` | Бизнес-логика |
| `internal/storage/pggorm/` | Репозитории (GORM) |
| `internal/adapter/` | Внешние интеграции (Kaspi, WB, Telegram, AI) |
| `internal/middleware/` | Парсинг параметров запроса |
| `internal/scheduler/` | Cron-задачи (robfig/cron) |
| `domain/` | Доменные типы и ошибки |
| `configs/config.yml` | Локальный конфиг |
| `db/migrations/` | SQL-миграции (goose) |

## Как добавить новый эндпоинт

1. Домен в `domain/` (если нужна новая сущность)
2. Репозиторий в `internal/storage/pggorm/`
3. Сервис в `internal/services/<feature>/`
4. Хэндлер в `internal/handler/<feature>/`
5. Зарегистрировать в `di/repositories.go`, `di/services.go`, `di/handlers.go`, `di/router.go`

## Моки

Генерируются через mockery (конфиг `.mockery.yaml`):

```bash
mockery
```

Моки кладутся в `mocks/`.

## Конфигурация

Читается через Viper. Структура в `internal/config/`. Для prod используется `configs/config.prod.yml`.
Чувствительные значения (JWT secret, API keys) передавать через env-переменные, не хардкодить в конфиг.

## Линтер

Используется `golangci-lint` с `enable-all: true`. Перед коммитом обязательно `make lint`. При добавлении кода следи за:
- длиной функций (`funlen`)
- именами переменных (`varnamelen`)
- порядком деклараций (`decorder`)

## GORM — критичные паттерны

**Delete** обязательно принимает модель с `TableName()` — иначе GORM не знает таблицу, ничего не удаляет, `RowsAffected = 0`:
```go
// ✅ правильно
r.DB.Delete(&model.Alert{}, "id = ? AND shop_id = ?", id, shopID)

// ❌ неправильно — удалит ноль строк без ошибки
r.DB.Where("id = ? AND shop_id = ?", id, shopID).Delete(&struct{}{})
```

После Delete всегда проверяй `result.RowsAffected == 0` → возвращай `gorm.ErrRecordNotFound`.

**GetByID** возвращает `nil, nil` (не ошибку) когда запись не найдена — вызывающий код обязан проверять nil:
```go
if err == gorm.ErrRecordNotFound { return nil, nil }  // паттерн в репозиториях
// в сервисе затем:
if result == nil { return nil, gorm.ErrRecordNotFound }
```

`First` — возвращает `ErrRecordNotFound` если нет строк. `Find` — всегда возвращает nil, даже если пусто.

## Обработка ошибок

В проекте **два пакета ошибок** (исторически):
- `domain/` — старый: `ServiceError`, `NotFoundError`, `ForbiddenError` и т.д.
- `domain/apperror/` — новый: `APPError`, `NotFoundError`, `ValidationError`

Оба обрабатываются в `internal/handler/error.go` (`handleError`). В новом коде предпочитай `apperror`.

Хэндлер возвращает `echo.NewHTTPError(code, msg)` только для простых случаев. Для доменных ошибок — возвращай `*apperror.APPError` или `*domain.ServiceError` — `handleError` сам проставит правильный HTTP-статус.

## Логирование в хэндлерах

Хэндлеры получают логгер через конструктор и тегируют его:
```go
Handler struct {
    log *slog.Logger
    // ...
}

func NewHandler(..., log *slog.Logger) *Handler {
    return &Handler{log: log.With("handler", "feature")}
}
```

`prov.log` передаётся из `di/handlers.go`. Request ID берётся из Echo:
```go
c.Response().Header().Get(echo.HeaderXRequestID)
```

В prod логи идут в JSON (`slog.NewJSONHandler`), в local — pretty. В GCP фильтруй по `jsonPayload.<field>`.

## Схема БД

- Все первичные ключи: `BIGSERIAL PRIMARY KEY` (int64 на уровне БД, `int` в Go-моделях)
- Временны́е метки: `TIMESTAMPTZ NOT NULL DEFAULT now()`
- JSONB-поля хранятся как `[]byte` в Go-моделях (`Payload []byte`)
- Внешние ключи с `ON DELETE CASCADE` — удаление магазина каскадно удаляет всё связанное

## Как добавить cron-задачу

1. Создать функцию `func(ctx context.Context) error` в `internal/scheduler/`
2. Зарегистрировать в `Provider` через `prov.scheduler.Register(ctx, spec, name, fn)`
3. Файл планировщика — `internal/app/di/sheduller.go` (опечатка в имени файла — это норма)

## Моки — добавление нового интерфейса

При добавлении нового интерфейса под мок нужно обновить `.mockery.yaml`, добавив пакет, затем запустить `mockery`. Мок-пакет называется `<packagename>_mocks`.

---

## WB API — лимиты (базовый токен, апрель 2026)

Используется token-bucket rate limiter (`golang.org/x/time/rate`) через `LimiterManager` в `internal/adapter/wb/client/rate_limiter.go`. Каждый API-ключ получает свои бакеты при регистрации через `RegisterKey(apiKey, tokenType)`.

| Категория (`RateCategory`) | Эндпоинты | Лимит (basic) | Spacing |
|---|---|---|---|
| `RateContent` | GET/POST карточки | 1/hr | ~1 час |
| `RateMarketplace` | GET заказы FBS/DBW/DBS | 10/hr | ~6 мин |
| `RateFBSOrders` | GET новые сборочные задания | 30/hr | ~2 мин |
| `RateDBOrders` | GET DB-заказы | 10/hr | ~6 мин |
| `RatePrices` | GET/POST цены | 1/hr | ~1 час |
| `RateStatistics` | GET остатки, продажи | 1/min | ~1 мин |
| `RateFeedbacks` | GET отзывы/вопросы | 5/hr | ~12 мин |
| `RateFeedbacksWrite` | POST ответы на отзывы/вопросы | 5/hr | ~12 мин |

Для **сервисного токена** (`TokenTypeService = "service"`) — отдельная политика `servicePolicies` с повышенными лимитами (WB официально не публикует точные значения).

Spacing реализован через `limiter.Wait(ctx)` — запрос блокируется до разрешения, не ретраит.

## Синк магазинов — архитектура

Планировщик: `internal/scheduler/jobs/sync_shops.go` → `SyncActiveShopsJob.Run`.

Каждый магазин: `mon.SyncShop` → `alerts.RunShop` (сразу после синка, в рамках одного тика).

### Fast sync (каждые 30 мин, `Spec: "0/30 * * * *"`)
`runFastSync` в `internal/services/monitoring/sync_shop.go`:
- FBS/DBW/DBS заказы (`RateFBSOrders` / `RateMarketplace`)
- Отзывы и вопросы (`RateFeedbacks`)
- Продажи и отказы (`RateStatistics`)

### Heavy sync (не чаще 1 раза в час, `heavySyncInterval = time.Hour`)
`runHeavySync` — только если `time.Since(LastHeavySyncAt) >= 1h`:
- Склады продавца
- Карточки товаров — инкрементально по курсору (`RateContent`, 1/hr)
- Полный синк карточек раз в 7 дней (`fullCardsSyncInterval`)
- Снапшоты цен и остатков (`RatePrices` + `RateStatistics`)

`ShopTimeout: 8m` — таймаут на один магазин.

## Типы токенов магазина

`domain.Shop.TokenType` — `"basic"` (по умолчанию) или `"service"`.

Устанавливается при создании/обновлении магазина (`api/shop.go`). Вызов `wbClient.SetTokenType(shop.APIToken, shop.TokenType)` в начале `SyncShop` применяет нужную rate-политику к ключу.

Миграция: `db/migrations/0021_shops_token_type.sql`.

## Восстановление цен — статусная машина

Таблица `price_restore_batches`, поля: `shop_id`, `restore_status`, `restore_started_at`, `restored_at`.

Переходы: `pending → restoring → restored | failed`.

Ключевые правила:
- `TrySetRestoring` — атомарный `UPDATE WHERE restore_status='pending'`. `RowsAffected == 0` = кто-то уже начал, возвращаем `AlreadyRestoring: true`.
- **Все** пути ошибок после захвата блокировки (`TrySetRestoring` вернул `true`) обязаны вызывать `SetRestoreFailed` — иначе батч зависнет в `restoring` навсегда.
- При успехе `SetRestoredAt` одновременно ставит `restore_status = 'restored'`.
- Повторный клик «Вернуть цены» при `status = failed` — показывает диалог подтверждения снова (позволяет ретрай).
- При успехе — broadcast во все включённые Telegram-эндпоинты магазина (`endpointsRepo.ListByShopID`).

Реализация: `internal/services/price_restore/`, репозиторий: `internal/storage/pggorm/repository/pricerestorebatch/`.

Миграция с индексами: `db/migrations/0022_price_restore_batch_status.sql`.

## Alert Engine — паттерны

`internal/services/alert/engine.go` → `RunShop(ctx, shopID)`:
1. Загружает активные алерты магазина
2. Загружает продукты, снапшоты (2 последних), размеры
3. Вычисляет: `stock_low`, `stock_zero`, `price_drop` (по снапшотам)
4. Отправляет уведомления об заказах FBS/DBW/DBS, продажах, отзывах, вопросах

Дедупликация алертов: `CreateIfNotExists` с `ON CONFLICT (alert_id, product_id, snapshot_id, product_size_key) DO NOTHING`.

### Ранняя детекция нулевого остатка (stock_zero from orders)

`internal/services/alert/stock_zero_orders.go` → `buildStockWarningByNmID`:
- Вызывается в `runShopFBSOrders` перед форматированием сообщения
- Считает pending FBS-заказы (`supplier_status='new'`) из DB — **без WB API**
- Если `pending_count >= snapshot.StockTotal > 0` → ставит флаг `StockWarning` на строку заказа
- В уведомлении о заказе добавляется `⚠️ последний товар!`
- Не создаёт отдельный `AlertTriggered` — нет дублей с последующим снапшотным алертом

## Уведомления — формат сообщений

### Дайджест алертов (`internal/services/notification/digest_ru.go`)

Секция `price_drop` — агрегированная, не per-item:
- Парсит `tr.Message` формата `"Цена упала: %.2f → %.2f (-%.2f, %.2f%%)"` через `fmt.Sscanf`
- Сортирует по % снижения (от большего к меньшему)
- Показывает до 15 товаров (`maxPriceDropShow = 15`)
- Строка `... и ещё N товаров` если > 15
- Итоговая строка: `Итого: X товаров, ср. снижение Y%, суммарно −Z ₽`

Секции `stock_zero` и `stock_low` — per-item с группировкой по размерам.

### Уведомления о заказах (`internal/services/notification/assembly_orders.go`)

`AssemblyOrderRow.StockWarning = true` → добавляет `⚠️ последний товар!` к строке заказа.
FBS: `FormatFBSOrdersMessageRU(orders, titleByNmID, stockWarningByNmID, loc)`.

### Telegram callback — кнопка «Вернуть цены назад»

Поток в `internal/services/telegrambot/handle_callback.go`:
1. Клик → `handleRestorePrices`: проверяет `RestoreStatus`:
   - `restoring` → «⏳ уже восстанавливается»
   - `restored` → «ℹ️ восстановлено в XX:XX»
   - `failed` → показывает диалог подтверждения (ретрай)
   - `pending` → показывает диалог подтверждения
2. «Да» → `handleConfirmRestore` → `priceRestore.Restore` → broadcast при успехе
