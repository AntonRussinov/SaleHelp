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
