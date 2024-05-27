package main

import (
	"context"
	"fmt"
	routes "go-microservice-demo/internal/route"
	"go-microservice-demo/pkg/config"
	"go-microservice-demo/pkg/controller"
	"go-microservice-demo/pkg/repository"
	"go-microservice-demo/pkg/service"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	config.LoadConfig("configs/config.yaml")

	// MySQL Connection
	db, err := gorm.Open("mysql", config.AppConfig.MySQL.URI)
	if err != nil {
		log.Fatalf("无法连接到MySQL数据库: %v", err)
	}
	defer db.Close()

	// MongoDB Connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig.Mongo.URI))
	if err != nil {
		log.Fatalf("无法连接到MongoDB数据库: %v", err)
	}
	defer client.Disconnect(ctx)
	mongoDatabase := client.Database(config.AppConfig.Mongo.Database)

	// Initialize Repositories
	orderRepository := repository.NewOrderRepository(db)
	productRepository := repository.NewProductRepository(mongoDatabase)

	// Initialize Services
	orderService := service.NewOrderService(orderRepository, resty.New())
	productService := service.NewProductService(productRepository)

	// Initialize Controllers
	orderController := controller.NewOrderController(orderService)
	productController := controller.NewProductController(productService)

	// Setup Router
	router := gin.Default()
	routes.SetupRouter(router, orderController, productController)

	// Start Server
	log.Printf("Server running on port %d", config.AppConfig.Server.Port)
	if err := router.Run(fmt.Sprintf(":%d", config.AppConfig.Server.Port)); err != nil {
		log.Fatalf("无法启动服务器: %v", err)
	}
}
