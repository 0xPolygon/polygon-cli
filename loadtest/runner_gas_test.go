package loadtest

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/0xPolygon/polygon-cli/loadtest/config"
	"github.com/ethereum/go-ethereum/ethclient"
)

type rpcReq struct {
	Method string `json:"method"`
	ID     int    `json:"id"`
}

func TestRunnerGasPriceDecrease(t *testing.T) {
	var baseFee int64 = 100
	var blockNum uint64 = 1000

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req rpcReq
		json.NewDecoder(r.Body).Decode(&req)
		w.Header().Set("Content-Type", "application/json")
		switch req.Method {
		case "eth_blockNumber":
			fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"result":"0x%x"}`, req.ID, blockNum)
		case "eth_getBlockByNumber":
			fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"result":{"number":"0x%x","baseFeePerGas":"0x%x","parentHash":"0x0000000000000000000000000000000000000000000000000000000000000000","difficulty":"0x0","extraData":"0x","gasLimit":"0x0","gasUsed":"0x0","miner":"0x0000000000000000000000000000000000000000","mixHash":"0x0000000000000000000000000000000000000000000000000000000000000000","nonce":"0x0000000000000000","receiptsRoot":"0x0000000000000000000000000000000000000000000000000000000000000000","sha3Uncles":"0x0000000000000000000000000000000000000000000000000000000000000000","stateRoot":"0x0000000000000000000000000000000000000000000000000000000000000000","timestamp":"0x0","transactionsRoot":"0x0000000000000000000000000000000000000000000000000000000000000000","logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"}}`, req.ID, blockNum, baseFee)
		case "eth_feeHistory":
			// Return a fee history where the latest base fee is our current baseFee
			fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"result":{"baseFeePerGas":["0x%x","0x%x"],"reward":[["0x%x"]]}}`, req.ID, baseFee, baseFee, 10)
		case "eth_maxPriorityFeePerGas":
			fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"result":"0xa"}`, req.ID)
		case "eth_chainId":
			fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"result":"0x89"}`, req.ID)
		default:
			fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"result":"0x0"}`, req.ID)
		}
	}))
	defer ts.Close()

	client, _ := ethclient.Dial(ts.URL)
	r := &Runner{
		cfg: &config.Config{
			ChainSupportBaseFee:   true,
			BigGasPriceMultiplier: big.NewFloat(1.0),
		},
		client: client,
	}

	ctx := context.Background()

	// 1. Initial price
	gp1, _ := r.getSuggestedGasPrices(ctx)
	// baseFee 100 * 2 + priority 10 = 210
	if gp1.Int64() != 210 {
		t.Errorf("Expected 210, got %v", gp1)
	}

	// 2. Decrease base fee and advance block
	baseFee = 50
	blockNum = 1001
	// Reset the cache to force a new fetch
	r.cachedLatestBlockTime = time.Time{}

	gp2, _ := r.getSuggestedGasPrices(ctx)
	// new baseFee 50 * 2 + priority 10 = 110
	// Before the fix, this would have returned 210 because of the broken blocksToWait logic.
	if gp2.Int64() != 110 {
		t.Errorf("Expected 110, got %v", gp2)
	}
}
