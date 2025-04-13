package main

import (
    "log"
    "time"
    "github.com/streadway/amqp"
    "encoding/json"
)

func main() {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("Failed to open a channel: %v", err)
    }
    defer ch.Close()

    // Создание очереди с приоритетом
    args := amqp.Table{"x-max-priority": int32(10)}
    _, err = ch.QueueDeclare(
        "rzd_queue", // имя очереди
        true,        // durable
        false,       // delete when unused
        false,       // exclusive
        false,       // no-wait
        args,        // arguments
    )
    if err != nil {
        log.Fatalf("Failed to declare a queue: %v", err)
    }

    // Функция для отправки сообщения
    sendMessage := func(requestType string, priority int, correlationId string, payload interface{}) {
        message := map[string]interface{}{
            "requestType":   requestType,
            "priority":      priority,
            "correlationId": correlationId,
            "payload":       payload,
        }

        body, err := json.Marshal(message)
        if err != nil {
            log.Fatalf("Failed to marshal JSON: %v", err)
        }

        err = ch.Publish(
            "",           // exchange
            "rzd_queue",  // routing key
            false,        // mandatory
            false,        // immediate
            amqp.Publishing{
                ContentType: "application/json",
                Body:        body,
                Priority:    uint8(priority),
            })
        if err != nil {
            log.Fatalf("Failed to publish a message: %v", err)
        }

        log.Printf("Sent message: %s", body)
    }

    // Отправка тестовых сообщений
    sendMessage("GetTrainRoutes", 5, "uuid-getroutes-001", map[string]interface{}{
        "fromCode":   2004000,
        "toCode":     2000000,
        "direction":  0,
        "trainType":  1,
        "checkSeats": false,
        "fromDate":   "2024-10-15T00:00:00Z",
        "withChange": false,
    })

    sendMessage("GetTrainCarriages", 1, "uuid-getcarriages-001", map[string]interface{}{
        "trainNumber": "119А",
        "direction":   0,
        "fromCode":    2004000,
        "fromTime":    "2024-10-15T00:00:00Z",
        "toCode":      2000000,
    })

    sendMessage("SearchStation", 10, "uuid-searchstation-001", map[string]interface{}{
        "query":      "ЧЕБ",
        "compactMode": true,
        "lang":       "ru",
    })

    // Задержка для отправки сообщений
    time.Sleep(5 * time.Second)
}
