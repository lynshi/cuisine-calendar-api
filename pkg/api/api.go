package api

import (
	"net/http"

	"github.com/justinas/alice"

	"github.com/lynshi/cuisine-calendar-api/pkg/database"
	"github.com/lynshi/cuisine-calendar-api/pkg/router"
)

type App struct {
	db     *database.DB
	router *router.Router
}

// New returns an app with the provided router and database.
func New(db *database.DB, router *router.Router) *App {
	return &App{
		db: db,
		router: router,
	}
}

// Runs a server with App, returning an error if it returns.
func (a *App) Run(address string) error {
	return http.ListenAndServe(address, a.router)
}

// SetupRouter adds routes to the app.
func (app *App) SetupRouter() {
	commonHandlers := alice.New(loggingHandler)
	app.router.Get("/getRecipe", commonHandlers.ThenFunc(app.getRecipe))
	app.router.Post("/putRecipe", commonHandlers.ThenFunc(app.putRecipe))
}
