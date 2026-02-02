# Используем официальный образ Go для сборки
FROM golang:1.22.2-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o HabitFlow .

# Используем легковесный образ для запуска
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /app

# Устанавливаем tzdata для поддержки часовых поясов
RUN apk add --no-cache tzdata

# Копируем скомпилированный бинарник из builder
COPY --from=builder /app/HabitFlow .

# Копируем необходимые файлы (шаблоны и конфиги)
COPY templates ./templates

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./HabitFlow"]
