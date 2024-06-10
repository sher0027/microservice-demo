package main

import (
	"fmt"
	"log"
	"order-service/config"
	"order-service/controller"
	"order-service/repository"
	"order-service/route"
	"order-service/service"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Error load configuration: %v", err)
	}

	// MySQL Connection for Order Service
	db, err := gorm.Open("mysql", cfg.MySQL.URI)
	if err != nil {
		log.Fatalf("Error connect to DB: %v", err)
	}
	defer db.Close()

	// Initialize Kafka Producer
	kafkaProducer, err := initKafkaProducer(cfg.Kafka.Brokers)
	if err != nil {
		log.Fatalf("Error initializing Kafka producer: %v", err)
	}
	defer kafkaProducer.AsyncClose()

	// Initialize Repositories
	orderRepository := repository.NewOrderRepository(db)

	// Initialize Services
	orderService := service.NewOrderService(orderRepository, resty.New(), kafkaProducer)

	// Initialize Controllers
	orderController := controller.NewOrderController(orderService)

	// Setup Router
	router := gin.Default()
	route.SetupRouter(router, orderController)

	// Register the service
	err = config.RegisterService(cfg)
	if err != nil {
		log.Fatalf("Error registering service: %v", err)
	}

	// Start Server
	log.Printf("Order Service running on port %d", cfg.Service.Port)
	if err := router.Run(fmt.Sprintf(":%d", cfg.Service.Port)); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func initKafkaProducer(brokers []string) (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case success := <-producer.Successes():
				log.Printf("Message produced successfully: %v", success)
			case err := <-producer.Errors():
				log.Printf("Failed to produce message: %v", err)
			}
		}
	}()

	return producer, nil
}
