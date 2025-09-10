package service

import (
	"context"
	"fmt"

	"github.com/Neimess/food_tracker/internal/domain"
	"github.com/Neimess/food_tracker/internal/repository"
)

type cartRepoI interface {
	Add(context.Context, *domain.CartItem) error
	Remove(context.Context, int64) error
	Clear(context.Context) error
	List(context.Context) ([]domain.CartItem, error)
}
type PlannerService struct {
	fRepo    *repository.FoodsRepo
	fiRepo   *repository.FoodIngredientsRepo
	ingRepo  *repository.IngredientsRepo
	cartRepo cartRepoI

	selectedFoods []int64
}

func NewPlannerService(fRepo *repository.FoodsRepo, fiRepo *repository.FoodIngredientsRepo,
	ingRepo *repository.IngredientsRepo, cartRepo cartRepoI) *PlannerService {
	return &PlannerService{
		fRepo:    fRepo,
		fiRepo:   fiRepo,
		ingRepo:  ingRepo,
		cartRepo: cartRepo,
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
			Name: ing.Name,
			Department: ing.Department,
			Qty: fing.Quantity,
			Unit: fing.Unit,
		}

		if err := s.cartRepo.Add(ctx, item); err != nil {
			return fmt.Errorf("failed to add to cart: %w", err)
		}
	}
	s.selectedFoods = append(s.selectedFoods, foodID)
	return nil
}

func (s *PlannerService) RemoveFood(ctx context.Context, foodID int64) error {
	rows, err := s.fiRepo.ListByFood(ctx, foodID)
	if err != nil {
		return fmt.Errorf("failed to list ingredients for food %d: %w", foodID, err)
	}
	for _, row := range rows {
		if err := s.cartRepo.Remove(ctx, row.IngredientID); err != nil {
			return fmt.Errorf("failed to remove ingredient %d from cart: %w", row.IngredientID, err)
		}
	}
	newList := make([]int64, 0, len(s.selectedFoods))
    for _, id := range s.selectedFoods {
        if id != foodID {
            newList = append(newList, id)
        }
    }
    s.selectedFoods = newList
    return nil
}

func (s *PlannerService) ListFoods(ctx context.Context) ([]domain.Food, error) {
	return s.fRepo.List(ctx)
}

func (s *PlannerService) ListSelectedFoods(ctx context.Context) ([]domain.Food, error) {
	res := make([]domain.Food, 0, len(s.selectedFoods))
    for _, id := range s.selectedFoods {
        food, err := s.fRepo.Get(ctx, id)
        if err != nil {
            return nil, err
        }
        res = append(res, *food)
    }
    return res, nil
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
	s.selectedFoods = nil
	return s.cartRepo.Clear(ctx)
}