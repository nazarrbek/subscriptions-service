# subscriptions-service

````markdown
# Subscription Service

REST API сервис для управления подписками пользователей.

## Стек

- Go
- Chi Router
- PostgreSQL
- pgx/v5
- Docker Compose
- golang-migrate
- Swagger
- Viper

---

## Запуск проекта

### 1. Клонировать репозиторий

```bash
git clone https://github.com/nazarrbek/subscriptions-service.git
cd subscriptions-service
```

### 2. Создать .env

Создайте файл `.env` на основе `.env.example`

Пример:

```env
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=subscriptions
DB_SSLMODE=disable
```

### 3. Запустить PostgreSQL

```bash
docker compose up -d
```

Проверить контейнер:

```bash
docker ps
```

---

### 4. Применить миграции

```bash
migrate \
-path migrations \
-database "postgres://postgres:postgres@localhost:5432/subscriptions?sslmode=disable" \
up
```

---

### 5. Запустить приложение

```bash
go run cmd/server/main.go
```

Сервис будет доступен по адресу

```
http://localhost:8080
```

---

## Swagger

После запуска приложения документация доступна по адресу

```
http://localhost:8080/swagger/index.html
```

---

## API

### Создать подписку

```
POST /subscriptions
```

### Получить список подписок

```
GET /subscriptions
```

### Получить подписку по ID

```
GET /subscriptions/{id}
```

### Обновить подписку

```
PUT /subscriptions/{id}
```

### Удалить подписку

```
DELETE /subscriptions/{id}
```

### Подсчитать общую стоимость подписок

```
GET /subscriptions/total
```

Параметры:

| Параметр | Описание |
|----------|----------|
| user_id | UUID пользователя |
| service_name | Название сервиса |
| from | Начало периода (MM-YYYY) |
| to | Конец периода (MM-YYYY) |

Пример:

```
GET /subscriptions/total?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&service_name=Yandex%20Plus&from=01-2025&to=12-2025
```

---

## Структура проекта

```
cmd/
    server/

internal/
    config/
    dto/
    handler/
    middleware/
    models/
    repository/
    service/

migrations/

docs/
```

---

## Логирование

Реализовано middleware логирование HTTP-запросов.

Логируется:

- HTTP метод
- URL
- Время выполнения запроса

---

## Автор

Nazarbek Amanbek
````



Создай файл `.env.example`:

```env
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=subscriptions
DB_SSLMODE=disable
```


