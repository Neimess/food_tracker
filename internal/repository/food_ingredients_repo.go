package repository

import (
	"context"
	"database/sql"

	"github.com/Neimess/food_tracker/internal/domain"
)

type FoodIngredientsRepo struct{ DB *sql.DB }

func NewFoodIngredientsRepo(db *sql.DB) *FoodIngredientsRepo {
	return &FoodIngredientsRepo{DB: db}
}

func (r *FoodIngredientsRepo) ListByFood(ctx context.Context, foodID int64) ([]domain.FoodIngredient, error) {
	const q = `
	SELECT fi.food_id, fi.ingredient_id, i.name, fi.quantity, fi.unit
	FROM food_ingredients fi
	JOIN ingredients i ON fi.ingredient_id = i.id
	WHERE fi.food_id = ?
	ORDER BY i.name`
	rows, err := r.DB.QueryContext(ctx, q, foodID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.FoodIngredient
	for rows.Next() {
		var fi domain.FoodIngredient
		if err := rows.Scan(&fi.FoodID, &fi.IngredientID, &fi.IngredientName, &fi.Quantity, &fi.Unit); err != nil {
			return nil, err
		}
		res = append(res, fi)
	}
	return res, rows.Err()
}

func (r *FoodIngredientsRepo) Upsert(ctx context.Context, foodID, ingID int64, qty float64, unit string) error {
	_, err := r.DB.ExecContext(ctx, `
	INSERT INTO food_ingredients(food_id, ingredient_id, quantity, unit)
	VALUES (?,?,?,?)
	ON CONFLICT(food_id, ingredient_id) DO UPDATE SET quantity=excluded.quantity, unit=excluded.unit
	`, foodID, ingID, qty, unit)
	return err
}

func (r *FoodIngredientsRepo) Delete(ctx context.Context, foodID, ingID int64) error {
	_, err := r.DB.ExecContext(ctx,
		`DELETE FROM food_ingredients WHERE food_id=? AND ingredient_id=?`,
		foodID, ingID)
	return err
}

func (r *FoodIngredientsRepo) DeleteAllForFood(ctx context.Context, foodID int64) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM food_ingredients WHERE food_id = ?`, foodID)
	return err
}
