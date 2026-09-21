package p2p

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/consensus/misc/eip1559"
	"github.com/ethereum/go-ethereum/core/types"
)

// TestCalcBaseFeeAgainstChain measures whether go-ethereum's EIP-1559 base fee
// formula reproduces the next block's base fee on a live chain, using the chain
// config the transaction validator would use.
//
// It is skipped unless BASEFEE_RPC is set, because it talks to the network.
//
//	BASEFEE_RPC=https://polygon-bor-rpc.publicnode.com go test ./p2p/ -run CalcBaseFee -v
//
// The question it answers: bor gates its re-broadcast of stuck transactions on
// eip1559.CalcBaseFee(config, head), so copying that here would need the same
// number. But bor's CalcBaseFee is not geth's -- it derives the gas target from
// a Dandeli percentage that itself varies with the base fee, takes the change
// denominator from a Bhilai-era schedule, and caps the per-block change at 5%
// post-Lisovo. All three read config.Bor, which a sensor does not have.
func TestCalcBaseFeeAgainstChain(t *testing.T) {
	url := os.Getenv("BASEFEE_RPC")
	if url == "" {
		t.Skip("set BASEFEE_RPC to run against a live chain")
	}

	chainID := uint64(137)
	if v := os.Getenv("BASEFEE_CHAIN_ID"); v != "" {
		if _, err := fmt.Sscanf(v, "%d", &chainID); err != nil {
			t.Fatalf("BASEFEE_CHAIN_ID: %v", err)
		}
	}
	config := broadcastChainConfig(chainID)

	head, err := blockNumber(url)
	if err != nil {
		t.Fatalf("eth_blockNumber: %v", err)
	}

	const samples = 30
	var exact, off int
	var worstPct float64

	for i := samples; i > 0; i-- {
		parent, err := headerByNumber(url, head-uint64(i))
		if err != nil {
			t.Fatalf("fetching parent: %v", err)
		}
		child, err := headerByNumber(url, head-uint64(i)+1)
		if err != nil {
			t.Fatalf("fetching child: %v", err)
		}
		if parent.BaseFee == nil || child.BaseFee == nil {
			t.Skip("chain has no base fee")
		}

		got := eip1559.CalcBaseFee(config, parent)
		if got.Cmp(child.BaseFee) == 0 {
			exact++
			continue
		}
		off++

		diff := new(big.Int).Sub(got, child.BaseFee)
		pct, _ := new(big.Float).Quo(
			new(big.Float).SetInt(new(big.Int).Abs(diff)),
			new(big.Float).SetInt(child.BaseFee),
		).Float64()
		if pct*100 > worstPct {
			worstPct = pct * 100
		}
		if off <= 3 {
			t.Logf("block %d: want %v, geth formula gives %v (%+v, %.2f%%), parent gasUsed %d of %d",
				child.Number, child.BaseFee, got, diff, pct*100, parent.GasUsed, parent.GasLimit)
		}
	}

	t.Logf("RESULT: %d/%d exact, %d off, worst error %.2f%%", exact, exact+off, off, worstPct)
	if off > 0 {
		t.Logf("go-ethereum's formula does NOT reproduce this chain's base fee; a sensor cannot " +
			"compute the next base fee without the chain's consensus parameters")
	}
}

func rpcCall(url, method string, params []any, out any) error {
	body, err := json.Marshal(map[string]any{
		"jsonrpc": "2.0", "id": 1, "method": method, "params": params,
	})
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	var envelope struct {
		Result json.RawMessage `json:"result"`
		Error  *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return err
	}
	if envelope.Error != nil {
		return fmt.Errorf("rpc: %s", envelope.Error.Message)
	}
	return json.Unmarshal(envelope.Result, out)
}

func blockNumber(url string) (uint64, error) {
	var hex string
	if err := rpcCall(url, "eth_blockNumber", []any{}, &hex); err != nil {
		return 0, err
	}
	n := new(big.Int)
	if _, ok := n.SetString(hex[2:], 16); !ok {
		return 0, fmt.Errorf("bad block number %q", hex)
	}
	return n.Uint64(), nil
}

// headerByNumber fetches only the fields CalcBaseFee reads.
func headerByNumber(url string, number uint64) (*types.Header, error) {
	var raw struct {
		Number   string `json:"number"`
		GasLimit string `json:"gasLimit"`
		GasUsed  string `json:"gasUsed"`
		BaseFee  string `json:"baseFeePerGas"`
	}
	if err := rpcCall(url, "eth_getBlockByNumber", []any{fmt.Sprintf("0x%x", number), false}, &raw); err != nil {
		return nil, err
	}

	parse := func(s string) (*big.Int, error) {
		if s == "" {
			return nil, nil
		}
		n, ok := new(big.Int).SetString(s[2:], 16)
		if !ok {
			return nil, fmt.Errorf("bad hex %q", s)
		}
		return n, nil
	}
	num, err := parse(raw.Number)
	if err != nil {
		return nil, err
	}
	gasLimit, err := parse(raw.GasLimit)
	if err != nil {
		return nil, err
	}
	gasUsed, err := parse(raw.GasUsed)
	if err != nil {
		return nil, err
	}
	baseFee, err := parse(raw.BaseFee)
	if err != nil {
		return nil, err
	}
	return &types.Header{
		Number:     num,
		GasLimit:   gasLimit.Uint64(),
		GasUsed:    gasUsed.Uint64(),
		BaseFee:    baseFee,
		Difficulty: big.NewInt(1),
	}, nil
}
