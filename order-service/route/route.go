package route

import (
	"order-service/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, orderController *controller.OrderController) {
	api := router.Group("/api")
	{
		api.POST("/order", orderController.PlaceOrder)
	}
}
