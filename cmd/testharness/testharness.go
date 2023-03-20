package testharness

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
)

var (
	harnessPort *uint16
	listenAddr  *string
)

const (
	packetSize        = 2<<15 - 1
	HarnessIdentifier = "Test Harness"
)

type (
	Handler500  struct{}
	Handler400  struct{}
	HandlerHuge struct{}
	HandlerJunk struct{}
)

var TestHarnessCmd = &cobra.Command{
	Use:   "testharness --mode [mode] --port [portnumber]",
	Short: "Run a simple test harness on the given port",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info().Uint16("port", *harnessPort).Str("ip", *listenAddr).Msg("Starting server")
		return startHarness()
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		parsedIp := net.ParseIP(*listenAddr)
		if parsedIp == nil {
			return fmt.Errorf("the ip %s could not be parsed", *listenAddr)
		}
		return nil
	},
}

func startHarness() error {
	http.Handle("/500", new(Handler500))
	http.Handle("/400", new(Handler400))
	http.Handle("/huge", new(HandlerHuge))
	http.Handle("/junk", new(HandlerJunk))

	return http.ListenAndServe(fmt.Sprintf("%s:%d", *listenAddr, *harnessPort), nil)
}

func init() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	harnessPort = TestHarnessCmd.PersistentFlags().Uint16("port", 11235, "If the mode is tcp level or higher this will set the port that the server listens on ")
	listenAddr = TestHarnessCmd.PersistentFlags().String("listen-ip", "127.0.0.1", "The IP that we'll use to listen")
}

func (m Handler500) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("handling request")
	w.WriteHeader(500)
	_, _ = w.Write([]byte(HarnessIdentifier))
}

func (m Handler400) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("handling request")
	w.WriteHeader(400)
	_, _ = w.Write([]byte(HarnessIdentifier))
}

func (m HandlerJunk) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idxS := r.URL.Query().Get("idx")
	idx, err := strconv.Atoi(idxS)
	if err != nil {
		idx = rand.Intn(len(junkJSONRPC))
	}

	w.Header().Set("Content-Type", junkContentTypeHeader[rand.Intn(len(junkContentTypeHeader))])

	log.Debug().Int("idx", idx).Msg("handling request")

	junkResponse := junkJSONRPC[idx%len(junkJSONRPC)]
	w.WriteHeader(200)
	_, _ = w.Write([]byte(junkResponse))
}

func (m HandlerHuge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b := r.URL.Query().Get("bytes")
	bc, err := strconv.Atoi(b)
	if err != nil {
		bc = 1024 * 1024
	}
	hugeResponse := make([]byte, bc)
	for i := 0; i+len(HarnessIdentifier) < len(hugeResponse); i = i + len(HarnessIdentifier) {
		for j := 0; j < len(HarnessIdentifier); j = j + 1 {
			hugeResponse[i+j] = HarnessIdentifier[j]
		}
	}
	hugeResponseStr := `{"jsonrpc": "2.0", "result": ` + string(hugeResponse) + `, "id": 1}`
	log.Debug().Msg("handling request")
	w.WriteHeader(200)
	_, _ = w.Write([]byte(hugeResponseStr))
}

//A server that's listening on UDP rather tha TCP
//A server that responds with valid data but with a delay (e.g. a 2-second delay)
//A server that returns data with invalid or malformed JSON syntax
//A server that returns data in a different character encoding than expected (e.g. ISO-8859-1 instead of UTF-8)
//A server that responds with a different HTTP status code than expected (e.g. 301 instead of 200)
//A server that sends back a response that exceeds the content-length specified in the response header
//A server that sends back a response with missing headers
//A server that sends back a response with extra headers
//A server that requires an authentication header, but fails if it is not provided or if it is incorrect
//A server that requires a specific content type header, and fails if it is not provided or if it is incorrect
//A server that has a firewall that blocks certain IP addresses, causing the request to fail.
