# Заметки консультантов (Consultant Notes)

## Обзор

Модуль для ведения профессиональных заметок консультантами о собаках, с которыми они работают. Заметки поддерживают Markdown-форматирование и доступны только консультантам и администраторам.

## Структура данных

### ConsultantNote (Заметка)

```go
type ConsultantNote struct {
    ID           uint      // Уникальный идентификатор
    ConsultantID uint      // ID консультанта-автора
    Consultant   *User     // Связь с консультантом
    DogID        uint      // ID собаки
    Dog          *Dog      // Связь с собакой
    Title        string    // Заголовок заметки
    Content      string    // Содержимое (Markdown)
    CreatedAt    time.Time // Дата создания
    UpdatedAt    time.Time // Дата последнего обновления
}
```

## Бизнес-процессы

### 1. Создание заметки

**Endpoint**: `POST /api/v1/consultant-notes`

**Права доступа**: Consultant (с доступом к собаке), Admin

**Бизнес-логика**:
1. Консультант создаёт заметку для собаки
2. Система проверяет доступ к собаке через `consultant_access`
3. Если доступа нет - возвращается 403 Forbidden
4. Заметка создаётся с текущим временем
5. `Content` сохраняется как есть (Markdown)

**Валидация**:
- `dog_id`: обязательное поле, должен быть доступ
- `title`: обязательное, макс 255 символов
- `content`: обязательное, текст в формате Markdown

**Пример запроса**:
```json
{
  "dog_id": 1,
  "title": "Первая тренировка",
  "content": "# Отчёт о тренировке\n\n## Прогресс\n- Освоена команда \"сидеть\"\n- Хорошо реагирует на поощрение\n\n## Рекомендации\nПродолжить закрепление команд."
}
```

**Пример ответа**:
```json
{
  "id": 25,
  "consultant_id": 5,
  "dog_id": 1,
  "title": "Первая тренировка",
  "content": "# Отчёт о тренировке...",
  "created_at": "2025-11-23T14:00:00Z",
  "updated_at": "2025-11-23T14:00:00Z"
}
```

**Ошибки**:
- 403 - Консультант не имеет доступа к собаке
- 400 - Невалидные данные

### 2. Получение заметки по ID

**Endpoint**: `GET /api/v1/consultant-notes/:id`

**Права доступа**: 
- Consultant - только свои заметки
- Admin - все заметки

**Бизнес-логика**:
1. Система загружает заметку с информацией о собаке и владельце
2. Проверяется RBAC:
   - Консультант видит только свои заметки
   - Админ видит все заметки
3. Возвращается полная информация

**Пример ответа**:
```json
{
  "id": 25,
  "consultant_id": 5,
  "dog_id": 1,
  "dog_name": "Рекс",
  "owner_id": 3,
  "owner_name": "Иван Петров",
  "title": "Первая тренировка",
  "content": "# Отчёт о тренировке\n\n## Прогресс...",
  "created_at": "2025-11-23T14:00:00Z",
  "updated_at": "2025-11-23T14:00:00Z"
}
```

**Ошибки**:
- 404 - Заметка не найдена
- 403 - Попытка получить чужую заметку (не админ)

### 3. Обновление заметки

**Endpoint**: `PUT /api/v1/consultant-notes/:id`

**Права доступа**:
- Consultant - только свои заметки
- Admin - все заметки

**Бизнес-логика**:
1. Проверяется владение заметкой (consultant_id или admin)
2. Обновляются поля `title` и/или `content`
3. Поле `updated_at` обновляется автоматически
4. Частичное обновление поддерживается

**Валидация**:
- `title`: опционально, если указано - макс 255 символов
- `content`: опционально, если указано - не пустое

**Пример запроса**:
```json
{
  "title": "Первая тренировка (обновлено)",
  "content": "# Отчёт о тренировке\n\n## Обновлённая информация\n..."
}
```

**Ошибки**:
- 404 - Заметка не найдена
- 403 - Попытка обновить чужую заметку

### 4. Удаление заметки

**Endpoint**: `DELETE /api/v1/consultant-notes/:id`

**Права доступа**:
- Consultant - только свои заметки
- Admin - все заметки

**Бизнес-логика**:
1. Проверяется владение заметкой
2. Заметка удаляется безвозвратно
3. Возвращается 204 No Content

**Ошибки**:
- 404 - Заметка не найдена
- 403 - Попытка удалить чужую заметку

### 5. Список заметок с фильтрацией

**Endpoint**: `GET /api/v1/consultant-notes`

**Права доступа**:
- Consultant - только свои заметки
- Admin - все заметки

**Query параметры**:
- `search` - Поиск в title и content (ILIKE)
- `dog_id` - Фильтр по собаке
- `owner_id` - Фильтр по владельцу собаки
- `from_date` - Начало периода (RFC3339)
- `to_date` - Конец периода (RFC3339)
- `sort_by` - Поле сортировки: `created_at`, `updated_at`, `dog_name`, `owner_name`
- `order` - Порядок: `asc`, `desc` (default: `desc`)
- `page` - Номер страницы (default: 1)
- `page_size` - Размер страницы (default: 20, max: 100)

**Бизнес-логика (RBAC)**:
```go
// Консультант видит только свои заметки
if role == "consultant" {
    query.Where("consultant_id = ?", userID)
}

// Админ видит все заметки (без фильтрации)
```

**Фильтрация**:
```go
// Поиск по тексту
if search != "" {
    query.Where("title ILIKE ? OR content ILIKE ?", "%"+search+"%", "%"+search+"%")
}

// По собаке
if dogID > 0 {
    query.Where("dog_id = ?", dogID)
}

// По владельцу
if ownerID > 0 {
    query.Joins("LEFT JOIN dogs ON dogs.id = consultant_notes.dog_id").
          Where("dogs.owner_id = ?", ownerID)
}

// Период
if fromDate != "" {
    query.Where("created_at >= ?", fromDate)
}
if toDate != "" {
    query.Where("created_at <= ?", toDate)
}
```

**Сортировка**:
```go
sortField := "consultant_notes.created_at" // по умолчанию

switch sortBy {
case "updated_at":
    sortField = "consultant_notes.updated_at"
case "dog_name":
    sortField = "dogs.name"  // требует JOIN
case "owner_name":
    sortField = "users.name" // требует JOIN
}

order := "DESC"
if orderParam == "asc" {
    order = "ASC"
}

query.Order(sortField + " " + order)
```

**Пример запроса**:
```
GET /api/v1/consultant-notes?dog_id=1&search=тренировка&sort_by=updated_at&order=desc
```

**Пример ответа**:
```json
{
  "notes": [
    {
      "id": 25,
      "consultant_id": 5,
      "dog_id": 1,
      "dog_name": "Рекс",
      "owner_id": 3,
      "owner_name": "Иван Петров",
      "title": "Первая тренировка",
      "content": "# Отчёт о тренировке...",
      "created_at": "2025-11-23T14:00:00Z",
      "updated_at": "2025-11-23T15:30:00Z"
    }
  ],
  "page": 1,
  "page_size": 20,
  "total_count": 1,
  "total_pages": 1
}
```

## Markdown поддержка

### Сохранение
- На бэкенде Markdown сохраняется как есть (plain text)
- Никакой обработки или валидации Markdown на бэкенде
- Фронтенд отвечает за рендеринг

### Поддерживаемые элементы (фронтенд)
```markdown
# Заголовки
## H2
### H3

**Жирный текст**
*Курсив*

- Списки
- Пункты

1. Нумерованные
2. Списки

`inline code`

```
code block
```

[Ссылки](https://example.com)
```

### Рекомендации для консультантов
1. Использовать заголовки для структуры
2. Списки для прогресса и рекомендаций
3. Код-блоки для расписаний
4. Жирный текст для важного

## Контроль доступа (RBAC)

### Таблица прав

| Операция | Owner | Consultant (автор) | Consultant (другой) | Admin |
|----------|-------|-------------------|---------------------|-------|
| Создать заметку | ❌ | ✅ Для собак с доступом | ❌ | ✅ |
| Список заметок | ❌ | ✅ Только своих | ❌ | ✅ Всех |
| Получить заметку | ❌ | ✅ Только свою | ❌ | ✅ Любую |
| Обновить заметку | ❌ | ✅ Только свою | ❌ | ✅ Любую |
| Удалить заметку | ❌ | ✅ Только свою | ❌ | ✅ Любую |

### Важно: Владельцы НЕ видят заметки

Заметки консультантов - это **внутренняя профессиональная документация**:
- Владельцы **не могут** просматривать заметки
- Заметки видны только автору и админам
- Для общения с владельцем используйте **события** с типом `note`

## База данных

### Схема таблицы

```sql
CREATE TABLE consultant_notes (
    id SERIAL PRIMARY KEY,
    consultant_id INTEGER NOT NULL REFERENCES users(id),
    dog_id INTEGER NOT NULL REFERENCES dogs(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_consultant_notes_consultant_id ON consultant_notes(consultant_id);
CREATE INDEX idx_consultant_notes_dog_id ON consultant_notes(dog_id);
CREATE INDEX idx_consultant_notes_created_at ON consultant_notes(created_at);
CREATE INDEX idx_consultant_notes_updated_at ON consultant_notes(updated_at);
```

### Индексы
- `idx_consultant_notes_consultant_id` - для фильтрации по консультанту
- `idx_consultant_notes_dog_id` - для фильтрации по собаке
- `idx_consultant_notes_created_at` - для сортировки
- `idx_consultant_notes_updated_at` - для сортировки

### Каскадные операции
- `ON DELETE CASCADE` - при удалении собаки удаляются все заметки о ней

## Примеры использования

### Ведение профессионального дневника

```bash
# 1. Создать заметку после первой тренировки
curl -X POST http://localhost:8080/api/v1/consultant-notes \
  -H "Authorization: Bearer $CONSULTANT_TOKEN" \
  -d '{
    "dog_id": 1,
    "title": "Первая тренировка - 23.11.2025",
    "content": "# Отчёт о тренировке\n\n## Прогресс\n- Собака активная, внимательная\n- Успешно освоена команда \"сидеть\"\n- Реагирует на кличку\n\n## План на следующую встречу\n- Закрепить \"сидеть\"\n- Начать \"лежать\"\n- Работа с поводком"
  }'

# 2. Получить все заметки о собаке
curl -X GET "http://localhost:8080/api/v1/consultant-notes?dog_id=1" \
  -H "Authorization: Bearer $CONSULTANT_TOKEN"

# 3. Поиск по ключевому слову
curl -X GET "http://localhost:8080/api/v1/consultant-notes?search=поводок" \
  -H "Authorization: Bearer $CONSULTANT_TOKEN"

# 4. Обновить заметку
curl -X PUT http://localhost:8080/api/v1/consultant-notes/25 \
  -H "Authorization: Bearer $CONSULTANT_TOKEN" \
  -d '{
    "content": "# Отчёт о тренировке (обновлено)\n\n## Дополнение\n- Собака показала отличный результат"
  }'

# 5. Удалить заметку
curl -X DELETE http://localhost:8080/api/v1/consultant-notes/25 \
  -H "Authorization: Bearer $CONSULTANT_TOKEN"
```

### Расширенная фильтрация

```bash
# Все заметки за ноябрь 2025
curl -X GET "http://localhost:8080/api/v1/consultant-notes?from_date=2025-11-01T00:00:00Z&to_date=2025-11-30T23:59:59Z" \
  -H "Authorization: Bearer $CONSULTANT_TOKEN"

# Заметки о собаках конкретного владельца
curl -X GET "http://localhost:8080/api/v1/consultant-notes?owner_id=3" \
  -H "Authorization: Bearer $CONSULTANT_TOKEN"

# Сортировка по последнему обновлению
curl -X GET "http://localhost:8080/api/v1/consultant-notes?sort_by=updated_at&order=desc" \
  -H "Authorization: Bearer $CONSULTANT_TOKEN"

# Комбинированный фильтр
curl -X GET "http://localhost:8080/api/v1/consultant-notes?dog_id=1&search=команда&sort_by=created_at&order=asc&page=1&page_size=10" \
  -H "Authorization: Bearer $CONSULTANT_TOKEN"
```

## Шаблоны заметок

### Отчёт о тренировке
```markdown
# Тренировка {дата}

## Цели на сегодня
- [ ] Закрепить команду "сидеть"
- [ ] Начать команду "лежать"

## Выполнено
- [x] Команда "сидеть" - отлично
- [ ] Команда "лежать" - в процессе

## Наблюдения
- Собака активная, хорошо концентрируется
- Отвлекается на других собак

## План на следующую встречу
1. Продолжить "лежать"
2. Начать работу с отвлекающими факторами
```

### Медицинские наблюдения
```markdown
# Ветеринарный осмотр {дата}

## Общее состояние
Собака в хорошей форме, активная.

## Осмотр
- **Вес**: 25 кг
- **Температура**: 38.5°C (норма)
- **Зубы**: Требуется чистка

## Рекомендации
- Регулярная чистка зубов
- Контроль веса
- Следующий осмотр через 6 месяцев

## Назначения
`Таблетки от глистов - 1 раз в 3 месяца`
```

### Поведенческая консультация
```markdown
# Консультация по поведению {дата}

## Проблема
Собака лает на прохожих во время прогулок.

## Анализ
- Реакция на незнакомцев
- Защитное поведение
- Недостаток социализации

## План коррекции
1. **Неделя 1-2**: Контроль на дистанции от раздражителя
2. **Неделя 3-4**: Постепенное сближение
3. **Неделя 5+**: Положительное подкрепление спокойного поведения

## Домашнее задание владельцу
- Прогулки в спокойных местах
- Награда за спокойное поведение
- Избегать конфликтных ситуаций
```

## Бизнес-правила

1. **Доступ**: Заметку можно создать только для собаки с доступом
2. **Приватность**: Заметки видны только автору и админу
3. **Владение**: Консультант не может редактировать чужие заметки
4. **Каскад**: При удалении собаки удаляются все заметки о ней
5. **Markdown**: Поддерживается любой валидный Markdown
6. **Поиск**: Case-insensitive поиск по title и content
7. **История**: UpdatedAt обновляется при каждом изменении

## Use Cases

### 1. Кинолог ведёт дневник тренировок
- Создаёт заметку после каждой тренировки
- Отслеживает прогресс собаки
- Планирует следующие занятия
- Анализирует проблемные моменты

### 2. Ветеринар ведёт медицинскую карту
- Записывает результаты осмотров
- Ведёт историю назначений
- Отслеживает хронические заболевания
- Планирует профилактику

### 3. Грумер отмечает особенности
- Записывает предпочтения собаки
- Отмечает проблемные зоны
- Планирует следующие процедуры
- Рекомендации по уходу

## Связанные модули

- [Консультанты](./consultants.md) - система доступа к собакам
- [Собаки](./dogs.md) - привязка заметок к собакам
- [События](./events.md) - публичные записи (в отличие от приватных заметок)
- [Пользователи](./users.md) - авторы заметок
