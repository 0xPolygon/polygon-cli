package util

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// VerbosityLevel represents the verbosity levels.
// https://pkg.go.dev/github.com/rs/zerolog#readme-leveled-logging
type VerbosityLevel int

const (
	Silent VerbosityLevel = 0
	Panic  VerbosityLevel = 100
	Fatal  VerbosityLevel = 200
	Error  VerbosityLevel = 300
	Warn   VerbosityLevel = 400
	Info   VerbosityLevel = 500
	Debug  VerbosityLevel = 600
	Trace  VerbosityLevel = 700
)

// SetLogLevel sets the log level based on the flags.
// https://logging.apache.org/log4j/2.x/manual/customloglevels.html
func SetLogLevel(verbosity int) {
	switch {
	case verbosity == int(Silent):
		zerolog.SetGlobalLevel(zerolog.NoLevel)
	case verbosity <= int(Panic):
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case verbosity <= int(Fatal):
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case verbosity <= int(Error):
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case verbosity <= int(Warn):
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case verbosity <= int(Info):
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case verbosity <= int(Debug):
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}
}

// LogMode represents the logger mode.
type LogMode string

const (
	Console LogMode = "console"
	JSON    LogMode = "json"
)

// SetLogMode updates the log format.
func SetLogMode(mode LogMode) error {
	switch mode {
	case Console:
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Trace().Msg("Starting logger in console mode")
	case JSON:
		log.Trace().Msg("Starting logger in JSON mode")
	default:
		return fmt.Errorf("unsupported log mode: %s", mode)
	}
	return nil
}
