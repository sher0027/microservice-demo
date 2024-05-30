package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"notification-service/config"

	"github.com/IBM/sarama"
)

type OrderPlacedEvent struct {
	OrderNumber string `json:"orderNumber"`
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Kafka configuration
	brokers := cfg.Kafka.Brokers
	group := "notificationGroup"
	topic := "notificationTopic"

	// Initialize Kafka Consumer
	kafkaConsumer, err := initKafkaConsumer(brokers, group)
	if err != nil {
		log.Fatalf("Error initializing Kafka consumer: %v", err)
	}
	defer kafkaConsumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			if err := kafkaConsumer.Consume(ctx, []string{topic}, &Consumer{}); err != nil {
				log.Fatalf("Error consuming messages: %v", err)
			}
		}
	}()

	// Graceful shutdown
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm
	cancel()
}

func initKafkaConsumer(brokers []string, group string) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

	consumerGroup, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		return nil, err
	}

	return consumerGroup, nil
}

type Consumer struct{}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var event OrderPlacedEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}
		log.Printf("Received order placed event: %+v", event)
		session.MarkMessage(message, "")
	}
	return nil
}
