package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"runtime"
	"time"

	"github.com/Chaika-Team/ChaikaRzdScraper/internal/domain"
	"github.com/Chaika-Team/ChaikaRzdScraper/internal/infrastructure/rzd/rabbitmq"
	"github.com/Chaika-Team/ChaikaRzdScraper/internal/service"
)

func StartRabbitMQHandler(svc service.Service, rabbitMQURL string) {
	consumer, err := rabbitmq.NewConsumer(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ consumer: %v", err)
	}
	defer consumer.Close()

	msgs, err := consumer.Consume("rzd_queue")
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Println("Successfully registered consumer and subscribed to queue 'rzd_queue'")

	log.Println("Waiting for messages...")

	// Количество горутин в зависимости от числа ядер процессора
	numWorkers := runtime.NumCPU() * 2

	// Канал для блокировки выполнения
	stopChan := make(chan struct{})

	// Запуск горутин
	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			for d := range msgs {
				log.Printf("Worker %d: Received a message: %s", workerID, d.Body)
				var message map[string]interface{}
				if err := json.Unmarshal(d.Body, &message); err != nil {
					log.Printf("Worker %d: Error decoding JSON: %v", workerID, err)
					continue
				}

				// Обработка сообщения в зависимости от типа запроса
				requestType := message["requestType"].(string)
				switch requestType {
				case "GetTrainRoutes":
					log.Printf("Worker %d: Handling GetTrainRoutes request", workerID)
					fromDateStr := message["payload"].(map[string]interface{})["fromDate"].(string)
					fromDate, err := time.Parse(time.RFC3339, fromDateStr)
					if err != nil {
						log.Printf("Worker %d: Error parsing fromDate: %v", workerID, err)
						continue
					}

					params := domain.GetTrainRoutesParams{
						FromCode:   int(message["payload"].(map[string]interface{})["fromCode"].(float64)),
						ToCode:     int(message["payload"].(map[string]interface{})["toCode"].(float64)),
						Direction:  domain.Direction(message["payload"].(map[string]interface{})["direction"].(float64)),
						TrainType:  domain.TrainSearchType(message["payload"].(map[string]interface{})["trainType"].(float64)),
						CheckSeats: message["payload"].(map[string]interface{})["checkSeats"].(bool),
						FromDate:   fromDate,
						WithChange: message["payload"].(map[string]interface{})["withChange"].(bool),
					}

					routes, err := svc.GetTrainRoutes(context.Background(), params)
					if err != nil {
						log.Printf("Worker %d: Failed to get train routes: %v", workerID, err)
						continue
					}

					response := map[string]interface{}{
						"requestType":   "GetTrainRoutes",
						"correlationId": message["correlationId"],
						"payload":       routes,
					}
					sendResponse(response)

				case "GetTrainCarriages":
					log.Printf("Worker %d: Handling GetTrainCarriages request", workerID)
					fromTimeStr := message["payload"].(map[string]interface{})["fromTime"].(string)
					fromTime, err := time.Parse(time.RFC3339, fromTimeStr)
					if err != nil {
						log.Printf("Worker %d: Error parsing fromTime: %v", workerID, err)
						continue
					}

					params := domain.GetTrainCarriagesParams{
						TrainNumber: message["payload"].(map[string]interface{})["trainNumber"].(string),
						Direction:   domain.Direction(message["payload"].(map[string]interface{})["direction"].(float64)),
						FromCode:    int(message["payload"].(map[string]interface{})["fromCode"].(float64)),
						FromTime:    fromTime,
						ToCode:      int(message["payload"].(map[string]interface{})["toCode"].(float64)),
					}

					carriages, err := svc.GetTrainCarriages(context.Background(), params)
					if err != nil {
						log.Printf("Worker %d: Failed to get train carriages: %v", workerID, err)
						continue
					}

					response := map[string]interface{}{
						"requestType":   "GetTrainCarriages",
						"correlationId": message["correlationId"],
						"payload":       carriages,
					}
					sendResponse(response)

				case "SearchStation":
					log.Printf("Worker %d: Handling SearchStation request", workerID)
					params := domain.SearchStationParams{
						Query:       message["payload"].(map[string]interface{})["query"].(string),
						CompactMode: message["payload"].(map[string]interface{})["compactMode"].(bool),
					}

					stations, err := svc.SearchStation(context.Background(), params)
					if err != nil {
						log.Printf("Worker %d: Failed to search stations: %v", workerID, err)
						continue
					}

					response := map[string]interface{}{
						"requestType":   "SearchStation",
						"correlationId": message["correlationId"],
						"payload":       stations,
					}
					sendResponse(response)

				default:
					log.Printf("Worker %d: Unknown request type: %s", workerID, requestType)
				}
			}
		}(i)
	}

	// Блокировка выполнения, чтобы горутины продолжали работать
	<-stopChan
}

func sendResponse(response map[string]interface{}) {
	// Здесь вы можете реализовать отправку ответа обратно в очередь или другой механизм доставки
	log.Printf("Sending response: %v", response)
}
