# Chaika RZD Scraper (В разработке)

## Описание проекта

Chaika RZD Scraper — это инструмент для получения информации о железнодорожных маршрутах, составах поездов и станциях на
основе данных, предоставляемых РЖД через их API.

Проект использует gRPC для обмена данными и предоставляет API для следующих функций:

- Получение маршрутов поездов.
- Получение информации о вагонах поезда.
- Поиск станций по части названия.

## Установка и настройка

### Требования

- Go 1.18 или выше.

### Установка

1. Склонируйте репозиторий:

    ```bash
    git clone https://github.com/Chaika-Team/ChaikaRzdScraper.git
    cd rzd-scraper
    ```

2. Установите зависимости:

    ```bash
    go mod tidy
    ```

3. Настройте конфигурационный файл `config.yml` для подключения к API РЖД. Пример конфигурации:

    ```yaml
    RZD:
      LANGUAGE: "ru"
      BASE_PATH: "https://pass.rzd.ru/"
      USER_AGENT: "Mozilla/5.0 (compatible; RzdClient/1.0)"
      TIMEOUT: 2000
      MAX_RETRIES: 10
      RID_LIFETIME: 300000
      PROXY: ""
    GRPC:
      PORT: "50051"
    ```

4. Запустите сервер gRPC:

    ```bash
    go run cmd/rzd-scraper/main.go
    ```

## Примеры использования

После запуска сервера можно делать запросы к его gRPC API.

### Пример запроса маршрутов

Запрос для получения маршрутов поездов между станциями с кодами `2004000` и `2000000`:

```protobuf
// Пример запроса для получения маршрутов поездов
    service.RzdService.GetTrainRoutes({
FromCode: 2004000,
    ToCode: 2000000,
    Direction: 0, // OneWay
    TrainType: 1, // AllTrains
    CheckSeats: false,
FromDate: "2025-04-14",
    WithChange: false
    });
```

### Пример запроса информации о вагонах

Запрос для получения информации о вагонах для поезда с номером `119А`:

```protobuf
// Пример запроса для получения информации о вагонах
    service.RzdService.GetTrainCarriages({
TrainNumber: "119А",
    Direction: 0, // OneWay
    FromCode: 2004000,
    FromTime: "2025-04-14T10:00:00",
    ToCode: 2000000
    });
```

### Пример поиска станции

Запрос для поиска станций, содержащих строку "ЧЕБ":

```protobuf
// Пример запроса для поиска станций
    service.RzdService.SearchStation({
Query: "ЧЕБ",
    CompactMode: true,
    Lang: "ru"
    });
```

## Тестирование

В проекте предусмотрены e2e тесты для проверки функциональности API. Для запуска тестов выполните следующую команду:

```bash
go test ./...
```

## Лицензия

Этот проект распространяется под лицензией [GPL-3.0](LICENSE).

## Контакты

- Email: chaika.contact@yandex.ru
