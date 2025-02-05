# Используем официальный образ Go
FROM golang:1.23.1 as builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарник
RUN go build -o main .

# Финальный образ
FROM debian:latest

WORKDIR /app

# Копируем бинарник из сборочного контейнера
COPY --from=builder /app/main .

# Добавляем entrypoint-скрипт
COPY docker-entrypoint.sh .

# Даем права на выполнение
RUN chmod +x docker-entrypoint.sh

# Запуск через entrypoint
ENTRYPOINT ["/app/docker-entrypoint.sh"]
