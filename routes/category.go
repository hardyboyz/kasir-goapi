package routes

import (
	"myshop-api/handlers"
	"myshop-api/repositories"
	"myshop-api/services"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CategoryRoutes(r *gin.Engine, db *pgxpool.Pool) {
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	r.GET("/api/category", gin.WrapF(categoryHandler.HandleCategories))
	r.POST("/api/category", gin.WrapF(categoryHandler.HandleCategories))
	r.GET("/api/category/:id", gin.WrapF(categoryHandler.HandleCategoryByID))
	r.PUT("/api/category/:id", gin.WrapF(categoryHandler.HandleCategoryByID))
	r.DELETE("/api/category/:id", gin.WrapF(categoryHandler.HandleCategoryByID))
}
