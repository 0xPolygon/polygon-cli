package p2p

import (
	"math/big"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Gas price oracle constants (matching Bor/geth defaults)
const (
	// gpoSampleNumber is the number of transactions to sample per block
	gpoSampleNumber = 3
	// gpoCheckBlocks is the number of blocks to check for gas price estimation
	gpoCheckBlocks = 20
	// gpoPercentile is the percentile to use for gas price estimation
	gpoPercentile = 60
)

var (
	// gpoMaxPrice is the maximum gas price to suggest (500 gwei)
	gpoMaxPrice = big.NewInt(500_000_000_000)
	// gpoIgnorePrice is the minimum tip to consider (2 gwei, lower than Bor's 25 gwei for broader network compatibility)
	gpoIgnorePrice = big.NewInt(2_000_000_000)
	// gpoDefaultPrice is the default gas price when no data is available (1 gwei)
	gpoDefaultPrice = big.NewInt(1_000_000_000)
)

// GasPriceOracle estimates gas prices based on recent block data.
// It follows Bor/geth's gas price oracle approach.
type GasPriceOracle struct {
	conns *Conns

	mu       sync.RWMutex
	lastHead common.Hash
	lastTip  *big.Int
}

// NewGasPriceOracle creates a new gas price oracle that uses the given Conns for block data.
func NewGasPriceOracle(conns *Conns) *GasPriceOracle {
	return &GasPriceOracle{
		conns: conns,
	}
}

// SuggestGasPrice estimates the gas price based on recent blocks.
// For EIP-1559 networks, this returns baseFee + suggestedTip.
// For legacy networks, this returns the 60th percentile of gas prices.
func (o *GasPriceOracle) SuggestGasPrice() *big.Int {
	head := o.conns.HeadBlock()
	if head.Block == nil {
		return gpoDefaultPrice
	}

	// For EIP-1559: return baseFee + suggested tip
	if baseFee := head.Block.BaseFee(); baseFee != nil {
		tip := o.SuggestGasTipCap()
		if tip == nil {
			tip = gpoDefaultPrice
		}
		return new(big.Int).Add(baseFee, tip)
	}

	// Legacy: return percentile of gas prices
	return o.suggestLegacyGasPrice()
}

// suggestLegacyGasPrice estimates gas price for pre-EIP-1559 networks.
func (o *GasPriceOracle) suggestLegacyGasPrice() *big.Int {
	keys := o.conns.blocks.Keys()
	if len(keys) == 0 {
		return gpoDefaultPrice
	}

	if len(keys) > gpoCheckBlocks {
		keys = keys[:gpoCheckBlocks]
	}

	var prices []*big.Int
	for _, hash := range keys {
		cache, ok := o.conns.blocks.Peek(hash)
		if !ok || cache.Body == nil {
			continue
		}

		txs, err := cache.Body.Transactions.Items()
		if err != nil {
			continue
		}
		for _, tx := range txs {
			if price := tx.GasPrice(); price != nil && price.Sign() > 0 {
				prices = append(prices, new(big.Int).Set(price))
			}
		}
	}

	if len(prices) == 0 {
		return gpoDefaultPrice
	}

	sort.Slice(prices, func(i, j int) bool {
		return prices[i].Cmp(prices[j]) < 0
	})

	price := prices[(len(prices)-1)*gpoPercentile/100]
	if price.Cmp(gpoMaxPrice) > 0 {
		return new(big.Int).Set(gpoMaxPrice)
	}
	return price
}

// SuggestGasTipCap estimates a gas tip cap (priority fee) based on recent blocks.
// This implementation follows Bor/geth's gas price oracle approach:
// - Samples the lowest N tips from each of the last M blocks
// - Ignores tips below a threshold
// - Returns the configured percentile of collected tips
// - Caches results until head changes
func (o *GasPriceOracle) SuggestGasTipCap() *big.Int {
	head := o.conns.HeadBlock()
	if head.Block == nil {
		return nil
	}
	headHash := head.Block.Hash()

	// Check cache first
	o.mu.RLock()
	if headHash == o.lastHead && o.lastTip != nil {
		tip := new(big.Int).Set(o.lastTip)
		o.mu.RUnlock()
		return tip
	}
	lastTip := o.lastTip
	o.mu.RUnlock()

	// Collect tips from recent blocks
	keys := o.conns.blocks.Keys()
	if len(keys) == 0 {
		return lastTip
	}

	// Limit to checkBlocks most recent
	if len(keys) > gpoCheckBlocks {
		keys = keys[:gpoCheckBlocks]
	}

	var results []*big.Int
	for _, hash := range keys {
		tips := o.getBlockTips(hash, gpoSampleNumber, gpoIgnorePrice)
		if len(tips) == 0 && lastTip != nil {
			// Empty block or all tips below threshold, use last tip
			tips = []*big.Int{lastTip}
		}
		results = append(results, tips...)
	}

	if len(results) == 0 {
		return lastTip
	}

	// Sort and get percentile
	sort.Slice(results, func(i, j int) bool {
		return results[i].Cmp(results[j]) < 0
	})
	tip := results[(len(results)-1)*gpoPercentile/100]

	// Apply max price cap
	if tip.Cmp(gpoMaxPrice) > 0 {
		tip = new(big.Int).Set(gpoMaxPrice)
	}

	// Cache result
	o.mu.Lock()
	o.lastHead = headHash
	o.lastTip = tip
	o.mu.Unlock()

	return new(big.Int).Set(tip)
}

// getBlockTips returns the lowest N tips from a block that are above the ignore threshold.
// Transactions are sorted by effective tip ascending, and the first N valid tips are returned.
func (o *GasPriceOracle) getBlockTips(hash common.Hash, limit int, ignoreUnder *big.Int) []*big.Int {
	cache, ok := o.conns.blocks.Peek(hash)
	if !ok || cache.Body == nil || cache.Header == nil {
		return nil
	}

	baseFee := cache.Header.BaseFee
	if baseFee == nil {
		return nil // Pre-EIP-1559 block
	}

	// Calculate tips for all transactions
	txs, err := cache.Body.Transactions.Items()
	if err != nil {
		return nil
	}
	var allTips []*big.Int
	for _, tx := range txs {
		tip := effectiveGasTip(tx, baseFee)
		if tip != nil && tip.Sign() > 0 {
			allTips = append(allTips, tip)
		}
	}

	if len(allTips) == 0 {
		return nil
	}

	// Sort by tip ascending (lowest first, like Bor)
	sort.Slice(allTips, func(i, j int) bool {
		return allTips[i].Cmp(allTips[j]) < 0
	})

	// Collect tips above threshold, up to limit
	var tips []*big.Int
	for _, tip := range allTips {
		if ignoreUnder != nil && tip.Cmp(ignoreUnder) < 0 {
			continue
		}
		tips = append(tips, tip)
		if len(tips) >= limit {
			break
		}
	}

	return tips
}

// effectiveGasTip returns the effective tip (priority fee) for a transaction.
// For EIP-1559 transactions: min(maxPriorityFeePerGas, maxFeePerGas - baseFee)
// For legacy transactions: gasPrice - baseFee (the implicit tip)
// Returns nil if the tip cannot be determined or is negative.
func effectiveGasTip(tx *types.Transaction, baseFee *big.Int) *big.Int {
	switch tx.Type() {
	case types.DynamicFeeTxType, types.BlobTxType:
		tip := tx.GasTipCap()
		if tip == nil {
			return nil
		}
		// Effective tip is min(maxPriorityFeePerGas, maxFeePerGas - baseFee)
		if tx.GasFeeCap() != nil {
			effectiveTip := new(big.Int).Sub(tx.GasFeeCap(), baseFee)
			if effectiveTip.Cmp(tip) < 0 {
				tip = effectiveTip
			}
		}
		if tip.Sign() <= 0 {
			return nil
		}
		return new(big.Int).Set(tip)
	default:
		// Legacy/AccessList transactions: tip is gasPrice - baseFee
		if price := tx.GasPrice(); price != nil {
			tip := new(big.Int).Sub(price, baseFee)
			if tip.Sign() > 0 {
				return tip
			}
		}
		return nil
	}
}
