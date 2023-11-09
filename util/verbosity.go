package util

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// VerbosityLevel represents the verbosity levels.
// https://pkg.go.dev/github.com/rs/zerolog#readme-leveled-logging
type verbosityLevel int

const (
	Silent verbosityLevel = 0
	Panic  verbosityLevel = 100
	Fatal  verbosityLevel = 200
	Error  verbosityLevel = 300
	Warn   verbosityLevel = 400
	Info   verbosityLevel = 500
	Debug  verbosityLevel = 600
	Trace  verbosityLevel = 700
)

// setLogLevel sets the log level based on the flags.
// https://logging.apache.org/log4j/2.x/manual/customloglevels.html
func SetLogLevel(verbosity int, pretty bool) {
	switch {
	case verbosity == int(Silent):
		zerolog.SetGlobalLevel(zerolog.NoLevel)
	case verbosity < int(Panic):
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case verbosity < int(Fatal):
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case verbosity < int(Error):
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case verbosity < int(Warn):
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case verbosity < int(Info):
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case verbosity < int(Debug):
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	if pretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Debug().Msg("Starting logger in console mode")
	} else {
		log.Debug().Msg("Starting logger in JSON mode")
	}
}
