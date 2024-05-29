package service

import (
	"log"
	"errors"
	"strings"
	"order-service/dto"
	"order-service/model"
	"order-service/repository"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

type OrderService struct {
	OrderRepo *repository.OrderRepository
	Client    *resty.Client
}

func NewOrderService(orderRepo *repository.OrderRepository, client *resty.Client) *OrderService {
	return &OrderService{
		OrderRepo: orderRepo,
		Client:    client,
	}
}

func (service *OrderService) PlaceOrder(orderRequest *dto.OrderRequest) error {
	log.Println("Starting PlaceOrder process")
	order := model.Order{
		OrderNumber: uuid.New().String(),
	}

	for _, item := range orderRequest.OrderLineItems {
		orderLineItem := model.OrderLineItems{
			SkuCode:  item.SkuCode,
			Price:    item.Price,
			Quantity: item.Quantity,
			OrderId:  order.Id,
		}
		order.OrderLineItems = append(order.OrderLineItems, orderLineItem)
	}
	log.Printf("Order Line Items: %+v", order.OrderLineItems)

	if len(order.OrderLineItems) == 0 {
		log.Println("No order line items found")
		return errors.New("no order line items found")
	}

	skuCodes := make([]string, len(order.OrderLineItems))
	for i, item := range order.OrderLineItems {
		skuCodes[i] = item.SkuCode
	}
	skuCodesParam := strings.Join(skuCodes, ",")

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

	if err := service.OrderRepo.Save(&order); err != nil {
		log.Printf("Error saving order: %v", err)
		return err
	}

	log.Println("Order placed successfully")
	return nil
}
