package web

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/Neimess/food_tracker/internal/repository"
)

type DepartmentsHandlers struct {
	tpl  *template.Template
	Repo *repository.DepartmentsRepo
}

func (h *DepartmentsHandlers) Index(w http.ResponseWriter, r *http.Request) {
	items, _ := h.Repo.List(r.Context())
	render(h.tpl, w, "departments_index.tmpl", map[string]any{"Items": items})
}

func (h *DepartmentsHandlers) New(w http.ResponseWriter, r *http.Request) {
	render(h.tpl, w, "departments_form.tmpl", nil)
}

func (h *DepartmentsHandlers) Create(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	name := r.Form.Get("name")
	if _, err := h.Repo.Create(r.Context(), name); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/admin/departments", http.StatusSeeOther)
}

func (h *DepartmentsHandlers) Edit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	item, err := h.Repo.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	render(h.tpl, w, "departments_form.tmpl", map[string]any{"Item": item})
}

func (h *DepartmentsHandlers) Save(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	id, _ := strconv.ParseInt(r.Form.Get("id"), 10, 64)
	name := r.Form.Get("name")
	if err := h.Repo.Update(r.Context(), id, name); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/admin/departments", http.StatusSeeOther)
}
