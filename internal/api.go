package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// RunCuisineCalendarAPI is the entry point for the API.
func RunCuisineCalendarAPI() {
	app := App{}
	app.initialize()
}

type App struct {
	router *mux.Router
}

func (a *App) initialize() {
	a.router = mux.NewRouter()
	a.router.HandleFunc("/recipe/{recipeId:[0-9]+}", getRecipe).Methods(http.MethodGet)
}
