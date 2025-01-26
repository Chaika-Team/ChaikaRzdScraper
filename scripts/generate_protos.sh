#!/bin/bash

# Получаем путь до директории, где расположен этот скрипт
SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# Переходим в корневую директорию проекта
PROJECT_DIR=$(cd "$SCRIPT_DIR/.." && pwd)

# Указываем относительный путь до proto-файла
PROTO_DIR="$PROJECT_DIR/proto"
OUT_DIR="$PROJECT_DIR/internal/transports/grpc/pb"

# Генерация Protobuf файлов
protoc -I="$PROTO_DIR" \
    --go_out="$OUT_DIR" --go_opt=paths=source_relative \
    --go-grpc_out="$OUT_DIR" --go-grpc_opt=paths=source_relative \
    "$PROTO_DIR/rzd_service.proto"
