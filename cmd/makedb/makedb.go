package main

import (
	"flag"

	"github.com/rs/zerolog/log"

	"github.com/lynshi/cuisine-calendar-api/pkg/database"
	"github.com/lynshi/cuisine-calendar-api/internal/dbmodels"
	"github.com/lynshi/cuisine-calendar-api/internal/logsetup"
)

func main() {
	debug := flag.Bool("debug", false, "sets log level to debug")

	flag.Parse()

	logsetup.SetupZerolog(*debug)

	db, err := database.New("", "", "", "", 5432, true)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("could not connect to database")
	}

	dbmodels.CreateRecipeTable(db)
}
