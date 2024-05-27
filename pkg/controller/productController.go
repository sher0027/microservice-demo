package controller

import (
	dto "go-microservice-demo/pkg/dto/product"
	"go-microservice-demo/pkg/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	ProductService *service.ProductService
}

func NewProductController(service *service.ProductService) *ProductController {
	return &ProductController{
		ProductService: service,
	}
}

func (ctrl *ProductController) CreateProduct(c *gin.Context) {
	var request dto.ProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.ProductService.CreateProduct(&request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (ctrl *ProductController) GetAllProducts(c *gin.Context) {
	products, err := ctrl.ProductService.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}
