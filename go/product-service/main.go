package main

import (
	"context"
	"fmt"
	"log"
	"product-service/config"
	"product-service/controller"
	"product-service/repository"
	"product-service/route"
	"product-service/service"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Error load configuration: %v", err)
	}

	// MongoDB Connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongo.URI))
	if err != nil {
		log.Fatalf("Error connect to DB: %v", err)
	}
	defer client.Disconnect(ctx)
	mongoDatabase := client.Database(cfg.Mongo.Database)

	// Initialize Repositories
	productRepository := repository.NewProductRepository(mongoDatabase)

	// Initialize Services
	productService := service.NewProductService(productRepository)

	// Initialize Controllers
	productController := controller.NewProductController(productService)

	// Setup Router
	router := gin.Default()
	route.SetupRouter(router, productController)

	// Register the service
	err = config.RegisterService(cfg)
	if err != nil {
		log.Fatalf("Error registering service: %v", err)
	}

	// Start Server
	log.Printf("Product Service running on port %d", cfg.Service.Port)
	if err := router.Run(fmt.Sprintf(":%d", cfg.Service.Port)); err != nil {
		log.Fatalf("Error start server: %v", err)
	}
}
