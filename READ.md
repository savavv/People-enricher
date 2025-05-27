# People Enricher API

Сервис для обогащения данных о людях (ФИО) с использованием публичных API для определения вероятного возраста, пола и национальности.

---

## Описание

API принимает данные о человеке с полями:

```json
{
  "name": "Dmitriy",
  "surname": "Ushakov",
  "patronymic": "Vasilevich" // необязательно
}
```

Быстрый старт
1. Клонировать репозиторий
   ```
   https://github.com/savavv/People-enricher.git
   ```
   ```
   cd People-enricher
   ```

2. Создать и настроить .env файл
  ```
  DB_HOST=localhost
  DB_PORT=5432
  DB_USER=postgres
  DB_PASSWORD=your_password
  DB_NAME=people_db
  SERVER_PORT=8080
  ```
3. Запустить PostgreSQL (локально или в Docker)
   ```
   docker run --name people-postgres -e POSTGRES_PASSWORD=your_password -p 5432:5432 -d postgres
    ```
4. Запустить сервис
   ```
   go run ./cmd/main.go
   ```
   ```
   Или с Docker:
   docker build -t people-enricher .
   docker run -p 8080:8080 --env-file .env people-enricher
   ```
   ## API Endpoints
   ```
   POST /people — Добавить нового человека (обогатит и сохранит)

   GET /people — Получить список людей с фильтрами и пагинацией (name, surname, page, limit)

   PUT /people/{id} — Обновить данные человека по ID

   DELETE /people/{id} — Удалить человека по ID

   GET /swagger/index.html — Swagger UI с документацией API
   ```
