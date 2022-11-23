/*
Copyright © 2022 Polygon <engineering@polygon.technology>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	edgedummy "github.com/0xPolygon/polygon-edge/consensus/dummy"
	"github.com/hashicorp/go-hclog"
	"github.com/maticnetwork/polygon-cli/rpctypes"
	"io"
	"os"
	"path/filepath"

	edgeblockchain "github.com/0xPolygon/polygon-edge/blockchain"
	edgechain "github.com/0xPolygon/polygon-edge/chain"
	edgeconsensus "github.com/0xPolygon/polygon-edge/consensus"
	edgecrypto "github.com/0xPolygon/polygon-edge/crypto"
	edgestate "github.com/0xPolygon/polygon-edge/state"
	edgeitrie "github.com/0xPolygon/polygon-edge/state/immutable-trie"
	edgeevm "github.com/0xPolygon/polygon-edge/state/runtime/evm"
	edgeprecompiled "github.com/0xPolygon/polygon-edge/state/runtime/precompiled"
	edgetxpool "github.com/0xPolygon/polygon-edge/txpool"
	edgetypes "github.com/0xPolygon/polygon-edge/types"

	"github.com/spf13/cobra"
)

type (
	forgeParams struct {
		Client         *string
		DataDir        *string
		GenesisFile    *string
		Verifier       *string
		Mode           *string
		Count          *uint64
		JSONBlocksFile *string

		GenesisData []byte
	}
	BlockReader struct {
		scanner *bufio.Scanner
	}
)

var (
	inputForgeParams forgeParams
	BlockReadEOF     = errors.New("no more blocks to read")
)

// forgeCmd represents the forge command
var forgeCmd = &cobra.Command{
	Use:   "forge",
	Short: "A utility for generating blockchain data either for testing or migration",
	Long: `
go run main.go dumpblocks http://172.26.26.12:8545/ 0 50000 > eth.50k
cat eth.50k | grep '"difficulty"' > eth.50k.blocks
go run main.go forge --genesis ../polygon-edge/genesis.json --mode json --json-blocks eth.50k.blocks --count 50000`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("forge called")
		bc, err := NewEdgeBlockchain()
		if err != nil {
			return err
		}

		br, err := OpenJSONBlockReader(*inputForgeParams.JSONBlocksFile)
		if err != nil {
			return err
		}
		// in the future add a different type of reader potentially?

		err = readAllBlocksToChain(bc, br)

		return err
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if *inputForgeParams.Client != "edge" {
			return fmt.Errorf("the client %s is not supported. Only Edge is supported", *inputForgeParams.Client)
		}
		if *inputForgeParams.Mode != "json" {
			return fmt.Errorf("the mode %s is not suported yet. Only json is supported", *inputForgeParams.Mode)
		}
		f, err := os.Open(*inputForgeParams.GenesisFile)
		if err != nil {
			return fmt.Errorf("unable to open genesis file: %w", err)
		}
		genesisData, err := io.ReadAll(f)
		if err != nil {
			return fmt.Errorf("unable to read genesis file data: %w", err)
		}
		inputForgeParams.GenesisData = genesisData
		return nil
	},
}

type edgeTxpoolHub struct {
	state edgestate.State
	*edgeblockchain.Blockchain
}

type edgeBlockchainHandle struct {
	Blockchain   *edgeblockchain.Blockchain
	Executor     *edgestate.Executor
	StateStorage *edgeitrie.Storage
}

func NewEdgeBlockchain() (*edgeBlockchainHandle, error) {

	var chainConfig edgechain.Chain
	err := json.Unmarshal(inputForgeParams.GenesisData, &chainConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to parse genesis data: %w", err)
	}
	logger := hclog.Default()

	stateStorage, err := edgeitrie.NewLevelDBStorage(filepath.Join(*inputForgeParams.DataDir, "trie"), logger)
	if err != nil {
		return nil, fmt.Errorf("unable to open leveldb storage: %w")
	}
	state := edgeitrie.NewState(stateStorage)

	// consensus := edgedummy.Dummy{}
	executor := edgestate.NewExecutor(chainConfig.Params, state, logger)
	executor.SetRuntime(edgeprecompiled.NewPrecompiled())
	executor.SetRuntime(edgeevm.NewEVM())

	genesisRoot := executor.WriteGenesis(chainConfig.Genesis.Alloc)
	// should I override the state root here?
	chainConfig.Genesis.StateRoot = genesisRoot
	signer := edgecrypto.NewEIP155Signer(uint64(chainConfig.Params.ChainID))
	bc, err := edgeblockchain.NewBlockchain(logger, *inputForgeParams.DataDir, &chainConfig, nil, executor, signer)
	if err != nil {
		return nil, fmt.Errorf("unable to setup blockchain: %w", err)
	}
	executor.GetHash = bc.GetHashHelper
	//hub := &edgeTxpoolHub{
	//	state:      state,
	//	Blockchain: bc,
	//}
	txpool, err := edgetxpool.NewTxPool(
		logger,
		chainConfig.Params.Forks.At(0),
		nil,
		nil,
		nil,
		nil,
		&edgetxpool.Config{
			MaxSlots:            1000,
			PriceLimit:          1000,
			MaxAccountEnqueued:  1000,
			DeploymentWhitelist: nil,
		},
	)
	txpool.SetSigner(signer)
	dummyConsensus, err := edgedummy.Factory(&edgeconsensus.Params{
		TxPool:     txpool,
		Blockchain: bc,
		Executor:   executor,
		Logger:     logger,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create dummy consensus: %w", err)
	}
	bc.SetConsensus(dummyConsensus)
	if err := bc.ComputeGenesis(); err != nil {
		return nil, err
	}
	if err := dummyConsensus.Initialize(); err != nil {
		return nil, err
	}
	bh := &edgeBlockchainHandle{
		Blockchain:   bc,
		Executor:     executor,
		StateStorage: &stateStorage,
	}
	return bh, nil
}

func OpenJSONBlockReader(file string) (*BlockReader, error) {
	blockFile, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("unable to open json block file: %w", err)
	}
	scanner := bufio.NewScanner(blockFile)

	br := BlockReader{
		scanner: scanner,
	}

	return &br, nil
}

func (b *BlockReader) ReadBlock() (rpctypes.PolyBlock, error) {
	if !b.scanner.Scan() {
		return nil, BlockReadEOF
	}
	rawBlockBytes := b.scanner.Bytes()
	var raw rpctypes.RawBlockResponse
	err := json.Unmarshal(rawBlockBytes, &raw)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal block: %w - %s", err, string(rawBlockBytes))
	}
	return rpctypes.NewPolyBlock(&raw), nil
}
func readAllBlocksToChain(bh *edgeBlockchainHandle, br *BlockReader) error {
	bc := bh.Blockchain
	blocksToRead := *inputForgeParams.Count
	parentHash := bc.Genesis()
	fmt.Println(bc.Header().StateRoot)
	var i uint64
	for i = 0; i < blocksToRead; i = i + 1 {
		b, err := br.ReadBlock()
		if err != nil {
			return fmt.Errorf("could not read block %d due to error: %w", i, err)
		}
		edgeBlock := PolyblockToEdge(b)
		edgeBlock.Header.Number = i + 1
		edgeBlock.Header.ParentHash = parentHash
		edgeBlock.Header.Hash = edgeBlock.Header.ComputeHash().Hash
		// edgeBlock.Header.StateRoot = bh.StateStorage.
		parentHash = edgeBlock.Header.Hash

		// fmt.Println(edgeBlock.Header().Number)
		// h := edgeBlock.Header()
		// fmt.Println(h.Number)
		err = bc.WriteBlock(edgeBlock, "polycli")
		if err != nil {
			return fmt.Errorf("unable to write block: %w", err)
		}
	}
	return nil
}
func PolyblockToEdge(pb rpctypes.PolyBlock) *edgetypes.Block {
	h := new(edgetypes.Header)
	h.ParentHash = edgetypes.Hash(pb.ParentHash())
	h.Sha3Uncles = edgetypes.Hash(pb.UncleHash())
	h.Miner = pb.Miner().Bytes()
	h.StateRoot = edgetypes.Hash(pb.Root())
	h.TxRoot = edgetypes.Hash(pb.TxHash())
	h.ReceiptsRoot = edgetypes.Hash(pb.ReceiptsRoot())
	lb := pb.LogsBloom()
	l := edgetypes.Bloom{}
	copy(l[:], lb)
	h.LogsBloom = l
	h.Difficulty = pb.Difficulty().Uint64()
	h.Number = pb.Number().Uint64()
	h.GasLimit = pb.GasLimit()
	h.GasUsed = pb.GasUsed()
	h.Timestamp = pb.Time()
	h.ExtraData = pb.Extra()
	h.MixHash = edgetypes.Hash{}
	var nonce [8]byte
	binary.LittleEndian.PutUint64(nonce[:], pb.Nonce())
	h.Nonce = edgetypes.Nonce(nonce)
	h.Hash = edgetypes.Hash(pb.Hash())

	txs := pb.Transactions()
	etxs := make([]*edgetypes.Transaction, 0)
	for _, tx := range txs {
		etx := edgetypes.Transaction{}
		etx.Nonce = tx.Nonce()
		etx.GasPrice = tx.GasPrice()
		etx.Gas = tx.Gas()
		addr := edgetypes.Address(tx.To())
		etx.To = &addr
		etx.Value = tx.Value()
		etx.Input = tx.Data()
		etx.V = tx.V()
		etx.R = tx.R()
		etx.S = tx.S()
		etx.Hash = edgetypes.Hash(tx.Hash())
		etx.From = edgetypes.Address(tx.From())
		etxs = append(etxs, &etx)
	}

	b := edgetypes.Block{
		Header:       h,
		Transactions: etxs,
		// At some point we might want to include uncles?
	}
	return &b
}

func GenerateRandomBlock(number uint64) *edgetypes.Block {
	return nil
}

func init() {
	rootCmd.AddCommand(forgeCmd)

	inputForgeParams.Client = forgeCmd.PersistentFlags().StringP("client", "c", "edge", "Specify which blockchain client should be use to forge the data")
	inputForgeParams.DataDir = forgeCmd.PersistentFlags().StringP("data-dir", "d", "./forged-data", "Specify a folder to be used to store the chain data")
	inputForgeParams.GenesisFile = forgeCmd.PersistentFlags().StringP("genesis", "g", "genesis.json", "Specify a file to be used for genesis configuration")
	inputForgeParams.Verifier = forgeCmd.PersistentFlags().StringP("verifier", "v", "dummy", "Specify a consensus engin to use for forging")
	inputForgeParams.Mode = forgeCmd.PersistentFlags().StringP("mode", "m", "json", "The forge mode indicates how we should get the transactions for our blocks")
	inputForgeParams.Count = forgeCmd.PersistentFlags().Uint64P("count", "C", 100, "The number of blocks to try to forge")
	inputForgeParams.JSONBlocksFile = forgeCmd.PersistentFlags().String("json-blocks", "", "a file of json encoded blocks")

}
