package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/lynshi/cuisine-calendar-api/internal/router"
)

type getRecipeResponse struct {
	RecipeID    int         `json:"recipe_id"`
	Name        string      `json:"name"`
	Servings    int         `json:"servings"`
	Ingredients interface{} `json:"ingredients"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Owner       string      `json:"owner"`
}

func (app *appContext) getRecipe(w http.ResponseWriter, r *http.Request) {
	params := router.GetParams(r)

	recipeID, err := strconv.Atoi(router.GetParamByName(&params, "recipeId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	response := app.retrieveRecipeByID(recipeID)
	respondWithJSON(w, http.StatusOK, response)
}

func (app *appContext) retrieveRecipeByID(id int) *getRecipeResponse {
	recipe := app.db.GetRecipeByID(id)
	response := &getRecipeResponse{
		RecipeID:    recipe.ID,
		Name:        recipe.Name,
		Servings:    recipe.Servings,
		Ingredients: recipe.Ingredients,
		CreatedAt:   recipe.CreatedAt,
		UpdatedAt:   recipe.UpdatedAt,
		Owner:       recipe.Owner,
	}

	return response
}
