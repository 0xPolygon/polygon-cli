package gasmanager

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

// GasProvider defines the interface for gas providers.
type GasProvider interface {
	Start(ctx context.Context)
}

// GasProviderBase provides common functionality for gas providers.
type GasProviderBase struct {
	client      *ethclient.Client
	vault       *GasVault
	onNewHeader func(header *types.Header)
}

// NewGasProviderBase creates a new GasProviderBase with the given Ethereum client and gas vault.
func NewGasProviderBase(client *ethclient.Client, vault *GasVault) *GasProviderBase {
	return &GasProviderBase{
		client: client,
		vault:  vault,
	}
}

// Start begins the operation of the GasProviderBase by starting to watch for new block headers.
func (o *GasProviderBase) Start(ctx context.Context) {
	log.Trace().Msg("Starting GasProviderBase")
	go o.watchNewHeaders(ctx)
}

// watchNewHeaders continuously monitors for new block headers and triggers the onNewHeader callback when a new header is detected.
func (o *GasProviderBase) watchNewHeaders(ctx context.Context) {
	if o.onNewHeader == nil {
		return
	}

	const pollInterval = 1 * time.Second
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	log.Trace().Msg("Starting to watch for new block headers")
	var lastHeader *types.Header

	for {
		select {
		case <-ctx.Done():
			log.Trace().Msg("Stopping block header watch")
			return
		case <-ticker.C:
			header, err := o.client.HeaderByNumber(ctx, nil)
			if err != nil {
				log.Warn().Err(err).Msg("Failed to fetch latest block header, retrying")
				continue
			}

			if lastHeader == nil || header.Number.Cmp(lastHeader.Number) > 0 {
				log.Trace().Uint64("block_number", header.Number.Uint64()).Msg("New block header detected")
				o.onNewHeader(header)
				lastHeader = header
			}
		}
	}
}
