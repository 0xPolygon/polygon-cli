package gaslimiter

import (
	"context"
	"math"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type OscillatingGasProvider struct {
	GasProviderBase
	oscillationCurve Curve
}

func NewOscillatingGasProvider(client *ethclient.Client, vault *GasVault, oscillationCurve Curve) *OscillatingGasProvider {
	p := &OscillatingGasProvider{
		GasProviderBase:  *NewGasProviderBase(client, vault),
		oscillationCurve: oscillationCurve,
	}

	p.GasProviderBase.onStart = p.onStart
	p.GasProviderBase.onNewHeader = p.onNewHeader
	return p
}

func (o *OscillatingGasProvider) Start(ctx context.Context) {
	o.GasProviderBase.Start(ctx)
}

func (o *OscillatingGasProvider) onStart() {
	o.vault.AddGas(uint64(math.Floor(o.oscillationCurve.Y())))
}

func (o *OscillatingGasProvider) onNewHeader(header *types.Header) {
	o.oscillationCurve.MoveNext()
	if o.vault != nil {
		o.vault.AddGas(uint64(math.Floor(o.oscillationCurve.Y())))
	}
}
