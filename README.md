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
make run
# или
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

## Тестирование

```bash
make test
# или
go test ./... -v -count=1
```

## Линтинг

```bash
make lint
# или
go vet ./...
```

## Swagger

После запуска сервиса документация доступна здесь:

```text
http://localhost:8080/swagger/index.html
```

## API

### Создание подписки

`POST /subscriptions`

**Тело запроса:**
```json
{
  "service_name": "Netflix",
  "price": 999,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "01-2025",
  "end_date": "12-2025"
}
```

**Валидация:**
- `service_name` — обязательное, не пустое
- `price` — обязательное, >= 0
- `user_id` — обязательное, валидный UUID
- `start_date` — обязательное, формат `MM-YYYY`
- `end_date` — опционально, формат `MM-YYYY`

**Ответы:**
- `201 Created` — `{"id": "uuid"}`
- `400 Bad Request` — ошибка валидации
- `500 Internal Server Error`

---

### Список подписок (с пагинацией)

`GET /subscriptions?limit=10&offset=0`

| Параметр | Описание | По умолчанию |
|----------|----------|:------------:|
| `limit`  | Количество записей (макс. 100) | 10 |
| `offset` | Смещение от начала | 0 |

**Ответ:**
```json
{
  "data": [...],
  "total": 42,
  "limit": 10,
  "offset": 0
}
```

**Коды:**
- `200 OK`
- `400 Bad Request` — невалидные limit/offset
- `500 Internal Server Error`

---

### Получение подписки по ID

`GET /subscriptions/{id}`

**Ответы:**
- `200 OK` — объект подписки
- `400 Bad Request` — невалидный UUID
- `404 Not Found` — подписка не найдена

---

### Обновление подписки

`PUT /subscriptions/{id}`

**Тело запроса:**
```json
{
  "service_name": "Updated Name",
  "price": 1500,
  "start_date": "01-2025",
  "end_date": "12-2025"
}
```

**Ответы:**
- `200 OK` — обновлённый объект подписки
- `400 Bad Request` — ошибка валидации
- `404 Not Found` — подписка не найдена
- `500 Internal Server Error`

---

### Удаление подписки

`DELETE /subscriptions/{id}`

**Ответы:**
- `204 No Content`
- `400 Bad Request` — невалидный UUID
- `404 Not Found` — подписка не найдена

---

### Подсчёт стоимости подписок

`GET /subscriptions/total`

| Параметр | Описание |
|----------|----------|
| `user_id` | UUID пользователя, опционально |
| `service_name` | Название сервиса, опционально |
| `from` | Начало периода в формате `MM-YYYY` |
| `to` | Конец периода в формате `MM-YYYY` |

Пример:

```text
GET /subscriptions/total?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&service_name=Yandex%20Plus&from=01-2025&to=12-2025
```

## Формат ошибок

Все ошибки возвращаются в формате JSON:

```json
{
  "error": "описание ошибки"
}
```

## Структура проекта

```text
cmd/server/         — точка входа
internal/apperror/  — кастомные ошибки приложения
internal/config/    — конфигурация (Viper)
internal/dto/       — DTO (запросы, ответы, валидация)
internal/handler/   — HTTP-обработчики
internal/middleware/ — middleware (логирование)
internal/models/    — доменные модели
internal/repository/— слой работы с БД
internal/service/   — бизнес-логика
migrations/         — SQL-миграции
docs/               — Swagger-документация
```

## Логирование

Используется structured logging через `slog`.

## Автор

Nazarbek Amanbek
