package p2p

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	ds "github.com/0xPolygon/polygon-cli/p2p/datastructures"
)

// validatorSetResponse is the subset of GET /stake/validators-set we need.
type validatorSetResponse struct {
	ValidatorSet struct {
		Validators []struct {
			Signer string `json:"signer"`
		} `json:"validators"`
	} `json:"validator_set"`
}

// ValidatorSet maintains the set of authorized block signers by periodically
// fetching the Heimdall validator set. It is used to decide whether a received
// block should be rebroadcast to peers.
type ValidatorSet struct {
	rest    *client.RESTClient
	refresh time.Duration
	signers ds.Locked[map[common.Address]struct{}]
}

// NewValidatorSet creates a ValidatorSet that fetches the validator set from
// the Heimdall REST gateway at restURL, refreshing every refresh interval.
func NewValidatorSet(restURL string, refresh time.Duration) *ValidatorSet {
	return &ValidatorSet{
		rest:    client.NewRESTClient(restURL, 30*time.Second, nil, false),
		refresh: refresh,
	}
}

// fetch retrieves the current validator set from Heimdall and returns the set
// of authorized signer addresses. It fails if no validators are parsed so we
// never operate on an empty/partial set.
func (v *ValidatorSet) fetch(ctx context.Context) (map[common.Address]struct{}, error) {
	body, _, err := v.rest.Get(ctx, "/stake/validators-set", nil)
	if err != nil {
		return nil, fmt.Errorf("fetching validator set: %w", err)
	}

	var resp validatorSetResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decoding validator set: %w", err)
	}

	signers := make(map[common.Address]struct{}, len(resp.ValidatorSet.Validators))
	for _, val := range resp.ValidatorSet.Validators {
		if val.Signer == "" {
			continue
		}
		signers[common.HexToAddress(val.Signer)] = struct{}{}
	}
	if len(signers) == 0 {
		return nil, fmt.Errorf("validator set contained no signers")
	}

	return signers, nil
}

// Start performs a blocking initial fetch of the validator set and then
// launches a background goroutine that refreshes it on the configured
// interval until ctx is cancelled. It returns an error if the initial fetch
// fails so the caller can abort startup rather than serve with no validators.
func (v *ValidatorSet) Start(ctx context.Context) error {
	signers, err := v.fetch(ctx)
	if err != nil {
		return err
	}
	v.signers.Set(signers)
	log.Info().Int("validators", len(signers)).Msg("Loaded validator set")

	go v.refreshLoop(ctx)
	return nil
}

// refreshLoop periodically refreshes the validator set. On a fetch error it
// logs a warning and keeps the previous set rather than wiping it.
func (v *ValidatorSet) refreshLoop(ctx context.Context) {
	ticker := time.NewTicker(v.refresh)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			signers, err := v.fetch(ctx)
			if err != nil {
				log.Warn().Err(err).Msg("Failed to refresh validator set, keeping previous set")
				continue
			}
			v.signers.Set(signers)
			log.Debug().Int("validators", len(signers)).Msg("Refreshed validator set")
		}
	}
}

// HasSigner reports whether addr is the signer address of a known validator.
func (v *ValidatorSet) HasSigner(addr common.Address) bool {
	signers := v.signers.Get()
	_, ok := signers[addr]
	return ok
}
