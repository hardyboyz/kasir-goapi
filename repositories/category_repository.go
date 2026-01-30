package repositories

import (
	"context"
	"errors"
	"myshop-api/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (repo *CategoryRepository) GetAll() ([]models.Category, error) {
	query := "SELECT id, name, description FROM categories"
	rows, err := repo.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]models.Category, 0)
	for rows.Next() {
		var p models.Category
		err := rows.Scan(&p.ID, &p.Name, &p.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, p)
	}

	return categories, nil
}

func (repo *CategoryRepository) Create(category *models.Category) error {
	query := "INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id"
	err := repo.db.QueryRow(context.Background(), query, category.Name, category.Description).Scan(&category.ID)
	return err
}

// GetByID - ambil Category by ID
func (repo *CategoryRepository) GetByID(id int) (*models.Category, error) {
	query := "SELECT id, name, description FROM categories WHERE id = $1"

	var p models.Category
	err := repo.db.QueryRow(context.Background(), query, id).Scan(&p.ID, &p.Name, &p.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("Category tidak ditemukan")
		}
		return nil, err
	}

	return &p, nil
}

func (repo *CategoryRepository) Update(category *models.Category) error {
	query := "UPDATE categories SET name = $1, description = $2 WHERE id = $3"
	result, err := repo.db.Exec(context.Background(), query, category.Name, category.Description, category.ID)
	if err != nil {
		return err
	}

	rows := result.RowsAffected()
	if rows == 0 {
		return errors.New("Category tidak ditemukan")
	}

	return nil
}

func (repo *CategoryRepository) Delete(id int) error {
	query := "DELETE FROM categories WHERE id = $1"
	result, err := repo.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	rows := result.RowsAffected()
	if rows == 0 {
		return errors.New("Category tidak ditemukan")
	}

	return nil
}
