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
echo "1) build"
echo "2) up"
echo "3) down"

read -p "Введите соответствующий номер операции над файлом (1-3): " command

case $command in
  1) COMMAND="build" ;;
  2) COMMAND="up" ;;
  3) COMMAND="down" ;;
  *) echo "Неверный номер операции"; exit 1 ;;
esac

if [ "$COMMAND" = "build" ]; then
  docker-compose -f docker-compose.yml -f "$C_FILE_PATH" build
elif [ "$COMMAND" = "up" ]; then
  docker-compose -f docker-compose.yml -f "$C_FILE_PATH" up
elif [ "$COMMAND" = "down" ]; then
  docker-compose -f docker-compose.yml -f "$C_FILE_PATH" down
fi