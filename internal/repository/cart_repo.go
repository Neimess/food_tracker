package repository

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/Neimess/food_tracker/internal/domain"
)

type cartKey struct {
	IngredientID int64
	Unit         string
}

type CartRepoInMemory struct {
	d        map[cartKey]domain.CartItem
	mu       sync.RWMutex
	filePath string
}

func NewCartRepo(filepath string) *CartRepoInMemory {
	return &CartRepoInMemory{
		d:        make(map[cartKey]domain.CartItem),
		filePath: filepath,
	}
}

func (r *CartRepoInMemory) AddDelta(ctx context.Context, item *domain.CartItem) error {
	if item == nil {
		return errors.New("item is nil")
	}

	key := cartKey{
		IngredientID: item.IngredientID,
		Unit:         strings.TrimSpace(item.Unit),
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	cur, ok := r.d[key]
	if !ok {
		cur = domain.CartItem{
			IngredientID: item.IngredientID,
			Name:         item.Name,
			Department:   item.Department,
			Unit:         key.Unit,
			Qty:          0,
		}
	} else {
		if cur.Name == "" && item.Name != "" {
			cur.Name = item.Name
		}
		if cur.Department == "" && item.Department != "" {
			cur.Department = item.Department
		}
	}

	cur.Qty += item.Qty
	if cur.Qty <= 0 {
		delete(r.d, key)
		return nil
	}

	r.d[key] = cur
	return nil
}

func (r *CartRepoInMemory) List(ctx context.Context) ([]domain.CartItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]domain.CartItem, 0, len(r.d))
	for _, v := range r.d {
		out = append(out, v)
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Department != out[j].Department {
			return out[i].Department < out[j].Department
		}
		if out[i].Name != out[j].Name {
			return out[i].Name < out[j].Name
		}
		if out[i].IngredientID != out[j].IngredientID {
			return out[i].IngredientID < out[j].IngredientID
		}
		return out[i].Unit < out[j].Unit
	})

	return out, nil
}

func (r *CartRepoInMemory) Clear(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.d = make(map[cartKey]domain.CartItem)
	return nil
}

func (r *CartRepoInMemory) ToggleChecked(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for k, v := range r.d {
		if v.IngredientID == int64(id) {
			v.Checked = !v.Checked
			r.d[k] = v
			return nil
		}
	}
	return fmt.Errorf("item %d not found", id)
}
