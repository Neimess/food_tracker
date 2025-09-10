package web

import (
	"html/template"
	"net/http"

	"github.com/Neimess/food_tracker/internal/repository"
)

type CategoriesHandlers struct {
	tpl  *template.Template
	Repo *repository.FoodCategoriesRepo
}

func (h *CategoriesHandlers) Index(w http.ResponseWriter, r *http.Request) {
	items, err := h.Repo.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	render(h.tpl, w, "categories_index.tmpl", map[string]any{"Items": items})
}

func (h *CategoriesHandlers) New(w http.ResponseWriter, r *http.Request) {
	render(h.tpl, w, "categories_form.tmpl", nil)
}

func (h *CategoriesHandlers) Create(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	name := r.Form.Get("name")
	if name == "" {
		http.Error(w, "name is required", 400)
		return
	}
	if _, err := h.Repo.Create(r.Context(), name); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
}

func (h *CategoriesHandlers) Edit(w http.ResponseWriter, r *http.Request) {
	// для простоты — список + форма создания; редактирование можно сделать позже
	http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
}

func (h *CategoriesHandlers) Save(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
}
