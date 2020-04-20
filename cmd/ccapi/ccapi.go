package main

import (
	"flag"

	api "github.com/lynshi/cuisine-calendar-api/internal/api"
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

	api.RunCuisineCalendarAPI("", "", "", "", 5432, *debug)
}
