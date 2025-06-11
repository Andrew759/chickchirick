#!/bin/bash

echo "Выберите окружение:"
echo "1) prod"
echo "2) dev"
echo "3) stage"

read -p "Введите соответствующий ему номер (1-3): " choice

case $choice in
  1)
    FILE="docker-compose.dev.yml"
    ;;
  2)
    FILE="docker-compose.prod.yml"
    ;;
  3)
    FILE="docker-compose.stage.yml"
    ;;
  *)
    echo "Неверный номер окружения"
    exit 1
    ;;
esac

echo "Выберите действие с файлом $FILE..."
echo "1) build"
echo "2) up"

read -p "Введите соответствующий номер операции над файлом (1-2): " command

case $command in
  1)
    COMMAND="build"
    ;;
  2)
    COMMAND="up"
    ;;
  *)
    echo "Неверный номер операции"
    exit 1
    ;;
esac

docker-compose -f docker-compose.yml -f "$FILE" "$COMMAND"