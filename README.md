# PVZ Server

**PVZ Server** — это backend-сервис для сотрудников пунктов выдачи заказов (ПВЗ), который позволяет вести учёт приёмок товаров и управления ими. Сервис обеспечивает создание ПВЗ, регистрацию приёмок, добавление и удаление товаров, а также закрытие приёмок. 

## Технологии и зависимости

- **Язык программирования:** Go 1.24
- **Фреймворк для HTTP-сервера:** [Gin](https://github.com/gin-gonic/gin)
- **Работа с базой данных:** `database/sql` + [lib/pq](https://github.com/lib/pq)
- **Миграции базы данных:** [golang-migrate/migrate](https://github.com/golang-migrate/migrate)
- **JWT авторизация:** [github.com/golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt)
- **Тестирование:** `testing`, [stretchr/testify](https://github.com/stretchr/testify)
- **База данных:** PostgreSQL

---

## Сборка и запуск
## docker-compose - TBA SOON

### Необходимо сделать перед сборкой

1. Создать базу даннных:

```bash
createdb -U postgres pvz
```

2. Применить миграции для создания необходимых таблиц:
```bash
make migrate-up
```

### Сборка проекта

Для сборки используйте команду:

```bash
make build
```

Бинарный файл будет создан в директории `bin/`.

### Запуск сервера

Для запуска используйте:

```bash
make run
```

После этого сервис будет доступен на порту `:8080`

### Команды Makefile

| Команда                 | Назначение                                                          |
|-------------------------|---------------------------------------------------------------------|
| `make build`            | Собирает бинарный файл в директорию `bin/`                         |
| `make run`              | Запускает приложение с использованием `.env`                       |
| `make test`             | Запускает unit-тесты в `internal/handlers/`                        |
| `make integration_test` | Запускает интеграционный тест в `internal/handlers/integration_test` |
| `make migrate-up`       | Применяет миграции к базе данных                                   |
| `make migrate-down`     | Откатывает миграции                                                 |
| `make clean`            | Удаляет собранные бинарники из `bin/`                              |


---

## Возможности API

### 1. Авторизация 

### POST /dummyLogin
Для доступа к эндпоинтам требуется авторизация через JWT токен.
Токен можно получить следующим образом. 

```http
POST /dummyLogin
```

Пример тела запроса:

```json

{
  "role": "moderator"
}

```

Ответ:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6..."
}
```

Полученный токен необходимо передавать в заголовке каждого запроса в формате:

```http
Authorization: Bearer <ваш_токен>
```
Это обязательное условие для всех защищённых эндпоинтов, включая:

- `/pvz`

- `/receptions`

- `/products`

- `/pvz/{pvzId}/delete_last_product`

- `/pvz/{pvzId}/close_last_reception`

### 2. Создание ПВЗ (только модератор)

### POST /pvz

Создаёт новый ПВЗ в одном из допустимых городов: Москва, Санкт-Петербург, Казань.

Пример запроса:

```http
Authorization: Bearer <токен_модератора>
Content-Type: application/json
```
```json
{
  "city": "Москва"
}
```

### 3. Создание приёмки (только сотрудник ПВЗ)

#### POST /receptions

Создание новой приёмки товаров. Должна быть только одна активная приёмка.

Пример запроса:

```http
Authorization: Bearer <токен_модератора>
Content-Type: application/json
```
```json
{
  "pvzId": "pvz_id"
}
```

### 4. Добавление товара в приёмку

### POST /products

Добавляет товар (электроника, одежда, обувь) в текущую активную приёмку.

Пример запроса:

```http
Authorization: Bearer <токен_модератора>
Content-Type: application/json
```
```json
{
  "type": "одежда",
  "pvzId": "pvz_id"
}
```

### 5. Удаление последнего товара

### POST /pvz/{pvzId}/delete_last_product

Удаляет последний добавленный товар из текущей приёмки. Только пока приёмка не закрыта.

Пример запроса:

```http
POST /pvz/pvz_id/delete_last_product
```

```http
Authorization: Bearer <токен_сотрудника>
Content-Type: application/json
```

### 6. Закрытие текущей приёмки

### POST /pvz/{pvzId}/close_last_reception

Закрывает последнюю активную приёмку.

Пример запроса:

```http
POST /pvz/pvz_id/close_last_reception
```

```http
Authorization: Bearer <токен_сотрудника>
Content-Type: application/json
```

### 7. Получение списка ПВЗ

### GET /pvz

Возвращает список ПВЗ, приёмок и товаров. Можно указать фильтр по дате и пагинации.

```http
Authorization: Bearer <токен_сотрудника или модератора>
Content-Type: application/json
```
Query-параметры:
| Параметр    | Тип       | Описание                               | Пример                  |
|-------------|-----------|----------------------------------------|--------------------------|
| `startDate` | `datetime`| Начало интервала фильтрации приёмок    | `2025-04-13T00:00:00Z`  |
| `endDate`   | `datetime`| Конец интервала фильтрации приёмок     | `2025-04-14T23:59:59Z`  |
| `page`      | `int`     | Номер страницы                         | `1`                      |
| `limit`     | `int`     | Кол-во элементов на странице    | `10`                     |


### Сущности и их взаимосвязи

