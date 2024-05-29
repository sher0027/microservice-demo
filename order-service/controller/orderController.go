package controller

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"order-service/dto"
	"order-service/service"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	OrderService *service.OrderService
}

func NewOrderController(orderService *service.OrderService) *OrderController {
	return &OrderController{
		OrderService: orderService,
	}
}

func (controller *OrderController) PlaceOrder(c *gin.Context) {
	var orderRequest dto.OrderRequest

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Received OrderRequest: %+v", orderRequest)

	if err := controller.OrderService.PlaceOrder(&orderRequest); err != nil {
		log.Printf("Error placing order: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Order placed successfully"})
}
