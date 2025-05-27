# syntax=docker/dockerfile:1
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum, чтобы кешировать зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Сборка бинарника
RUN CGO_ENABLED=0 GOOS=linux go build -o people-enricher ./cmd/main.go

# Финальный образ - минимальный
FROM alpine:latest

WORKDIR /app

# Копируем бинарь из builder-а
COPY --from=builder /app/people-enricher .

# Копируем .env (если есть)
COPY .env .

# Порт, который слушает сервис
EXPOSE 8080

# Запуск сервиса
CMD ["./people-enricher"]
