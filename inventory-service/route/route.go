package routes

import (
	"inventory-service/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, inventoryController *controller.InventoryController) {
	api := router.Group("/api")
	{
		api.GET("/inventory", inventoryController.IsInStock)
	}
}
