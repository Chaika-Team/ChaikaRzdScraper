// internal/service/service.go
package service

import (
	"context"
	"log"
	"time"
	"runtime"

	"github.com/Chaika-Team/ChaikaRzdScraper/internal/domain"
	"github.com/Chaika-Team/ChaikaRzdScraper/internal/infrastructure/rzd"
	"github.com/streadway/amqp"
	"encoding/json"
	"github.com/Chaika-Team/ChaikaRzdScraper/pkg/config"
)

// mainService реализует интерфейс Service
type mainService struct {
	rzdClient *rzd.Client
}

// New возвращает новый экземпляр сервиса
func New(rzdClient *rzd.Client, cfg *config.Config) Service {
	svc := &mainService{rzdClient: rzdClient}
	go svc.startConsuming(cfg.RabbitMQ.URL) // Передаем URL из конфигурации
	return svc
}

// GetTrainRoutes получение маршрутов поездов
func (s *mainService) GetTrainRoutes(ctx context.Context, params domain.GetTrainRoutesParams) ([]domain.TrainRoute, error) {
	return s.rzdClient.GetTrainRoutes(ctx, params)
}

// GetTrainCarriages получение информации о вагонах
func (s *mainService) GetTrainCarriages(ctx context.Context, params domain.GetTrainCarriagesParams) ([]domain.Car, error) {
	return s.rzdClient.GetTrainCarriages(ctx, params)
}

// SearchStation получение кодов станций по поисковому запросу
func (s *mainService) SearchStation(ctx context.Context, params domain.SearchStationParams) ([]domain.Station, error) {
	return s.rzdClient.SearchStation(ctx, params)
}

func (s *mainService) startConsuming(rabbitMQURL string) {
	log.Println("Connecting to RabbitMQ...")

	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	log.Println("Connected to RabbitMQ")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	log.Println("Channel opened")

	msgs, err := ch.Consume(
		"rzd_queue", // имя очереди
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Println("Waiting for messages...")

	// Количество горутин в зависимости от числа ядер процессора
	numWorkers := runtime.NumCPU() * 2

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
					s.handleGetTrainRoutes(message)
				case "GetTrainCarriages":
					log.Printf("Worker %d: Handling GetTrainCarriages request", workerID)
					s.handleGetTrainCarriages(message)
				case "SearchStation":
					log.Printf("Worker %d: Handling SearchStation request", workerID)
					s.handleSearchStation(message)
				default:
					log.Printf("Worker %d: Unknown request type: %s", workerID, requestType)
				}
			}
		}(i)
	}

	// Блокировка основного потока, чтобы горутины продолжали работать
	select {}
}

func (s *mainService) handleGetTrainRoutes(message map[string]interface{}) {
	params := domain.GetTrainRoutesParams{
		FromCode:   int(message["payload"].(map[string]interface{})["fromCode"].(float64)),
		ToCode:     int(message["payload"].(map[string]interface{})["toCode"].(float64)),
		Direction:  domain.Direction(message["payload"].(map[string]interface{})["direction"].(float64)),
		TrainType:  domain.TrainSearchType(message["payload"].(map[string]interface{})["trainType"].(float64)),
		CheckSeats: message["payload"].(map[string]interface{})["checkSeats"].(bool),
		FromDate:   time.Now(), // Пример, замените на реальное значение
		WithChange: message["payload"].(map[string]interface{})["withChange"].(bool),
	}

	routes, err := s.rzdClient.GetTrainRoutes(context.Background(), params)
	if err != nil {
		log.Printf("Failed to get train routes: %v", err)
		return
	}

	response := map[string]interface{}{
		"requestType":   "GetTrainRoutes",
		"correlationId": message["correlationId"],
		"payload":       routes,
	}
	s.sendResponse(response)
}

func (s *mainService) handleGetTrainCarriages(message map[string]interface{}) {
	params := domain.GetTrainCarriagesParams{
		TrainNumber: message["payload"].(map[string]interface{})["trainNumber"].(string),
		Direction:   domain.Direction(message["payload"].(map[string]interface{})["direction"].(float64)),
		FromCode:    int(message["payload"].(map[string]interface{})["fromCode"].(float64)),
		FromTime:    time.Now(), // Пример, замените на реальное значение
		ToCode:      int(message["payload"].(map[string]interface{})["toCode"].(float64)),
	}

	carriages, err := s.rzdClient.GetTrainCarriages(context.Background(), params)
	if err != nil {
		log.Printf("Failed to get train carriages: %v", err)
		return
	}

	response := map[string]interface{}{
		"requestType":   "GetTrainCarriages",
		"correlationId": message["correlationId"],
		"payload":       carriages,
	}
	s.sendResponse(response)
}

func (s *mainService) handleSearchStation(message map[string]interface{}) {
	params := domain.SearchStationParams{
		Query:       message["payload"].(map[string]interface{})["query"].(string),
		CompactMode: message["payload"].(map[string]interface{})["compactMode"].(bool),
	}

	stations, err := s.rzdClient.SearchStation(context.Background(), params)
	if err != nil {
		log.Printf("Failed to search stations: %v", err)
		return
	}

	response := map[string]interface{}{
		"requestType":   "SearchStation",
		"correlationId": message["correlationId"],
		"payload":       stations,
	}
	s.sendResponse(response)
}

func (s *mainService) sendResponse(response map[string]interface{}) {
	// Здесь вы можете реализовать отправку ответа обратно в очередь или другой механизм доставки
	log.Printf("Sending response: %v", response)
}
