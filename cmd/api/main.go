package main

import (
	"log"
	"product-service/config"
	"product-service/internal/delivery/handler"
	"product-service/internal/delivery/router"
	elasticrepository "product-service/internal/repository/elasticsearch"
	"product-service/internal/usecase"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/labstack/echo/v4"
)

// @title Product Service API
// @version 1.0
// @description Product Service API with Elasticsearch.

// @host localhost:8080
func main() {
	// Load config
	cfg := config.Load()

	// Initialize Elasticsearch client
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: cfg.Elasticsearch.URLs,
	})
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	// Initialize repositories
	productRepo := elasticrepository.NewProductRepository(es)

	// Initialize use cases
	productUseCase := usecase.NewProductUseCase(productRepo)

	// Initialize handlers
	productHandler := handler.NewProductHandler(productUseCase)

	// Setup Echo
	e := echo.New()
	router.SetupRouter(e, productHandler)

	// Start server
	e.Logger.Fatal(e.Start(":" + cfg.Server.Port))
}
