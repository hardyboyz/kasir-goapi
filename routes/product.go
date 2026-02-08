package routes

import (
	"myshop-api/handlers"
	"myshop-api/repositories"
	"myshop-api/services"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ProductRoutes(r *gin.Engine, db *pgxpool.Pool) {
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	r.GET("/api/product", productHandler.HandleProducts)
	r.POST("/api/product", productHandler.HandleProducts)
	r.GET("/api/product/:id", productHandler.HandleProductByID)
	r.PUT("/api/product/:id", productHandler.HandleProductByID)
	r.DELETE("/api/product/:id", productHandler.HandleProductByID)
}
