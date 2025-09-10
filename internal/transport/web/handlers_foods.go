package web

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/Neimess/food_tracker/internal/repository"
)

type FoodsHandlers struct {
	tpl  *template.Template
	Repo *repository.FoodsRepo
	Cats *repository.FoodCategoriesRepo
	FIng *repository.FoodIngredientsRepo
	Ings *repository.IngredientsRepo
}

func (h *FoodsHandlers) Index(w http.ResponseWriter, r *http.Request) {
	items, err := h.Repo.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	render(h.tpl, w, "foods_index.tmpl", map[string]any{"Items": items})
}

func (h *FoodsHandlers) New(w http.ResponseWriter, r *http.Request) {
	cats, _ := h.Cats.List(r.Context())
	ings, _ := h.Ings.List(r.Context())
	render(h.tpl, w, "foods_form.tmpl", map[string]any{
		"Cats": cats,
		"Ings": ings,
	})
}

func (h *FoodsHandlers) Create(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	name := r.Form.Get("name")
	isComplex := r.Form.Get("is_complex") == "on"
	catID, _ := strconv.ParseInt(r.Form.Get("category_id"), 10, 64)
	foodID, err := h.Repo.Create(r.Context(), name, isComplex, catID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if isComplex {
		if err := h.FIng.DeleteAllForFood(r.Context(), foodID); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if err := h.saveCompositionFromForm(w, r, foodID); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	http.Redirect(w, r, "/admin/foods", http.StatusSeeOther)
}

func (h *FoodsHandlers) Edit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	item, err := h.Repo.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	cats, _ := h.Cats.List(r.Context())
	ings, _ := h.Ings.List(r.Context())
	rows, _ := h.FIng.ListByFood(r.Context(), id)
	render(h.tpl, w, "foods_form.tmpl", map[string]any{
		"Item": item,
		"Cats": cats,
		"Ings": ings,
		"Rows": rows,
	})
}

func (h *FoodsHandlers) Save(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	foodID, _ := strconv.ParseInt(r.Form.Get("id"), 10, 64)
	name := r.Form.Get("name")
	isComplex := r.Form.Get("is_complex") == "on"
	catID, _ := strconv.ParseInt(r.Form.Get("category_id"), 10, 64)
	if err := h.Repo.Update(r.Context(), foodID, name, isComplex, catID); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if isComplex {
		if err := h.FIng.DeleteAllForFood(r.Context(), foodID); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if err := h.saveCompositionFromForm(w, r, foodID); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	} else {
		if err := h.FIng.DeleteAllForFood(r.Context(), foodID); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	http.Redirect(w, r, "/admin/foods", http.StatusSeeOther)
}

func (h *FoodsHandlers) Compose(w http.ResponseWriter, r *http.Request) {
	foodID, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	food, err := h.Repo.Get(r.Context(), foodID)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	rows, _ := h.FIng.ListByFood(r.Context(), foodID)
	ings, _ := h.Ings.List(r.Context())
	render(h.tpl, w, "foods_compose.tmpl", map[string]any{
		"Food": food, "Rows": rows, "AllIngs": ings,
	})
}

func (h *FoodsHandlers) ComposeAdd(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	foodID, _ := strconv.ParseInt(r.Form.Get("food_id"), 10, 64)
	ingID, _ := strconv.ParseInt(r.Form.Get("ingredient_id"), 10, 64)
	qty, _ := strconv.ParseFloat(r.Form.Get("quantity"), 64)
	unit := r.Form.Get("unit")
	if err := h.FIng.Upsert(r.Context(), foodID, ingID, qty, unit); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/admin/foods/compose?id="+strconv.FormatInt(foodID, 10), http.StatusSeeOther)
}

func (h *FoodsHandlers) ComposeDel(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	foodID, _ := strconv.ParseInt(r.Form.Get("food_id"), 10, 64)
	ingID, _ := strconv.ParseInt(r.Form.Get("ingredient_id"), 10, 64)
	if err := h.FIng.Delete(r.Context(), foodID, ingID); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/admin/foods/compose?id="+strconv.FormatInt(foodID, 10), http.StatusSeeOther)
}

func (h *FoodsHandlers) saveCompositionFromForm(w http.ResponseWriter, r *http.Request, foodID int64) error {
	ingIDs := r.Form["ingredient_id[]"]
	qtys := r.Form["quantity[]"]
	units := r.Form["unit[]"]

	for i := range ingIDs {
		if ingIDs[i] == "" {
			continue
		}
		ingID, _ := strconv.ParseInt(ingIDs[i], 10, 64)
		qty, _ := strconv.ParseFloat(qtys[i], 64)
		unit := units[i]

		if err := h.FIng.Upsert(r.Context(), foodID, ingID, qty, unit); err != nil {
			return err
		}
	}
	return nil
}
