package fund

import (
	"errors"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// addr turns a small int into a distinct, non-zero address.
func addr(i int) common.Address {
	var out common.Address
	out[0] = 0xaa
	out[18] = byte(i >> 8)
	out[19] = byte(i)
	return out
}

func makeWallets(n int) []common.Address {
	wallets := make([]common.Address, n)
	for i := range wallets {
		wallets[i] = addr(i + 1)
	}
	return wallets
}

func TestSplitIntoBatches(t *testing.T) {
	funder := common.HexToAddress("0xfeed")

	tests := []struct {
		name      string
		wallets   []common.Address
		batchSize uint64
		wantSizes []int
	}{
		{name: "empty", wallets: nil, batchSize: 4, wantSizes: nil},
		{name: "zero batch size", wallets: makeWallets(3), batchSize: 0, wantSizes: nil},
		{name: "exact multiple", wallets: makeWallets(8), batchSize: 4, wantSizes: []int{4, 4}},
		{name: "remainder", wallets: makeWallets(9), batchSize: 4, wantSizes: []int{4, 4, 1}},
		{name: "smaller than batch", wallets: makeWallets(2), batchSize: 4, wantSizes: []int{2}},
		{
			name:      "funder skipped in the middle",
			wallets:   append(append(makeWallets(4), funder), makeWallets(4)...),
			batchSize: 4,
			wantSizes: []int{4, 4},
		},
		{name: "only funder", wallets: []common.Address{funder}, batchSize: 4, wantSizes: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			batches := splitIntoBatches(tt.wallets, funder, tt.batchSize)
			if len(batches) != len(tt.wantSizes) {
				t.Fatalf("got %d batches, want %d", len(batches), len(tt.wantSizes))
			}
			for i, b := range batches {
				if b.index != i {
					t.Errorf("batch %d has index %d", i, b.index)
				}
				if len(b.accounts) != tt.wantSizes[i] {
					t.Errorf("batch %d has %d accounts, want %d", i, len(b.accounts), tt.wantSizes[i])
				}
				for _, a := range b.accounts {
					if a == funder {
						t.Errorf("batch %d contains the funder address", i)
					}
				}
			}
		})
	}
}

func TestSplitIntoBatchesPreservesOrder(t *testing.T) {
	wallets := makeWallets(7)
	batches := splitIntoBatches(wallets, common.HexToAddress("0xfeed"), 3)
	var got []common.Address
	for _, b := range batches {
		got = append(got, b.accounts...)
	}
	if len(got) != len(wallets) {
		t.Fatalf("got %d accounts, want %d", len(got), len(wallets))
	}
	for i := range wallets {
		if got[i] != wallets[i] {
			t.Fatalf("account %d is %s, want %s", i, got[i], wallets[i])
		}
	}
}

func TestFirstFailedBatch(t *testing.T) {
	batches := splitIntoBatches(makeWallets(12), common.HexToAddress("0xfeed"), 4)
	if got := firstFailedBatch(batches); got != -1 {
		t.Fatalf("no failures: got %d, want -1", got)
	}
	batches[2].err = errors.New("boom")
	batches[1].err = errors.New("boom")
	if got := firstFailedBatch(batches); got != 1 {
		t.Fatalf("got %d, want 1", got)
	}
}

func TestSummarizeSendFailures(t *testing.T) {
	batches := splitIntoBatches(makeWallets(10), common.HexToAddress("0xfeed"), 4)
	for _, b := range batches {
		b.nonce = 100 + uint64(b.index)
	}
	sendErr := errors.New("replacement transaction underpriced")
	batches[1].err = sendErr

	err := summarizeSendFailures(batches, firstFailedBatch(batches))
	if err == nil {
		t.Fatal("expected an error")
	}
	if !errors.Is(err, sendErr) {
		t.Errorf("summary does not wrap the send error: %v", err)
	}
	msg := err.Error()
	for _, want := range []string{
		"1 of 3 funding transactions were not sent",
		"4 accounts were not funded",
		"1 transaction(s) with nonces above 101",
		"fix-nonce-gap",
		"batch 2 (nonce 101, 4 accounts)",
	} {
		if !strings.Contains(msg, want) {
			t.Errorf("summary missing %q:\n%s", want, msg)
		}
	}
}

func TestSummarizeSendFailuresLastBatchNoStranded(t *testing.T) {
	batches := splitIntoBatches(makeWallets(10), common.HexToAddress("0xfeed"), 4)
	batches[2].err = errors.New("boom")

	msg := summarizeSendFailures(batches, firstFailedBatch(batches)).Error()
	if strings.Contains(msg, "stay pending") {
		t.Errorf("no batches were sent above the gap, but summary mentions stranded txs:\n%s", msg)
	}
	if !strings.Contains(msg, "2 accounts were not funded") {
		t.Errorf("summary has wrong unfunded count:\n%s", msg)
	}
}
