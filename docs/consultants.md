# Консультанты (Consultants)

## Обзор

Модуль для профессиональных консультантов (кинологов, грумеров, ветеринаров). Позволяет создавать профили, искать специалистов и управлять системой приглашений для получения доступа к собакам клиентов.

## Структура данных

### ConsultantProfile (Профиль консультанта)

```go
type ConsultantProfile struct {
    UserID      uint   // ID пользователя (1:1 с users)
    User        *User  // Связь с пользователем
    Description string // Описание опыта и услуг
    Services    string // Список услуг (через запятую или свободный текст)
    Breeds      string // Специализация по породам
    Location    string // Город/регион работы
    Surname     string // Фамилия (для поиска)
}
```

### Invite (Приглашение)

```go
type Invite struct {
    ID           uint         // Уникальный идентификатор
    OwnerID      uint         // ID владельца собаки
    ConsultantID uint         // ID консультанта
    DogID        uint         // ID собаки
    Token        string       // Уникальный токен приглашения
    Status       InviteStatus // pending, accepted, rejected
    CreatedAt    time.Time    // Дата создания
    ExpiresAt    time.Time    // Дата истечения (по умолчанию +24ч)
}

type InviteStatus string

const (
    InvitePending  InviteStatus = "pending"
    InviteAccepted InviteStatus = "accepted"
    InviteRejected InviteStatus = "rejected"
)
```

## Бизнес-процессы

### 1. Создание/обновление профиля консультанта

**Endpoint**: `PUT /api/v1/consultants/profile`

**Права доступа**: Consultant, Admin

**Бизнес-логика**:
1. Консультант создаёт/обновляет свой профиль
2. Если профиль не существует - создаётся новый
3. Если существует - обновляются поля
4. Профиль связан 1:1 с записью в `users`

**Валидация**:
- Все поля опциональные
- `description`, `services`, `breeds`, `location` - текстовые поля
- `surname` - для улучшения поиска

**Пример запроса**:
```json
{
  "description": "Профессиональный кинолог с 10-летним опытом. Специализация: воспитание и дрессировка собак.",
  "services": "Дрессировка, Коррекция поведения, Консультации",
  "breeds": "Немецкая овчарка, Лабрадор, Золотистый ретривер",
  "location": "Москва",
  "surname": "Иванов"
}
```

**Пример ответа**:
```json
{
  "user_id": 5,
  "name": "Алексей",
  "surname": "Иванов",
  "description": "Профессиональный кинолог с 10-летним опытом...",
  "services": "Дрессировка, Коррекция поведения, Консультации",
  "breeds": "Немецкая овчарка, Лабрадор, Золотистый ретривер",
  "location": "Москва"
}
```

### 2. Поиск консультантов

**Endpoint**: `GET /api/v1/consultants`

**Права доступа**: All authenticated

**Query параметры**:
- `query` - Поиск по имени, фамилии, описанию (ILIKE)
- `services` - Фильтр по услугам (подстрока)
- `breeds` - Фильтр по породам (подстрока)
- `location` - Фильтр по локации (подстрока)
- `page` - Номер страницы (default: 1)
- `page_size` - Размер страницы (default: 20, max: 100)

**Бизнес-логика**:
1. Поиск по профилям консультантов с JOIN к `users`
2. Фильтрация по нескольким критериям одновременно
3. Case-insensitive поиск (ILIKE)
4. Пагинация результатов

**Реализация поиска**:
```go
// Поиск по тексту
if query != "" {
    search := "%" + query + "%"
    query.Where("users.name ILIKE ? OR surname ILIKE ? OR description ILIKE ?", 
                search, search, search)
}

// Фильтры
if services != "" {
    query.Where("services ILIKE ?", "%"+services+"%")
}

if breeds != "" {
    query.Where("breeds ILIKE ?", "%"+breeds+"%")
}

if location != "" {
    query.Where("location ILIKE ?", "%"+location+"%")
}
```

**Пример запроса**:
```
GET /api/v1/consultants?query=дрессировка&breeds=овчарка&location=Москва
```

**Пример ответа**:
```json
{
  "data": [
    {
      "id": 5,
      "name": "Алексей Иванов",
      "description": "Профессиональный кинолог...",
      "services": "Дрессировка, Коррекция поведения",
      "breeds": "Немецкая овчарка, Лабрадор",
      "location": "Москва"
    }
  ],
  "total_count": 1,
  "page": 1,
  "page_size": 20
}
```

### 3. Получение профиля консультанта

**Endpoint**: `GET /api/v1/consultants/:id`

**Права доступа**: All authenticated

**Бизнес-логика**:
1. Любой авторизованный пользователь может просмотреть профиль консультанта
2. Возвращается полная информация о консультанте
3. Используется для предпросмотра перед приглашением

**Пример ответа**:
```json
{
  "user_id": 5,
  "name": "Алексей",
  "surname": "Иванов",
  "description": "Профессиональный кинолог с 10-летним опытом работы...",
  "services": "Дрессировка, Коррекция поведения, Консультации",
  "breeds": "Немецкая овчарка, Лабрадор, Золотистый ретривер",
  "location": "Москва"
}
```

### 4. Приглашение консультанта

**Endpoint**: `POST /api/v1/consultants/:id/invite`

**Права доступа**: Owner (владелец собаки)

**Бизнес-логика**:
1. Владелец приглашает консультанта для работы с конкретной собакой
2. Генерируется уникальный токен приглашения (32 символа, hex)
3. Создаётся запись в таблице `invites` со статусом `pending`
4. Устанавливается время истечения (по умолчанию +24 часа)
5. **Stub**: Отправка email с ссылкой приглашения (не реализовано)
6. Возвращается информация о приглашении (включая токен для тестирования)

**Валидация**:
- Только владелец может приглашать для своих собак
- Нельзя пригласить на уже предоставленный доступ

**Пример запроса**:
```json
{
  "dog_id": 1
}
```

**Пример ответа**:
```json
{
  "id": 15,
  "token": "a3f5e8d2c9b1...",
  "status": "pending",
  "consultant_id": 5,
  "dog_id": 1,
  "expires_at": "2025-11-24T10:00:00Z"
}
```

**Email stub**:
```
From: noreply@pawtrack.com
To: consultant@example.com
Subject: Приглашение для работы с собакой

Владелец пригласил вас для работы с собакой.
Примите приглашение: https://pawtrack.com/invites/accept?token=a3f5e8d2c9b1...

Ссылка действует 24 часа.
```

### 5. Принятие приглашения

**Endpoint**: `POST /api/v1/invites/accept?token=xxx`

**Права доступа**: Consultant

**Бизнес-логика**:
1. Консультант переходит по ссылке с токеном
2. Система проверяет:
   - Токен существует
   - Статус = `pending`
   - Не истёк срок действия
   - Токен предназначен для текущего консультанта
3. При успешной валидации:
   - Статус меняется на `accepted`
   - Создаётся запись в `consultant_access`:
     ```sql
     INSERT INTO consultant_access (consultant_id, dog_id, granted_at)
     VALUES (5, 1, NOW());
     ```
   - Консультант получает доступ к собаке
4. Теперь консультант может:
   - Просматривать информацию о собаке
   - Просматривать события собаки
   - Создавать события
   - Создавать заметки

**Пример запроса**:
```
POST /api/v1/invites/accept?token=a3f5e8d2c9b1...
```

**Пример ответа**:
```json
{
  "message": "invite accepted"
}
```

**Ошибки**:
- 400 - Токен невалидный, истёк или уже использован
- 400 - Приглашение не для этого консультанта
- 401 - Не авторизован

## Система доступа

### ConsultantAccess (Доступ консультанта)

```go
type ConsultantAccess struct {
    ID           uint       // Уникальный идентификатор
    ConsultantID uint       // ID консультанта
    DogID        uint       // ID собаки
    GrantedAt    time.Time  // Когда предоставлен доступ
    RevokedAt    *time.Time // Когда отозван (NULL если активен)
}
```

### Проверка доступа

Используется во всех модулях для проверки прав консультанта:

```go
// В DogRepository
func HasConsultantAccess(consultantID, dogID uint) (bool, error) {
    var count int64
    err := db.Model(&ConsultantAccess{}).
        Where("consultant_id = ? AND dog_id = ?", consultantID, dogID).
        Where("revoked_at IS NULL").
        Count(&count).Error
    return count > 0, err
}
```

### Что даёт доступ

После принятия приглашения консультант может:

✅ **Просмотр**:
- Информация о собаке (GET /dogs/:id)
- События собаки (GET /events?dog_id=X)
- Свои заметки о собаке (GET /consultant-notes?dog_id=X)

✅ **Создание**:
- События для собаки (POST /events)
- Заметки о собаке (POST /consultant-notes)

❌ **Запрещено**:
- Редактирование информации о собаке
- Удаление собаки
- Удаление событий
- Просмотр заметок других консультантов

## Контроль доступа (RBAC)

| Операция | Owner | Consultant | Admin |
|----------|-------|------------|-------|
| Создать профиль | ❌ | ✅ Свой | ✅ |
| Обновить профиль | ❌ | ✅ Свой | ✅ Любой |
| Поиск консультантов | ✅ | ✅ | ✅ |
| Просмотр профиля | ✅ | ✅ | ✅ |
| Пригласить консультанта | ✅ Для своих собак | ❌ | ✅ |
| Принять приглашение | ❌ | ✅ Своё | ✅ |

## База данных

### Схема таблиц

```sql
-- Профили консультантов
CREATE TABLE consultant_profiles (
    user_id INTEGER PRIMARY KEY REFERENCES users(id),
    description TEXT,
    services TEXT,
    breeds TEXT,
    location VARCHAR(255),
    surname VARCHAR(100)
);

-- Приглашения
CREATE TABLE invites (
    id SERIAL PRIMARY KEY,
    owner_id INTEGER NOT NULL REFERENCES users(id),
    consultant_id INTEGER NOT NULL REFERENCES users(id),
    dog_id INTEGER NOT NULL REFERENCES dogs(id),
    token VARCHAR(64) UNIQUE NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'accepted', 'rejected')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_invites_token ON invites(token);
CREATE INDEX idx_invites_consultant_id ON invites(consultant_id);

-- Доступы консультантов
CREATE TABLE consultant_access (
    id SERIAL PRIMARY KEY,
    consultant_id INTEGER NOT NULL REFERENCES users(id),
    dog_id INTEGER NOT NULL REFERENCES dogs(id),
    granted_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(consultant_id, dog_id)
);

CREATE INDEX idx_consultant_access_consultant ON consultant_access(consultant_id);
CREATE INDEX idx_consultant_access_dog ON consultant_access(dog_id);
```

## Примеры использования

### Flow: От поиска до доступа

```bash
# 1. Владелец ищет консультанта
curl -X GET "http://localhost:8080/api/v1/consultants?services=дрессировка&location=Москва" \
  -H "Authorization: Bearer $OWNER_TOKEN"

# 2. Владелец приглашает консультанта
curl -X POST http://localhost:8080/api/v1/consultants/5/invite \
  -H "Authorization: Bearer $OWNER_TOKEN" \
  -d '{"dog_id": 1}'

# Ответ: {"token": "a3f5e8d2c9b1..."}

# 3. Консультант принимает приглашение
curl -X POST "http://localhost:8080/api/v1/invites/accept?token=a3f5e8d2c9b1..." \
  -H "Authorization: Bearer $CONSULTANT_TOKEN"

# 4. Консультант получает доступ к собаке
curl -X GET http://localhost:8080/api/v1/dogs/1 \
  -H "Authorization: Bearer $CONSULTANT_TOKEN"
```

## Бизнес-правила

1. **Профиль 1:1**: Один консультант = один профиль
2. **Истечение приглашений**: По умолчанию 24 часа
3. **Уникальность токена**: Каждое приглашение имеет уникальный токен
4. **Единственное использование**: Токен можно использовать только один раз
5. **Целевой консультант**: Приглашение можно принять только тем консультантом, кому оно адресовано
6. **Активный доступ**: Доступ остаётся до отзыва (revoked_at = NULL)
7. **Email stub**: Email уведомления пока не отправляются (stub)

## Будущие улучшения

- [ ] Реальная отправка email через SMTP/SendGrid
- [ ] Отзыв доступа владельцем
- [ ] История приглашений
- [ ] Рейтинг консультантов
- [ ] Отзывы от владельцев
- [ ] Портфолио консультантов (фото, сертификаты)

## Связанные модули

- [Пользователи](./users.md) - базовая информация консультанта
- [Собаки](./dogs.md) - доступ к собакам клиентов
- [События](./events.md) - создание событий для собак
- [Заметки консультантов](./consultant-notes.md) - профессиональные записи
