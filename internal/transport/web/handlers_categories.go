package web

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/Neimess/food_tracker/internal/repository"
)

type CategoriesHandlers struct {
	tpl  *template.Template
	Repo *repository.FoodCategoriesRepo
}

func (h *CategoriesHandlers) New(w http.ResponseWriter, r *http.Request) {
	render(h.tpl, w, "categories_form.tmpl", nil)
}
func (h *CategoriesHandlers) Index(w http.ResponseWriter, r *http.Request) {
	items, err := h.Repo.List(r.Context())
	if err != nil {
		handleError(w, "Не удалось получить список категорий", err, http.StatusInternalServerError)
		return
	}
	render(h.tpl, w, "categories_index.tmpl", map[string]any{"Items": items})
}

func (h *CategoriesHandlers) Create(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	name := r.Form.Get("name")
	if name == "" {
		handleError(w, "Поле name обязательно", nil, http.StatusBadRequest)
		return
	}
	if _, err := h.Repo.Create(r.Context(), name); err != nil {
		handleError(w, "Не удалось создать категорию", err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
}

func (h *CategoriesHandlers) Edit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	item, err := h.Repo.Get(r.Context(), id)
	if err != nil {
		handleError(w, "Категория не найдена", err, http.StatusNotFound)
		return
	}
	render(h.tpl, w, "categories_form.tmpl", map[string]any{"Item": item})
}

func (h *CategoriesHandlers) Save(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	id, _ := strconv.ParseInt(r.Form.Get("id"), 10, 64)
	name := r.Form.Get("name")
	if err := h.Repo.Update(r.Context(), id, name); err != nil {
		handleError(w, "Не удалось обновить категорию", err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
}
