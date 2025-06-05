# Dockerfile References: https://docs.docker.com/engine/reference/builder/

ARG GOLANG_VERSION='1.23.5'
ARG USER_UID=1000

FROM golang:${GOLANG_VERSION}-alpine as builder

#TODO Починить
#RUN apk update && apk upgrade && \
#    apk add --no-cache bash git openssh


# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum из корневой директории
COPY ./go.mod ./go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY ./ ./

# Компилируем приложение
RUN go build -o main ./cmd/main.go


# Финальный минимальный образ
FROM alpine:latest

WORKDIR /root/

# Копируем бинарник из стадии сборки
COPY --from=builder . .

# Открываем порт
EXPOSE 8080

# Команда по умолчанию
CMD ["./app/main"]
