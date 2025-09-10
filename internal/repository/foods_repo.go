package repository

import (
	"context"
	"database/sql"

	"github.com/Neimess/food_tracker/internal/domain"
)

type FoodsRepo struct{ DB *sql.DB }

func NewFoodsRepo(db *sql.DB) *FoodsRepo {
	return &FoodsRepo{DB: db}
}

func (r *FoodsRepo) List(ctx context.Context) ([]domain.Food, error) {
	const q = `
	SELECT f.id, f.name, f.is_complex, f.food_category_id, c.name
	FROM foods f
	JOIN food_categories c ON f.food_category_id = c.id
	ORDER BY f.name`
	rows, err := r.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.Food
	for rows.Next() {
		var f domain.Food
		if err := rows.Scan(&f.ID, &f.Name, &f.IsComplex, &f.CategoryID, &f.CategoryName); err != nil {
			return nil, err
		}
		res = append(res, f)
	}
	return res, rows.Err()
}

func (r *FoodsRepo) Get(ctx context.Context, id int64) (*domain.Food, error) {
	const q = `
	SELECT f.id, f.name, f.is_complex, f.food_category_id, c.name
	FROM foods f
	JOIN food_categories c ON f.food_category_id = c.id
	WHERE f.id = ?`
	var f domain.Food
	err := r.DB.QueryRowContext(ctx, q, id).Scan(&f.ID, &f.Name, &f.IsComplex, &f.CategoryID, &f.CategoryName)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *FoodsRepo) Create(ctx context.Context, name string, isComplex bool, categoryID int64) (int64, error) {
	res, err := r.DB.ExecContext(ctx,
		`INSERT INTO foods(name, is_complex, food_category_id) VALUES(?,?,?)`,
		name, isComplex, categoryID)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *FoodsRepo) Update(ctx context.Context, id int64, name string, isComplex bool, categoryID int64) error {
	_, err := r.DB.ExecContext(ctx,
		`UPDATE foods SET name=?, is_complex=?, food_category_id=? WHERE id=?`,
		name, isComplex, categoryID, id)
	return err
}

func (r *FoodsRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM foods WHERE id=?`, id)
	return err
}
