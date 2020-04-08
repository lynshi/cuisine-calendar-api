package api

import (
	"github.com/gorilla/mux"
	"github.com/lynshi/cuisine-calendar-api/internal/handlers"
)

GET_METHOD := "GET"

// RunCuisineCalendarAPI is the entry point for the API.
func RunCuisineCalendarAPI() {
	app := App{}
	app.initialize()
}

// App holds application data.
type App struct {
	router *mux.Router
}

func (a *App) initialize() {
	a.router = mux.NewRouter()
	a.router.HandleFunc("/recipe/{recipeId:[0-9]+}", handlers.GetRecipe).Methods("GET_METHOD")
}
