package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"myshop-api/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow(context.Background(), "SELECT name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec(context.Background(), "UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = tx.QueryRow(context.Background(), "INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	for i := range details {
		details[i].TransactionID = transactionID
		_, err = tx.Exec(context.Background(), "INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)",
			transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}

func (repo *TransactionRepository) GetTodayReport() (*models.Report, error) {
	query := `
		SELECT
			COALESCE(SUM(t.total_amount), 0) AS total_revenue,
			COUNT(t.id) AS total_transaksi,
			COALESCE(p.name, '') AS produk_nama,
			COALESCE(SUM(td.quantity), 0) AS qty_terjual
		FROM transactions t
		LEFT JOIN transaction_details td ON t.id = td.transaction_id
		LEFT JOIN products p ON td.product_id = p.id
		WHERE DATE(t.created_at) = CURRENT_DATE
		GROUP BY p.id, p.name
		ORDER BY qty_terjual DESC
		LIMIT 1
	`

	var report models.Report
	var produkNama string
	var qtyTerjual int

	err := repo.db.QueryRow(context.Background(), query).Scan(&report.TotalRevenue, &report.TotalTransaksi, &produkNama, &qtyTerjual)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		report.TotalRevenue = 0
		report.TotalTransaksi = 0
		report.ProdukTerlaris = models.ProdukTerlaris{Nama: "", QtyTerjual: 0}
	} else {
		report.ProdukTerlaris = models.ProdukTerlaris{Nama: produkNama, QtyTerjual: qtyTerjual}
	}

	return &report, nil
}

func (repo *TransactionRepository) GetReportByDateRange(startDate, endDate string) (*models.Report, error) {
	query := `
		SELECT
			COALESCE(SUM(t.total_amount), 0) AS total_revenue,
			COUNT(DISTINCT t.id) AS total_transaksi,
			COALESCE(p.name, '') AS produk_nama,
			COALESCE(SUM(td.quantity), 0) AS qty_terjual
		FROM transactions t
		LEFT JOIN transaction_details td ON t.id = td.transaction_id
		LEFT JOIN products p ON td.product_id = p.id
		WHERE DATE(t.created_at) BETWEEN $1 AND $2
		GROUP BY p.id, p.name
		ORDER BY qty_terjual DESC
		LIMIT 1
	`

	var report models.Report
	var produkNama string
	var qtyTerjual int

	err := repo.db.QueryRow(context.Background(), query, startDate, endDate).Scan(&report.TotalRevenue, &report.TotalTransaksi, &produkNama, &qtyTerjual)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		report.TotalRevenue = 0
		report.TotalTransaksi = 0
		report.ProdukTerlaris = models.ProdukTerlaris{Nama: "", QtyTerjual: 0}
	} else {
		report.ProdukTerlaris = models.ProdukTerlaris{Nama: produkNama, QtyTerjual: qtyTerjual}
	}

	return &report, nil
}
