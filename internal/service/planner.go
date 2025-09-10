package service

import (
	"context"
	"fmt"

	"github.com/Neimess/food_tracker/internal/domain"
	"github.com/Neimess/food_tracker/internal/repository"
)

type cartRepoI interface {
	AddDelta(context.Context, *domain.CartItem) error
	Clear(context.Context) error
	List(context.Context) ([]domain.CartItem, error)
}
type PlannerService struct {
	fRepo    *repository.FoodsRepo
	fiRepo   *repository.FoodIngredientsRepo
	ingRepo  *repository.IngredientsRepo
	cartRepo cartRepoI

	selectedFoods map[int64]int64
}

func NewPlannerService(fRepo *repository.FoodsRepo, fiRepo *repository.FoodIngredientsRepo,
	ingRepo *repository.IngredientsRepo, cartRepo cartRepoI) *PlannerService {
	return &PlannerService{
		fRepo:         fRepo,
		fiRepo:        fiRepo,
		ingRepo:       ingRepo,
		cartRepo:      cartRepo,
		selectedFoods: make(map[int64]int64),
	}
}

func (s *PlannerService) AddFood(ctx context.Context, foodID int64) error {
	_, err := s.fRepo.Get(ctx, foodID)
	if err != nil {
		return fmt.Errorf("food not found: %w", err)
	}
	fings, err := s.fiRepo.ListByFood(ctx, foodID)
	if err != nil {
		return fmt.Errorf("failed to list ingredients for food %d: %w", foodID, err)
	}
	for _, fing := range fings {
		ing, err := s.ingRepo.Get(ctx, fing.IngredientID)
		if err != nil {
			return fmt.Errorf("failed to load ingredient %d: %w", fing.IngredientID, err)
		}

		item := &domain.CartItem{
			IngredientID: ing.ID,
			Name:         ing.Name,
			Department:   ing.Department,
			Qty:          fing.Quantity,
			Unit:         fing.Unit,
		}

		if err := s.cartRepo.AddDelta(ctx, item); err != nil {
			return fmt.Errorf("failed to add to cart: %w", err)
		}
	}
	s.selectedFoods[foodID] = s.selectedFoods[foodID] + 1
	return nil
}

func (s *PlannerService) RemoveFood(ctx context.Context, foodID int64) error {
	count, ok := s.selectedFoods[foodID]
	if !ok || count <= 0 {
		return fmt.Errorf("food %d is not in selection", foodID)
	}

	rows, err := s.fiRepo.ListByFood(ctx, foodID)
	if err != nil {
		return fmt.Errorf("failed to list ingredients for food %d: %w", foodID, err)
	}

	for _, row := range rows {
		ing, err := s.ingRepo.Get(ctx, row.IngredientID)
		if err != nil {
			return fmt.Errorf("failed to load ingredient %d: %w", row.IngredientID, err)
		}

		item := &domain.CartItem{
			IngredientID: ing.ID,
			Name:         ing.Name,
			Department:   ing.Department,
			Qty:          -row.Quantity,
			Unit:         row.Unit,
		}

		if err := s.cartRepo.AddDelta(ctx, item); err != nil {
			return fmt.Errorf("failed to subtract ingredient %d from cart: %w", row.IngredientID, err)
		}
	}

	if count == 1 {
		delete(s.selectedFoods, foodID)
	} else {
		s.selectedFoods[foodID] = count - 1
	}
	return nil
}

func (s *PlannerService) ListFoods(ctx context.Context) ([]domain.Food, error) {
	return s.fRepo.List(ctx)
}

func (s *PlannerService) ListSelectedFoods(ctx context.Context) ([]domain.Food, error) {
	res := make([]domain.Food, 0, len(s.selectedFoods))
	for id := range s.selectedFoods {
		food, err := s.fRepo.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		res = append(res, *food)
	}
	return res, nil
}

type FoodWithCount struct {
	Food  domain.Food
	Count int64
}

func (s *PlannerService) ListSelectedFoodsWithCount(ctx context.Context) ([]FoodWithCount, error) {
	res := make([]FoodWithCount, 0, len(s.selectedFoods))
	for id, count := range s.selectedFoods {
		food, err := s.fRepo.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		res = append(res, FoodWithCount{Food: *food, Count: count})
	}
	return res, nil
}

func (s *PlannerService) SelectedFoodCounts() map[int64]int64 {
	out := make(map[int64]int64, len(s.selectedFoods))
	for k, v := range s.selectedFoods {
		out[k] = v
	}
	return out
}

func (s *PlannerService) CountForFood(foodID int64) int64 {
	return s.selectedFoods[foodID]
}

func (s *PlannerService) GetComposition(ctx context.Context, foodID int64) (*domain.Food, []domain.CompositionItem, error) {
	food, err := s.fRepo.Get(ctx, foodID)
	if err != nil {
		return nil, nil, fmt.Errorf("food not found: %w", err)
	}

	fings, err := s.fiRepo.ListByFood(ctx, foodID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list ingredients for food %d: %w", foodID, err)
	}

	res := make([]domain.CompositionItem, 0, len(fings))
	for _, fi := range fings {
		ing, err := s.ingRepo.Get(ctx, fi.IngredientID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load ingredient %d: %w", fi.IngredientID, err)
		}

		res = append(res, domain.CompositionItem{
			Ingredient: *ing,
			Quantity:   fi.Quantity,
			Unit:       fi.Unit,
		})
	}

	return food, res, nil
}

func (s *PlannerService) BuildCart(ctx context.Context) ([]domain.CartItem, error) {
	return s.cartRepo.List(ctx)
}

func (s *PlannerService) Clear(ctx context.Context) error {
	s.selectedFoods = make(map[int64]int64)
	return s.cartRepo.Clear(ctx)
}
