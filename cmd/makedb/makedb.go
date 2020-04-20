package main

import (
	"flag"

	"github.com/lynshi/cuisine-calendar-api/internal/database"
	"github.com/rs/zerolog"
)

func main() {
	debug := flag.Bool("debug", false, "sets log level to debug")

	flag.Parse()

	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	database.BuildDatabase("", "", "", "", 5432, *debug)
}
