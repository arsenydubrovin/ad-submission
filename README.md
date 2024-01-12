# Сервис для работы с объявлениями

Решение тестового задания [github.com/avito-tech/adv-backend-trainee-assignment](https://github.com/avito-tech/adv-backend-trainee-assignment).

Задача — разработать API для создания и получения объявлений, с возможностью валидации данных, пагинацией, сортировкой, опциональными полями.

## Методы API

### Создать объявление

HTTP метод: `POST /advert`

Формат запроса:

```json
{
    "title": "...",
    "description": "...",
    "price": 0,
    "photoLinks": ["...", "...", "..."]
}
```

Валидация:

- `title`, `price`, `photoLinks` обязательны,
- `title` не длиннее 200 символов,
- `description` не длиннее 1000 символов,
- `price` — целое, неотрицательное число,
- от 1 до 3 ссылок в `photoLinks`.

Метод создаёт новое объявление и возвращает JSON с полями созданного объявления и его `id`.

### Получить объявление по `id`

HTTP метод: `GET /advert/:id`

Параметры:

| Параметр | Описание | Значение по умолчанию |
| --- | --- | --- |
| `fields` | Значения через запятую | — |

Параметр `fields` добавит к ответу описание и список всех фотографий.

### Получить объявления постранично

HTTP метод: `GET /adverts`

Параметры:

| Параметр | Описание | Значение по умолчанию |
| --- | --- | --- |
| `page` | Целое неотрицательное число от 1 до 10 000 000 | 1 |
| `page_size` | Целое неотрицательное число от 1 до 100 | 10 |
| `sort` | Одно из допусимых значений | `id` |

Параметр `sort` может быть равен `id`, `date`, `price`, `-id`, `-date`, `-price`. Префикс `-` означает сортировку по убыванию.

Метод возвращает список объявлений. Каждое объявление содержит поля `id`, `title`, `primaryPhotoLink`. В параметрах указывается номер страницы, количество записей на странице и порядок сортировки.

## Запуск в Docker

Соберите образ.

```bash
make docker-build
```

Запустите приложение при помощи docker-compose.

```bash
make docker-run
```

## Локальная разработка

Установите необходимые инструменты: `air`, `gofumpt`, `golangci-lint`, `pre-commit` и `pymarkdown`.

```bash
make init
```

Запустите контейнер с PostgresSQL и приложение в live-режиме.

```bash
make run
```

Дополнительно:

|   |   |
|---|---|
| `make deps` | Обновить зависимости |
| `make lint` | Запустить линтеры для всего проекта |
| `make` | Показать `help` |

## Стек технологий

- Go `1.21`
- Python `3.11`
- Docker
- `make` в качестве таск-раннера

Пакеты:

- [`echo`](https://github.com/labstack/echo) — легковесный веб-фреймворк
- [`godotenv`](https://github.com/joho/godotenv) — для конфигурации приложения через переменные окружения, чтоби приблизить его к 12-факторному
- [`pq`](https://github.com/lib/pq) — драйвер для `PostgreSQL`
- [`log/slog`](https://pkg.go.dev/log/slog) — для структурированного логгирования
- [`console-slog`](https://github.com/phsym/console-slog) — форматтер для логов для локальной разработки
- [`slog-echo`](https://github.com/samber/slog-echo) — мидлварь для `echo`, объединяющий логи в `slog`

Инструменты:

- [`air`](https://github.com/cosmtrek/air) — автоматическая пересборка приложения
- [`gofumpt`](https://github.com/mvdan/gofumpt) — форматтер, построже, чем `gofmt`
- [`golangci-lint`](https://github.com/golangci/golangci-lint) — металинтер для go
- [`pre-commit`](https://github.com/pre-commit/pre-commit) — хуки для валидации кода перед каждый коммитом
