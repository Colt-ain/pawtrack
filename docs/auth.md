# Аутентификация и авторизация

## Обзор

Модуль отвечает за регистрацию пользователей, вход в систему и управление JWT токенами.

## Роли пользователей

Система поддерживает три роли:

1. **Owner** - Владелец собак
2. **Consultant** - Профессиональный консультант (кинолог, грумер, ветеринар)
3. **Admin** - Администратор системы

## Бизнес-процессы

### 1. Регистрация владельца

**Endpoint**: `POST /api/v1/auth/register/owner`

**Бизнес-логика**:
1. Пользователь отправляет имя, email и пароль
2. Система проверяет уникальность email
3. Пароль хешируется с помощью bcrypt (cost 10)
4. Создаётся запись в таблице `users` с ролью `owner`
5. Возвращается информация о созданном пользователе (без пароля)

**Валидация**:
- Email должен быть валидным и уникальным
- Пароль: минимум 6 символов
- Имя: обязательное поле

**Пример запроса**:
```json
{
  "name": "Иван Петров",
  "email": "ivan@example.com",
  "password": "secret123"
}
```

**Пример ответа**:
```json
{
  "id": 1,
  "name": "Иван Петров",
  "email": "ivan@example.com",
  "role": "owner"
}
```

### 2. Регистрация консультанта

**Endpoint**: `POST /api/v1/auth/register/consultant`

**Бизнес-логика**:
Аналогична регистрации владельца, но роль устанавливается `consultant`.

**Особенности**:
- После регистрации консультант может создать профиль с описанием услуг
- Консультант получает доступ к собакам только после приглашения от владельца

### 3. Вход в систему

**Endpoint**: `POST /api/v1/auth/login`

**Бизнес-логика**:
1. Пользователь отправляет email и пароль
2. Система ищет пользователя по email
3. Проверяется пароль с помощью bcrypt.CompareHashAndPassword
4. При успехе генерируется JWT токен со следующими claims:
   - `user_id` - ID пользователя
   - `email` - Email пользователя
   - `role` - Роль пользователя
   - `exp` - Время истечения (по умолчанию +24 часа)
   - `iat` - Время создания
5. Токен подписывается секретом (HMAC SHA256)
6. Возвращается токен и информация о пользователе

**Пример запроса**:
```json
{
  "email": "ivan@example.com",
  "password": "secret123"
}
```

**Пример ответа**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "name": "Иван Петров",
    "email": "ivan@example.com",
    "role": "owner"
  }
}
```

**Ошибки**:
- 400 - Неверный формат данных
- 401 - Неверный email или пароль

## JWT токены

### Структура токена

```json
{
  "user_id": 1,
  "email": "ivan@example.com",
  "role": "owner",
  "exp": 1700000000,
  "iat": 1699913600
}
```

### Использование токена

Все защищённые эндпоинты требуют токен в заголовке:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Middleware авторизации

На каждый защищённый запрос:
1. Извлекается токен из заголовка `Authorization`
2. Проверяется формат: `Bearer <token>`
3. Токен валидируется (подпись, время истечения)
4. Извлекаются claims (user_id, email, role)
5. Данные сохраняются в контексте Gin для доступа в handlers

## Безопасность

### Хеширование паролей
- Используется bcrypt с cost factor 10
- Исходный пароль никогда не сохраняется в БД
- При проверке используется constant-time сравнение

### JWT безопасность
- Токены подписываются секретом (не передаётся клиенту)
- Алгоритм: HMAC SHA256
- Токены имеют ограниченное время жизни
- Нет механизма отзыва (stateless), при утечке нужно ждать истечения

### Рекомендации
1. Использовать HTTPS в production
2. Хранить JWT_SECRET в переменных окружения
3. Не передавать токены в URL
4. Использовать httpOnly cookies в production (опционально)
5. Регулярно обновлять токены (refresh tokens - не реализовано)

## Обработка ошибок

### Дублирование email
При регистрации с существующим email:
```json
{
  "error": "email already in use"
}
```
Status: 409 Conflict

### Неверные credentials
При вводе неверного пароля или несуществующего email:
```json
{
  "error": "invalid credentials"
}
```
Status: 401 Unauthorized

### Невалидный токен
При использовании истёкшего или подделанного токена:
```json
{
  "error": "invalid or expired token"
}
```
Status: 401 Unauthorized

### Отсутствие токена
При обращении к защищённому эндпоинту без токена:
```json
{
  "error": "missing authorization header"
}
```
Status: 401 Unauthorized

## Конфигурация

### Переменные окружения

- `JWT_SECRET` - Секрет для подписи токенов (обязательный в production)
- `JWT_EXPIRY_HOURS` - Время жизни токена в часах (default: 24)

### Пример конфигурации

```yaml
# docker-compose.yml
environment:
  JWT_SECRET: "your-super-secret-key-change-in-production"
  JWT_EXPIRY_HOURS: "24"
```

## Примеры использования

### Полный flow регистрации и входа

```bash
# 1. Регистрация
curl -X POST http://localhost:8080/api/v1/auth/register/owner \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Иван Петров",
    "email": "ivan@example.com",
    "password": "secret123"
  }'

# 2. Вход
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "ivan@example.com",
    "password": "secret123"
  }' | jq -r '.token')

# 3. Использование токена
curl -X GET http://localhost:8080/api/v1/dogs \
  -H "Authorization: Bearer $TOKEN"
```

## Связанные модули

- [Пользователи](./users.md) - управление профилями после регистрации
- [Консультанты](./consultants.md) - расширенный профиль для консультантов
