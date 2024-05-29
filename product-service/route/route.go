package routes

import (
	"product-service/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine,  productController *controller.ProductController) {
	api := router.Group("/api")
	{
		api.POST("/product", productController.CreateProduct)
		api.GET("/product", productController.GetAllProducts)
	}
}
