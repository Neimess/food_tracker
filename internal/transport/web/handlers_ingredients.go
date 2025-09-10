package web

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/Neimess/food_tracker/internal/repository"
)

type IngredientsHandlers struct {
	tpl  *template.Template
	Repo *repository.IngredientsRepo
	Deps *repository.DepartmentsRepo
}

func (h *IngredientsHandlers) Index(w http.ResponseWriter, r *http.Request) {
	items, err := h.Repo.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	render(h.tpl, w, "ingredients_index.tmpl", map[string]any{"Items": items})
}

func (h *IngredientsHandlers) New(w http.ResponseWriter, r *http.Request) {
	deps, _ := h.Deps.List(r.Context())
	render(h.tpl, w, "ingredients_form.tmpl", map[string]any{"Deps": deps})
}

func (h *IngredientsHandlers) Create(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	name := r.Form.Get("name")
	typ := r.Form.Get("type")
	depID, _ := strconv.ParseInt(r.Form.Get("department_id"), 10, 64)
	if _, err := h.Repo.Create(r.Context(), name, typ, depID); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/admin/ingredients", http.StatusSeeOther)
}

func (h *IngredientsHandlers) Edit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	it, err := h.Repo.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	deps, _ := h.Deps.List(r.Context())
	render(h.tpl, w, "ingredients_form.tmpl", map[string]any{"Item": it, "Deps": deps})
}

func (h *IngredientsHandlers) Save(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	id, _ := strconv.ParseInt(r.Form.Get("id"), 10, 64)
	name := r.Form.Get("name")
	typ := r.Form.Get("type")
	depID, _ := strconv.ParseInt(r.Form.Get("department_id"), 10, 64)
	if err := h.Repo.Update(r.Context(), id, name, typ, depID); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/admin/ingredients", http.StatusSeeOther)
}
