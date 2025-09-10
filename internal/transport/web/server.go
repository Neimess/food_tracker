package web

import (
	"html/template"
	"net/http"

	"github.com/Neimess/food_tracker/internal/repository"
)

type Server struct {
	mux *http.ServeMux
	tpl *template.Template

	Cat *CategoriesHandlers
	Ing *IngredientsHandlers
	Foo *FoodsHandlers
	Dep *DepartmentsHandlers
}

func NewServer(
	catRepo *repository.FoodCategoriesRepo,
	depRepo *repository.DepartmentsRepo,
	ingRepo *repository.IngredientsRepo,
	foodsRepo *repository.FoodsRepo,
	fiRepo *repository.FoodIngredientsRepo,
) *Server {
	tpl := template.Must(template.ParseGlob("internal/transport/web/templates/*.tmpl"))
	s := &Server{
		mux: http.NewServeMux(),
		tpl: tpl,
	}

	s.Cat = &CategoriesHandlers{tpl: tpl, Repo: catRepo}
	s.Ing = &IngredientsHandlers{tpl: tpl, Repo: ingRepo, Deps: depRepo}
	s.Foo = &FoodsHandlers{tpl: tpl, Repo: foodsRepo, Cats: catRepo, FIng: fiRepo, Ings: ingRepo}
	s.Dep = &DepartmentsHandlers{tpl: tpl, Repo: depRepo}

	// routes
	s.mux.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/admin/foods", http.StatusFound)
	})

	// categories
	s.mux.HandleFunc("/admin/categories", s.Cat.Index)         // GET
	s.mux.HandleFunc("/admin/categories/new", s.Cat.New)       // GET
	s.mux.HandleFunc("/admin/categories/create", s.Cat.Create) // POST
	s.mux.HandleFunc("/admin/categories/edit", s.Cat.Edit)     // GET ?id=
	s.mux.HandleFunc("/admin/categories/save", s.Cat.Save)     // POST

	// ingredients
	s.mux.HandleFunc("/admin/ingredients", s.Ing.Index)
	s.mux.HandleFunc("/admin/ingredients/new", s.Ing.New)
	s.mux.HandleFunc("/admin/ingredients/create", s.Ing.Create)
	s.mux.HandleFunc("/admin/ingredients/edit", s.Ing.Edit) // ?id=
	s.mux.HandleFunc("/admin/ingredients/save", s.Ing.Save)

	// foods
	s.mux.HandleFunc("/admin/foods", s.Foo.Index)
	s.mux.HandleFunc("/admin/foods/new", s.Foo.New)
	s.mux.HandleFunc("/admin/foods/create", s.Foo.Create)
	s.mux.HandleFunc("/admin/foods/edit", s.Foo.Edit)
	s.mux.HandleFunc("/admin/foods/save", s.Foo.Save)

	// food composition
	s.mux.HandleFunc("/admin/foods/compose", s.Foo.Compose)           // GET ?id=
	s.mux.HandleFunc("/admin/foods/compose/add", s.Foo.ComposeAdd)    // POST
	s.mux.HandleFunc("/admin/foods/compose/delete", s.Foo.ComposeDel) // POST

	s.mux.HandleFunc("/admin/departments", s.Dep.Index)
	s.mux.HandleFunc("/admin/departments/new", s.Dep.New)
	s.mux.HandleFunc("/admin/departments/create", s.Dep.Create)
	s.mux.HandleFunc("/admin/departments/edit", s.Dep.Edit)
	s.mux.HandleFunc("/admin/departments/save", s.Dep.Save)
	return s
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}
