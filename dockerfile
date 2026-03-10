# Этап 1: Сборка приложения
FROM golang:1.24-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum (если есть)
COPY go.mod go.sum* ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Компилируем приложение для Linux
RUN GOOS=linux GOARCH=amd64 go build -o myfinproject ./main.go

# Этап 2: Создание финального образа
FROM alpine:latest

# Устанавливаем необходимые пакеты
RUN apk --no-cache add ca-certificates tzdata

# Создаем пользователя для запуска приложения (безопасность)
RUN adduser -D -h /app appuser

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем скомпилированное приложение из этапа сборки
COPY --from=builder /app/myfinproject .

# Копируем веб-директорию с фронтендом
COPY --from=builder /app/web ./web

# Создаем директорию для данных и логов
RUN mkdir -p /app/data /app/logs && chown -R appuser:appuser /app

# Переключаемся на непривилегированного пользователя
USER appuser

# Указываем порт, который будет слушать приложение
EXPOSE 7540

# Переменные окружения (можно переопределить при запуске)
ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/data/scheduler.db
# ENV TODO_PASSWORD=your_password_here # раскомментируйте для включения аутентификации
# ENV TODO_SECRET=your_secret_key_here # секретный ключ для JWT

# Команда для запуска приложения
CMD ["./myfinproject"]