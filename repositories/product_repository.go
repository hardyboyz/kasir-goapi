package repositories

import (
	"context"
	"database/sql"
	"errors"
	"myshop-api/models"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) GetAll() ([]models.Product, error) {
	query := "SELECT products.id, products.name, price, stock, category_id, categories.id, categories.name, description FROM products LEFT JOIN categories ON products.category_id = categories.id"
	rows, err := repo.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &p.Category.ID, &p.Category.Name, &p.Category.Description)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (repo *ProductRepository) Create(product *models.Product) error {
	query := "INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	err := repo.db.QueryRow(context.Background(), query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
	return err
}

// GetByID - ambil produk by ID
func (repo *ProductRepository) GetByID(id int) (*models.Product, error) {
	query := "SELECT A.id, A.name, A.price, A.stock, A.category_id, B.name, B.description FROM products A LEFT JOIN categories B ON A.category_id = B.id WHERE A.id = $1"

	var p models.Product
	var c models.Category
	err := repo.db.QueryRow(context.Background(), query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &c.Name, &c.Description)
	if err == sql.ErrNoRows {
		return nil, errors.New("produk tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (repo *ProductRepository) Update(product *models.Product) error {
	query := "UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5"
	result, err := repo.db.Exec(context.Background(), query, product.Name, product.Price, product.Stock, product.CategoryID, product.ID)
	if err != nil {
		return err
	}

	rows := result.RowsAffected()
	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}

	return nil
}

func (repo *ProductRepository) Delete(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := repo.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	rows := result.RowsAffected()
	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}

	return nil
}

func (repo *ProductRepository) SearchProducts(name string, categoryID int) ([]models.Product, error) {
	query := "SELECT products.id, products.name, price, stock, category_id, categories.id, categories.name, description FROM products LEFT JOIN categories ON products.category_id = categories.id WHERE 1=1"
	args := []interface{}{}
	argCount := 0

	if name != "" {
		argCount++
		query += " AND products.name ILIKE $" + strconv.Itoa(argCount)
		args = append(args, "%"+name+"%")
	}

	if categoryID > 0 {
		argCount++
		query += " AND products.category_id = $" + strconv.Itoa(argCount)
		args = append(args, categoryID)
	}

	rows, err := repo.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &p.Category.ID, &p.Category.Name, &p.Category.Description)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}
