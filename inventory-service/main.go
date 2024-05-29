package main

import (
	"fmt"
	"inventory-service/route"
	"inventory-service/config"
	"inventory-service/controller"
	"inventory-service/repository"
	"inventory-service/service"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	cfg, err := config.LoadConfig("config/config.yaml")
    if err != nil {
        log.Fatalf("Error load configuration: %v", err)
    }

    // MySQL Connection for Inventory Service
    db, err := gorm.Open("mysql", cfg.MySQL.URI)
    if err != nil {
        log.Fatalf("Error connect to DB: %v", err)
    }
    defer db.Close()

    // Initialize Repositories
    inventoryRepository := repository.NewInventoryRepository(db)

    // Initialize Services
    inventoryService := service.NewInventoryService(inventoryRepository)

    // Initialize Controllers
    inventoryController := controller.NewInventoryController(inventoryService)

    // Setup Router
    router := gin.Default()
    routes.SetupRouter(router, inventoryController)

    // Start Server
    log.Printf("Inventory Service running on port %d", cfg.Server.Port)
    if err := router.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
        log.Fatalf("Error start server: %v", err)
    }
}
