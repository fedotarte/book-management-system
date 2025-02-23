# 📚 Book Management System

## Описание
API для управления книгами, пользователями и отзывами.

## технологии
 - go 1.23
 - gin
 - gorm
 - postgres
 - mongo
 - cron


## 📌 руководство к запуску
Миграции осуществляется через make-команды

### Установка
1. **Установите зависимости**  
   ```sh
   make install
   ```

2. **Запустите базу данных**  
   ```sh
   make start-db
   ```

3. **Примените миграции**  
   ```sh
   make migrate-up
   ```

4. **Запустите сервер**  
   - Локально:
     ```sh
     make run-backend-local
     ```
   - В Docker:
     ```sh
     make run-backend-docker
     ```

---

## 📌 Работа с миграциями
**Создать новую миграцию**  
```sh
make new-migration
```

**Применить миграции**  
```sh
make migrate-up
```

**Откатить последнюю миграцию**  
```sh
make migrate-down
```

---

## 📌 Тестирование
**Запуск тестов:**  
```sh
make test
```

---

## 📌 API Документация
Swagger доступен по адресу:  
```
http://localhost:8080/swagger/index.html
```

---


## 📂 Структура проекта

📦 **cmd/** – точка входа в приложение  
📦 **internal/** – основной код приложения  
📂 **handlers/** – обработчики HTTP-запросов  
📂 **services/** – бизнес-логика  
📂 **repositories/** – работа с базами данных  
📂 **database/** – подключение к PostgreSQL и MongoDB  
📂 **middleware/** – мидлвари на авторизацию, ролевку и логирование  
📂 **dto/** – DTO для API  
📂 **pkg/** – утилсы (логгер, конвертации, jwt)  
📂 **migrations/** – SQL миграции  
📂 **docs/** – Swagger  
📄 **.env** – переменные окружения(должны быть `env.development`, `env.staging`, `env.production` )  
📄 **Dockerfile** – инструкции для контейнера бэка
📄 **docker-compose.yml** – описание сервисов инфры  
📄 **Makefile** – команды для CLI  
📄 **README.md** – документация  

## Графическое представление системы

### Как выглядит взаимодействие компонент:
```mermaid
graph TD;
    A[main.go] -->|Loads| B[middleware]
    A -->|Loads| C[routes]
    B -->|Uses| D[pkg/utils]
    
    C -->|Registers| E[handlers]
    E -->|Calls| F[services]
    F -->|Uses| G[dto]
    F -->|Calls| H[repositories]
    H -->|Uses| I[models]
    
    subgraph Middleware with Utils
        B
        D
    end
    
    subgraph HTTP Layer
        C
        E
    end
    
    subgraph Business Logic
        F
        G
    end
    
    subgraph Data Layer
        H
        I
    end
```

### Верхнеуровневое устройство системы 
```mermaid
graph LR
  User[польз]-->|USES| Frontend((️BMS FRONT))-->|HTTP| Backend((BMS API))
  Backend -->|books/users/authors| Postgres[(PostgreSQL)]
  Backend -->|reviews/votes| MongoDB[(MongoDB)]
  Backend -->|images| FileStorage[(пока на fs тачке)]
```

### Контекст C4
```mermaid
C4Context
    title Book Management System - Context Diagram

    Person(User, "Пользователь", "Читает книги, оставляет отзывы, ставит оценки")
    Person(Moderator, "Модератор", "Подтверждает книги, редактирует контент")
    Person(Admin, "Админ", "Управляет пользователями и системой")

    System(BookManagementSystem, "Book Management System", "Обрабатывает книги, отзывы и пользователей")

    ContainerDb(PostgresDB, "PostgreSQL", "Хранит данные о книгах и пользователях")
    ContainerDb(MongoDB, "MongoDB", "Хранит отзывы пользователей")
    Rel(User, BookManagementSystem, "Читает книги, оставляет отзывы")
    Rel(Moderator, BookManagementSystem, "Модерирует контент")
    Rel(Admin, BookManagementSystem, "Управляет контентом, пользователями и ролями")
    
    Rel(BookManagementSystem, PostgresDB, "Записывает данные о книгах, пользователях")
    Rel(BookManagementSystem, MongoDB, "Хранит отзывы пользователей")

```

### Контейнеры C4
```mermaid
C4Container
    title Book Management System - Containers

    Person(User, "Пользователь", "Взаимодействует с системой через API")

    System_Boundary(BookManagementSystem, "Book Management System") {
        Container(API, "API Backend", "Gin + GORM", "Обрабатывает запросы пользователей")
        ContainerDb(PostgresDB, "PostgreSQL", "Хранит книги, пользователей, авторов")
        ContainerDb(MongoDB, "MongoDB", "Хранит отзывы и оценки")
        Container(FileStorage, "S3 Storage/fs сервера", "Хранит обложки книг")
    }

    Rel(User, API, "Отправляет запросы API")
    Rel(API, PostgresDB, "Сохраняет книги, пользователей")
    Rel(API, MongoDB, "Сохраняет отзывы")
    Rel(API, FileStorage, "Загружает обложки книг")

```

### ER логическая
```mermaid 
erDiagram
    users {
        UUID id PK
        STRING username
        STRING email
        STRING password
        STRING role
        DATETIME deleted_at
        DATETIME created_at
        DATETIME updated_at
    }

    books {
        UUID id PK
        STRING title
        STRING description
        STRING cover_image
        BOOLEAN confirmed
        DATETIME deleted_at
        DATETIME created_at
        DATETIME updated_at
    }

    authors {
        UUID id PK
        STRING name
        TEXT bio
        DATETIME deleted_at
        DATETIME created_at
        DATETIME updated_at
    }

    book_authors {
        UUID book_id FK
        UUID author_id FK
    }

    user_books {
        UUID user_id FK
        UUID book_id FK
        STRING status
        INT pages_read
        DATETIME created_at
        DATETIME updated_at
    }

    reading_progresses {
        UUID id PK
        UUID user_id FK
        UUID book_id FK
        STRING status
        INT pages_read
        DATETIME created_at
        DATETIME updated_at
    }

    book_ratings {
        UUID user_id FK
        UUID book_id FK
        INT rating
        DATETIME created_at
    }

    refresh_tokens {
        UUID id PK
        UUID user_id FK
        STRING token
        DATETIME expires_at
    }

    moderator_actions {
        UUID id PK
        UUID moderator_id FK
        STRING action
        UUID target_id
        STRING target_type
        DATETIME created_at
    }

    users ||--o{ user_books : "читает"
    users ||--o{ reading_progresses : "читаемый статус"
    users ||--o{ book_ratings : "ставит оценку"
    users ||--o{ refresh_tokens : "имеет сессии"
    users ||--o{ moderator_actions : "выполняет действия"

    books ||--o{ book_authors : "написана"
    books ||--o{ user_books : "добавлена в список"
    books ||--o{ reading_progresses : "прогресс"
    books ||--o{ book_ratings : "получает оценки"

    authors ||--o{ book_authors : "пишет книги"

```
---
## 📌 Контрибьютинг
1. Форкни репозиторий
2. Создай новую ветку (`feature/your-feature`)
3. Запусти `make test` перед коммитом
4. Сделай PR


---

## 📌 Контакты

[tg](https://t.me/fedtart)