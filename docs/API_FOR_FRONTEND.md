# API для фронта: эндпоинты и Google-авторизация

Документ описывает, какие эндпоинты нужны для фронта (создание пользователя через Google, магазины, алерты, AI-настройки) и чего не хватает в БД для Google-авторизации.

**Документация Swagger:** полное описание всех эндпоинтов (включая запросы/ответы и авторизацию) доступно в Swagger UI. В режиме разработки (не production) откройте в браузере **`/swagger/`**. Спецификация генерируется командой `make swagger` (файлы `docs/swagger.json`, `docs/swagger.yaml`).

---

## 1. Эндпоинты: что есть и что добавить

### Авторизация (Google)

| Метод | Путь | Назначение | Статус |
|-------|------|------------|--------|
| POST | `/api/v1/auth/google` | Вход по Google: тело `{ "idToken": "..." }`. Бэкенд проверяет токен, находит/создаёт user по `google_id` или `email`, возвращает JWT. | **Реализовано** |
| POST | `/api/v1/auth/register` | Регистрация по email и паролю: тело `{ "email", "password" }` (password не менее 8 символов). Возвращает JWT и user. При дубликате email — 409. | **Реализовано** |
| POST | `/api/v1/auth/login` | Вход по email и паролю: тело `{ "email", "password" }`. Возвращает JWT и user. | **Реализовано** |
| POST | `/api/v1/auth/refresh` | Обновление JWT (если используете refresh token). | По необходимости |

**Логика POST /auth/google (рекомендуемая):**

1. Фронт после успешного Google Sign-In присылает `id_token` (JWT от Google).
2. Бэкенд проверяет подпись id_token (Google JWKS), извлекает `sub` (google_id) и `email`.
3. Поиск пользователя: `GetByGoogleID( sub )` → если найден, возвращаем свои JWT/session.
4. Если не найден: `GetByEmail( email )` → если найден, привязываем к нему `google_id` (UPDATE), возвращаем JWT.
5. Если не найден ни по google_id, ни по email: создаём пользователя (email, google_id, без password), возвращаем JWT.

---

### Пользователь

| Метод | Путь | Назначение | Статус |
|-------|------|------------|--------|
| GET | `/api/v1/me` | Текущий пользователь (из JWT): id, email. Заголовок `Authorization: Bearer <token>`. | **Реализовано** |

Отдельный «создать пользователя» не нужен — пользователь создаётся при первом входе через Google (или при первом обращении в Telegram).

---

### Магазины

| Метод | Путь | Назначение | Статус |
|-------|------|------------|--------|
| GET | `/api/v1/shops` | Список магазинов текущего пользователя (owner). Query: `limit`, `offset`. Требует JWT. | **Реализовано** |
| POST | `/api/v1/shops` | Создать магазин: тело `{ "name": "...", "apiToken": "...", "marketplace": "wildberries" }`. Валидация токена WB, создание shop. Требует JWT. | **Реализовано** |
| GET | `/api/v1/shops/:id` | Детали магазина (название, статус, last_sync_at и т.д.). | Опционально |
| POST | `/api/v1/shop/:id/sync` | Запуск синка по магазину. | **Есть** (JWT + проверка владельца). |

**Важно:** все эндпоинты с `shopId`/`:id` проверяют, что магазин принадлежит текущему пользователю (JWT → `shop.UserID == userID`). Иначе — 401/403.

---

### Алерты

| Метод | Путь | Назначение | Статус |
|-------|------|------------|--------|
| POST | `/api/v1/shops/:shopId/alerts` | Создать алерт. Тело уже есть: `CreateAlertRequest` (type, severity, status, message, scope, payload). | **Есть** |
| GET | `/api/v1/shops/:shopId/alerts` | Список алертов с фильтрами (productId, type, status) и пагинацией. | **Есть** |
| GET | `/api/v1/shops/:shopId/alerts/:alertId` | Один алерт. | **Есть** |
| PUT | `/api/v1/shops/:shopId/alerts/:alertId` | Обновить алерт. | **Есть** |
| DELETE | `/api/v1/shops/:shopId/alerts/:alertId` | Удалить алерт. | **Есть** |

Проверка доступа включена: эндпоинты под защищённой группой (JWT). В хендлере проверяется, что текущий user — владелец магазина (`shop.UserID == userID`). При отсутствии/невалидном JWT — 401, при доступе к чужому магазину — 403.

---

### Настройки AI (автоответы)

| Метод | Путь | Назначение | Статус |
|-------|------|------------|--------|
| GET | `/api/v1/shops/:shopId/ai-settings` | Получить настройки AI магазина (toneOfVoice, brandName, customRules, isActive). | **Есть** |
| PUT | `/api/v1/shops/:shopId/ai-settings` | Создать/обновить настройки. Тело: `UpsertAISettingsRequest` (toneOfVoice, brandName, customRules, isActive). | **Есть** |

Проверка доступа включена: эндпоинты под JWT, в хендлере проверяется владелец магазина (`shop.UserID == userID`).

---

### Подписка

| Метод | Путь | Назначение | Статус |
|-------|------|------------|--------|
| GET | `/api/v1/shop/:id/subscription/status` | Статус подписки по магазину (активна/срок, лимиты чатов и т.д.). | **Есть** (JWT + владелец магазина). |
| POST | `/api/v1/shop/:id/subscription/grant` | Выдать/продлить подписку. | **Есть** (только роль **admin** в JWT). |

Итого по подписке: status — для владельца магазина; grant — только для пользователя с ролью `admin`. Роль хранится в таблице `users.role`, при логине попадает в JWT. Назначить себя админом: `UPDATE users SET role = 'admin' WHERE email = 'твой@email.com';` (после миграции `0016_users_role.sql`).

---

### Товары (для выбора в алертах)

| Метод | Путь | Назначение | Статус |
|-------|------|------------|--------|
| GET | `/api/v1/shops/:shopId/products` | Список товаров магазина. Query: `limit`, `offset`. Ответ: `{ "items", "total", "limit", "offset" }`. Требует JWT, проверка доступа к магазину. | **Реализовано** |

Рекомендуется единообразно: список товаров по магазину — `GET /shops/:shopId/products` с query-параметрами limit/offset (и при необходимости поиск по названию).

---

## 2. Сводка: что сделано

1. **Реализовано:**
   - `POST /api/v1/auth/google` — вход по Google id_token, выдача JWT.
   - `GET /api/v1/me` — текущий пользователь (заголовок `Authorization: Bearer <token>`).
   - `GET /api/v1/shops` — список магазинов пользователя (JWT).
   - `POST /api/v1/shops` — создание магазина (JWT, проверка WB-токена).
   - `GET /api/v1/shops/:shopId/products` — список товаров магазина с пагинацией (JWT, проверка доступа).
   - Health/версия: `GET /version` (или аналог) — в Swagger.

2. **Подписка (зарегистрированы и документированы в Swagger):**
   - `GET /api/v1/shop/:id/subscription/status` — статус подписки.
   - `POST /api/v1/shop/:id/subscription/grant` — выдача подписки (рекомендуется защитить только для админа).

3. **Swagger:** для всех перечисленных эндпоинтов (auth, shops, products, alerts, subscription, AI-настройки, health, telegram webhook) добавлена аннотированная документация; в dev UI доступен по `/swagger/`.

4. **Продукты:** основной способ — список по магазину `GET /api/v1/shops/:shopId/products`. Дополнительно в Swagger документирован `GET /api/v1/products` (query-параметр) для обратной совместимости.

5. **Защита доступа сделана:** все роуты с магазином (shops, products, alerts, ai-settings, subscription, sync) под middleware JWT и с проверкой владельца магазина (`shop.UserID == userID`). Участники по Telegram пока не проверяются — при необходимости добавить отдельно.

---

## 3. БД: хватает ли текущей структуры для Google-авторизации

**Текущая таблица `users`:**
- `id`, `telegram_chat_id`, `email`, `password_hash`, `created_at`, `updated_at`

**Для входа только по Google этого недостаточно:** нет поля, по которому однозначно идентифицировать пользователя после проверки Google id_token. У Google в id_token приходит `sub` (уникальный id пользователя в Google) и `email`.

**Что нужно добавить:**

| Поле | Тип | Назначение |
|------|-----|------------|
| `google_id` | `TEXT UNIQUE` (nullable) | Значение `sub` из Google id_token. По нему ищем/создаём пользователя при входе через Google. |

**Рекомендуемая миграция:**

```sql
-- +goose Up
ALTER TABLE users
  ADD COLUMN IF NOT EXISTS google_id TEXT UNIQUE;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id) WHERE google_id IS NOT NULL;
-- +goose Down
ALTER TABLE users DROP COLUMN IF EXISTS google_id;
```

**Логика:**

- Пользователь может войти только через Google (тогда `email` и `google_id` заполнены, `password_hash` и `telegram_chat_id` могут быть NULL).
- Или только через Telegram (как сейчас: `telegram_chat_id` есть, `google_id` NULL).
- Связка «один и тот же пользователь» при первом входе через Google и при первом заходе в Telegram — по **email**: если при Google-логине пришёл email, по которому уже есть user (например созданный позже через Telegram с тем же email), можно обновить этому user поле `google_id` и считать его тем же аккаунтом. Либо делать только привязку по google_id без слияния — тогда два входа (Google и Telegram) будут разными пользователями, пока не сделаете отдельный «привязать Telegram» в ЛК.

**Опционально (для отображения на фронте):**

- `name` / `display_name` — имя из Google (`name` в payload id_token).
- `avatar_url` — картинка профиля Google.

Итого: **текущей структуры БД недостаточно** — нужно добавить как минимум `google_id`. Остальное (email, password_hash, telegram_chat_id) уже есть и подходит для сценария «вход с фронта по Google, в боте — по Telegram с возможной привязкой по email».

---

## 4. Порядок внедрения (кратко)

**Уже сделано:** миграция `google_id`, репозиторий/сервис авторизации, роуты auth/shops/products/subscription/AI/алерты, документация Swagger (все эндпоинты, UI в dev по `/swagger/`, генерация — `make swagger`).

**Опционально для продакшена:** ограничить `POST /shop/:id/subscription/grant` только для админа (сейчас доступен владельцу магазина по JWT).

Фронт уже может: войти по Google, получить JWT, создать пользователя (неявно), создать магазин, настроить алерты и AI-автоответы, смотреть статус подписки. Актуальный список эндпоинтов и схем — в Swagger (`/swagger/`).
