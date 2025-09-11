package repository

import (
	"context"
	"database/sql"

	"github.com/Neimess/food_tracker/internal/domain"
)

type DepartmentsRepo struct {
	db *sql.DB
}

func NewDepartmentsRepo(db *sql.DB) *DepartmentsRepo {
	return &DepartmentsRepo{db: db}
}

func (r *DepartmentsRepo) List(ctx context.Context) ([]domain.IngredientDepartment, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name FROM ingredient_departments ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	var res []domain.IngredientDepartment
	for rows.Next() {
		var d domain.IngredientDepartment
		if err := rows.Scan(&d.ID, &d.Name); err != nil {
			return nil, err
		}
		res = append(res, d)
	}
	return res, rows.Err()
}

func (r *DepartmentsRepo) Create(ctx context.Context, name string) (int64, error) {
	res, err := r.db.ExecContext(ctx, `INSERT INTO ingredient_departments(name) VALUES(?)`, name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *DepartmentsRepo) Get(ctx context.Context, id int64) (*domain.IngredientDepartment, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, name FROM ingredient_departments WHERE id=?`, id)
	var d domain.IngredientDepartment
	if err := row.Scan(&d.ID, &d.Name); err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DepartmentsRepo) Update(ctx context.Context, id int64, name string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE ingredient_departments SET name=? WHERE id=?`, name, id)
	return err
}
