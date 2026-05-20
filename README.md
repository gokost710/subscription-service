# Subscription Service

REST-сервис для учета и агрегации онлайн-подписок пользователей.

Проект реализует CRUDL-операции над подписками и расчет суммарной стоимости подписок за выбранный период с фильтрацией по пользователю и названию сервиса.

## Возможности

- CRUDL для записей о подписках.
- Расчет суммарной стоимости подписок за период.
- Фильтрация по `user_id` и `service_name`.
- PostgreSQL в качестве хранилища.
- SQL-миграции для инициализации базы.
- Конфигурация через `.env`.
- JSON-логи через `slog`.
- Swagger UI.
- Запуск через Docker Compose.

## Стек

- Go
- Gin
- pgx
- PostgreSQL
- golang-migrate
- Docker Compose
- Swagger UI

## Быстрый старт

Создайте `.env` на основе примера:

```bash
cp .env.example .env
```

Запустите сервис:

```bash
docker compose up --build
```

После старта API будет доступен на:

```text
http://localhost:8080
```

Swagger UI:

```text
http://localhost:8080/swagger/index.html
```

OpenAPI JSON:

```text
http://localhost:8080/swagger/doc.json
```

Остановить сервис:

```bash
docker compose down
```

## Конфигурация

Пример переменных окружения находится в `.env.example`.

```env
HTTP_PORT=8080

DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=subscriptions
DB_SSL_MODE=disable

LOG_LEVEL=info
```

Для запуска через Docker Compose используйте:

```env
DB_HOST=postgres
```

Для локального запуска приложения без Docker Compose обычно нужен:

```env
DB_HOST=localhost
```

## API

### Healthcheck

```http
GET /health
```

### Создать подписку

```http
POST /subscriptions
```

```json
{
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025"
}
```

Поле `end_date` опционально:

```json
{
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025",
  "end_date": "12-2025"
}
```

### Получить подписку

```http
GET /subscriptions/{id}
```

### Получить список подписок

```http
GET /subscriptions
```

Поддерживаемые query-параметры:

- `user_id`
- `service_name`
- `limit`
- `offset`

Пример:

```http
GET /subscriptions?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&service_name=Yandex%20Plus&limit=10&offset=0
```

### Обновить подписку

```http
PUT /subscriptions/{id}
```

Тело запроса такое же, как при создании.

### Удалить подписку

```http
DELETE /subscriptions/{id}
```

### Посчитать суммарную стоимость

```http
GET /subscriptions/summary
```

Обязательные query-параметры:

- `from` в формате `MM-YYYY`
- `to` в формате `MM-YYYY`

Опциональные query-параметры:

- `user_id`
- `service_name`

Пример:

```http
GET /subscriptions/summary?from=07-2025&to=09-2025&user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&service_name=Yandex%20Plus
```

Ответ:

```json
{
  "total_price": 1200
}
```

## Примеры curl

Создать подписку:

```bash
curl -X POST http://localhost:8080/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }'
```

Получить список:

```bash
curl http://localhost:8080/subscriptions
```

Посчитать сумму:

```bash
curl "http://localhost:8080/subscriptions/summary?from=07-2025&to=09-2025&user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&service_name=Yandex%20Plus"
```

## Миграции

Миграции лежат в директории `migrations`.

При запуске через Docker Compose они применяются автоматически отдельным контейнером `migrate` после готовности PostgreSQL.

Текущая миграция создает таблицу `subscriptions` и индексы для фильтрации и расчета по периоду.

## Проверка

Юнит-тесты:

```bash
go test ./...
```

Статическая проверка:

```bash
go vet ./...
```

Локальная сборка:

```bash
mkdir -p bin
go build -o bin/subscription-service ./cmd/app
```

Полная проверка через Docker Compose:

```bash
docker compose up --build
```

## Формат даты

Даты подписок принимаются в формате:

```text
MM-YYYY
```

Пример:

```text
07-2025
```

В PostgreSQL дата хранится как `DATE` с первым днем месяца.
