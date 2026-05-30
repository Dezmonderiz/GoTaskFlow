# GoTaskFlow

Mini Task Tracker на Go.

## Планируемый стек

- Go
- Gin или Chi
- PostgreSQL
- Redis
- golang-migrate
- testing + httptest
- golangci-lint
- Docker + Docker Compose
- GitHub Actions

## Структура

```text
cmd/app        точка входа приложения
internal       внутренние пакеты API
web            статические файлы веб-интерфейса
migrations     SQL-миграции базы данных
documentation  проектная документация
artifacts      артефакты разработки
tests          дополнительные тестовые материалы
```

## Запуск

```bash
go run ./cmd/app
```

Сервер стартует на `http://localhost:8080`.

Минимальный frontend доступен на главной странице:

```text
http://localhost:8080/
```

## Docker Compose

Запуск приложения вместе с PostgreSQL и Redis:

```bash
docker compose up --build
```

Compose поднимает сервисы:

- `app` - Go API и frontend;
- `postgres` - база данных PostgreSQL;
- `redis` - кэш для `GET /api/stats`;
- `migrate` - одноразовый запуск миграций перед стартом приложения.

После запуска приложение доступно на:

```text
http://localhost:8080/
```

Остановка:

```bash
docker compose down
```

Остановка с удалением данных PostgreSQL:

```bash
docker compose down -v
```

По умолчанию приложение ожидает PostgreSQL по адресу:

```text
postgres://postgres:postgres@localhost:5432/gotaskflow?sslmode=disable
```

Можно переопределить настройки через переменные окружения:

```text
PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/gotaskflow?sslmode=disable
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
STATS_CACHE_TTL_SECONDS=60
```

## Миграции

SQL-миграции лежат в `migrations`.

Пример запуска через `golang-migrate`:

```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/gotaskflow?sslmode=disable" up
```

## Тесты

```bash
go test ./...
```

## CI/CD

GitHub Actions workflow находится в `.github/workflows/ci.yml`.

Pipeline выполняет:

- checkout кода;
- установку Go из `go.mod`;
- загрузку зависимостей;
- `gofmt` check;
- `go vet ./...`;
- `go test ./...`;
- security scan через `gosec`;
- сборку приложения.

## API

```text
GET    /health
GET    /api/tasks
POST   /api/tasks
GET    /api/tasks/{id}
PATCH  /api/tasks/{id}/status
DELETE /api/tasks/{id}
GET    /api/stats
```

### Создать задачу

```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Learn Go","description":"Build Mini Task Tracker"}'
```

### Изменить статус

```bash
curl -X PATCH http://localhost:8080/api/tasks/1/status \
  -H "Content-Type: application/json" \
  -d '{"status":"in_progress"}'
```

Доступные статусы: `todo`, `in_progress`, `done`.

### Статистика

```bash
curl http://localhost:8080/api/stats
```

Пример ответа:

```json
{
  "total": 10,
  "todo": 3,
  "done": 3,
  "in_progress": 4
}
```

`GET /api/stats` кэшируется в Redis на `STATS_CACHE_TTL_SECONDS`, по умолчанию на 60 секунд. При создании, удалении или изменении статуса задачи кэш статистики сбрасывается.
