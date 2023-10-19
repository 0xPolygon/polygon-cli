package monitor

import (
	"context"
	_ "embed"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/maticnetwork/polygon-cli/rpctypes"
	"github.com/maticnetwork/polygon-cli/util"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	//go:embed usage.md
	usage string

	// flags
	rpcUrl         string
	batchSizeValue string
	intervalStr    string
)

// MonitorCmd represents the monitor command
var MonitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor blocks using a JSON-RPC endpoint.",
	Long:  usage,
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkFlags()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return monitor(cmd.Context())
	},
}

func init() {
	MonitorCmd.PersistentFlags().StringVarP(&rpcUrl, "rpc-url", "r", "http://localhost:8545", "The RPC endpoint url")
	MonitorCmd.PersistentFlags().StringVarP(&batchSizeValue, "batch-size", "b", "auto", "Number of requests per batch")
	MonitorCmd.PersistentFlags().StringVarP(&intervalStr, "interval", "i", "5s", "Amount of time between batch block rpc calls")
}

func checkFlags() (err error) {
	// Check rpc-url flag.
	if err = util.ValidateUrl(rpcUrl); err != nil {
		return
	}

	// Check interval duration flag.
	interval, err = time.ParseDuration(intervalStr)
	if err != nil {
		return err
	}

	// Check batch-size flag.
	if batchSizeValue == "auto" {
		batchSize = -1
	} else {
		batchSize, err = strconv.Atoi(batchSizeValue)
		if batchSize == 0 {
			return fmt.Errorf("batch-size can't be equal to zero")
		}
		if err != nil {
			// Failed to convert to int, handle the error
			return fmt.Errorf("batch-size needs to be an integer")
		}
	}

	return nil
}

func monitor(ctx context.Context) error {
	rpc, err := rpc.DialContext(ctx, rpcUrl)
	if err != nil {
		log.Error().Err(err).Msg("Unable to dial rpc")
		return err
	}
	ec := ethclient.NewClient(rpc)

	ms := new(monitorStatus)
	ms.MaxBlockRetrieved = big.NewInt(0)
	ms.BlocksLock.Lock()
	ms.Blocks = make(map[string]rpctypes.PolyBlock, 0)
	ms.BlocksLock.Unlock()
	ms.ChainID = big.NewInt(0)
	ms.PendingCount = 0
	observedPendingTxs = make(historicalRange, 0)

	isUiRendered := false
	errChan := make(chan error)
	go func() {
		for {
			err = fetchBlocks(ctx, ec, ms, rpc, isUiRendered)
			if err != nil {
				continue
			}

			if !isUiRendered {
				go func() {
					errChan <- renderMonitorUI(ctx, ec, ms, rpc)
				}()
				isUiRendered = true
			}

			time.Sleep(interval)
		}
	}()

	err = <-errChan
	return err
}
