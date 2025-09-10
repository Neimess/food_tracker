package domain

type FoodCategory struct {
	ID   int64
	Name string
}

type IngredientDepartment struct {
	ID   int64
	Name string
}

type Food struct {
	ID           int64
	Name         string
	IsComplex    bool
	CategoryID   int64
	CategoryName string
}

type Ingredient struct {
	ID           int64
	Name         string
	DepartmentID int64
	Department   string
}

type FoodIngredient struct {
	FoodID         int64
	IngredientID   int64
	IngredientName string
	Quantity       float64
	Unit           string
}

type CartItem struct {
	IngredientID   int64
	Name 		   string
	Unit           string
	Qty            float64
	Department     string
	Checked        bool
}

type CompositionItem struct {
    Ingredient Ingredient
    Quantity   float64
    Unit       string
}