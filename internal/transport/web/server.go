package web

import (
	"context"
	"embed"
	"errors"
	"html/template"
	"net/http"

	"github.com/Neimess/food_tracker/internal/config"
	"github.com/Neimess/food_tracker/internal/repository"
)

//go:embed templates/*.tmpl
var tmplFS embed.FS

func loadTemplates() *template.Template {
	return template.Must(template.ParseFS(tmplFS, "templates/*.tmpl"))
}

type Server struct {
	mux *http.ServeMux
	tpl *template.Template
	srv *http.Server
	cfg *config.HTTPServer
	Cat *CategoriesHandlers
	Ing *IngredientsHandlers
	Foo *FoodsHandlers
	Dep *DepartmentsHandlers
}

func NewServer(
	cfg *config.HTTPServer,
	catRepo *repository.FoodCategoriesRepo,
	depRepo *repository.DepartmentsRepo,
	ingRepo *repository.IngredientsRepo,
	foodsRepo *repository.FoodsRepo,
	fiRepo *repository.FoodIngredientsRepo,
) *Server {
	tpl := loadTemplates()
	s := &Server{
		mux: http.NewServeMux(),
		tpl: tpl,
		cfg: cfg,
	}

	s.Cat = &CategoriesHandlers{tpl: tpl, Repo: catRepo}
	s.Ing = &IngredientsHandlers{tpl: tpl, Repo: ingRepo, Deps: depRepo}
	s.Foo = &FoodsHandlers{tpl: tpl, Repo: foodsRepo, Cats: catRepo, FIng: fiRepo, Ings: ingRepo}
	s.Dep = &DepartmentsHandlers{tpl: tpl, Repo: depRepo}

	s.mux.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/admin/foods", http.StatusFound)
	})

	s.mux.HandleFunc("/admin/categories", s.Cat.Index)         // GET
	s.mux.HandleFunc("/admin/categories/new", s.Cat.New)       // GET
	s.mux.HandleFunc("/admin/categories/create", s.Cat.Create) // POST
	s.mux.HandleFunc("/admin/categories/edit", s.Cat.Edit)     // GET ?id=
	s.mux.HandleFunc("/admin/categories/save", s.Cat.Save)     // POST

	s.mux.HandleFunc("/admin/ingredients", s.Ing.Index)
	s.mux.HandleFunc("/admin/ingredients/new", s.Ing.New)
	s.mux.HandleFunc("/admin/ingredients/create", s.Ing.Create)
	s.mux.HandleFunc("/admin/ingredients/edit", s.Ing.Edit) // ?id=
	s.mux.HandleFunc("/admin/ingredients/save", s.Ing.Save)

	s.mux.HandleFunc("/admin/foods", s.Foo.Index)
	s.mux.HandleFunc("/admin/foods/new", s.Foo.New)
	s.mux.HandleFunc("/admin/foods/create", s.Foo.Create)
	s.mux.HandleFunc("/admin/foods/edit", s.Foo.Edit)
	s.mux.HandleFunc("/admin/foods/save", s.Foo.Save)

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
	s.srv = &http.Server{
		Addr:         s.cfg.Address,
		Handler:      s.registerMiddlewares(s.mux),
		ReadTimeout:  s.cfg.ReadTimeout,
		WriteTimeout: s.cfg.WriteTimeout,
		IdleTimeout:  s.cfg.IdleTimeout,
	}
	err := s.srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.srv == nil {
		return nil
	}
	err := s.srv.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}
