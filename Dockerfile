# Этап 1: Сборка
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем главный сервер
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /app/server ./cmd/server/main.go

# Собираем embedder (для RAG)
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /app/embedder ./cmd/embedder/main.go


# Этап 2: Запуск
FROM alpine:latest

WORKDIR /root/

# Копируем бинарный файл сервера
COPY --from=builder /app/server .

# Копируем embedder (хотя он обычно запускается отдельно)
COPY --from=builder /app/embedder .

EXPOSE 8000

# Запускаем приложение
CMD ["./server"]