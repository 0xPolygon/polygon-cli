package monitor

import (
	_ "embed"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/0xPolygon/polygon-cli/flag"
	"github.com/0xPolygon/polygon-cli/util"
	"github.com/spf13/cobra"
)

var (
	//go:embed usage.md
	usage string

	// flags
	rpcUrl          string
	batchSizeValue  string
	subBatchSize    int
	blockCacheLimit int
	intervalStr     string

	defaultBatchSize = 100
)

type SafeBatchSize struct {
	value int
	auto  bool // true if batchSize should be set automatically based on the UI
	mutex sync.RWMutex
}

func (s *SafeBatchSize) Set(value int, auto bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.value = value
	s.auto = auto
}

func (s *SafeBatchSize) Get() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.value
}

func (s *SafeBatchSize) Auto() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.auto
}

var MonitorCmd = &cobra.Command{
	Use:          "monitor",
	Short:        "Monitor blocks using a JSON-RPC endpoint.",
	Long:         usage,
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		rpcUrl, err = flag.GetRPCURL(cmd)
		if err != nil {
			return err
		}

		// By default, hide logs from `polycli monitor`.
		verbosityFlag := cmd.Flag("verbosity")
		if verbosityFlag != nil && !verbosityFlag.Changed {
			util.SetLogLevel(util.Silent)
		}

		prettyFlag := cmd.Flag("pretty-logs")
		if prettyFlag != nil && prettyFlag.Value.String() == "true" {
			if err = util.SetLogMode(util.Console); err != nil {
				return err
			}
		}

		return checkFlags()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return monitor(cmd.Context())
	},
}

func init() {
	f := MonitorCmd.Flags()
	f.StringVarP(&rpcUrl, flag.RPCURL, "r", flag.DefaultRPCURL, "the RPC endpoint URL")
	f.StringVarP(&batchSizeValue, "batch-size", "b", "auto", "number of requests per batch")
	f.IntVarP(&subBatchSize, "sub-batch-size", "s", 50, "number of requests per sub-batch")
	f.IntVarP(&blockCacheLimit, "cache-limit", "c", 200, "number of cached blocks for the LRU block data structure (Min 100)")
	f.StringVarP(&intervalStr, "interval", "i", "5s", "amount of time between batch block RPC calls")
}

func checkFlags() (err error) {
	interval, err = time.ParseDuration(intervalStr)
	if err != nil {
		return err
	}

	if batchSizeValue == "auto" {
		batchSize.Set(defaultBatchSize, true) // -1 value and true for auto mode
	} else {
		batchSizeInt, err := strconv.Atoi(batchSizeValue)
		if batchSizeInt == 0 || err != nil {
			return fmt.Errorf("invalid batch-size provided")
		}
		batchSize.Set(batchSizeInt, false) // specific value and false for auto mode
	}

	// Check batch-size flag.
	if blockCacheLimit < 100 {
		return fmt.Errorf("block-cache can't be less than 100")
	}

	return nil
}
