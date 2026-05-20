# Subscription Service Effective Mobile

REST-сервис для учета онлайн-подписок пользователей.

## Технологии

- Go
- Gin
- pgx
- PostgreSQL
- golang-migrate
- Docker Compose
- Swagger UI
- slog

## Запуск

Запустите сервис вместе с PostgreSQL и миграциями:

```bash
docker compose up --build
```

Docker Compose поднимает три сервиса:

- `postgres` - база данных PostgreSQL;
- `migrate` - контейнер, который применяет SQL-миграции;
- `app` - REST API.

API после запуска:

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

Если нужно удалить volume с данными PostgreSQL:

```bash
docker compose down -v
```

## Пример конфигурации

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

Для запуска через Docker Compose используется:

```env
DB_HOST=postgres
```

Для локального запуска приложения без Docker Compose обычно нужен:

```env
DB_HOST=localhost
```

Поддерживаемые уровни логирования:

```text
debug, info, warn, warning, error
```

## Модель подписки

Пример тела запроса на создание подписки:

```json
{
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025"
}
```

С датой окончания:

```json
{
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025",
  "end_date": "12-2025"
}
```

В API даты передаются как `MM-YYYY`. В PostgreSQL они хранятся как `DATE` с первым днем месяца.

## API

### Healthcheck

```http
GET /health
```

Ответ:

```json
{
  "status": "ok"
}
```

### Создать подписку

```http
POST /subscriptions
```

### Получить подписку по ID

```http
GET /subscriptions/{id}
```

### Получить список подписок

```http
GET /subscriptions
```

Query-параметры:

- `user_id` - фильтр по пользователю;
- `service_name` - фильтр по названию сервиса;
- `limit` - количество записей;
- `offset` - смещение.

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

- `from` - начало периода в формате `MM-YYYY`;
- `to` - конец периода в формате `MM-YYYY`.

Опциональные query-параметры:

- `user_id`;
- `service_name`.

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

Удалить подписку:

```bash
curl -X DELETE http://localhost:8080/subscriptions/1
```

## Миграции

Миграции находятся в директории `migrations`.

При запуске через Docker Compose миграции применяются автоматически отдельным контейнером `migrate` после того, как PostgreSQL пройдет healthcheck.

Текущая миграция создает таблицу `subscriptions`, ограничения и индексы для фильтрации и расчета по периоду.

## Проверка проекта

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
go build -o bin/subscription-service ./cmd/app
```

Проверка Docker Compose:

```bash
docker compose config
docker compose up --build
```

