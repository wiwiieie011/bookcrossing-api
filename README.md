![Banner](https://images.unsplash.com/photo-1519682337058-a94d519337bc?q=80&w=1950&auto=format&fit=crop)

# Bookcrossing API ✨

![Go](https://img.shields.io/badge/Go-1.25%2B-00ADD8?logo=go) 
![Gin](https://img.shields.io/badge/Gin-1.11-00ADD8)
![GORM](https://img.shields.io/badge/GORM-ORM-blue) 
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?logo=postgresql) 
![Redis](https://img.shields.io/badge/Redis-7-DC382D?logo=redis)

Серверное REST API для обмена книгами: пользователи публикуют книги, инициируют обмены, пишут отзывы, авторизуются по JWT. Кэширование — Redis, БД — PostgreSQL, веб‑фреймворк — Gin.
Демо: [ссылка на демо] (если есть)

## Getting Started

Prerequisites:

- Go 1.25+
- Docker и Docker Compose (по желанию)
- PostgreSQL 16 и Redis 7 (локально или через Docker)
- Make (опционально), Air для hot‑reload (опционально)

Installation:

1. Клонируйте репозиторий и перейдите в папку проекта
2. Создайте .env с переменными окружения (пример ниже)

```env
PORT=1010
DB_HOST=localhost
DB_USER=postgres
DB_PASS=postgres
DB_NAME=bookcrossing
DB_PORT=5432
DB_SSLMODE=disable
REDIS_HOST=localhost
REDIS_PORT=6379
SUPER_SECRET_KEY=your-secret
```

Running:

- Через Docker

```bash
docker-compose up -d --build
```

Note: PORT в .env должен совпадать с портом в docker-compose.yml.

- Локально (Postgres и Redis уже запущены)

```bash
go run ./cmd/bookcrossing
# или
make run
```

Дополнительные команды:

```bash
make lint   # golangci-lint
make fmt    # форматирование
make vet    # статический анализ
make tidy   # очистка зависимостей
go test ./... -v
docker build -t bookcrossing/app .
```

## Contributors

- 👤 - wiwiieie011
- 👤 - dasler-fw
- 👤 - dzhambazbiev-ux
- 👤 - Bekkhanbs

## О технологиях

- Gin — минималистичный, быстрый HTTP‑фреймворк для Go.
- GORM — ORM над PostgreSQL с миграциями и ассоциациями.
- Redis — кэш ответов/данных и вспомогательные операции.
- JWT — аутентификация пользователей и защита маршрутов.

## Feedback

Откройте Issue: [ссылка на репозиторий]/issues или напишите: you@example.com.
