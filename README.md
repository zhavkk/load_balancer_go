# Load Balancer Go

HTTP балансировщик нагрузки, написанный на Go, с поддержкой rate limiting и расширяемой архитектурой.

## Требования

- Go 1.21 или выше
- Docker и Docker Compose
- PostgreSQL 15 (если запуск без Docker)
- Goose (для миграций базы данных)

## Установка

1. Клонируйте репозиторий:
```bash
git clone https://github.com/zhavkk/load_balancer_go.git
cd load_balancer_go
```

2. Установите зависимости:
```bash
go mod download
```

## Сборка и запуск

### Локальный запуск

1. Соберите проект:
```bash
make build
```

2. Запустите приложение:
```bash
make run
```

### Запуск с Docker

1. Запустите все сервисы с помощью Docker Compose:
```bash
make docker-up
```

2. Для остановки:
```bash
make docker-down
```

## API Эндпоинты

### Управление бэкендами

- `GET /api/backends` - Получить список всех бэкендов
- `POST /api/backends` - Добавить новый бэкенд
  ```json
  {
    "url": "http://backend:8080"
  }
  ```
- `DELETE /api/backends/{id}` - Удалить бэкенд по ID

### Управление rate limiting

- `GET /api/limits` - Получить текущие лимиты
- `POST /api/limits` - Установить лимиты для клиента
  ```json
  {
    "client_id": "client123",
    "rps": 100,
    "burst": 200
  }
  ```
- `DELETE /api/limits/{client_id}` - Удалить лимиты для клиента

### Статистика

- `GET /api/stats` - Получить статистику по всем бэкендам
- `GET /api/stats/{backend_id}` - Получить статистику по конкретному бэкенду

### Прокси

- `GET /*` - Проксирование запросов на бэкенды согласно выбранному алгоритму балансировки

## Миграции базы данных

1. Применить миграции:
```bash
make migrate-up
```

2. Откатить последнюю миграцию:
```bash
make migrate-down
```

## Конфигурация

Основные настройки находятся в файле `config/config.yml`:

```yaml
proxy:
  port: "8080"
  algorithm: "round-robin"
backends:
  - url: "http://10.0.0.1:8000"
  - url: "http://10.0.0.2:8000"
rate_limit:
  enabled: true
  default_rps: 10
  default_burst: 20
  use_ip: true
db:
  dsn: "postgres://user:pass@db:5432/limits_db"
  update_interval: "5m"
env: "local"
```

## Переменные окружения

- `CONFIG_PATH` - путь к файлу конфигурации (по умолчанию: `config/config.yml`)
- `DB_DSN` - строка подключения к базе данных (по умолчанию: `postgres://user:pass@localhost:5432/limits_db?sslmode=disable`)

## Тестирование

Запуск тестов:
```bash
make test
```

## Структура проекта

```
.
├── cmd/
│   └── load_balancer/    # Точка входа приложения
├── internal/
│   ├── app/             # Основная логика приложения
│   ├── balancer/        # Реализация балансировщика
│   ├── config/          # Конфигурация
│   ├── handlers/        # HTTP обработчики
│   ├── logger/          # Логирование
│   ├── proxy/           # Прокси-сервер
│   ├── ratelimiter/     # Ограничение запросов
│   ├── repository/      # Работа с базой данных
│   ├── server/          # HTTP сервер
│   └── storage/         # Хранилище данных
├── config/              # Конфигурационные файлы
├── migrations/          # Миграции базы данных
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── README.md
```