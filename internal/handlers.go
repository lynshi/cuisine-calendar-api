package api

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type getRecipeResponse struct {
	id   int
	name string
}

func getRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	response := &getRecipeResponse{
		id:   id,
		name: "test",
	}

	respondWithJSON(w, http.StatusOK, response)
}
