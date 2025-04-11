#!/bin/bash

# Скрипт для генерации Protobuf файлов на основе proto-файлов

set -euo pipefail
# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "Error: protoc is not installed"
    exit 1
fi
# Получаем путь до директории, где расположен этот скрипт
SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# Переходим в корневую директорию проекта
PROJECT_DIR=$(cd "$SCRIPT_DIR/.." && pwd)
# Указываем относительный путь до proto-файла
PROTO_DIR="$PROJECT_DIR/proto/rzd"
OUT_DIR="$PROJECT_DIR/internal/transports/grpc/pb"
# Validate directories exist
for dir in "$PROTO_DIR" "$OUT_DIR"; do
    if [ ! -d "$dir" ]; then
        echo "Error: Directory $dir does not exist"
        exit 1
    fi
done
# Генерация Protobuf файлов
if ! protoc -I="$PROTO_DIR" \
    --go_out="$OUT_DIR" --go_opt=paths=source_relative \
    --go-grpc_out="$OUT_DIR" --go-grpc_opt=paths=source_relative \
    "$PROTO_DIR/rzd_service.proto"; then
    echo "Error: Failed to generate protobuf files"
    exit 1
fi
echo "Successfully generated protobuf files"