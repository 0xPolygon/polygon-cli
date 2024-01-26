package ecrecover

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	_ "embed"

	"github.com/ethereum/go-ethereum/common"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/maticnetwork/polygon-cli/util"

	"github.com/rs/zerolog/log"
)

func ecrecover(ctx context.Context) (common.Address, error) {
	rpc, err := ethrpc.DialContext(ctx, rpcUrl)
	if err != nil {
		log.Error().Err(err).Msg("Unable to dial rpc")
		return common.Address{}, err
	}

	ec := ethclient.NewClient(rpc)
	if _, err = ec.BlockNumber(ctx); err != nil {
		return common.Address{}, err
	}

	block, err := ec.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
	if err != nil {
		log.Error().Err(err).Msg("Unable to retrieve block")
		return common.Address{}, err
	}

	if len(block.Transactions()) == 0 {
		return common.Address{}, fmt.Errorf("no transaction to derive public key from")
	}

	signerBytes, err := util.Ecrecover(block)
	if err != nil {
		log.Error().Err(err).Msg("Unable to recover signature")
		return common.Address{}, err
	}
	signerAddress := ethcommon.BytesToAddress(signerBytes)

	chainID, err := ec.NetworkID(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get network ID")
	}

	var publicKey *ecdsa.PublicKey
	for _, tx := range block.Transactions() {
		tx.Value()
		if msg, err := tx.AsMessage(types.LatestSignerForChainID(chainID), nil); err == nil {
			if msg.From().Hex() == signerAddress {
				publicKeyECDSA, err := crypto.SigToPub(signHash(msg), msg.Signature())
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to recover public key")
				}
				publicKey = publicKeyECDSA
				break
			}
		}
	}

}
