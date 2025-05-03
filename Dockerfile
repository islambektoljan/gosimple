# Сборка
FROM golang:1.24-alpine AS build
WORKDIR /app

# Копируем файлы go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod tidy

# Копируем весь проект
COPY . ./

# Строим приложение
RUN go build -o main .

# Финальный образ
FROM alpine:latest
RUN apk --no-cache add ca-certificates bash

# Копируем исполняемый файл из предыдущего этапа сборки
COPY --from=build /app/main /main

# Копируем папку с миграциями в контейнер
COPY internal/db/migrations /app/internal/db/migrations

# Копируем скрипт ожидания в контейнер
COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

# Устанавливаем рабочую директорию и указываем переменные окружения
WORKDIR /app

# Переменные окружения и порт
ENV PORT=8080
EXPOSE 8080

# Запуск приложения с ожиданием готовности базы данных
CMD ["/wait-for-it.sh", "db:5432", "--", "/main"]
