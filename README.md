# Subscription Service

REST API для управления подписками пользователей.

## Стек

- Go
- Chi Router
- PostgreSQL
- pgx/v5
- Docker / Docker Compose
- golang-migrate
- Swagger
- Viper

## Конфигурация

Скопируйте `.env.example` в `.env` и при необходимости измените значения.

```env
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=subscriptions
DB_SSLMODE=disable
```

## Запуск локально

1. Поднимите PostgreSQL:

```bash
docker compose up -d
```

2. Примените миграции:

```bash
migrate \
    -path migrations \
    -database "postgres://postgres:postgres@localhost:5432/subscriptions?sslmode=disable" \
    up
```

3. Запустите приложение:

```bash
go run cmd/server/main.go
```

Сервис будет доступен по адресу:

```text
http://localhost:8080
```

## Запуск через Docker

Сначала убедитесь, что Docker Desktop или другой Docker daemon запущен, затем соберите образ:

```bash
docker build -t subscriptions-service:test .
```

## Swagger

После запуска сервиса документация доступна здесь:

```text
http://localhost:8080/swagger/index.html
```

## API

- `POST /subscriptions` — создать подписку
- `GET /subscriptions` — получить список подписок
- `GET /subscriptions/{id}` — получить подписку по ID
- `PUT /subscriptions/{id}` — обновить подписку
- `DELETE /subscriptions/{id}` — удалить подписку
- `GET /subscriptions/total` — подсчитать общую стоимость подписок

Параметры для `GET /subscriptions/total`:

| Параметр | Описание |
|----------|----------|
| user_id | UUID пользователя, опционально |
| service_name | Название сервиса, опционально |
| from | Начало периода в формате `MM-YYYY` |
| to | Конец периода в формате `MM-YYYY` |

Пример:

```text
GET /subscriptions/total?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&service_name=Yandex%20Plus&from=01-2025&to=12-2025
```

## Как запушить в Git

```bash
git status
git add README.md Dockerfile cmd/server/main.go internal/config/config.go internal/dto/create_subscription_response.go internal/handler/handler.go internal/middleware/logger.go internal/repository/subscription.go internal/service/subscription.go docs/docs.go docs/swagger.json docs/swagger.yaml
git commit -m "Fix docker build, api docs and shutdown"
git push origin main
```

Если `main` защищен, пушьте в свою ветку и создавайте Pull Request.

## Структура проекта

```text
cmd/server/
internal/config/
internal/dto/
internal/handler/
internal/middleware/
internal/models/
internal/repository/
internal/service/
migrations/
docs/
```

## Логирование

Используется structured logging через `slog`.

## Автор

Nazarbek Amanbek


