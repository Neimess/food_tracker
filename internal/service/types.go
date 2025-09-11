package service

import "github.com/Neimess/food_tracker/internal/domain"

type FoodWithCount struct {
	Food  domain.Food
	Count int64
}

type persistedState struct {
	Version       int64             `json:"version"`
	SelectedFoods map[int64]int64   `json:"selected_foods"`
	Cart          []domain.CartItem `json:"cart"`
}
