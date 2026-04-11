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
