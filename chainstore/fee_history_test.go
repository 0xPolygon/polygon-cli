package chainstore

import (
	"context"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"
)

type feeHistoryRPCCall struct {
	blockCount        string
	newestBlock       string
	rewardPercentiles []float64
}

type feeHistoryRPCService struct {
	mu    sync.Mutex
	calls []feeHistoryRPCCall
}

func (s *feeHistoryRPCService) FeeHistory(_ context.Context, blockCount string, newestBlock string, rewardPercentiles []float64) (*FeeHistoryResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.calls = append(s.calls, feeHistoryRPCCall{
		blockCount:        blockCount,
		newestBlock:       newestBlock,
		rewardPercentiles: append([]float64(nil), rewardPercentiles...),
	})
	callNumber := int64(len(s.calls))

	return &FeeHistoryResult{
		OldestBlock:   big.NewInt(callNumber),
		BaseFeePerGas: []*big.Int{big.NewInt(callNumber)},
		GasUsedRatio:  []float64{float64(callNumber)},
	}, nil
}

func (s *feeHistoryRPCService) callCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.calls)
}

func newFeeHistoryTestStore(t *testing.T, ttl time.Duration) (*PassthroughStore, *feeHistoryRPCService) {
	t.Helper()

	server := rpc.NewServer()
	service := &feeHistoryRPCService{}
	require.NoError(t, server.RegisterName("eth", service))

	client := rpc.DialInProc(server)
	capabilities := NewCapabilityManager(client, time.Hour)
	capabilities.capabilities["eth_feeHistory"] = true
	capabilities.lastChecked = time.Now()

	store := &PassthroughStore{
		client:       client,
		cache:        NewChainCache(),
		capabilities: capabilities,
		config: &ChainStoreConfig{
			FrequentTTL: ttl,
		},
	}

	t.Cleanup(func() {
		client.Close()
		server.Stop()
	})

	return store, service
}

func TestGetFeeHistoryCachesIdenticalRequests(t *testing.T) {
	store, service := newFeeHistoryTestStore(t, time.Minute)
	percentiles := []float64{10, 50, 90}

	first, err := store.GetFeeHistory(t.Context(), 10, "latest", percentiles)
	require.NoError(t, err)

	percentiles[0] = 25
	second, err := store.GetFeeHistory(t.Context(), 10, "latest", []float64{10, 50, 90})
	require.NoError(t, err)

	require.Equal(t, 1, service.callCount())
	require.Equal(t, first, second)
}

func TestGetFeeHistorySeparatesRequestsByParameters(t *testing.T) {
	tests := []struct {
		name                    string
		firstBlockCount         int
		firstNewestBlock        string
		firstRewardPercentiles  []float64
		secondBlockCount        int
		secondNewestBlock       string
		secondRewardPercentiles []float64
	}{
		{
			name:                    "block count",
			firstBlockCount:         1,
			firstNewestBlock:        "latest",
			firstRewardPercentiles:  []float64{10, 50, 90},
			secondBlockCount:        100,
			secondNewestBlock:       "latest",
			secondRewardPercentiles: []float64{10, 50, 90},
		},
		{
			name:                    "newest block",
			firstBlockCount:         10,
			firstNewestBlock:        "latest",
			firstRewardPercentiles:  []float64{10, 50, 90},
			secondBlockCount:        10,
			secondNewestBlock:       "0x1234",
			secondRewardPercentiles: []float64{10, 50, 90},
		},
		{
			name:                    "reward percentiles",
			firstBlockCount:         10,
			firstNewestBlock:        "latest",
			firstRewardPercentiles:  []float64{10, 50, 90},
			secondBlockCount:        10,
			secondNewestBlock:       "latest",
			secondRewardPercentiles: []float64{25, 50, 75},
		},
		{
			name:                    "nil and empty reward percentiles",
			firstBlockCount:         10,
			firstNewestBlock:        "latest",
			firstRewardPercentiles:  nil,
			secondBlockCount:        10,
			secondNewestBlock:       "latest",
			secondRewardPercentiles: []float64{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store, service := newFeeHistoryTestStore(t, time.Minute)

			first, err := store.GetFeeHistory(t.Context(), tt.firstBlockCount, tt.firstNewestBlock, tt.firstRewardPercentiles)
			require.NoError(t, err)
			second, err := store.GetFeeHistory(t.Context(), tt.secondBlockCount, tt.secondNewestBlock, tt.secondRewardPercentiles)
			require.NoError(t, err)

			require.Equal(t, 2, service.callCount())
			require.NotEqual(t, first.OldestBlock, second.OldestBlock)
		})
	}
}

func TestGetFeeHistoryRefetchesAfterTTLExpires(t *testing.T) {
	const ttl = 10 * time.Millisecond
	store, service := newFeeHistoryTestStore(t, ttl)

	first, err := store.GetFeeHistory(t.Context(), 10, "latest", []float64{10, 50, 90})
	require.NoError(t, err)

	time.Sleep(2 * ttl)

	second, err := store.GetFeeHistory(t.Context(), 10, "latest", []float64{10, 50, 90})
	require.NoError(t, err)

	require.Equal(t, 2, service.callCount())
	require.NotEqual(t, first.OldestBlock, second.OldestBlock)
}

func TestSetFeeHistoryPrunesExpiredEntries(t *testing.T) {
	cache := NewChainCache()
	cache.SetFeeHistory(1, "latest", nil, &FeeHistoryResult{}, time.Nanosecond)

	require.Eventually(t, func() bool {
		_, valid := cache.GetFeeHistory(1, "latest", nil)
		return !valid
	}, time.Second, time.Millisecond)

	cache.SetFeeHistory(2, "latest", nil, &FeeHistoryResult{}, time.Minute)
	require.Len(t, cache.feeHistories, 1)
}
