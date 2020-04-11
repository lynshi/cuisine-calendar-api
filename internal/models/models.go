package models

import (
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

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

type TextInstruction struct {
	ID       int
	Text     string `gorm:"type:text"`
	RecipeID int
	Recipe   Recipe
}
