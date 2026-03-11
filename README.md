# Планировщик задач

Веб-приложение для управления задачами с поддержкой повторяющихся событий. Позволяет создавать, редактировать, удалять задачи и отмечать их как выполненные. Поддерживает различные правила повторения задач. Реализовано на Go с использованием SQLite для хранения данных и JWT для аутентификации.

## Функциональность

- ✅ Создание, редактирование и удаление задач
- ✅ Установка даты выполнения и правила повторения
- ✅ Поддержка различных правил повторения:
  - `d N` - каждые N дней
  - `y` - каждый год
  - `w дни` - по дням недели (1-7, где 1=пн, 7=вс)
  - `m дни [месяцы]` - по дням месяца
- ✅ Отметка задач как выполненных
- ✅ Поиск задач по заголовку, комментарию или дате
- ✅ Аутентификация (опционально)
- ✅ Сохранение данных в SQLite
- ✅ Docker-образ для простого развертывания

## Структура проекта

```
myfinproject/
├── main.go                 # Точка входа, конфигурация, запуск сервера
├── go.mod                  # Зависимости проекта
├── go.sum                  # Контрольные суммы зависимостей
├── README.md               # Документация
├── Dockerfile              # Многоступенчатая сборка Docker-образа
├── docker-compose.test.yml # Docker Compose для тестирования
├── web/                    # Статические файлы фронтенда
│   ├── index.html          # Главная страница
│   ├── css/                # Стили
│   └── js/                 # JavaScript
├── pkg/                    # Пакеты приложения
│   ├── api/                 # HTTP обработчики
│   │   ├── api.go           # Регистрация маршрутов
│   │   ├── middleware.go    # Аутентификация
│   │   ├── signin.go        # Вход в систему
│   │   ├── nextDate.go      # API для вычисления дат
│   │   ├── task_add.go      # Добавление задачи
│   │   ├── task_handlers.go # Получение/обновление задачи
│   │   ├── task_list.go     # Список задач
│   │   ├── task_done.go     # Отметка о выполнении
│   │   └── task_delete.go   # Удаление задачи
│   ├── auth/                 # Аутентификация и JWT
│   │   └── auth.go           # Генерация и проверка токенов
│   ├── db/                   # Работа с базой данных
│   │   ├── db.go             # Подключение к SQLite
│   │   └── task_db.go        # CRUD операции с задачами
│   ├── nextdate/             # Логика повторения задач
│   │   └── calc_nextDate.go  # Вычисление следующих дат
│   └── server/                # HTTP сервер
│       └── server.go          # Запуск и остановка сервера
└── tests/                    # Интеграционные тесты
    ├── addtask_4_test.go
    ├── app_1_test.go
    ├── db_2_test.go
    ├── nextdate_3_test.go
    ├── task_6_test.go
    ├── task_7_test.go
    ├── tasks_5_test.go
    └── settings.go           # Настройки тестов
```

## Запуск локально

### Предварительные требования
- Go 1.24 или выше
- SQLite

### Переменные окружения

Создайте файл `.env` (опционально):

```env
TODO_PORT=7540              # Порт сервера (по умолчанию 7540)
TODO_DBFILE=scheduler.db     # Путь к файлу БД
TODO_PASSWORD=secret         # Пароль для аутентификации
TODO_SECRET=secret-key       # Секретный ключ для JWT
TODO_WEB_DIR=./web           # Путь к веб-файлам
```

### Запуск

```bash
# Установка зависимостей
go mod download

# Запуск сервера
go run main.go

# Или с компиляцией
go build -o myfinproject ./main.go
./myfinproject
```

После запуска откройте браузер по адресу: **http://localhost:7540**

### Запуск контейнера

```bash
# Создание директорий для данных
mkdir -p data logs

# Запуск контейнера
docker run -d \
  --name myfin-scheduler \
  -p 7540:7540 \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/logs:/app/logs \
  myfinproject:latest

# Просмотр логов
docker logs -f myfin-scheduler
```

### Для Windows PowerShell

```powershell
docker run -d `
  --name myfin-scheduler `
  -p 7540:7540 `
  -v ${PWD}/data:/app/data `
  -v ${PWD}/logs:/app/logs `
  myfinproject:latest
```

### Запуск с аутентификацией

```bash
docker run -d \
  --name myfin-scheduler \
  -p 7540:7540 \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/logs:/app/logs \
  -e TODO_PASSWORD="your_password" \
  -e TODO_SECRET="your_secret_key" \
  myfinproject:latest
```

## Docker Hub

Готовый образ доступен на Docker Hub:

```bash
# Скачать образ
docker pull vgdryukov/myfinproject:latest

# Запустить
docker run -d --name myfin-scheduler -p 7540:7540 vgdryukov/myfinproject:latest
```

## API Endpoints

| Метод | Путь | Описание | Аутентификация |
|-------|------|----------|----------------|
| POST | `/api/signin` | Вход в систему | Нет |
| GET | `/api/nextdate` | Вычислить следующую дату | Да |
| POST | `/api/task` | Создать задачу | Да |
| GET | `/api/task?id=` | Получить задачу | Да |
| PUT | `/api/task` | Обновить задачу | Да |
| DELETE | `/api/task?id=` | Удалить задачу | Да |
| POST | `/api/task/done?id=` | Отметить выполненной | Да |
| GET | `/api/tasks?search=` | Список задач | Да |

## Автор

[vgdryukov](https://github.com/vgdryukov)

## Лицензия

MIT