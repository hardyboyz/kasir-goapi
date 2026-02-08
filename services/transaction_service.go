package services

import (
	"myshop-api/models"
	"myshop-api/repositories"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []models.CheckoutItem, useLock bool) (*models.Transaction, error) {
	return s.repo.CreateTransaction(items)
}

func (s *TransactionService) GetTodayReport() (*models.Report, error) {
	return s.repo.GetTodayReport()
}

func (s *TransactionService) GetReportByDateRange(startDate, endDate string) (*models.Report, error) {
	return s.repo.GetReportByDateRange(startDate, endDate)
}
