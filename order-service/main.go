package main

import (
	"fmt"
	"order-service/route"
	"order-service/config"
	"order-service/controller"
	"order-service/repository"
	"order-service/service"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	cfg, err := config.LoadConfig("config/config.yaml")
    if err != nil {
        log.Fatalf("Error load configuration: %v", err)
    }

    // MySQL Connection for Order Service
    db, err := gorm.Open("mysql", cfg.MySQL.URI)
    if err != nil {
        log.Fatalf("Error connect to DB: %v", err)
    }
    defer db.Close()

    // Initialize Repositories
    orderRepository := repository.NewOrderRepository(db)

    // Initialize Services
    orderService := service.NewOrderService(orderRepository, resty.New())

    // Initialize Controllers
    orderController := controller.NewOrderController(orderService)

    // Setup Router
    router := gin.Default()
    routes.SetupRouter(router, orderController)

    // Start Server
    log.Printf("Order Service running on port %d", cfg.Server.Port)
    if err := router.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
        log.Fatalf("Error start server: %v", err)
    }
}