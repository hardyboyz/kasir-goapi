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

	r.GET("/api/product", gin.WrapF(productHandler.HandleProducts))
	r.POST("/api/product", gin.WrapF(productHandler.HandleProducts))
	r.GET("/api/product/:id", gin.WrapF(productHandler.HandleProductByID))
	r.PUT("/api/product/:id", gin.WrapF(productHandler.HandleProductByID))
	r.DELETE("/api/product/:id", gin.WrapF(productHandler.HandleProductByID))
}
