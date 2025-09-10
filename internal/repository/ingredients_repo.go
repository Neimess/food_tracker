package repository

import (
	"context"
	"database/sql"

	"github.com/Neimess/food_tracker/internal/domain"
)

type IngredientsRepo struct{ DB *sql.DB }

func NewIngredientsRepo(db *sql.DB) *IngredientsRepo {
	return &IngredientsRepo{DB: db}
}

func (r *IngredientsRepo) List(ctx context.Context) ([]domain.Ingredient, error) {
	const q = `
	SELECT i.id, i.name, i.department_id, d.name
	FROM ingredients i
	JOIN ingredient_departments d ON i.department_id = d.id
	ORDER BY i.name`
	rows, err := r.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.Ingredient
	for rows.Next() {
		var ing domain.Ingredient
		if err := rows.Scan(&ing.ID, &ing.Name, &ing.DepartmentID, &ing.Department); err != nil {
			return nil, err
		}
		res = append(res, ing)
	}
	return res, rows.Err()
}

func (r *IngredientsRepo) Get(ctx context.Context, id int64) (*domain.Ingredient, error) {
	const q = `
	SELECT i.id, i.name, i.department_id, d.name
	FROM ingredients i
	JOIN ingredient_departments d ON i.department_id = d.id
	WHERE i.id = ?`
	var ing domain.Ingredient
	err := r.DB.QueryRowContext(ctx, q, id).
		Scan(&ing.ID, &ing.Name, &ing.DepartmentID, &ing.Department)
	if err != nil {
		return nil, err
	}
	return &ing, nil
}

func (r *IngredientsRepo) Create(ctx context.Context, name, t string, deptID int64) (int64, error) {
	res, err := r.DB.ExecContext(ctx,
		`INSERT INTO ingredients(name, department_id) VALUES(?,?)`,
		name, deptID)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *IngredientsRepo) Update(ctx context.Context, id int64, name, t string, deptID int64) error {
	_, err := r.DB.ExecContext(ctx,
		`UPDATE ingredients SET name=?, department_id=? WHERE id=?`,
		name, deptID, id)
	return err
}
