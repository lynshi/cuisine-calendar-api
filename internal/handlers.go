package api

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type getRecipeResponse struct {
	RecipeId int
	Name     string
}

func getRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	recipeId, err := strconv.Atoi(vars["recipeId"])
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
