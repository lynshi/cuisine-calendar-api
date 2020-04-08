package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type getRecipeResponse struct {
	id   int
	name string
}

// GetRecipe returns details for a recipe by id.
func GetRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("\"id\" could not be converted to int"))
		return
	}

	response := &getRecipeResponse{
		id:   id,
		name: "test",
	}

	json.NewEncoder(w).Encode(response)
}
