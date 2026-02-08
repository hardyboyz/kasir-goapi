package routes

import (
	"myshop-api/handlers"
	"myshop-api/repositories"
	"myshop-api/services"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TransactionRoutes(r *gin.Engine, db *pgxpool.Pool) {
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	r.POST("/api/checkout", transactionHandler.HandleCheckout)
	r.GET("/api/report/hari-ini", transactionHandler.HandleTodayReport)
	r.GET("/api/report", transactionHandler.HandleReport)
}
