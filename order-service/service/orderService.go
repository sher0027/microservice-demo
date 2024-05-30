package service

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"order-service/dto"
	"order-service/event"
	"order-service/model"
	"order-service/repository"

	"github.com/IBM/sarama"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

type OrderService struct {
	OrderRepo     *repository.OrderRepository
	Client        *resty.Client
	KafkaProducer sarama.AsyncProducer
}

func NewOrderService(orderRepo *repository.OrderRepository, client *resty.Client, kafkaProducer sarama.AsyncProducer) *OrderService {
	return &OrderService{
		OrderRepo:     orderRepo,
		Client:        client,
		KafkaProducer: kafkaProducer,
	}
}

func (service *OrderService) PlaceOrder(orderRequest *dto.OrderRequest) error {
	log.Println("Starting PlaceOrder process")
	order := model.Order{
		OrderNumber: uuid.New().String(),
	}
	log.Printf("Generated Order Number: %s", order.OrderNumber)
	log.Printf("OrderRequest: %+v", orderRequest)

	if len(orderRequest.OrderLineItems) == 0 {
		log.Println("No order line items found")
		return errors.New("no order line items found")
	}

	var orderLineItems []model.OrderLineItems
	skuCodes := make([]string, len(orderRequest.OrderLineItems))
	for i, item := range orderRequest.OrderLineItems {
		orderLineItem := model.OrderLineItems{
			SkuCode:  item.SkuCode,
			Price:    item.Price,
			Quantity: item.Quantity,
		}
		orderLineItems = append(orderLineItems, orderLineItem)
		skuCodes[i] = item.SkuCode
	}
	log.Printf("SKU Codes: %v", skuCodes)

	skuCodesParam := strings.Join(skuCodes, ",")
	log.Printf("SKU Codes Param: %s", skuCodesParam)

	var inventoryResponse []dto.InventoryResponse
	resp, err := service.Client.R().
		SetQueryParams(map[string]string{
			"skuCode": skuCodesParam,
		}).
		SetResult(&inventoryResponse).
		Get("http://localhost:8082/api/inventory")

	if err != nil {
		log.Printf("Error making inventory check request: %v", err)
		return errors.New("error checking inventory")
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("Unexpected response status code: %d", resp.StatusCode())
		return errors.New("error checking inventory")
	}

	log.Printf("Inventory Response: %+v", inventoryResponse)

	allInStock := true
	for _, inventory := range inventoryResponse {
		if !inventory.IsInStock {
			allInStock = false
			break
		}
	}
	if !allInStock {
		log.Println("Not all products are in stock")
		return errors.New("product is not in stock, please try again later")
	}

	if err := service.OrderRepo.SaveOrderWithItems(&order, orderLineItems); err != nil {
		log.Printf("Error saving order with items: %v", err)
		return err
	}

	orderPlacedEvent := event.OrderPlacedEvent{
		OrderNumber: order.OrderNumber,
	}

	eventMessage, err := json.Marshal(orderPlacedEvent)
	if err != nil {
		log.Printf("Error marshalling order placed event: %v", err)
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: "notificationTopic",
		Value: sarama.StringEncoder(eventMessage),
	}

	service.KafkaProducer.Input() <- msg

	log.Println("Order placed successfully and Kafka message sent")
	return nil
}
