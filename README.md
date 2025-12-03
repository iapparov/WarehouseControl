# WarehouseControl

## Описание

WarehouseControl — сервис управления складом: товары, история изменений и аудит.

Основные возможности:
- Управление товарами: создание, просмотр, редактирование, удаление.
- История изменений: фиксация операций (created/updated/deleted) с привязкой к пользователю и логину.
- Просмотр отличий между версиями товара (`item_diff`).
- Экспорт истории в CSV.
- JWT-аутентификация и роли (admin/manager/viewer).

## Состав репозитория

- **cmd/WarehouseControl/main.go** — точка входа (Uber FX DI).
- **internal/**
  - **app/** — бизнес-логика: `item`, `history`, `user`.
  - **auth/** — JWT аутентификация.
  - **config/** — загрузка конфигурации.
  - **di/** — регистрация зависимостей.
  - **domain/** — модели `item`, `history`, `user` (включая диффы).
  - **storage/postgres/** — репозитории PostgreSQL и сервис подключения.
  - **web/** — DTO, хэндлеры и роутер.
- **config/local.yaml** — пример конфигурации.
- **migrations/** — SQL-миграции (users, items, history, функции и триггеры).
- **docs/** — Swagger-документация.
- **web/index.html** — минимальный UI: аутентификация, CRUD товаров, история и CSV.
- **docker-compose.yml** — запуск PostgreSQL.


---

## Быстрый старт

### 1. Запуск инфраструктуры

```sh
docker-compose up -d
```
(Запустит контейнеры: postgres → порт 5433)

### 2. Конфигурация
Заполните `config/local.yaml` при необходимости.

### 3. Применить миграции

```sh
migrate -path migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" up
```

### 4. Запуск сервиса

```sh
go run ./cmd/WarehouseControl/main.go
```

Сервис стартует на порту 8080.

## API

Аутентификация:
- `POST /api/auth/register` — регистрация (login, password, role).
- `POST /api/auth/login` — вход, возвращает JWT.
- `POST /api/auth/refresh` — обновление токенов.

Товары:
- `GET /api/items` — список товаров.
- `GET /api/items/{id}` — товар по UUID.
- `POST /api/items` — создать товар (admin).
- `PUT /api/items/{id}` — обновить товар (admin/manager).
- `DELETE /api/items/{id}` — удалить товар (admin).

История:
- `GET /api/history?from=YYYY-MM-DD&to=YYYY-MM-DD&[id,action,login]` — история изменений.
- `GET /api/history/csv?from=YYYY-MM-DD&to=YYYY-MM-DD&[id,action,login]` — выгрузка CSV.

Swagger: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

Для защиты эндпоинтов используется JWT Bearer Auth:
- Укажите заголовок `Authorization: Bearer <token>` для защищённых эндпоинтов.
- Роли ограничивают операции с товарами.

## Веб-интерфейс
Откройте `web/index.html`: вход/регистрация, роль, CRUD товаров, история с диффами, экспорт CSV.

## Тесты
Юнит-тесты: `go test ./... -cover`

## Миграции

- `000001_create_user_table.*.sql`
- `000002_create_item_table.*.sql`
- `000003_create_history_table.*.sql`
- `000004_create_history_functions.*.sql`
- `000005_create_history_triggers.*.sql`

---

## Логирование
Через `wb-go/wbf/zlog`.

## Зависимости

- Go 1.21+
- PostgreSQL 16+
- Docker (для локального запуска инфраструктуры)

---

## Swagger

- Исходники в `docs/`
- Аннотации в `internal/web/handlers/*`
- Документация доступна по `/swagger/*`.