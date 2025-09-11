package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/Neimess/food_tracker/internal/domain"
)

type FoodCategoriesRepo struct {
	db *sql.DB
}

func NewFoodCategoriesRepo(db *sql.DB) *FoodCategoriesRepo {
	return &FoodCategoriesRepo{db: db}
}

func (r *FoodCategoriesRepo) List(ctx context.Context) ([]domain.FoodCategory, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id,name FROM food_categories ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}()

	var res []domain.FoodCategory
	for rows.Next() {
		var fc domain.FoodCategory
		if err := rows.Scan(&fc.ID, &fc.Name); err != nil {
			return nil, err
		}
		res = append(res, fc)
	}
	return res, rows.Err()
}

func (r *FoodCategoriesRepo) Create(ctx context.Context, name string) (int64, error) {
	result, err := r.db.ExecContext(ctx, "INSERT INTO food_categories (name) VALUES (?)", name)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *FoodCategoriesRepo) Get(ctx context.Context, id int64) (*domain.FoodCategory, error) {
	var fc domain.FoodCategory
	err := r.db.QueryRowContext(ctx, "SELECT id,name FROM food_categories WHERE id = ?", id).
		Scan(&fc.ID, &fc.Name)
	if err != nil {
		return nil, err
	}
	return &fc, nil
}

func (r *FoodCategoriesRepo) Update(ctx context.Context, id int64, name string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE food_categories SET name = ? WHERE id = ?", name, id)
	return err
}
