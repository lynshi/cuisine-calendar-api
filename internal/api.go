package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// RunCuisineCalendarAPI is the entry point for the API.
func RunCuisineCalendarAPI(debug bool) {
	log.Info().Msg("Starting Cuisine Calendar API")
	log.Info().Msg(fmt.Sprintf("Debug is %v", debug))

	app := App{
		debug: debug,
	}
	app.initialize()
	defer app.db.Close()
}

type App struct {
	router *mux.Router
	db     *sql.DB
	debug  bool
}

func (a *App) initialize() {
	a.initializeRouter()
	a.initializeDatabaseConnection("", "", "", "", 5432)
}

func (a *App) initializeRouter() {
	log.Info().Msg("Initializing router and adding handler functions")

	a.router = mux.NewRouter()

	log.Debug().Msg("Adding recipe GET handler")
	a.router.HandleFunc("/recipe/{recipeId:[0-9]+}", getRecipe).Methods(http.MethodGet)
}
