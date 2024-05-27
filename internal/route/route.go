package routes

import (
	"go-microservice-demo/pkg/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, orderController *controller.OrderController, productController *controller.ProductController) {
	api := router.Group("/api")
	{
		api.POST("/order", orderController.PlaceOrder)
		// api.GET("/order", orderController.GetAllOrders)
		api.POST("/product", productController.CreateProduct)
		api.GET("/product", productController.GetAllProducts)
	}
}
