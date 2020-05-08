package dbmodels

import (
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"

	DB "github.com/lynshi/cuisine-calendar-api/pkg/database"
)

// Recipe models the recipes table.
type Recipe struct {
	ID          int
	Name        string `gorm:"index:recipe_name"`
	Servings    int
	Ingredients postgres.Jsonb // Needs to have index created manually
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Owner       string `gorm:"index:owner_name"`
	Permissions string
}

// CreateRecipeTable creates the recipe table in the database. 
func CreateRecipeTable(db *DB.DB) {
	db.AutoMigrate(&Recipe{})
}

// AddRecipe adds an entry in the recipes table using `recipe`.
func AddRecipe(db *DB.DB, recipe *Recipe) {
	db.Create(recipe)
}

// UpdateRecipe updates an existing recipe.
func UpdateRecipe(db *DB.DB, recipe *Recipe) {
	db.Save(recipe)
}

// GetRecipeByID retrieves an entry in recipes by ID, and returns an error if no such entry exists.
func GetRecipeByID(db *DB.DB, id int) (Recipe, error) {
	var recipe Recipe
	err := db.First(&recipe, id).Error
	return recipe, err
}
