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

## Возможности API