package gasmanager

import (
	"context"
	"math"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

type OscillatingGasProvider struct {
	GasProviderBase
	wave Wave
}

func NewOscillatingGasProvider(client *ethclient.Client, vault *GasVault, wave Wave) *OscillatingGasProvider {
	p := &OscillatingGasProvider{
		GasProviderBase: *NewGasProviderBase(client, vault),
		wave:            wave,
	}

	p.GasProviderBase.onStart = p.onStart
	p.GasProviderBase.onNewHeader = p.onNewHeader
	return p
}

func (o *OscillatingGasProvider) Start(ctx context.Context) {
	o.GasProviderBase.Start(ctx)
}

func (o *OscillatingGasProvider) onStart() {
	o.vault.AddGas(uint64(math.Floor(o.wave.Y())))
}

func (o *OscillatingGasProvider) onNewHeader(header *types.Header) {
	log.Trace().Uint64("block_number", header.Number.Uint64()).Msg("oscillating gas provider processing new block header")
	o.wave.MoveNext()
	if o.vault != nil {
		log.Trace().Float64("new_gas_amount", o.wave.Y()).Msg("adding gas to vault based on oscillation wave")
		o.vault.AddGas(uint64(math.Floor(o.wave.Y())))
		log.Trace().Uint64("available_budget", o.vault.GetAvailableBudget()).Msg("updated gas vault budget after oscillation")
	}
}
