#!/bin/bash

if ! command -v docker-compose &> /dev/null; then
  echo "Ошибка: docker-compose не найден"
  exit 1
fi

echo "Выберите окружение:"
echo "1) prod"
echo "2) dev"
echo "3) stage"

read -p "Введите соответствующий ему номер (1-3): " choice

case $choice in
  1) C_FILE_PATH="build/prod/docker-compose.yml" ;;
  2) C_FILE_PATH="build/dev/docker-compose.yml" ;;
  3) C_FILE_PATH="build/stage/docker-compose.yml" ;;
  *) echo "Неверный номер окружения"; exit 1 ;;
esac

echo "Выберите действие с файлом $C_FILE_PATH..."
echo "1) build (no cache)"
echo "2) build (with cache)"
echo "3) up"
echo "4) down"
echo "5) build main file only"

read -p "Введите соответствующий номер операции над файлом (1-5): " command

case $command in
  1) COMMAND="build (no cache)" ;;
  2) COMMAND="build (with cache)" ;;
  3) COMMAND="up" ;;
  4) COMMAND="down" ;;
  5) COMMAND="build main file only" ;;
  *) echo "Неверный номер операции"; exit 1 ;;
esac

if [ "$COMMAND" = "build (no cache)" ]; then
  docker-compose -f docker-compose.yml -f "$C_FILE_PATH" build --no-cache
elif [ "$COMMAND" = "build (with cache)" ]; then
  docker-compose -f docker-compose.yml -f "$C_FILE_PATH" build
elif [ "$COMMAND" = "up" ]; then
  docker-compose -f docker-compose.yml -f "$C_FILE_PATH" up
elif [ "$COMMAND" = "down" ]; then
  docker-compose -f docker-compose.yml -f "$C_FILE_PATH" down
#TODO: временное решение с контейнером
elif [ "$COMMAND" = "build main file only" ]; then
   docker cp ./local-dir chickchirick_api_1:/app/target-dir
   docker-compose -f docker-compose.yml -f "$C_FILE_PATH" exec api go build -o main app/cmd/main.go
fi