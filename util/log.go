package util

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Verbosity levels.
// https://pkg.go.dev/github.com/rs/zerolog#readme-leveled-logging
const (
	Silent = 0
	Panic  = 100
	Fatal  = 200
	Error  = 300
	Warn   = 400
	Info   = 500
	Debug  = 600
	Trace  = 700
)

// ParseVerbosity parses a verbosity string (e.g., "info", "debug") or integer
// and returns the corresponding verbosity level.
func ParseVerbosity(v string) (int, error) {
	v = strings.TrimSpace(v)

	// Try parsing as an integer first
	if level, err := strconv.Atoi(v); err == nil {
		return level, nil
	}

	// Parse as string
	switch strings.ToLower(v) {
	case "silent":
		return Silent, nil
	case "panic":
		return Panic, nil
	case "fatal":
		return Fatal, nil
	case "error":
		return Error, nil
	case "warn", "warning":
		return Warn, nil
	case "info":
		return Info, nil
	case "debug":
		return Debug, nil
	case "trace":
		return Trace, nil
	default:
		return 0, fmt.Errorf(`invalid verbosity level: %s
valid options:
  0   - silent
  100 - panic
  200 - fatal
  300 - error
  400 - warn
  500 - info
  600 - debug
  700 - trace`, v)
	}
}

// SetLogLevel sets the log level based on the flags.
// https://logging.apache.org/log4j/2.x/manual/customloglevels.html
func SetLogLevel(verbosity int) {
	switch {
	case verbosity == Silent:
		zerolog.SetGlobalLevel(zerolog.NoLevel)
	case verbosity <= Panic:
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case verbosity <= Fatal:
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case verbosity <= Error:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case verbosity <= Warn:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case verbosity <= Info:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case verbosity <= Debug:
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
