# Load Balancer with Rate Limiting

## Описание

Это HTTP-балансировщик нагрузки с опциональным ограничителем трафика на основе алгоритма Token Bucket.  
Конфигурация хранится во внешнем YAML-файле и в PostgreSQL для управления лимитами клиентов.

- Балансирует входящие HTTP-запросы на пул бэкендов.
- Автоматически исключает недоступные бэкенды.
- Ограничивает трафик на основе IP или API-ключа.
- Управляет лимитами через API.
- Хранит лимиты в PostgreSQL.
- Контейнеризуется через Docker и управляется через Makefile.

---

## Структура проекта
```
├── cmd/load_balancer/ # Точка входа
├── config/ # Конфиг и docker-compose.yml
├── migrations/ # SQL-миграции Goose
├── internal/ # Вся внутренняя логика
├── Dockerfile
├── Makefile
└── README.md
```

---

## Конфигурация

Пример `config/config.yml`:

```yaml
proxy:
  port: "8080"
  algorithm: "round-robin"

backends:
  - url: "http://localhost:9001"
  - url: "http://localhost:9002"

rate_limit:
  default_rps:   10
  default_burst: 20

db:
  dsn: "postgres://user:pass@db:5432/limits_db?sslmode=disable"
```
Сборка и запуск
Локально

make build          # Сборка бинарника
make migrate-up     # Применение миграций
make run            # Запуск приложения

Через Docker

make docker-up      # Поднимает сервисы и запускает миграции
make docker-down    # Останавливает сервисы

HTTP API

Управление лимитами клиентов
POST /clients

Добавить или обновить лимит:
```json
{
  "ip": "127.0.0.1",
  "capacity": 5,
  "refill_every": "2s"
}
```
DELETE /clients

Удалить лимит:
```json
{
  "ip": "127.0.0.1"
}
```
Пока не протестировано ничего из за нехватки времени :(
