package service

import (
	"errors"
	"fmt"
	dto "go-microservice-demo/pkg/dto/order"
	"go-microservice-demo/pkg/model"
	"go-microservice-demo/pkg/repository"
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
	order := model.Order{
		OrderNumber: uuid.New().String(),
	}

	for _, item := range orderRequest.OrderLineItems {
		orderLineItem := model.OrderLineItems{
			SkuCode:  item.SkuCode,
			Price:    item.Price,
			Quantity: item.Quantity,
		}
		order.OrderLineItems = append(order.OrderLineItems, orderLineItem)
	}

	skuCodes := make([]string, len(order.OrderLineItems))
	for i, item := range order.OrderLineItems {
		skuCodes[i] = item.SkuCode
	}

	var inventoryResponse []dto.InventoryResponse
	resp, err := service.Client.R().
		SetQueryParams(map[string]string{
			"skuCode": fmt.Sprintf("%v", skuCodes),
		}).
		SetResult(&inventoryResponse).
		Get("http://inventory-service/api/inventory")

	if err != nil || resp.StatusCode() != http.StatusOK {
		return errors.New("error checking inventory")
	}

	allInStock := true
	for _, inventory := range inventoryResponse {
		if !inventory.InStock {
			allInStock = false
			break
		}
	}

	if !allInStock {
		return errors.New("product is not in stock, please try again later")
	}

	if err := service.OrderRepo.Save(&order); err != nil {
		return err
	}

	return nil
}
