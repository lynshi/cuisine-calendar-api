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

func getRecipe(w http.ResponseWriter, r *http.Request) {
	params := router.GetParams(r)

	recipeID, err := strconv.Atoi(router.GetParamByName(&params, "recipeId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	response := &getRecipeResponse{
		RecipeID: recipeID,
		Name:     "test",
	}

	respondWithJSON(w, http.StatusOK, response)
}
