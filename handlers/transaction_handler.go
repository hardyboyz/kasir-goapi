package handlers

import (
	"myshop-api/models"

	"myshop-api/services"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// multiple item apa aja, quantity nya
func (h *TransactionHandler) HandleCheckout(c *gin.Context) {
	switch c.Request.Method {
	case "POST":
		h.Checkout(c)
	default:
		c.JSON(405, gin.H{"error": "Method not allowed"})
	}
}

func (h *TransactionHandler) Checkout(c *gin.Context) {
	var req models.CheckoutRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	transaction, err := h.service.Checkout(req.Items, false)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, transaction)
}

func (h *TransactionHandler) HandleTodayReport(c *gin.Context) {
	report, err := h.service.GetTodayReport()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, report)
}

func (h *TransactionHandler) HandleReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		c.JSON(400, gin.H{"error": "start_date and end_date are required"})
		return
	}

	report, err := h.service.GetReportByDateRange(startDate, endDate)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, report)
}
