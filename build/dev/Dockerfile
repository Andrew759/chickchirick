# Dockerfile References: https://docs.docker.com/engine/reference/builder/

ARG GOLANG_VERSION='1.24.4'
ARG USER_UID=1000

FROM golang:${GOLANG_VERSION}-alpine as builder

#Устанавливаются нужные пакеты
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

# Устанавливается рабочая директория
WORKDIR /app

# Копируется go.mod и go.sum из корневой директории
COPY ./go.mod ./go.sum ./

# Загружаются зависимости go
RUN go mod download

# Копируется исходный код
COPY ./ ./
COPY ./.env.dev .env

# Компилируется приложение
RUN go build -o main ./cmd/main.go

# Создается финальный образ
FROM alpine:latest

ENV GOLANG_VERSION=1.24.4

# Устанавливается рабочая директория
WORKDIR /root/

# Копируется бинарник из стадии сборки
COPY --from=builder . .

# Открывается порт
EXPOSE 8080

# Запуск бинарника
CMD ["./app/main"]

USER root

#Подготовка финального образа
RUN apk add --no-cache curl tar

# Установка golang
RUN apk add --no-cache curl tar git build-base && \
    curl -LO https://go.dev/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    rm -rf /usr/local/go && \
    tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    rm go${GOLANG_VERSION}.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"