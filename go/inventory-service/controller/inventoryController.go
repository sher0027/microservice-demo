package controller

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"inventory-service/service"
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
	log.Println("IsInStock endpoint hit")
	skuCodes := c.Query("skuCode")
	if skuCodes == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "skuCode query parameter is required"})
		return
	}
	skuCodeList := strings.Split(skuCodes, ",")
	
	responses, err := controller.InventoryService.IsInStock(skuCodeList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, responses)
}