package repository

import (
	"context"
	"sync"

	"github.com/Neimess/food_tracker/internal/domain"
)

type CartRepoInMemory struct {
	d  map[int64]domain.CartItem
	mu sync.RWMutex
}

func NewCartRepo() *CartRepoInMemory {
	return &CartRepoInMemory{
		d:  make(map[int64]domain.CartItem),
		mu: sync.RWMutex{},
	}
}

func (r *CartRepoInMemory) Add(ctx context.Context, item *domain.CartItem) error {
	r.withLock(func(m map[int64]domain.CartItem) {
		if ex, ok := m[item.IngredientID]; ok {
			ex.Qty += item.Qty
			m[item.IngredientID] = ex
		} else {
			m[item.IngredientID] = *item
		}
	})
	return nil
}

func (r *CartRepoInMemory) Remove(ctx context.Context, ingredientID int64) error {
	r.withLock(func(m map[int64]domain.CartItem) {
		delete(m, ingredientID)
	})
	return nil
}

func (r *CartRepoInMemory) Clear(ctx context.Context) error {
	r.withLock(func(m map[int64]domain.CartItem) {
		r.d = make(map[int64]domain.CartItem)
	})
	return nil
}

func (r *CartRepoInMemory) List(ctx context.Context) ([]domain.CartItem, error) {
	var res []domain.CartItem
	r.withRLock(func(m map[int64]domain.CartItem) {
		res = make([]domain.CartItem, 0, len(m))
		for _, v := range m {
			res = append(res, v)
		}
	})
	return res, nil
}

func (r *CartRepoInMemory) withLock(fn func(m map[int64]domain.CartItem)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	fn(r.d)
}

func (r *CartRepoInMemory) withRLock(fn func(m map[int64]domain.CartItem)) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	fn(r.d)
}
