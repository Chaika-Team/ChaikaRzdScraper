# Этап 1: Сборка
FROM golang:1.23.4-alpine AS builder

# Установка необходимых зависимостей для работы с Go
RUN apk add --no-cache git

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod tidy

# Копируем весь исходный код в контейнер
COPY . .

# Компилируем приложение
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o rzd-scraper cmd/rzd-scraper/main.go

# Этап 2: Запуск
FROM alpine:latest

# Установка необходимых библиотек для работы с Go-программой
RUN apk add --no-cache ca-certificates

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем скомпилированный файл из предыдущего этапа
COPY --from=builder /app/rzd-scraper .

# Устанавливаем переменные окружения для конфигурации
ENV RZD_LANGUAGE=ru
ENV RZD_TIMEOUT=5
ENV RZD_MAX_RETRIES=10
ENV RZD_RID_LIFETIME=300000
ENV RZD_PROXY=""
ENV RZD_USER_AGENT="Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/113.0"
ENV RZD_BASE_PATH="https://pass.rzd.ru/"
ENV RZD_DEBUG_MODE=false
ENV GRPC_PORT=50051

# Открываем порт для gRPC сервера
EXPOSE 50051

# Запускаем приложение
CMD ["./rzd-scraper"]
