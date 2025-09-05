# 🧩 Task Manager — микросервисная платформа (Core Service на Go)

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15%2B-336791?logo=postgresql)](https://www.postgresql.org/)
[![Build](https://img.shields.io/badge/tests-go%20test-success)](#-тестирование)

**Task Manager** — платформа для управления задачами, уведомлениями и рейтингами пользователей.  
Этот репозиторий содержит **ядро (Task Service)** — REST‑сервис на Go, хранящий пользователей и задачи в **PostgreSQL**.

> 🎯 Видение: система из нескольких микросервисов (Task, Notification, Rating), связанных **Kafka + gRPC** и интегрированных с **Telegram**. Текущая версия — монолитный сервис задач с чистой структурой и готовностью к выделению по сервисам.

---

## 🧭 Зачем это нужно

- Автоматизировать бизнес‑процессы: постановка задач, контроль сроков, исполнители.
- Гибкая эволюция до микросервисной архитектуры (Kafka/gRPC, Telegram‑уведомления, рейтинги).

---



## 🛠️ Используемые технологии (кратко)

- **Go 1.24+**: `net/http`, `context`, `encoding/json`, graceful shutdown.
- **API**: REST + JSON, `http.ServeMux`, унифицированные коды/ошибки.
- **Config**: `.env` через `godotenv`, ключ — `DATABASE_URL`.
- **БД**: PostgreSQL 15+, `pgx` через `database/sql`, пул соединений в `internal/db/postgres.go`.
- **Слои**: `handlers` → `service` → `storage` → БД; доменные сущности в `internal/entity`.
- **Миграции**: SQL в `internal/db/migrations`.
- **Тесты**: `go test ./...`, моки через `gomock` (`//go:generate`).
- **Дальше**: Docker Compose, Kafka + gRPC, JWT, CI/CD, Prometheus/Grafana.


---

## ⚙️ Установка и запуск

### 1) Клонирование
```bash
git clone https://github.com/HDBOOMONE12/TaskManager.git
cd TaskManager
```

### 2) Конфигурация окружения
Создайте `.env`, ориентируясь на `.env.example`:
```dotenv
DATABASE_URL=postgres://USER:PASSWORD@localhost:5432/taskmanager?sslmode=disable
```
> Переменные окружения и секреты не коммитим: это уже отражено в `.gitignore`.

### 3) Подготовка БД (миграции)
```bash
psql "$DATABASE_URL" -f internal/db/migrations/0001_create_users.sql
psql "$DATABASE_URL" -f internal/db/migrations/0002_create_tasks.sql
psql "$DATABASE_URL" -f internal/db/migrations/0003_indexes.sql
```

### 4) Запуск
```bash
go run ./cmd/taskmanager
```
По умолчанию сервис слушает `http://localhost:8080`.

---

## 📡 API (REST)

Формат ошибок (пример):
```json
{
  "error": "message"
}
```

### Пользователи

#### `GET /users`
Список пользователей.

**Успех — 200 OK**
```json
[
  { "id": 1, "name": "alice", "email": "alice@example.com" },
  { "id": 2, "name": "bob",   "email": "bob@example.com"   }
]
```

#### `POST /users`
Создать пользователя.

**Тело**
```json
{
  "name": "alice",
  "email": "alice@example.com"
}
```

**Успех — 201 Created**
```json
{ "id": 1, "name": "alice", "email": "alice@example.com" }
```

#### `GET /users/{id}`
Получить пользователя по ID.  
Коды: `200`, `404`.

#### `PATCH /users/{id}`
Частичное обновление имени/email.

**Тело (любой из полей опционален)**
```json
{ "name": "Alice Cooper", "email": "alice@company.com" }
```
Коды: `200`, `400`, `404`.

#### `DELETE /users/{id}`
Удалить пользователя. Коды: `204`, `404`.

---

### Задачи (привязаны к пользователю)

Статусы: **`todo` | `doing` | `done`**.  
Приоритет: целое число (**0..5**, по умолчанию `3`).  
Дата срока: ISO8601, поле `due_at` (опционально).

#### `GET /users/{user_id}/tasks`
Список задач пользователя. Коды: `200`, `400`, `404`.

#### `POST /users/{user_id}/tasks`
Создать задачу.

**Тело**
```json
{
  "title": "Сделать отчёт",
  "description": "К пятнице",
  "status": "todo",
  "priority": 3,
  "due_at": "2025-09-05T18:00:00Z"
}
```
**Успех — 201 Created**
```json
{
  "id": 10,
  "user_id": {user_id},
  "title": "Сделать отчёт",
  "description": "К пятнице",
  "status": "todo",
  "priority": 3,
  "due_at": "2025-09-05T18:00:00Z",
  "created_at": "2025-08-30T12:00:00Z",
  "updated_at": "2025-08-30T12:00:00Z"
}
```

#### `GET /users/{user_id}/tasks/{task_id}`
Получить задачу. Коды: `200`, `400`, `404`.

#### `PUT /users/{user_id}/tasks/{task_id}`
Полное обновление полей задачи. Коды: `200`, `400`, `404`.

#### `PATCH /users/{user_id}/tasks/{task_id}`
Частичное обновление (например, только `status` или `title`). Коды: `200`, `400`, `404`.

#### `DELETE /users/{user_id}/tasks/{task_id}`
Удалить задачу. Коды: `204`, `400`, `404`.

---

## 🧪 Тестирование
```bash
go test ./...
```
- В проекте есть юнит‑тесты сервиса задач (`internal/service/task_test.go`).
- Для репозитория задач предусмотрена генерация моков (`//go:generate mockgen` в `internal/storage/tasks_repo.go`).

---

## 🗂️ Структура проекта
```text
.
├── cmd/
│   └── taskmanager/
│       └── main.go
├── internal/
│   ├── db/
│   │   ├── postgres.go
│   │   └── migrations/
│   │       ├── 0001_create_users.sql
│   │       ├── 0002_create_tasks.sql
│   │       └── 0003_indexes.sql
│   ├── entity/
│   │   ├── user.go
│   │   └── task.go
│   ├── handlers/
│   │   ├── users.go
│   │   ├── tasks.go
│   │   ├── router_users.go
│   │   └── errors_tasks.go
│   ├── service/
│   │   ├── users.go
│   │   ├── tasks.go
│   │   └── task_test.go
│   ├── storage/
│   │   ├── users_repo.go
│   │   └── tasks_repo.go
│   └── mocks/
│       └── mock_task_repo.go
├── .env.example
├── .gitignore
├── go.mod
└── go.sum
```

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
- индексы по `user_id`, `user_id,status`, частичный индекс для активных задач

---
