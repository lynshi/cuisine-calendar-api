package main

import (
	"flag"
	
	"github.com/rs/zerolog/log"

	"github.com/lynshi/cuisine-calendar-api/internal/logsetup"
	"github.com/lynshi/cuisine-calendar-api/pkg/api"
	"github.com/lynshi/cuisine-calendar-api/pkg/database"
	"github.com/lynshi/cuisine-calendar-api/pkg/router"
)

func main() {
	debug := flag.Bool("debug", false, "sets log level to debug")
	flag.Parse()

	logsetup.SetupZerolog(*debug)

	router := router.New()
	db, err := database.New("", "", "", "", 5432, true)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("could not create DB.DB")
	}
	defer db.Close()

	app := api.New(db, router)
	app.SetupRouter()

	err = app.Run(":8080")
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("app exited")
	}
}
