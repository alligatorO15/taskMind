# Этап 1: Сборка Go-приложения
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Копируем файлы зависимостей для кеширования слоёв Docker
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код и собираем бинарник
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /taskmind ./cmd/server

# Этап 2: Минимальный образ для запуска
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Копируем собранный бинарник и конфигурацию
COPY --from=builder /taskmind .
COPY config.yml .

EXPOSE 8080

CMD ["./taskmind"]
