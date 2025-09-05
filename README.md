# 🧩 TaskMesh — микросервисная платформа (Core Service на Go)

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15%2B-336791?logo=postgresql)](https://www.postgresql.org/)
[![Build](https://img.shields.io/badge/tests-go%20test-success)](#-тестирование)
[![gRPC](https://img.shields.io/badge/gRPC-50051-8A2BE2?logo=grpc)](#-grpc-task-service--кратко)
[![Telegram](https://img.shields.io/badge/Telegram-bot-26A5E4?logo=telegram)](#-notification-service-привязка-telegram--taskmesh)

**TaskMesh** — современная микросервисная платформа для управления задачами, уведомлениями и рейтингами пользователей. Система автоматизирует постановку задач, контроль сроков и уведомления исполнителей, а также рассчитывает рейтинг по результатам работы.

---

## 🧱 Микросервисы (коротко)

1) **Task Service (ядро)**
   - CRUD пользователей и задач (REST), статусы `todo/doing/done`, дедлайны, приоритеты.
   - Экспортирует **gRPC** метод проверки пользователя по email.
   - Хранение: PostgreSQL. Покрыт юнит-тестами и моками.

2) **Notification Service**
   - Принимает Telegram **webhook** и отправляет ответы в чат.
   - Привязывает **email пользователя** к **Telegram chat_id** и проверяет существование пользователя через gRPC к Task Service.
   - В будущем получает события через Kafka для асинхронных уведомлений.

3) **Rating Service (в планах)**
   - Подсчёт рейтинга: +за раннее/своевременное выполнение, −за просрочку.
   - Влияет на доступ к задачам и приоритеты. Связи: gRPC + Kafka.

---

## 🛠️ Технологии

- **Go 1.24+**, `net/http`, `context`, graceful shutdown.
- **REST + JSON** (Task Service), **gRPC** (`google.golang.org/grpc`).
- **PostgreSQL 15+**, миграции в `/internal/**/db/migrations`.
- **Config**: `.env` через `godotenv`.
- Тесты: `go test ./...`, моки `gomock`.
- Дальше: Docker Compose, Kafka, JWT, CI/CD, Prometheus/Grafana.

---

## ⚙️ Установка и запуск

### 1) Клонирование
```bash
git clone https://github.com/HDBOOMONE12/TaskManager.git
cd TaskManager
```

### 2) Конфигурация окружения

**Task Service (`cmd/taskmanager/.env`)**
```dotenv
DATABASE_URL=postgres://USER:PASSWORD@localhost:5432/taskmanager?sslmode=disable
```

**Notification Service (`cmd/notification-service/.env`)**
```dotenv
TELEGRAM_BOT_TOKEN=123456:PUT_YOUR_TOKEN_HERE
PORT=8081
DATABASE_URL=postgres://USER:PASSWORD@localhost:5432/notify_db?sslmode=disable
TASK_SERVICE_GRPC_URL=localhost:50051
```
> Секреты не коммитим.

### 3) Миграции БД

**Task Service**
```bash
psql "$DATABASE_URL" -f internal/taskmanager/db/migrations/0001_create_users.sql
psql "$DATABASE_URL" -f internal/taskmanager/db/migrations/0002_create_tasks.sql
psql "$DATABASE_URL" -f internal/taskmanager/db/migrations/0003_indexes.sql
```

**Notification Service**
```bash
psql "$DATABASE_URL" -f internal/notification-service/db/migrations/0001_telegram_bindings.sql
```

### 4) Запуск

**Task Service (HTTP :8080, gRPC :50051)**
```bash
go run ./cmd/taskmanager
```

**Notification Service (HTTP :8081)**
```bash
go run ./cmd/notification-service
```

---

## 📡 API (REST, Task Service)

Формат ошибки:
```json
{ "error": "message" }
```

### Пользователи

`GET /users` — список (поддерживает `?email=`).  
`POST /users` — создать:
```json
{ "name": "alice", "email": "alice@example.com" }
```
`GET /users/{id}` — получить.  
`PATCH /users/{id}` — частичное обновление:
```json
{ "name": "Alice Cooper", "email": "alice@company.com" }
```
`DELETE /users/{id}` — удалить.

### Задачи (на пользователя)

Статусы: `todo | doing | done`. Приоритет: `0..5` (по умолчанию `3`). `due_at` — ISO8601.

`GET /users/{user_id}/tasks` — список.  
`POST /users/{user_id}/tasks` — создать:
```json
{
  "title": "Сделать отчёт",
  "description": "К пятнице",
  "status": "todo",
  "priority": 3,
  "due_at": "2025-09-05T18:00:00Z"
}
```
`GET /users/{user_id}/tasks/{task_id}` — получить.  
`PUT /users/{user_id}/tasks/{task_id}` — полное обновление.  
`PATCH /users/{user_id}/tasks/{task_id}` — частичное обновление.  
`DELETE /users/{user_id}/tasks/{task_id}` — удалить.

---

## 🔌 gRPC (Task Service) — кратко

- Порт: **:50051**
- Proto: `internal/taskmanager/proto/user.proto`
- Метод: `UserService.HasUserWithEmail(EmailRequest) → UserExistsResponse`

Генерация:
```bash
protoc --go_out=. --go-grpc_out=. internal/taskmanager/proto/user.proto
```

---

## 🔔 Notification Service: привязка Telegram ↔ TaskMesh

**Назначение:** мгновенные уведомления и идентификация пользователя по email.

**Как привязать аккаунт (5 шагов):**
1. Найдите своего бота в Telegram и отправьте `/start`.
2. Бот попросит email — пришлите адрес, который зарегистрирован в Task Service.
3. Notification Service вызывает gRPC `HasUserWithEmail` в Task Service и проверяет наличие пользователя.
4. Если пользователь найден — создаётся запись в БД: `telegram_bindings(email, chat_id)`, бот ответит «Привязка выполнена».
5. Готово: при создании/изменении задач (и последующей интеграции с Kafka) вы будете получать уведомления в этот чат.

**Webhook**
- Endpoint: `POST /webhook` (ожидает Telegram Update).
- Установка webhook (пример с ngrok):
  ```bash
  ngrok http 8081   # либо через ngrok-config.yml
  curl -X POST "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/setWebhook"        -d "url=https://<PUBLIC_HTTPS_URL>/webhook"
  ```
- Сброс:
  ```bash
  curl -X POST "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/deleteWebhook"
  ```

---

## 🗂️ Структура проекта
```text
.
├── cmd/
│   ├── notification-service/
│   │   ├── .env.example
│   │   ├── env.go
│   │   └── main.go
│   └── taskmanager/
│       ├── .env.example
│       ├── env.go
│       └── main.go
├── internal/
│   ├── notification-service/
│   │   ├── db/
│   │   │   ├── migrations/
│   │   │   │   └── 0001_telegram_bindings.sql
│   │   │   └── postgres.go
│   │   ├── entity/
│   │   │   └── telegram.go
│   │   ├── handlers/
│   │   │   ├── helpers.go
│   │   │   └── webhook.go
│   │   ├── notifyerrors/
│   │   │   └── errors.go
│   │   ├── senders/
│   │   │   ├── payload.go
│   │   │   └── telegram.go
│   │   ├── service/
│   │   │   └── binding_service.go
│   │   ├── storage/
│   │   │   └── telegram_binding_repo.go
│   │   └── taskclient/
│   │       ├── client.go
│   │       └── grpc.go
│   └── taskmanager/
│       ├── db/
│       │   ├── migrations/
│       │   │   ├── 0001_create_users.sql
│       │   │   ├── 0002_create_tasks.sql
│       │   │   └── 0003_indexes.sql
│       │   └── postgres.go
│       ├── entity/
│       │   ├── task.go
│       │   └── user.go
│       ├── grpcs/
│       │   └── server.go
│       ├── handlers/
│       │   ├── errors_tasks.go
│       │   ├── helpers.go
│       │   ├── router_users.go
│       │   ├── tasks.go
│       │   └── users.go
│       ├── mocks/
│       │   └── mock_task_repo.go
│       ├── proto/
│       │   ├── user.pb.go
│       │   ├── user.proto
│       │   └── user_grpc.pb.go
│       ├── service/
│       │   ├── task_test.go
│       │   ├── tasks.go
│       │   └── users.go
│       └── storage/
│           ├── tasks_repo.go
│           └── users_repo.go
├── .gitignore
├── README.md
├── go.mod
└── go.sum
```

---

## 🧪 Тестирование
```bash
go test ./...
```
- Юнит-тесты: `internal/taskmanager/service/task_test.go`
- Моки: `mockgen` для репозиториев.

---

## 🧱 Схема БД (миграции)

**users**
- `id BIGINT IDENTITY PRIMARY KEY`
- `username VARCHAR(50) UNIQUE NOT NULL`
- `email VARCHAR(50) UNIQUE NOT NULL`
- `created_at`, `updated_at`

**tasks**
- `id BIGINT IDENTITY PRIMARY KEY`
- `user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE`
- `title VARCHAR(50) NOT NULL`
- `description VARCHAR(50)`
- `status VARCHAR(50) CHECK (status IN ('todo','doing','done')) DEFAULT 'todo'`
- `priority INT NOT NULL DEFAULT 0`
- `due_date TIMESTAMPTZ`
- `created_at`, `updated_at`
- индексы: `user_id`, `(user_id,status)`, частичный индекс для активных

**telegram_bindings** (Notification Service)
- `email TEXT PRIMARY KEY`
- `chat_id BIGINT NOT NULL`

---

**GitHub:** https://github.com/HDBOOMONE12/TaskManager
