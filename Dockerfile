# ---------- Стадия сборки ----------
FROM golang:1.24-alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED=0
ENV GOOS=linux

# Установка таймзоны и зависимостей
RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /app

# Копируем зависимости
COPY go.mod .
COPY go.sum .
RUN go mod download

# Копируем исходники
COPY . .

# Сборка бинарника из cmd/main.go
RUN go build -ldflags="-s -w" -o /app/server ./cmd/main.go

# ---------- Финальный образ ----------
FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates

# Настраиваем часовой пояс
COPY --from=builder /usr/share/zoneinfo/America/New_York /usr/share/zoneinfo/America/New_York
ENV TZ=America/New_York

WORKDIR /app

# Копируем собранный бинарник
COPY --from=builder /app/server .
COPY .env .

# Запуск
CMD ["./server"]


