# Q&A Service API с авторизацией

REST API сервис для управления вопросами и ответами.

## Содержание

- [Особенности](#особенности)
- [Технологический стек](#технологический-стек)
- [Запуск](#запуск)
- [API Endpoints](#api-endpoints)
- [Примеры запросов](#примеры-запросов)
- [Тестирование](#тестирование)


---

## Особенности

### Функциональность

- Создание, получение, удаление вопросов
- Создание, получение, удаление ответов
- Каскадное удаление ответов при удалении вопроса
- Один пользователь может оставлять несколько ответов на один вопрос
- Валидация входных данных (email, пароль, текст)

---

## Технологический стек

| Компонент | Версия |
|---|---|
| **Язык** | Go 1.25+ |
| **Web-фреймворк** | net/http | 
| **ORM** | GORM |
| **База данных** | PostgreSQL 16 |
| **Миграции** | Goose |
| **Авторизация** | JWT |
| **Хеширование** | bcrypt |
| **Логирование** | log/slog |
| **Тестирование** | testify |
| **Контейнеризация** | Docker |

---


## Запуск

### Docker Compose

#### 1. Клонирование репозитория

```bash
git clone https://github.com/m4kson/hitalent-test.git
cd hitalent-test
```


#### 2. Запуск приложения

```bash
docker-compose up --build
```

Приложение будет доступно на `http://localhost:8080`

#### 4. Проверка здоровья

```bash
curl http://localhost:8080/health

```

#### 5. Просмотр логов

```bash
docker-compose logs -f app

docker-compose logs -f postgres
```

#### 6. Остановка

```bash
docker-compose down

```

---

## API Endpoints

API задокументированно в формате openapi, см. файл hitalent-test/docs/api/api.yaml

---

## Примеры запросов

### 1. Регистрация

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123"
  }'
```

**Успешный ответ (201):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "is_active": true,
  "created_at": "2025-01-15T10:30:45Z"
}
```

### 2. Логин

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123"
  }'
```

**Успешный ответ (200):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "is_active": true,
    "created_at": "2025-01-15T10:30:45Z"
  }
}
```

### 3. Создание вопроса

```bash
curl -X POST http://localhost:8080/questions/ \
  -H "Content-Type: application/json" \
  -d '{
    "text": "What is the capital of France?"
  }'
```

**Успешный ответ (201):**
```json
{
  "id": 1,
  "text": "What is the capital of France?",
  "created_at": "2025-01-15T10:30:45Z"
}
```


## Тестирование

### Unit-тесты

Запуск всех тестов:

```bash
go test ./...
```