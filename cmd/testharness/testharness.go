package testharness

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"net"
	"net/http"
	"os"
)

var (
	harnessPort *uint16
	harnessMode *string
	activeMode  HarnessMode
	listenAddr  *string
)

const (
	packetSize        = 2<<15 - 1
	HarnessIdentifier = "Test Harness"

	HarnessMode500    HarnessMode = "500"
	HarnessMode400    HarnessMode = "400"
	HarnessModeClosed HarnessMode = "closed"
	HarnessModeHang   HarnessMode = "hang"
)

var modeList = []HarnessMode{
	HarnessMode500,
	HarnessMode400,
	HarnessModeClosed,
	HarnessModeHang,
}

type (
	HarnessHandler interface {
		StartListener(port uint16) error
	}
	HarnessMode   string
	Handler500    struct{}
	Handler400    struct{}
	HandlerClosed struct{}
	HandlerHang   struct{}
)

func ListenerFactory(mode HarnessMode) (HarnessHandler, error) {
	switch mode {
	case HarnessMode500:
		return Handler500{}, nil
	case HarnessMode400:
		return Handler400{}, nil
	case HarnessModeClosed:
		return HandlerClosed{}, nil
	case HarnessModeHang:
		return HandlerHang{}, nil
	default:
		return nil, fmt.Errorf("the mode %s isn't supported yet", mode)
	}

}

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
			return err
		}
		activeMode = parsedMode
		parsedIp := net.ParseIP(*listenAddr)
		if parsedIp == nil {
			return fmt.Errorf("the ip %s could not be parsed", *listenAddr)
		}
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
	listenAddr = TestHarnessCmd.PersistentFlags().String("listen-ip", "127.0.0.1", "The IP that we'll use to listen")

}

func (m Handler500) StartListener(port uint16) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Msg("handling request")
		w.WriteHeader(500)
		_, _ = w.Write([]byte(HarnessIdentifier))
	})
	return http.ListenAndServe(fmt.Sprintf("%s:%d", *listenAddr, port), nil)
}
func (m Handler400) StartListener(port uint16) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Msg("handling request")
		w.WriteHeader(400)
		_, _ = w.Write([]byte(HarnessIdentifier))
	})
	return http.ListenAndServe(fmt.Sprintf("%s:%d", *listenAddr, port), nil)
}

func (m HandlerClosed) StartListener(port uint16) error {
	addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", *listenAddr, port))
	if err != nil {
		log.Error().Err(err).Msg("unable to resolve address")
		return err
	}
	listener, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		log.Error().Err(err).Msg("unable to start listening")
		return err
	}
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Error().Err(err).Msg("error accepting")
			continue
		}
		log.Debug().
			Str("RemoteAddr", conn.RemoteAddr().String()).
			Msg("accepted connection")
		err = conn.Close()
		if err != nil {
			log.Error().Err(err).Msg("Error closing connection")
		}
		//sysConn, err := conn.SyscallConn()
		//if err != nil {
		//	log.Error().Err(err).Msg("unable to open sysconn")
		//	conn.Close()
		//	continue
		//}
		//err = sysConn.Read(func(fd uintptr) (done bool) {
		//	log.Debug().Uint("fd", uint(fd)).Msg("reading?")
		//	rawConnFile := os.NewFile(uintptr(fd), "test file")
		//
		//	data := make([]byte, packetSize)
		//	var dataAt int64 = 0
		//	n, err := rawConnFile.ReadAt(data, dataAt)
		//	if err != nil {
		//		log.Error().Err(err).Msg("unable to read!?")
		//		return false
		//	}
		//	log.Debug().Int("n", n).Msg("reading")
		//	return false
		//})
		//if err != nil {
		//	log.Error().Err(err).Msg("error reading sysconn")
		//}
		//conn.Close()
	}

}
func (m HandlerHang) StartListener(port uint16) error {
	addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", *listenAddr, port))
	if err != nil {
		log.Error().Err(err).Msg("unable to resolve address")
		return err
	}
	listener, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		log.Error().Err(err).Msg("unable to start listening")
		return err
	}
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Error().Err(err).Msg("error accepting")
			continue
		}
		log.Debug().
			Str("RemoteAddr", conn.RemoteAddr().String()).
			Msg("accepted connection")
	}
}
