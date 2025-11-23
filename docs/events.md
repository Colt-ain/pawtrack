# События (Events)

## Обзор

Модуль для отслеживания событий в жизни собак. Владельцы и консультанты могут записывать различные события: прогулки, кормление, приём лекарств, тренировки, посещения ветеринара и прочее.

## Структура данных

### Event (Событие)

```go
type Event struct {
    ID        uint       // Уникальный идентификатор
    DogID     *uint      // ID собаки (опционально для общих событий)
    Dog       *Dog       // Связь с собакой
    Type      string     // Тип события (walk, feed, meds, training, vet, note)
    Note      string     // Описание события
    At        time.Time  // Время события
    CreatedAt time.Time  // Дата создания записи
    UpdatedAt time.Time  // Дата обновления записи
}
```

## Типы событий

Система поддерживает любые типы событий (свободный ввод), но рекомендуются:

- `walk` - Прогулка
- `feed` - Кормление
- `meds` - Приём лекарств
- `training` - Тренировка
- `vet` - Посещение ветеринара
- `grooming` - Груминг
- `note` - Общая заметка

## Бизнес-процессы

### 1. Создание события

**Endpoint**: `POST /api/v1/events`

**Права доступа**: Owner, Consultant (с доступом к собаке), Admin

**Бизнес-логика**:
1. Пользователь отправляет данные события
2. Если указан `dog_id`:
   - Владелец: может создавать для своих собак
   - Консультант: может создавать только для собак с доступом
   - Админ: может создавать для любых собак
3. Если `at` не указано - используется текущее время
4. Создаётся запись в таблице `events`

**Валидация**:
- `type`: обязательное поле, макс 50 символов
- `note`: опционально, макс 255 символов
- `dog_id`: опционально (для событий без привязки к собаке)
- `at`: опционально, формат RFC3339

**Пример запроса**:
```json
{
  "dog_id": 1,
  "type": "walk",
  "note": "Утренняя прогулка в парке, 30 минут",
  "at": "2025-11-23T08:00:00Z"
}
```

**Пример ответа**:
```json
{
  "id": 15,
  "dog_id": 1,
  "type": "walk",
  "note": "Утренняя прогулка в парке, 30 минут",
  "at": "2025-11-23T08:00:00Z",
  "created_at": "2025-11-23T08:05:00Z",
  "updated_at": "2025-11-23T08:05:00Z"
}
```

### 2. Получение списка событий с фильтрацией

**Endpoint**: `GET /api/v1/events`

**Права доступа**: All authenticated

**Query параметры**:
- `dog_id` - Фильтр по собаке (только собаки с доступом)
- `types` - Фильтр по типам (через запятую): `?types=walk,feed`
- `search` - Поиск по содержимому note (ILIKE)
- `from_date` - Начало периода (RFC3339)
- `to_date` - Конец периода (RFC3339)
- `page` - Номер страницы (default: 1)
- `page_size` - Размер страницы (default: 20, max: 100)

**Бизнес-логика (RBAC)**:
- **Owner**: Видит события только своих собак
- **Consultant**: Видит события собак с доступом
- **Admin**: Видит все события

**Фильтрация**:
```go
// По собаке
if dogID > 0 {
    query.Where("dog_id = ?", dogID)
}

// По типам
if len(types) > 0 {
    query.Where("type IN ?", types)
}

// Поиск по тексту
if search != "" {
    query.Where("note ILIKE ?", "%"+search+"%")
}

// Период
if fromDate != "" {
    query.Where("at >= ?", fromDate)
}
if toDate != "" {
    query.Where("at <= ?", toDate)
}
```

**Пример запроса**:
```
GET /api/v1/events?dog_id=1&types=walk,feed&from_date=2025-11-01T00:00:00Z
```

**Пример ответа**:
```json
{
  "events": [
    {
      "id": 15,
      "dog_id": 1,
      "type": "walk",
      "note": "Утренняя прогулка",
      "at": "2025-11-23T08:00:00Z"
    },
    {
      "id": 14,
      "dog_id": 1,
      "type": "feed",
      "note": "Завтрак, 200г",
      "at": "2025-11-23T07:00:00Z"
    }
  ],
  "page": 1,
  "page_size": 20,
  "total_count": 2,
  "total_pages": 1
}
```

### 3. Получение события по ID

**Endpoint**: `GET /api/v1/events/:id`

**Права доступа**: Owner (своих собак), Consultant (с доступом), Admin (всех)

**Бизнес-логика**:
1. Проверяется доступ к собаке, с которой связано событие
2. Если событие без собаки (`dog_id = NULL`) - доступно всем
3. Возвращается событие с полной информацией о собаке

**Пример ответа**:
```json
{
  "id": 15,
  "dog_id": 1,
  "dog": {
    "id": 1,
    "name": "Рекс",
    "breed": "Немецкая овчарка"
  },
  "type": "walk",
  "note": "Утренняя прогулка",
  "at": "2025-11-23T08:00:00Z",
  "created_at": "2025-11-23T08:05:00Z"
}
```

### 4. Удаление события

**Endpoint**: `DELETE /api/v1/events/:id`

**Права доступа**: Owner (владелец собаки), Admin (любые)

**Бизнес-логика**:
1. Находится событие по ID
2. Проверяется доступ:
   - Владелец: только события своих собак
   - Админ: любые события
   - Консультант: **НЕ МОЖЕТ** удалять события
3. Событие удаляется из БД
4. Возвращается статус 204 No Content

**Важно**: Консультанты могут создавать события, но не могут их удалять. Это защита от случайного удаления важной информации.

**Ошибки**:
- 404 - Событие не найдено или нет доступа
- 403 - Консультант пытается удалить событие

## Контроль доступа (RBAC)

### Таблица прав

| Операция | Owner | Consultant | Admin |
|----------|-------|------------|-------|
| Создать событие | ✅ Для своих собак | ✅ Для собак с доступом | ✅ Для любых |
| Список событий | ✅ Своих собак | ✅ Собак с доступом | ✅ Всех |
| Получить событие | ✅ Своих собак | ✅ Собак с доступом | ✅ Любое |
| Удалить событие | ✅ Своих собак | ❌ | ✅ Любое |

### Фильтрация по ролям

```go
// В EventRepository
func (r *eventRepository) List(filters *EventFilterParams) ([]Event, int64, error) {
    query := r.db.Model(&models.Event{})
    
    // RBAC фильтрация
    if filters.UserRole == "owner" {
        // Только события собак владельца
        query = query.Joins("INNER JOIN dogs ON dogs.id = events.dog_id").
                     Where("dogs.owner_id = ?", filters.UserID)
    } else if filters.UserRole == "consultant" {
        // Только события собак с доступом
        query = query.Joins("INNER JOIN dogs ON dogs.id = events.dog_id").
                     Joins("INNER JOIN consultant_access ON consultant_access.dog_id = dogs.id").
                     Where("consultant_access.consultant_id = ?", filters.UserID)
    }
    // Admin - без фильтрации
    
    return query.Find(&events).Error
}
```

## Расширенная фильтрация

### Примеры использования

#### Все прогулки за последнюю неделю
```
GET /api/v1/events?types=walk&from_date=2025-11-16T00:00:00Z
```

#### Поиск событий с лекарствами
```
GET /api/v1/events?search=антибиотик
```

#### События конкретной собаки за день
```
GET /api/v1/events?dog_id=1&from_date=2025-11-23T00:00:00Z&to_date=2025-11-23T23:59:59Z
```

#### Комбинированный фильтр
```
GET /api/v1/events?dog_id=1&types=feed,meds&page=1&page_size=50
```

## Пагинация

### Параметры
- `page` - Номер страницы (начинается с 1)
- `page_size` - Количество записей на странице (default: 20, max: 100)

### Расчёт страниц
```go
totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
offset := (page - 1) * pageSize
```

### Ответ с пагинацией
```json
{
  "events": [...],
  "page": 2,
  "page_size": 20,
  "total_count": 156,
  "total_pages": 8
}
```

## База данных

### Схема таблицы

```sql
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    dog_id INTEGER REFERENCES dogs(id) ON DELETE SET NULL,
    type VARCHAR(50) NOT NULL,
    note VARCHAR(255),
    at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_events_dog_id ON events(dog_id);
CREATE INDEX idx_events_type ON events(type);
CREATE INDEX idx_events_at ON events(at);
```

### Индексы
- `idx_events_dog_id` - для фильтрации по собаке
- `idx_events_type` - для фильтрации по типу
- `idx_events_at` - для сортировки и фильтрации по дате

### Каскадные операции
- `ON DELETE SET NULL` - при удалении собаки `dog_id` становится NULL (события сохраняются)

## Примеры использования

### Ведение дневника собаки

```bash
# Утренняя прогулка
curl -X POST http://localhost:8080/api/v1/events \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "dog_id": 1,
    "type": "walk",
    "note": "Прогулка 30 мин, активная игра",
    "at": "2025-11-23T08:00:00Z"
  }'

# Кормление
curl -X POST http://localhost:8080/api/v1/events \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "dog_id": 1,
    "type": "feed",
    "note": "Утренний приём пищи, 200г корма",
    "at": "2025-11-23T09:00:00Z"
  }'

# Приём лекарств
curl -X POST http://localhost:8080/api/v1/events \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "dog_id": 1,
    "type": "meds",
    "note": "Антибиотик, 1 таблетка",
    "at": "2025-11-23T09:30:00Z"
  }'
```

### Получение истории

```bash
# Все события за сегодня
curl -X GET "http://localhost:8080/api/v1/events?dog_id=1&from_date=2025-11-23T00:00:00Z&to_date=2025-11-23T23:59:59Z" \
  -H "Authorization: Bearer $TOKEN"

# Только прогулки
curl -X GET "http://localhost:8080/api/v1/events?dog_id=1&types=walk" \
  -H "Authorization: Bearer $TOKEN"
```

## Бизнес-правила

1. **Временная метка**: Если не указана `at`, используется текущее время
2. **Привязка к собаке**: События могут быть без привязки к собаке (`dog_id = NULL`)
3. **Свободные типы**: Система не ограничивает типы событий (любая строка до 50 символов)
4. **Защита от удаления**: Консультанты не могут удалять события (только создавать)
5. **Сохранение истории**: При удалении собаки события сохраняются (`ON DELETE SET NULL`)

## Связанные модули

- [Собаки](./dogs.md) - основная сущность для привязки событий
- [Пользователи](./users.md) - создатели событий
- [Консультанты](./consultants.md) - доступ к событиям собак клиентов
