package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm/dialects/postgres"

	"github.com/lynshi/cuisine-calendar-api/internal/models"
	"github.com/lynshi/cuisine-calendar-api/internal/router"
)

func (app *appContext) getRecipe(w http.ResponseWriter, r *http.Request) {
	params := router.GetURLParams(r)

	idString, err := router.GetURLParamByName(&params, "id")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	var recipeID int
	recipeID, err = strconv.Atoi(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	var response models.GetRecipeResponse
	response, err = app.retrieveRecipeByID(recipeID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (app *appContext) retrieveRecipeByID(id int) (models.GetRecipeResponse, error) {
	recipe, err := app.db.GetRecipeByID(id)
	if err != nil {
		return models.GetRecipeResponse{}, err
	}

	var ingredients map[string]string
	ingredients, err = parseIngredientsJSONB(&recipe.Ingredients)
	if err != nil {
		return models.GetRecipeResponse{}, err
	}

	response := models.GetRecipeResponse{
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

func (app *appContext) putRecipe(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var putRecipeRequest models.PutRecipeRequest
	err := decoder.Decode(&putRecipeRequest)

	if err != nil {
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}
	}

	var recipe models.Recipe
	if putRecipeRequest.ID == nil {
		recipe, err = app.db.GetRecipeByID(*putRecipeRequest.ID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}
	}

	recipe.Name = putRecipeRequest.Name
	recipe.Servings = putRecipeRequest.Servings
	recipe.CreatedAt = putRecipeRequest.CreatedAt
	recipe.UpdatedAt = putRecipeRequest.UpdatedAt

	var ingredients []byte
	ingredients, err = json.Marshal(putRecipeRequest.Ingredients)
	recipe.Ingredients = postgres.Jsonb{RawMessage: ingredients}

	app.db.PutRecipe(&recipe)

	response := models.PutRecipeResponse{
		RecipeID: recipe.ID,
	}

	respondWithJSON(w, http.StatusOK, response)
}

func parseIngredientsJSONB(jsonb *postgres.Jsonb) (map[string]string, error) {
	var ingredients map[string]string
	err := json.Unmarshal(jsonb.RawMessage, &ingredients)
	return ingredients, err
}
