package gasmanager

import (
	"context"
	"math"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

// OscillatingGasProvider is a gas provider that adds gas to the vault based on an oscillating wave pattern.
type OscillatingGasProvider struct {
	GasProviderBase
	wave Wave
}

// NewOscillatingGasProvider creates a new OscillatingGasProvider with the given Ethereum client, gas vault, and wave pattern.
// It also sets up the necessary callbacks for starting and processing new block headers.
func NewOscillatingGasProvider(client *ethclient.Client, vault *GasVault, wave Wave) *OscillatingGasProvider {
	p := &OscillatingGasProvider{
		GasProviderBase: *NewGasProviderBase(client, vault),
		wave:            wave,
	}

	p.GasProviderBase.onNewHeader = p.onNewHeader
	return p
}

// Start begins the operation of the OscillatingGasProvider by invoking the Start method of its base class.
func (o *OscillatingGasProvider) Start(ctx context.Context) {
	o.GasProviderBase.Start(ctx)
}

// onNewHeader is called when a new block header is received.
// It advances the wave and adds gas to the vault based on the new Y value of the wave.
func (o *OscillatingGasProvider) onNewHeader(header *types.Header) {
	log.Trace().Uint64("block_number", header.Number.Uint64()).Msg("Processing new block header")
	if o.vault != nil {
		gasAmount := uint64(math.Floor(o.wave.Y()))
		o.vault.AddGas(gasAmount)
		log.Trace().
			Uint64("gas_added", gasAmount).
			Uint64("available_budget", o.vault.GetAvailableBudget()).
			Msg("Gas added from oscillation wave")
	}
	o.wave.MoveNext()
}
