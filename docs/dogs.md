# Собаки (Dogs)

## Обзор

Модуль управления информацией о собаках. Каждая собака привязана к владельцу (owner) и может иметь доступы для консультантов.

## Структура данных

### Dog (Собака)

```go
type Dog struct {
    ID        uint      // Уникальный идентификатор
    Name      string    // Кличка собаки
    Breed     string    // Порода
    BirthDate time.Time // Дата рождения
    OwnerID   uint      // ID владельца
    Owner     *User     // Связь с владельцем
    CreatedAt time.Time // Дата создания записи
    UpdatedAt time.Time // Дата обновления записи
}
```

## Бизнес-процессы

### 1. Создание собаки

**Endpoint**: `POST /api/v1/dogs`

**Права доступа**: Owner, Admin

**Бизнес-логика**:
1. Пользователь отправляет данные о собаке
2. Система автоматически привязывает собаку к текущему пользователю (из JWT токена)
3. Создаётся запись в таблице `dogs`
4. Возвращается созданная собака

**Валидация**:
- `name`: обязательное поле, макс 100 символов
- `breed`: обязательное поле, макс 100 символов
- `birth_date`: обязательное поле, формат RFC3339

**Пример запроса**:
```json
{
  "name": "Рекс",
  "breed": "Немецкая овчарка",
  "birth_date": "2020-05-15T00:00:00Z"
}
```

**Пример ответа**:
```json
{
  "id": 1,
  "name": "Рекс",
  "breed": "Немецкая овчарка",
  "birth_date": "2020-05-15T00:00:00Z",
  "owner_id": 5,
  "created_at": "2025-11-23T10:00:00Z",
  "updated_at": "2025-11-23T10:00:00Z"
}
```

### 2. Получение списка собак

**Endpoint**: `GET /api/v1/dogs`

**Права доступа**: All authenticated

**Бизнес-логика (RBAC)**:
- **Owner**: Видит только своих собак
- **Consultant**: Видит только собак, к которым есть доступ (через `consultant_access`)
- **Admin**: Видит всех собак

**Реализация**:
```go
// Для владельца
query.Where("owner_id = ?", userID)

// Для консультанта
query.Joins("INNER JOIN consultant_access ON consultant_access.dog_id = dogs.id").
     Where("consultant_access.consultant_id = ?", userID)

// Для админа
// Без фильтрации
```

**Пример ответа**:
```json
[
  {
    "id": 1,
    "name": "Рекс",
    "breed": "Немецкая овчарка",
    "birth_date": "2020-05-15T00:00:00Z",
    "owner_id": 5
  },
  {
    "id": 2,
    "name": "Бобик",
    "breed": "Лабрадор",
    "birth_date": "2019-03-20T00:00:00Z",
    "owner_id": 5
  }
]
```

### 3. Получение собаки по ID

**Endpoint**: `GET /api/v1/dogs/:id`

**Права доступа**: Owner (своих), Consultant (с доступом), Admin (всех)

**Бизнес-логика**:
1. Извлекается ID собаки из URL
2. Проверяется право доступа:
   - Владелец: `dog.owner_id == user_id`
   - Консультант: проверка в таблице `consultant_access`
   - Админ: без проверки
3. При отсутствии прав - **404 Not Found** (не 403, чтобы не раскрывать существование ресурса)
4. Возвращается собака с информацией о владельце

**Пример ответа**:
```json
{
  "id": 1,
  "name": "Рекс",
  "breed": "Немецкая овчарка",
  "birth_date": "2020-05-15T00:00:00Z",
  "owner_id": 5,
  "owner": {
    "id": 5,
    "name": "Иван Петров",
    "email": "ivan@example.com"
  }
}
```

**Ошибки**:
- 404 - Собака не найдена или нет доступа
- 401 - Не авторизован

### 4. Обновление собаки

**Endpoint**: `PUT /api/v1/dogs/:id`

**Права доступа**: Owner (только свои), Admin (все)

**Бизнес-логика**:
1. Проверяется право доступа (только владелец или админ)
2. Консультанты **НЕ МОГУТ** редактировать собак (даже с доступом)
3. Обновляются поля: `name`, `breed`, `birth_date`
4. Поле `UpdatedAt` обновляется автоматически

**Валидация**:
- Можно обновить любое поле или все сразу
- Все поля опциональные (частичное обновление)

**Пример запроса**:
```json
{
  "name": "Рекс Великолепный",
  "breed": "Немецкая овчарка"
}
```

**Ошибки**:
- 404 - Собака не найдена или нет доступа
- 403 - Консультант пытается обновить собаку
- 400 - Неверный формат данных

### 5. Удаление собаки

**Endpoint**: `DELETE /api/v1/dogs/:id`

**Права доступа**: Owner (только свои), Admin (все)

**Бизнес-логика**:
1. Проверяется право доступа (только владелец или админ)
2. Консультанты **НЕ МОГУТ** удалять собак
3. При удалении собаки **каскадно удаляются**:
   - Все события (`events`)
   - Все записи доступа консультантов (`consultant_access`)
   - Все заметки консультантов (`consultant_notes`)
4. Возвращается статус 204 No Content

**Каскадное удаление** (на уровне БД):
```sql
-- events
ON DELETE CASCADE

-- consultant_access
ON DELETE CASCADE

-- consultant_notes
ON DELETE CASCADE
```

**Ошибки**:
- 404 - Собака не найдена или нет доступа
- 403 - Консультант пытается удалить собаку

## Контроль доступа (RBAC)

### Таблица прав доступа

| Операция | Owner | Consultant | Admin |
|----------|-------|------------|-------|
| Создать собаку | ✅ Свою | ❌ | ✅ |
| Список собак | ✅ Своих | ✅ С доступом | ✅ Всех |
| Получить собаку | ✅ Свою | ✅ С доступом | ✅ Любую |
| Обновить собаку | ✅ Свою | ❌ | ✅ Любую |
| Удалить собаку | ✅ Свою | ❌ | ✅ Любую |

### Предоставление доступа консультанту

Доступ предоставляется через систему приглашений (см. [Консультанты](./consultants.md)):

1. Владелец создаёт приглашение для консультанта
2. Консультант принимает приглашение
3. Создаётся запись в `consultant_access`:
```sql
INSERT INTO consultant_access (consultant_id, dog_id, granted_at)
VALUES (10, 1, NOW());
```
4. Консультант получает доступ на чтение и создание заметок

### Проверка доступа консультанта

```go
// В DogRepository
func (r *dogRepository) HasConsultantAccess(consultantID, dogID uint) (bool, error) {
    var count int64
    err := r.db.Model(&models.ConsultantAccess{}).
        Where("consultant_id = ? AND dog_id = ?", consultantID, dogID).
        Where("revoked_at IS NULL").
        Count(&count).Error
    return count > 0, err
}
```

## Связанные сущности

### События (Events)
- Каждое событие привязано к собаке через `dog_id`
- При удалении собаки события удаляются каскадно
- См. [События](./events.md)

### Заметки консультантов (Consultant Notes)
- Консультанты создают заметки для собак с доступом
- При удалении собаки заметки удаляются каскадно
- См. [Заметки](./consultant-notes.md)

### Доступы консультантов (Consultant Access)
- Таблица связи many-to-many между консультантами и собаками
- Поля: `consultant_id`, `dog_id`, `granted_at`, `revoked_at`
- См. [Консультанты](./consultants.md)

## База данных

### Схема таблицы

```sql
CREATE TABLE dogs (
    id SERIAL PRIMARY KEY,
    owner_id INTEGER NOT NULL REFERENCES users(id),
    name VARCHAR(100) NOT NULL,
    breed VARCHAR(100) NOT NULL,
    birth_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_dogs_owner_id ON dogs(owner_id);
```

### Индексы
- `idx_dogs_owner_id` - для быстрого поиска собак владельца

## Примеры использования

### Создание и управление собакой

```bash
# 1. Создать собаку
curl -X POST http://localhost:8080/api/v1/dogs \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Рекс",
    "breed": "Немецкая овчарка",
    "birth_date": "2020-05-15T00:00:00Z"
  }'

# 2. Получить список своих собак
curl -X GET http://localhost:8080/api/v1/dogs \
  -H "Authorization: Bearer $TOKEN"

# 3. Обновить собаку
curl -X PUT http://localhost:8080/api/v1/dogs/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Рекс Великолепный"
  }'

# 4. Удалить собаку
curl -X DELETE http://localhost:8080/api/v1/dogs/1 \
  -H "Authorization: Bearer $TOKEN"
```

## Бизнес-правила

1. **Владение**: Собака может иметь только одного владельца
2. **Доступ**: Консультант может иметь доступ к нескольким собакам
3. **Неизменность владельца**: Владельца собаки нельзя изменить (только через админа в БД)
4. **Защита данных**: При отсутствии прав возвращается 404 (не 403) для защиты от enumerate
5. **Каскадные удаления**: При удалении собаки удаляются все связанные данные

## Связанные модули

- [События](./events.md) - трекинг событий собак
- [Консультанты](./consultants.md) - система доступа
- [Заметки консультантов](./consultant-notes.md) - профессиональные записи
