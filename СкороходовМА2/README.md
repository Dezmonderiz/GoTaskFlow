# Практика 202: Fast Feedback Loop

## Проект

GoTaskFlow - учебное веб-приложение/API для управления задачами. Проект развивается как сквозной DevOps-проект: Go backend, PostgreSQL, Redis, Docker Compose, тесты, pre-commit и GitHub Actions.

## Ссылки

- Репозиторий проекта: https://github.com/Dezmonderiz/GoTaskFlow.git
- CI/CD workflow: `.github/workflows/ci.yml`
- Конфигурация pre-commit: `.pre-commit-config.yaml`
- Копия конфигурации для сдачи: `artifacts/.pre-commit-config.yaml`

## Что настроено

### Локальная обратная связь

Настроен `pre-commit`, который запускается перед коммитом и проверяет:

- пробелы в конце строк;
- наличие пустой строки в конце файлов;
- корректность YAML/JSON;
- отсутствие больших случайно добавленных файлов;
- наличие потенциальных секретов через `detect-secrets`;
- форматирование Go-кода через `gofmt`;
- статический анализ `go vet`;
- unit-тесты `go test ./...`.

### Серверная обратная связь

GitHub Actions запускается на `push` в `main/master` и на каждый `pull_request`.

Pipeline включает:

- Lint: `pre-commit`, `gofmt`, `go vet`;
- Test: `go test ./...`;
- Security Scan: `gosec ./...`;
- Build: `go build -buildvcs=false -v ./cmd/app`.

## Артефакты

- `documentation/git_hooks_config.md` - описание локальных Git hooks и лог успешного запуска.
- `documentation/ci_pipeline_logic.md` - логика CI pipeline и Mermaid-схема.
- `artifacts/.pre-commit-config.yaml` - конфигурация pre-commit.
- `artifacts/pipeline_fail_example.png` - пример упавшего pipeline.
- `artifacts/pipeline_success_example.png` - пример успешного pipeline.

## Сравнение с VSM из практики 201

В VSM AS-IS этап проверки кода включал локальное тестирование, ручную проверку и ожидание review. Оценочное время активной проверки составляло около 2 часов 40 минут.

После настройки Fast Feedback Loop автоматическая проверка занимает примерно 3-5 минут:

```text
5 минут / 160 минут * 100% = около 3%
```

Итог: автоматизированная проверка занимает около 3% от прежнего ручного этапа проверки кода. Это сокращает Lead Time и быстрее показывает ошибки разработчику.

## Уведомления

Для GitHub Actions используются стандартные уведомления GitHub:

- статус проверки отображается прямо в Pull Request;
- автор PR получает уведомление в GitHub;
- при включенных email notifications GitHub отправляет письмо о результате workflow.

В дальнейшем можно добавить интеграцию с Telegram, Slack или Discord через webhook.
