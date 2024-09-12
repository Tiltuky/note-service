FROM golang:1.22-alpine as builder

# Установка PostgreSQL клиента и bash
RUN apk update && apk add --no-cache postgresql-client bash

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем все файлы в контейнер
COPY . .

# Делаем скрипт wait-for-postgres.sh исполняемым
RUN chmod +x wait-for-postgres.sh

# Загружаем зависимости
RUN go mod download

# Сборка приложения
RUN go build -o note-service ./cmd/main.go

# Начинаем новую стадию сборки на основе минимального образа
FROM alpine:latest

# Установка минимальных необходимых утилит для запуска приложения
RUN apk --no-cache add postgresql-client

COPY --from=builder /app/note-service /note-service
COPY ./config/config.yml /config/config.yml
COPY ./docs/swagger.json /docs/swagger.json
COPY ./.env /.env
COPY ./wait-for-postgres.sh /wait-for-postgres.sh

# Устанавливаем права на скрипт
RUN chmod +x /wait-for-postgres.sh

# Запуск приложения
CMD ["./note-service"]
