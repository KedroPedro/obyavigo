# Обновления функционала

## 1. Управление ролями администраторов и модераторов

### Существующая функция
В `internal/api/adminHandlers.go` уже есть функция `UpdateUserRole()`, которая позволяет администраторам:
- Назначать роли: `admin`, `moderator`, `user`
- Снимать администраторов и модераторов с должности (устанавливать роль `user`)

### Использование
```
POST /api/admin/users/{id}/role/
Authorization: Bearer {token}
Content-Type: application/json

{
  "role": "user"
}
```

**Ограничения доступа:**
- Только администраторы (`role == "admin"`) могут менять роли пользователей
- Модераторы не имеют доступа к этой функции

## 2. Загрузка и управление аватарками пользователей

### Добавленные функции

#### 2.1. Загрузка аватарки
**Endpoint:** `POST /api/profile/upload-avatar/`

**Описание:** Загружает аватарку пользователя в MongoDB GridFS (в bucket "avatars"). При загрузке новой аватарки старая автоматически удаляется.

**Запрос:**
```bash
curl -X POST https://your-domain/api/profile/upload-avatar/ \
  -H "Authorization: Bearer {token}" \
  -F "avatar=@/path/to/avatar.jpg"
```

**Ответ:**
```json
{
  "message": "Аватар успешно загружен",
  "avatar_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Ограничения:**
- Максимальный размер файла: 5 МБ
- Требуется авторизация

#### 2.2. Получение аватарки
**Endpoint:** `GET /api/avatars/{id}/`

**Описание:** Получает изображение аватарки по ID из MongoDB.

**Пример:**
```
GET /api/avatars/550e8400-e29b-41d4-a716-446655440000/
```

**Ответ:** Изображение в формате JPEG/PNG с соответствующими заголовками

#### 2.3. Профиль пользователя с аватаркой
**Endpoint:** `GET /api/user/profile/`

**Описание:** Возвращает данные профиля с добавленным полем `avatar_id`.

**Ответ:**
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "phone_number": "+1234567890",
  "role": "user",
  "avatar_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

## 3. Удаление изображений при удалении аккаунта

### Обновленная функция DeleteAccountHandler

При удалении аккаунта теперь автоматически удаляются:
1. **Аватар пользователя** из bucket "avatars" в MongoDB
2. **Все изображения объявлений** пользователя из bucket "fs" в MongoDB
3. **Запись пользователя** из PostgreSQL

**Endpoint:** `POST /api/profile/delete-account/`

**Процесс удаления:**
```
1. Удаление аватара (DeleteUserAvatar)
2. Удаление всех изображений объявлений (поиск по metadata.user_id)
3. Удаление аккаунта из PostgreSQL
4. Очистка cookie авторизации
```

**Примечание:** Если удаление изображений не удается, процесс продолжается с логированием ошибки, чтобы не блокировать удаление аккаунта.

## Технические детали

### База данных

#### PostgreSQL
- Поле `profile_picture_id` в таблице `site.users` хранит UUID аватарки
- Методы:
  - `UpdateUserAvatar(userId, avatarID)` - обновление ID аватарки
  - `GetUserAvatar(userId)` - получение ID аватарки

#### MongoDB GridFS
- **Bucket "avatars"** - хранит аватарки пользователей
- **Bucket "fs"** - хранит изображения объявлений
- Metadata для аватарок:
  ```json
  {
    "user_id": "uuid",
    "length": 12345,
    "uploadDate": "2024-11-03T13:00:00Z"
  }
  ```

### Методы MongoDB

#### UploadUserAvatar
```go
func (m *Mongo) UploadUserAvatar(ctx context.Context, file *multipart.FileHeader, userID string) (string, error)
```

#### DeleteUserAvatar
```go
func (m *Mongo) DeleteUserAvatar(ctx context.Context, userID string) error
```

#### DeleteUserImages
```go
func (m *Mongo) DeleteUserImages(ctx context.Context, userID string) error
```

## Безопасность

1. **Авторизация**: Все endpoints требуют авторизации через JWT токен
2. **Размер файлов**: Ограничение в 5 МБ для аватарок
3. **Права доступа**: 
   - Пользователи могут управлять только своей аватаркой
   - Только администраторы могут менять роли
4. **Кэширование**: Изображения кэшируются на 1 год (Cache-Control: max-age=31536000)

## Миграция данных

Поле `profile_picture_id` уже существует в модели `User`, поэтому дополнительная миграция БД не требуется.
