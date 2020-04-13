package api

import (
	"net/http"
	"strconv"

	"github.com/lynshi/cuisine-calendar-api/internal/router"
)

type getRecipeResponse struct {
	RecipeId     int
	Name         string
	Instructions string
	Ingredients  string
	Quantity     float64
	QuantityUnit string
}

func getRecipe(w http.ResponseWriter, r *http.Request) {
	params := router.GetParams(r)

	recipeId, err := strconv.Atoi(params.ByName("recipeId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	response := &getRecipeResponse{
		RecipeId: recipeId,
		Name:     "test",
	}

	respondWithJSON(w, http.StatusOK, response)
}
