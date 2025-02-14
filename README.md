# RZD API Client (Разработка)

Этот проект представляет собой клиент для API РЖД, написанный на Go.
Запросы были вдохновлены репозиторием [visavi/rzd-api](https://github.com/visavi/rzd-api).

## Возможности
- Получение списка маршрутов поездов
- Получение информации о вагонах
- Поддержка фильтрации по типу поезда и наличию мест

## Установка
```sh
go get github.com/yourusername/rzd-api-client
```

## Использование

Пример запроса списка маршрутов:

```go
ctx := context.Background()
client := rzd.NewClient()

params := domain.GetTrainRoutesParams{
    FromCode:   2004000,          // Санкт-Петербург
    ToCode:     2000000,          // Москва
    Direction:  domain.OneWay,    // Только туда
    TrainType:  domain.AllTrains, // Поезда и электрички
    CheckSeats: false,            // Не проверять наличие мест
    FromDate:   time.Now().Add(24 * 2 * time.Hour),
    WithChange: false, // Без пересадок
}

routes, err := client.GetTrainRoutes(ctx, params)
if err != nil {
    log.Fatalf("failed to get train routes: %v", err)
}

for _, route := range routes {
    fmt.Printf("Поезд %s типа %d из %s в %s отправляется в %s и прибывает в %s\n",
        route.TrainNumber, route.TrainType, route.From.Name, route.To.Name,
        route.Departure.Format("15:04"), route.Arrival.Format("15:04"))
    for _, car := range route.CarTypes {
        fmt.Printf("\tВагон %s %s класса, свободных мест: %d, стоимость: %d руб.\n",
            car.TypeShortLabel, car.Class, car.FreeSeats, car.Tariff)
    }
}
```

Пример запроса списка вагонов для выбранного маршрута:

```go
route := routes[1] // Выбираем первый маршрут для примера

cartParams := domain.GetTrainCarriagesParams{
    TrainNumber: route.TrainNumber,
    Direction:   domain.OneWay,
    FromCode:    route.From.Code,
    FromTime:    route.Departure,
    ToCode:      route.To.Code,
}

carriages, err := client.GetTrainCarriages(ctx, cartParams)
if err != nil {
    log.Fatalf("failed to get train carriages: %v", err)
}

for _, car := range carriages {
    fmt.Printf("Вагон %s %s, стоимость: %d руб., перевозчик: %s, свободных мест: %d\n",
        car.CategoryLabelLocal, car.CarNumber, car.Tariff, car.Carrier.Name, car.FreeSeats)
}
```

## Лицензия
Проект распространяется под лицензией **GPLv3**.

