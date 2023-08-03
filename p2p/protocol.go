package p2p

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	ethp2p "github.com/ethereum/go-ethereum/p2p"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/forkid"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"

	"github.com/maticnetwork/polygon-cli/rpctypes"
)

func NewEth66Protocol(genesis *core.Genesis, genesisHash common.Hash, url string, networkID uint64) ethp2p.Protocol {
	return ethp2p.Protocol{
		Name:    "eth",
		Version: 66,
		Length:  17,
		Run: func(p *ethp2p.Peer, rw ethp2p.MsgReadWriter) error {
			log.Info().Interface("peer", p.Info().Enode).Send()

			block, err := getLatestBlock(url)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get latest block")
				return err
			}

			err = statusExchange(rw, &eth.StatusPacket{
				ProtocolVersion: 66,
				NetworkID:       networkID,
				Genesis:         genesisHash,
				ForkID:          forkid.NewID(genesis.Config, genesisHash, block.Number.ToUint64()),
				Head:            block.Hash.ToHash(),
				TD:              block.TotalDifficulty.ToBigInt(),
			})
			if err != nil {
				return err
			}

			for {
				handleMessage(rw)
			}
		},
	}
}

func getLatestBlock(url string) (*rpctypes.RawBlockResponse, error) {
	client, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	var block rpctypes.RawBlockResponse
	err = client.Call(&block, "eth_getBlockByNumber", "latest", true)
	if err != nil {
		return nil, err
	}

	return &block, nil
}

func statusExchange(rw ethp2p.MsgReadWriter, packet *eth.StatusPacket) error {
	err := ethp2p.Send(rw, eth.StatusMsg, &packet)
	if err != nil {
		return err
	}

	msg, err := rw.ReadMsg()
	if err != nil {
		return err
	}

	var status eth.StatusPacket
	err = msg.Decode(&status)
	if err != nil {
		return err
	}

	if status.NetworkID != packet.NetworkID {
		return errors.New("network IDs mismatch")
	}

	log.Info().Interface("status", status).Msg("New peer")

	return nil
}

func handleMessage(rw ethp2p.MsgReadWriter) error {
	msg, err := rw.ReadMsg()
	if err != nil {
		return err
	}
	defer msg.Discard()

	switch msg.Code {
	case eth.TransactionsMsg:
		var txs eth.TransactionsPacket
		err = msg.Decode(&txs)
		log.Info().Interface("txs", txs).Err(err).Send()
	case eth.BlockHeadersMsg:
		var request eth.GetBlockHeadersPacket66
		err = msg.Decode(&request)
		log.Info().Interface("request", request).Err(err).Send()
	case eth.NewBlockMsg:
		var block eth.NewBlockPacket
		err = msg.Decode(&block)
		log.Info().Interface("block", block.Block.Number()).Err(err).Send()
	case eth.NewPooledTransactionHashesMsg:
		var txs eth.NewPooledTransactionHashesPacket
		err = msg.Decode(&txs)
		log.Info().Interface("txs", txs).Err(err).Send()
	default:
		log.Info().Interface("msg", msg).Send()
	}

	return nil
}
