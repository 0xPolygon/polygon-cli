package sendtx

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	_ "embed"

	"github.com/0xPolygon/polygon-cli/flag"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

//go:embed usage.md
var usage string

type sendtxParams struct {
	rpcURL      string
	file        string
	batchSize   int
	concurrency int
}

var params sendtxParams

var SendtxCmd = &cobra.Command{
	Use:   "sendtx",
	Short: "Send raw transactions to a JSON-RPC endpoint in batches.",
	Long:  usage,
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, _ []string) error {
		var err error
		params.rpcURL, err = flag.GetRequiredRPCURL(cmd)
		if err != nil {
			return err
		}
		if params.concurrency <= 0 {
			params.concurrency = runtime.NumCPU()
		}
		if params.batchSize <= 0 {
			return fmt.Errorf("batch-size must be positive")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		return runSendTx(cmd.Context())
	},
}

func init() {
	f := SendtxCmd.Flags()
	f.StringVarP(&params.rpcURL, flag.RPCURL, "r", flag.DefaultRPCURL, "RPC endpoint URL")
	f.StringVarP(&params.file, "file", "f", "", "file containing raw transactions, one per line")
	f.IntVarP(&params.batchSize, "batch-size", "b", 1000, "transactions per batch request")
	f.IntVarP(&params.concurrency, "concurrency", "c", 0, "concurrent batch requests (default: number of CPUs)")
	flag.MarkFlagsRequired(SendtxCmd, "file")
}

type jsonRPCRequest struct {
	JSONRPC string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	ID      int      `json:"id"`
}

func runSendTx(ctx context.Context) error {
	start := time.Now()

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        params.concurrency,
			MaxIdleConnsPerHost: params.concurrency,
			IdleConnTimeout:     90 * time.Second,
		},
		Timeout: 120 * time.Second,
	}

	batches := make(chan []string, params.concurrency*2)

	var totalTxs atomic.Int64
	var batchesSent atomic.Int64
	var batchesFailed atomic.Int64

	// Producer: read file and send batches to channel.
	producerErr := make(chan error, 1)
	go func() {
		defer close(batches)
		err := produceBatches(ctx, params.file, params.batchSize, batches)
		if err != nil {
			producerErr <- err
		}
	}()

	// Workers: consume batches and POST them.
	var wg sync.WaitGroup
	for i := 0; i < params.concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case batch, ok := <-batches:
					if !ok {
						return
					}
					n := int64(len(batch))
					if err := sendBatch(ctx, client, params.rpcURL, batch); err != nil {
						batchesFailed.Add(1)
						log.Error().Err(err).Int64("txs", n).Msg("Batch send failed")
					} else {
						batchesSent.Add(1)
					}
					totalTxs.Add(n)
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	wg.Wait()

	// Check if producer encountered an error.
	select {
	case err := <-producerErr:
		return err
	default:
	}

	elapsed := time.Since(start)
	total := totalTxs.Load()
	sent := batchesSent.Load()
	failed := batchesFailed.Load()

	txPerSec := float64(0)
	if elapsed.Seconds() > 0 {
		txPerSec = float64(total) / elapsed.Seconds()
	}

	log.Info().
		Int64("txs", total).
		Int64("batches-sent", sent).
		Int64("batches-failed", failed).
		Str("elapsed", elapsed.Round(time.Millisecond).String()).
		Float64("tx/s", txPerSec).
		Msg("Send complete")

	return nil
}

func produceBatches(ctx context.Context, filename string, batchSize int, out chan<- []string) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	batch := make([]string, 0, batchSize)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if !strings.HasPrefix(line, "0x") {
			line = "0x" + line
		}
		batch = append(batch, line)

		if len(batch) >= batchSize {
			out <- batch
			batch = make([]string, 0, batchSize)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	// Send remaining transactions.
	if len(batch) > 0 {
		out <- batch
	}

	return nil
}

func sendBatch(ctx context.Context, client *http.Client, rpcURL string, txs []string) error {
	requests := make([]jsonRPCRequest, len(txs))
	for i, tx := range txs {
		requests[i] = jsonRPCRequest{
			JSONRPC: "2.0",
			Method:  "eth_sendRawTransaction",
			Params:  []string{tx},
			ID:      i + 1,
		}
	}

	body, err := json.Marshal(requests)
	if err != nil {
		return fmt.Errorf("marshaling batch: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rpcURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body) //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	return nil
}
