# GoTaskFlow

## Краткое описание проекта

GoTaskFlow - учебное веб-приложение для управления задачами. Пользователь может создавать задачи, просматривать список, менять статус и удалять задачи. Для контроля состояния проекта доступна статистика по задачам: всего, `todo`, `in_progress`, `done`.

Проект используется как сквозной DevOps-проект для отработки практик: Git, CI/CD, контейнеризация, миграции, кэширование, тестирование и анализ потока создания ценности.

## Технологический стек

| Компонент | Технология |
|---|---|
| Язык | Go |
| Web framework | Gin |
| База данных | PostgreSQL |
| Кэш | Redis |
| Миграции | golang-migrate |
| Тестирование | testing, httptest |
| CI/CD | GitHub Actions |
| Контейнеризация | Docker, Docker Compose |
| Предполагаемое облако | GitHub Actions + контейнерный деплой в VPS или облачный container service |

## Репозиторий проекта

https://github.com/Dezmonderiz/GoTaskFlow.git

## Итоги VSM-аудита

AS-IS поток поставки фичи для проекта GoTaskFlow показывает, что основная доля Lead Time уходит не на чистое написание кода, а на ожидание проверки, ручную подготовку окружения и деплой.

Оценочный Lead Time: около 3 рабочих дней.

Оценочный Processing Time: около 13 часов.

Основные узкие места:

- ожидание проверки и ручного code review;
- ручной деплой и проверка окружения PostgreSQL/Redis/Docker.

## Git Governance

Для проекта выбрана trunk-based стратегия:

- основная ветка: `main`;
- новые изменения выполняются в короткоживущих feature-ветках;
- перед слиянием выполняются проверки CI;
- изменения в `main` должны быть небольшими и регулярно интегрироваться.

Рекомендуемые Branch Protection Rules для `main`:

- запрет прямого push в `main`;
- обязательный Pull Request перед merge;
- обязательное прохождение GitHub Actions workflow `CI`;
- запрет merge при падающих тестах;
- требование актуальности ветки перед merge.

## Управление задачами

Для управления задачами используется GitHub Projects или GitHub Issues Board со статусами:

- Backlog;
- In Progress;
- QA/Review;
- Done.

Артефакт доски расположен в `artifacts/project_board_setup.png`.

## Документация

- `documentation/value_stream_map.md` - Value Stream Mapping AS-IS.
- `documentation/slo_sli_definitions.md` - SLI/SLO и метрики успеха.
- `documentation/architecture_overview.md` - архитектура приложения.
