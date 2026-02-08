package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRoutes(r *gin.Engine, db *pgxpool.Pool) {
	CategoryRoutes(r, db)
	ProductRoutes(r, db)
	TransactionRoutes(r, db)
}
