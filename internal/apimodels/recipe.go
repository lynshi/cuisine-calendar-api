package apimodels

import (
	"time"
)

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