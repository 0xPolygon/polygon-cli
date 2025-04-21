package monitor_tview

import (
	"context"
	_ "embed"
	"fmt"
	"sync"

	"github.com/0xPolygon/polygon-cli/cmd/flag_loader"
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

var Cmd = &cobra.Command{
	Use:          "monitor-tview",
	Short:        "Monitor network using a JSON-RPC endpoint.",
	Long:         usage,
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		rpcUrlFlagValue := flag_loader.GetRpcUrlFlagValue(cmd)
		rpcUrl = *rpcUrlFlagValue
		// By default, hide logs from `polycli monitor`.
		verbosityFlag := cmd.Flag("verbosity")
		if verbosityFlag != nil && !verbosityFlag.Changed {
			util.SetLogLevel(int(util.Silent))
		}
		prettyFlag := cmd.Flag("pretty-logs")
		if prettyFlag != nil && prettyFlag.Value.String() == "true" {
			return util.SetLogMode(util.Console)
		}
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkFlags()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return monitor(cmd.Context())
	},
}

func monitor(cmd context.Context) error {
	app := NewMonitorApp(rpcUrl)
	if err := app.Run(); err != nil {
		return err
	}
	return nil
}

func init() {
	Cmd.PersistentFlags().StringVarP(&rpcUrl, "rpc-url", "r", "http://localhost:8545", "The RPC endpoint url")
	Cmd.PersistentFlags().StringVarP(&batchSizeValue, "batch-size", "b", "auto", "Number of requests per batch")
	Cmd.PersistentFlags().IntVarP(&subBatchSize, "sub-batch-size", "s", 50, "Number of requests per sub-batch")
	Cmd.PersistentFlags().IntVarP(&blockCacheLimit, "cache-limit", "c", 200, "Number of cached blocks for the LRU block data structure (Min 100)")
	Cmd.PersistentFlags().StringVarP(&intervalStr, "interval", "i", "5s", "Amount of time between batch block rpc calls")
}

func checkFlags() (err error) {
	if err = util.ValidateUrl(rpcUrl); err != nil {
		return
	}

	// interval, err = time.ParseDuration(intervalStr)
	// if err != nil {
	// 	return err
	// }

	// if batchSizeValue == "auto" {
	// 	batchSize.Set(defaultBatchSize, true) // -1 value and true for auto mode
	// } else {
	// 	batchSizeInt, err := strconv.Atoi(batchSizeValue)
	// 	if batchSizeInt == 0 || err != nil {
	// 		return fmt.Errorf("invalid batch-size provided")
	// 	}
	// 	batchSize.Set(batchSizeInt, false) // specific value and false for auto mode
	// }

	// Check batch-size flag.
	if blockCacheLimit < 100 {
		return fmt.Errorf("block-cache can't be less than 100")
	}

	return nil
}
