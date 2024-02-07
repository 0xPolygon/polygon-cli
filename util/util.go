package util

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/ethereum/go-ethereum/consensus/clique"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"

	"github.com/ethereum/go-ethereum/core/types"
)

type (
	simpleRPCTransaction struct {
		Hash string `json:"hash"`
	}
	simpleRPCBlock struct {
		Number       string                 `json:"number"`
		Transactions []simpleRPCTransaction `json:"transactions"`
	}
	txpoolStatus struct {
		Pending any `json:"pending"`
		Queued  any `json:"queued"`
	}
)

func Ecrecover(block *types.Block) ([]byte, error) {
	header := block.Header()
	sigStart := len(header.Extra) - ethcrypto.SignatureLength
	if sigStart < 0 || sigStart > len(header.Extra) {
		return nil, fmt.Errorf("unable to recover signature")
	}
	signature := header.Extra[sigStart:]
	pubkey, err := ethcrypto.Ecrecover(clique.SealHash(header).Bytes(), signature)
	if err != nil {
		return nil, err
	}
	signer := ethcrypto.Keccak256(pubkey[1:])[12:]

	return signer, nil
}

func GetBlockRange(ctx context.Context, from, to uint64, c *ethrpc.Client) ([]*json.RawMessage, error) {
	blms := make([]ethrpc.BatchElem, 0)
	for i := from; i <= to; i = i + 1 {
		r := new(json.RawMessage)
		var err error
		blms = append(blms, ethrpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{"0x" + strconv.FormatUint(i, 16), true},
			Result: r,
			Error:  err,
		})
	}
	log.Trace().Uint64("start", from).Uint64("end", to).Msg("Fetching block range")

	err := c.BatchCallContext(ctx, blms)
	if err != nil {
		log.Error().Err(err).Msg("RPC issue fetching blocks")
		return nil, err
	}
	blocks := make([]*json.RawMessage, 0)

	for _, b := range blms {
		if b.Error != nil {
			return nil, b.Error
		}
		blocks = append(blocks, b.Result.(*json.RawMessage))

	}

	return blocks, nil
}

func GetBlockRangeInPages(ctx context.Context, from, to, pageSize uint64, c *ethrpc.Client) ([]*json.RawMessage, error) {
	var allBlocks []*json.RawMessage

	for i := from; i <= to; i += pageSize {
		end := i + pageSize - 1
		if end > to {
			end = to
		}

		blocks, err := GetBlockRange(ctx, i, end, c)
		if err != nil {
			return nil, err
		}

		allBlocks = append(allBlocks, blocks...)
	}

	return allBlocks, nil
}

func GetReceipts(ctx context.Context, rawBlocks []*json.RawMessage, c *ethrpc.Client, batchSize uint64) ([]*json.RawMessage, error) {
	txHashes := make([]string, 0)
	txHashMap := make(map[string]string, 0)
	for _, rb := range rawBlocks {
		var block simpleRPCBlock
		err := json.Unmarshal(*rb, &block)
		if err != nil {
			return nil, err

		}
		for _, tx := range block.Transactions {
			txHashes = append(txHashes, tx.Hash)
			txHashMap[tx.Hash] = block.Number
		}

	}
	if len(txHashes) == 0 {
		return nil, nil
	}

	blms := make([]ethrpc.BatchElem, 0)
	blmsBlockMap := make(map[int]string, 0)
	for i, tx := range txHashes {
		r := new(json.RawMessage)
		var err error
		blms = append(blms, ethrpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []interface{}{tx},
			Result: r,
			Error:  err,
		})
		blmsBlockMap[i] = txHashMap[tx]
	}

	var start uint64 = 0
	for {
		last := false
		end := start + batchSize
		if int(end) > len(blms) {
			last = true
			end = uint64(len(blms))
		}

		log.Trace().Str("startblock", blmsBlockMap[int(start)]).Uint64("start", start).Uint64("end", end).Msg("Fetching tx receipt range")
		// json: cannot unmarshal object into Go value of type []rpc.jsonrpcMessage
		// The error occurs when we call batchcallcontext with a single transaction for some reason.
		// polycli dumpblocks -c 1 http://127.0.0.1:9209/ 34457958 34458108
		// To handle this I'm making an exception when start and end are equal to make a single call.
		if start == end {
			log.Trace().Int("length", len(blmsBlockMap)).Msg("Test Jesse")
			if len(blmsBlockMap) == int(start) {
				start = start - 1
			}
			err := c.CallContext(ctx, &blms[start].Result, "eth_getTransactionReceipt", blms[start].Args[0])
			if err != nil {
				log.Error().Err(err).Uint64("start", start).Uint64("end", end).Msg("RPC issue fetching single receipt")
				return nil, err
			}
			break
		}

		err := c.BatchCallContext(ctx, blms[start:end])
		if err != nil {
			log.Error().Err(err).Str("randtx", txHashes[0]).Uint64("start", start).Uint64("end", end).Msg("RPC issue fetching receipts")
			return nil, err
		}
		start = end
		if last {
			break
		}
	}

	receipts := make([]*json.RawMessage, 0)

	for _, b := range blms {
		if b.Error != nil {
			log.Error().Err(b.Error).Msg("Block response err")
			return nil, b.Error
		}
		receipts = append(receipts, b.Result.(*json.RawMessage))
	}
	log.Info().Int("hashes", len(txHashes)).Int("receipts", len(receipts)).Msg("Fetched tx receipts")
	return receipts, nil
}

func GetTxPoolStatus(rpc *ethrpc.Client) (uint64, uint64, error) {
	var status = new(txpoolStatus)
	err := rpc.Call(status, "txpool_status")
	if err != nil {
		return 0, 0, err
	}
	pendingCount, err := tryCastToUint64(status.Pending)
	if err != nil {
		return 0, 0, err
	}
	queuedCount, err := tryCastToUint64(status.Queued)
	if err != nil {
		return pendingCount, 0, err
	}

	return pendingCount, queuedCount, nil
}

func tryCastToUint64(val any) (uint64, error) {
	switch t := val.(type) {
	case float64:
		return uint64(t), nil
	case string:
		return convHexToUint64(t)
	default:
		return 0, fmt.Errorf("the value %v couldn't be marshalled to uint64", t)

	}
}
func convHexToUint64(hexString string) (uint64, error) {
	hexString = strings.TrimPrefix(hexString, "0x")
	if len(hexString)%2 != 0 {
		hexString = "0" + hexString
	}

	result, err := strconv.ParseUint(hexString, 16, 64)
	if err != nil {
		return 0, err
	}
	return uint64(result), nil
}

// BlockUntilSuccessfulFn is designed to wait until a specified number of Ethereum blocks have been
// mined, periodically checking for the completion of a given function within each block interval.
type BlockUntilSuccessfulFn func(ctx context.Context, c *ethclient.Client, f func() error) error

func BlockUntilSuccessful(ctx context.Context, c *ethclient.Client, retryable func() error) error {
	// this function use to be very complicated (and not work). I'm dumbing this down to a basic time based retryable which should work 99% of the time
	b := backoff.WithContext(backoff.WithMaxRetries(backoff.NewConstantBackOff(5*time.Second), 24), ctx)
	return backoff.Retry(retryable, b)
}
