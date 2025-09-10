package repository

import (
	"context"
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

type Case struct {
	Name string
	Run  func(t *testing.T)
}

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	schema := `
	CREATE TABLE food_categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);
	`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("failed to create test schema: %v", err)
	}

	return db
}

func TestFoodCategoriesRepo(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewFoodCategoriesRepo(db)

	cases := []Case{
		{
			Name: "Create Category",
			Run: func(t *testing.T) {
				id, err := repo.Create(context.Background(), "Fruits")
				if err != nil {
					t.Fatalf("Create failed: %v", err)
				}
				if id != 1 {
					t.Errorf("want id=1, got %d", id)
				}
			},
		},
		{
			Name: "List Categories",
			Run: func(t *testing.T) {
				items, err := repo.List(context.Background())
				if err != nil {
					t.Fatalf("List failed: %v", err)
				}
				if len(items) != 1 {
					t.Fatalf("want 1 category, got %d", len(items))
				}
				if items[0].Name != "Fruits" {
					t.Errorf("want Name=Fruits, got %s", items[0].Name)
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, c.Run)
	}
}
