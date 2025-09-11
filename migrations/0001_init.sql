-- +goose Up

CREATE TABLE IF NOT EXISTS food_categories (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS ingredient_departments (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS foods (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  food_category_id INTEGER NOT NULL,
  FOREIGN KEY (food_category_id) REFERENCES food_categories(id)
    ON UPDATE RESTRICT ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS ingredients (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  department_id INTEGER NOT NULL,
  FOREIGN KEY (department_id) REFERENCES ingredient_departments(id)
    ON UPDATE RESTRICT ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS food_ingredients (
  food_id INTEGER NOT NULL,
  ingredient_id INTEGER NOT NULL,
  quantity NUMERIC NOT NULL,
  unit TEXT NOT NULL,                             -- 'g','ml','pcs'...
  PRIMARY KEY (food_id, ingredient_id),
  FOREIGN KEY (food_id) REFERENCES foods(id) ON DELETE CASCADE,
  FOREIGN KEY (ingredient_id) REFERENCES ingredients(id) ON DELETE RESTRICT
);



CREATE INDEX IF NOT EXISTS idx_foods_category       ON foods(food_category_id);
CREATE INDEX IF NOT EXISTS idx_ingr_department      ON ingredients(department_id);
CREATE INDEX IF NOT EXISTS idx_fi_food              ON food_ingredients(food_id);
CREATE INDEX IF NOT EXISTS idx_fi_ingr              ON food_ingredients(ingredient_id);

-- +goose Down
DROP TABLE IF EXISTS cart;
DROP TABLE IF EXISTS food_ingredients;
DROP TABLE IF EXISTS ingredients;
DROP TABLE IF EXISTS foods;
DROP TABLE IF EXISTS ingredient_departments;
DROP TABLE IF EXISTS food_categories;
