package models

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

// GetRecipeResponse models the response to return for getRecipe.
type GetRecipeResponse struct {
	RecipeID    int               `json:"recipeId"`
	Name        string            `json:"name"`
	Servings    int               `json:"servings"`
	Ingredients map[string]string `json:"ingredients"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
	Owner       string            `json:"owner"`
}

// PutRecipeRequest models a request to add or update a recipe.
type PutRecipeRequest struct {
	ID          *int              `json:"id"`
	Name        string            `json:"name"`
	Servings    int               `json:"servings"`
	Ingredients map[string]string `json:"ingredients"`
}

// PutRecipeResponse models a response to a PutRecipe request.
type PutRecipeResponse struct {
	RecipeID int `json:"recipeID"`
}
