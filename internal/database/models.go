package database

import (
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
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

// TextInstruction models the text_instructions table.
type TextInstruction struct {
	ID       int
	Text     string `gorm:"type:text"`
	RecipeID int
	Recipe   Recipe
}
