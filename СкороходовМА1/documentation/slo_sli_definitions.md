# SLI/SLO для GoTaskFlow

## Цель метрик

Метрики нужны, чтобы оценивать не только факт работоспособности приложения, но и качество поставки: стабильность API, скорость ответа, успешность деплоев и надежность процесса разработки.

## Пользовательские SLI

| SLI | Что измеряется | Источник данных |
|---|---|---|
| Availability | Доля успешных ответов API без 5xx ошибок | HTTP logs, reverse proxy, monitoring |
| Latency | Время ответа API endpoints | Application metrics, HTTP logs |
| Error rate | Доля запросов с ошибками 5xx | HTTP logs |
| Task operation success rate | Доля успешных операций создания, обновления и удаления задач | API logs, tests |
| Stats cache hit readiness | Наличие Redis и успешность работы `/api/stats` | Redis healthcheck, API logs |

## DevOps SLI

| SLI | Что измеряется | Источник данных |
|---|---|---|
| CI success rate | Доля успешных запусков GitHub Actions | GitHub Actions |
| Build time | Время выполнения pipeline | GitHub Actions |
| Test pass rate | Доля успешных тестов | `go test ./...` |
| Security scan result | Количество критичных проблем gosec | GitHub Actions |
| Deployment reproducibility | Возможность поднять проект одной командой | Docker Compose |

## SLO

| SLO | Целевое значение | Период |
|---|---:|---|
| Availability API | 99% успешных запросов без 5xx | 30 дней |
| Latency для `GET /api/tasks` | p95 меньше 300 ms | 7 дней |
| Latency для `GET /api/stats` | p95 меньше 150 ms при доступном Redis | 7 дней |
| CI success rate | не ниже 95% для ветки `main` | 30 дней |
| Build time | меньше 3 минут | каждый запуск CI |
| Security scan | 0 high severity issues | каждый Pull Request |
| Recovery после локального сбоя окружения | `docker compose up --build` поднимает проект меньше чем за 5 минут | каждый запуск |

## Error budget

Для учебного проекта принимается простой error budget:

- если availability API ниже 99%, новые фичи временно откладываются;
- приоритет получает исправление ошибок, тестов, миграций или Docker Compose;
- если CI падает в `main`, исправление pipeline становится задачей с максимальным приоритетом.

## Как использовать метрики

На текущем этапе часть метрик проверяется вручную и через CI:

- `go test ./...`;
- `go vet ./...`;
- `gosec ./...`;
- `docker compose up --build`;
- ручная проверка `/health`, `/api/tasks`, `/api/stats`.

На следующих этапах можно добавить Prometheus/Grafana или structured logging, чтобы собирать latency, error rate и cache behavior автоматически.
