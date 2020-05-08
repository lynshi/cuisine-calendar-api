package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"

	"github.com/lynshi/cuisine-calendar-api/internal/apimodels"
	"github.com/lynshi/cuisine-calendar-api/internal/dbmodels"
	"github.com/lynshi/cuisine-calendar-api/pkg/router"
)

func (app *App) getRecipe(w http.ResponseWriter, r *http.Request) {
	params := router.GetURLParams(r)

	idString, err := router.GetURLParamByName(&params, "id")
	if err != nil {
		err = errors.Wrap(err, "could not retrieve id from URL params")
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	var recipeID int
	recipeID, err = strconv.Atoi(idString)
	if err != nil {
		err = errors.Wrap(err, "could not convert id parameter to int")
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	var response apimodels.GetRecipeResponse
	response, err = app.retrieveRecipeByID(recipeID)
	if err != nil {
		err = errors.Wrap(err, "could not retrieve recipe by id")
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (app *App) retrieveRecipeByID(id int) (apimodels.GetRecipeResponse, error) {
	recipe, err := dbmodels.GetRecipeByID(app.db, id)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("recipe with %d could not be found", id))
		return apimodels.GetRecipeResponse{}, err
	}

	var ingredients map[string]string
	ingredients, err = parseIngredientsJSONB(&recipe.Ingredients)
	if err != nil {
		err = errors.Wrap(err, "ingredients could not be parsed")
		return apimodels.GetRecipeResponse{}, err
	}

	response := apimodels.GetRecipeResponse{
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

func (app *App) putRecipe(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var putRecipeRequest apimodels.PutRecipeRequest
	err := decoder.Decode(&putRecipeRequest)
	if err != nil {
		err = errors.Wrap(err, "could not decode put request")
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	isNewRecipe := putRecipeRequest.ID == nil

	var recipe dbmodels.Recipe
	if !isNewRecipe {
		recipe, err = dbmodels.GetRecipeByID(app.db, *putRecipeRequest.ID)
		if err != nil {
			err = errors.Wrap(err, "could not retrieve recipe")
			respondWithError(w, http.StatusInternalServerError, err)
			return
		}
	}

	recipe.Name = putRecipeRequest.Name
	recipe.Servings = putRecipeRequest.Servings

	var ingredients []byte
	ingredients, err = json.Marshal(putRecipeRequest.Ingredients)
	recipe.Ingredients = postgres.Jsonb{RawMessage: ingredients}

	if isNewRecipe {
		dbmodels.AddRecipe(app.db, &recipe)
	} else {
		dbmodels.UpdateRecipe(app.db, &recipe)
	}

	response := apimodels.PutRecipeResponse{
		RecipeID: recipe.ID,
	}

	respondWithJSON(w, http.StatusOK, response)
}

func parseIngredientsJSONB(jsonb *postgres.Jsonb) (map[string]string, error) {
	var ingredients map[string]string
	err := json.Unmarshal(jsonb.RawMessage, &ingredients)
	if err != nil {
		err = errors.Wrap(err, "error parsing ingredients from jsonb")
	}

	return ingredients, err
}
