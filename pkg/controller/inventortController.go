package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go-microservice-demo/pkg/service"
)

type InventoryController struct {
	InventoryService *service.InventoryService
}

func NewInventoryController(inventoryService *service.InventoryService) *InventoryController {
	return &InventoryController{
		InventoryService: inventoryService,
	}
}

func (controller *InventoryController) IsInStock(c *gin.Context) {
	skuCodes := c.QueryArray("skuCode")
	if len(skuCodes) == 1 {
		skuCodes = strings.Split(skuCodes[0], ",")
	}

	responses, err := controller.InventoryService.IsInStock(skuCodes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, responses)
}