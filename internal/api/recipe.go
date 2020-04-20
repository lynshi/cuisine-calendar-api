package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"

	"github.com/lynshi/cuisine-calendar-api/internal/router"
)

type getRecipeResponse struct {
	RecipeID    int            `json:"recipe_id"`
	Name        string         `json:"name"`
	Servings    int            `json:"servings"`
	Ingredients map[string]int `json:"ingredients"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Owner       string         `json:"owner"`
}

func (app *appContext) getRecipe(w http.ResponseWriter, r *http.Request) {
	params := router.GetParams(r)

	recipeID, err := strconv.Atoi(router.GetParamByName(&params, "recipeId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	var response getRecipeResponse
	response, err = app.retrieveRecipeByID(recipeID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (app *appContext) retrieveRecipeByID(id int) (getRecipeResponse, error) {
	recipe, err := app.db.GetRecipeByID(id)
	if err != nil {
		return getRecipeResponse{}, err
	}

	var ingredients map[string]int
	ingredients, err = parseIngredientsJSONB(&recipe.Ingredients)
	if err != nil {
		return getRecipeResponse{}, err
	}

	response := getRecipeResponse{
		RecipeID:    recipe.ID,
		Name:        recipe.Name,
		Servings:    recipe.Servings,
		Ingredients: ingredients,
		CreatedAt:   recipe.CreatedAt,
		UpdatedAt:   recipe.UpdatedAt,
		Owner:       recipe.Owner,
	}

	return response, nil
}

func parseIngredientsJSONB(jsonb *postgres.Jsonb) (map[string]int, error) {
	var ingredients map[string]int
	err := json.Unmarshal(jsonb.RawMessage, &ingredients)
	return ingredients, err
}
