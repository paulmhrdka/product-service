package router

import (
	_ "product-service/api/swagger"
	"product-service/internal/delivery/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func SetupRouter(e *echo.Echo, productHandler *handler.ProductHandler) {
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Swagger docs
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// API v1 routes
	v1 := e.Group("/api/v1")

	// Product routes
	products := v1.Group("/products")
	products.GET("/search", productHandler.Search)
}
