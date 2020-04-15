package api

import (
	"fmt"
	"net/http"

	"github.com/justinas/alice"
	"github.com/rs/zerolog/log"

	"github.com/lynshi/cuisine-calendar-api/internal/database"
	"github.com/lynshi/cuisine-calendar-api/internal/router"
)

type appContext struct {
	router *router.Router
	db     *database.DB
	debug  bool
}

// RunCuisineCalendarAPI is the entry point for the API.
func RunCuisineCalendarAPI(debug bool) {
	log.Info().Msg("Starting Cuisine Calendar API")
	log.Info().Msg(fmt.Sprintf("Debug is %v", debug))

	app := appContext{
		router: router.NewRouter(),
		db:     database.InitializeDatabaseConnection("", "", "", "", 5432),
		debug:  debug,
	}
	defer app.db.Close()

	app.setupRouter()

	log.Fatal().Err(http.ListenAndServe(":8080", app.router)).Msg("Server exited")
}

func (app *appContext) setupRouter() {
	log.Info().Msg("Initializing router and adding handler functions")

	commonHandlers := alice.New(loggingHandler)
	app.router.Get("/recipe/:recipeId", commonHandlers.ThenFunc(app.getRecipe))
}
