package testharness

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	harnessPort *uint16
	harnessMode *string
	activeMode  HarnessMode
)

const (
	HarnessIdentifier = "Test Harness"

	HarnessMode500 HarnessMode = "500"
	HarnessMode400 HarnessMode = "400"
)

var modeList = []HarnessMode{
	HarnessMode500,
	HarnessMode400,
}

type (
	HarnessHandler interface {
		StartListener(port uint16) error
	}
	HarnessMode string
	Handler500  struct{}
	Handler400  struct{}
)

var TestHarnessCmd = &cobra.Command{
	Use:   "testharness --mode [mode] --port [portnumber]",
	Short: "Run a simple test harness on the given port",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info().Uint16("port", *harnessPort).Str("mode", *harnessMode).Msg("Starting server")

		l, err := ListenerFactory(activeMode)
		if err != nil {
			log.Error().Err(err).Msg("Could not start listener")
			return err
		}
		return l.StartListener(*harnessPort)
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		parsedMode, err := stringToHarnessMode(*harnessMode)
		if err != nil {
			log.Error().Err(err).Msg("unable to parse mode")
		}
		activeMode = parsedMode
		return nil
	},
}

func stringToHarnessMode(inputMode string) (HarnessMode, error) {
	for _, m := range modeList {
		if inputMode == string(m) {
			return m, nil
		}
	}
	return "", fmt.Errorf("the mode %s is unrecognized", inputMode)
}

func init() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	harnessPort = TestHarnessCmd.PersistentFlags().Uint16("port", 11235, "If the mode is tcp level or higher this will set the port that the server listens on ")
	harnessMode = TestHarnessCmd.PersistentFlags().String("mode", "500", "The mode type that will be used for the harness")
}

func ListenerFactory(mode HarnessMode) (HarnessHandler, error) {
	switch mode {
	case HarnessMode500:
		return Handler500{}, nil
	case HarnessMode400:
		return Handler400{}, nil
	default:
		return nil, fmt.Errorf("the mode %s isn't supported yet", mode)
	}

}

func (m Handler500) StartListener(port uint16) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Msg("handling request")
		w.WriteHeader(500)
		w.Write([]byte(HarnessIdentifier))
	})
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
func (m Handler400) StartListener(port uint16) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Msg("handling request")
		w.WriteHeader(400)
		w.Write([]byte(HarnessIdentifier))
	})
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
