package controller

import (
	dto "go-microservice-demo/pkg/dto/order"
	"go-microservice-demo/pkg/service"
	"net/http"

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

func (ctrl *OrderController) PlaceOrder(c *gin.Context) {
    var orderRequest dto.OrderRequest
    if err := c.ShouldBindJSON(&orderRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := ctrl.OrderService.PlaceOrder(&orderRequest); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusCreated)
}
