package cmdsetup

import (
	"fmt"
	
	"github.com/rs/zerolog"
)

// SetupZerolog setups up zerolog to print stack traces and with the 
// given debug output level.
func SetupZerolog(debug bool) {
	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	zerolog.ErrorStackMarshaler = func(err error) interface{} {
		return fmt.Sprintf("%+v", err)
	}
}