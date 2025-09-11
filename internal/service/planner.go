package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Neimess/food_tracker/internal/domain"
	"github.com/Neimess/food_tracker/internal/repository"
)

type cartRepoI interface {
	AddDelta(context.Context, *domain.CartItem) error
	Clear(context.Context) error
	List(context.Context) ([]domain.CartItem, error)
	ToggleChecked(context.Context, int64) error
}

type PlannerService struct {
	fRepo    *repository.FoodsRepo
	fiRepo   *repository.FoodIngredientsRepo
	ingRepo  *repository.IngredientsRepo
	cartRepo cartRepoI

	muSelected sync.RWMutex
	selected   map[int64]int64

	statePath   string
	preloadOnce sync.Once
	preloadErr  error

	cacheDur time.Duration

	InvalidateChan chan struct{}
}

func NewPlannerService(
	statePath string,
	fRepo *repository.FoodsRepo,
	fiRepo *repository.FoodIngredientsRepo,
	ingRepo *repository.IngredientsRepo,
	cartRepo cartRepoI,
) *PlannerService {
	return &PlannerService{
		statePath:      statePath,
		fRepo:          fRepo,
		fiRepo:         fiRepo,
		ingRepo:        ingRepo,
		cartRepo:       cartRepo,
		selected:       make(map[int64]int64),
		InvalidateChan: make(chan struct{}, 1),
	}
}

func (s *PlannerService) Preload(ctx context.Context) error {
	s.preloadOnce.Do(func() {
		if s.statePath == "" {
			return
		}
		st, err := loadStateJSON(s.statePath)
		if err != nil {
			s.preloadErr = err
			return
		}
		if st == nil {
			return
		}

		if st.SelectedFoods != nil {
			s.withSelectedW(func(m map[int64]int64) {
				for k, v := range st.SelectedFoods {
					m[k] = v
				}
			})
		}

		for i := range st.Cart {
			_ = s.cartRepo.AddDelta(ctx, &st.Cart[i])
		}
	})
	go s.invalidateCart(ctx)
	return s.preloadErr
}

func (s *PlannerService) Close(ctx context.Context) error {
	if s.statePath == "" {
		return nil
	}

	cart, _ := s.cartRepo.List(ctx)

	st := &persistedState{
		Version:       1,
		SelectedFoods: s.cloneSelected(),
		Cart:          cart,
	}
	return saveStateJSONAtomic(s.statePath, st, 0o600)
}

func (s *PlannerService) AddFood(ctx context.Context, foodID int64) error {
	return s.applyFoodDelta(ctx, foodID, +1)
}

func (s *PlannerService) RemoveFood(ctx context.Context, foodID int64) error {
	if s.CountForFood(foodID) <= 0 {
		return fmt.Errorf("food %d is not in selection", foodID)
	}
	return s.applyFoodDelta(ctx, foodID, -1)
}

func (s *PlannerService) ListFoods(ctx context.Context) ([]domain.Food, error) {
	return s.fRepo.List(ctx)
}

func (s *PlannerService) ListSelectedFoods(ctx context.Context) ([]domain.Food, error) {
	snap := s.cloneSelected()
	res := make([]domain.Food, 0, len(snap))
	for id := range snap {
		food, err := s.fRepo.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		res = append(res, *food)
	}
	return res, nil
}

func (s *PlannerService) ListSelectedFoodsWithCount(ctx context.Context) ([]FoodWithCount, error) {
	snap := s.cloneSelected()
	res := make([]FoodWithCount, 0, len(snap))
	for id, count := range snap {
		food, err := s.fRepo.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		res = append(res, FoodWithCount{Food: *food, Count: count})
	}
	return res, nil
}

func (s *PlannerService) SelectedFoodCounts() map[int64]int64 {
	return s.cloneSelected()
}

func (s *PlannerService) CountForFood(foodID int64) int64 {
	var v int64
	s.withSelectedR(func(m map[int64]int64) { v = m[foodID] })
	return v
}

func (s *PlannerService) ToggleChecked(ctx context.Context, id int64) error {
	return s.cartRepo.ToggleChecked(ctx, id)
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

	items, err := s.compositionFromFings(ctx, fings)
	if err != nil {
		return nil, nil, err
	}
	return food, items, nil
}

func (s *PlannerService) BuildCart(ctx context.Context) ([]domain.CartItem, error) {
	return s.cartRepo.List(ctx)
}

func (s *PlannerService) Clear(ctx context.Context) error {
	s.withSelectedW(func(m map[int64]int64) {
		for k := range m {
			delete(m, k)
		}
	})
	return s.cartRepo.Clear(ctx)
}

func (s *PlannerService) applyFoodDelta(ctx context.Context, foodID int64, mult int64) error {
	if mult == 0 {
		return errors.New("applyFoodDelta: mult must be non-zero")
	}
	if _, err := s.fRepo.Get(ctx, foodID); err != nil {
		return fmt.Errorf("food not found: %w", err)
	}

	fings, err := s.fiRepo.ListByFood(ctx, foodID)
	if err != nil {
		return fmt.Errorf("failed to list ingredients for food %d: %w", foodID, err)
	}

	items, err := s.cartItemsFromFings(ctx, fings, mult)
	if err != nil {
		return err
	}
	for i := range items {
		if err := s.cartRepo.AddDelta(ctx, &items[i]); err != nil {
			return fmt.Errorf("failed to update cart: %w", err)
		}
	}

	if mult > 0 {
		s.withSelectedW(func(m map[int64]int64) { m[foodID] = m[foodID] + 1 })
	} else {
		s.withSelectedW(func(m map[int64]int64) {
			if cur := m[foodID]; cur <= 1 {
				delete(m, foodID)
			} else {
				m[foodID] = cur - 1
			}
		})
	}

	return nil
}

func (s *PlannerService) cartItemsFromFings(ctx context.Context, fings []domain.FoodIngredient, mult int64) ([]domain.CartItem, error) {
	res := make([]domain.CartItem, 0, len(fings))
	for _, fi := range fings {
		ing, err := s.ingRepo.Get(ctx, fi.IngredientID)
		if err != nil {
			return nil, fmt.Errorf("failed to load ingredient %d: %w", fi.IngredientID, err)
		}
		res = append(res, domain.CartItem{
			IngredientID: ing.ID,
			Name:         ing.Name,
			Department:   ing.Department,
			Qty:          fi.Quantity * mult,
			Unit:         fi.Unit,
		})
	}
	return res, nil
}

func (s *PlannerService) compositionFromFings(ctx context.Context, fings []domain.FoodIngredient) ([]domain.CompositionItem, error) {
	res := make([]domain.CompositionItem, 0, len(fings))
	for _, fi := range fings {
		ing, err := s.ingRepo.Get(ctx, fi.IngredientID)
		if err != nil {
			return nil, fmt.Errorf("failed to load ingredient %d: %w", fi.IngredientID, err)
		}
		res = append(res, domain.CompositionItem{
			Ingredient: *ing,
			Quantity:   fi.Quantity,
			Unit:       fi.Unit,
		})
	}
	return res, nil
}

func (s *PlannerService) withSelectedR(fn func(map[int64]int64)) {
	s.muSelected.RLock()
	fn(s.selected)
	s.muSelected.RUnlock()
}

func (s *PlannerService) withSelectedW(fn func(map[int64]int64)) {
	s.muSelected.Lock()
	fn(s.selected)
	s.muSelected.Unlock()
}

func (s *PlannerService) cloneSelected() map[int64]int64 {
	out := make(map[int64]int64)
	s.withSelectedR(func(m map[int64]int64) {
		for k, v := range m {
			out[k] = v
		}
	})
	return out
}

func (s *PlannerService) invalidateCart(ctx context.Context) {
	now := time.Now()
	daysUntilSaturday := (time.Saturday - now.Weekday() + 7) % 7
	if daysUntilSaturday == 0 {
		daysUntilSaturday = 7
	}

	nextSaturday := time.Date(
		now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location(),
	).AddDate(0, 0, int(daysUntilSaturday))

	d := time.Until(nextSaturday)
	timer := time.NewTimer(d)
	for {
		select {
		case <-ctx.Done():
		case <-timer.C:
			s.muSelected.Lock()
			s.selected = make(map[int64]int64)
			s.muSelected.Unlock()

			select {
			case s.InvalidateChan <- struct{}{}:
				log.Println("Cache invalidated")
			default:

			}
		}
	}
}
