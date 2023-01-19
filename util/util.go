package util

import (
	"context"
	"encoding/json"
	"strconv"

	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
)

type (
	simpleRPCTransaction struct {
		Hash string `json:"hash"`
	}
	simpleRPCBlock struct {
		Number       string                 `json:"number"`
		Transactions []simpleRPCTransaction `json:"transactions"`
	}
)

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
			log.Error().Err(b.Error).Msg("Block res err")
			return nil, b.Error
		}
		receipts = append(receipts, b.Result.(*json.RawMessage))
	}
	log.Info().Int("hashes", len(txHashes)).Int("receipts", len(receipts)).Msg("Fetched tx receipts")
	return receipts, nil
}
