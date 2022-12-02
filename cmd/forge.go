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
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-hclog"
	"github.com/maticnetwork/polygon-cli/rpctypes"
	"io"
	"math/big"
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
	edgebuildroot "github.com/0xPolygon/polygon-edge/types/buildroot"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

		zerolog.SetGlobalLevel(zerolog.TraceLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Debug().Msg("Starting logger in console mode")

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
	State        *edgeitrie.State
	Consensus    *edgeconsensus.Consensus
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

	executor := edgestate.NewExecutor(chainConfig.Params, state, logger)
	executor.SetRuntime(edgeprecompiled.NewPrecompiled())
	executor.SetRuntime(edgeevm.NewEVM())

	genesisRoot := executor.WriteGenesis(chainConfig.Genesis.Alloc)

	chainConfig.Genesis.StateRoot = genesisRoot
	signer := edgecrypto.NewEIP155Signer(uint64(chainConfig.Params.ChainID))
	bc, err := edgeblockchain.NewBlockchain(logger, *inputForgeParams.DataDir, &chainConfig, nil, executor, signer)
	if err != nil {
		return nil, fmt.Errorf("unable to setup blockchain: %w", err)
	}
	executor.GetHash = bc.GetHashHelper
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
	// eventually we should allow for different consensus. It would be better to use some private PoA consensus for all
	// of the forged blocks then switch to PoS or something like that at the last block
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
		State:        state,
		Consensus:    &dummyConsensus,
	}
	return bh, nil
}

func OpenJSONBlockReader(file string) (*BlockReader, error) {
	blockFile, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("unable to open json block file: %w", err)
	}
	maxCapacity := 5 * 1024 * 1024
	buf := make([]byte, maxCapacity)
	scanner := bufio.NewScanner(blockFile)
	scanner.Buffer(buf, maxCapacity)

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
	genesisBlock, _ := bc.GetBlockByHash(bc.Genesis(), true)

	// the block reward should probably be configurable depending on the needs
	blockReward := big.NewInt(2_000_000_000_000_000_000)
	parentBlock := genesisBlock

	// this should probably be based on a flag, but in our current use case, we're going to assume the 0th block is
	// a copy of the genesis block, so there's no point inserting it again
	_, err := br.ReadBlock()
	if err != nil {
		return fmt.Errorf("could not read off the genesis block from input: %w", err)
	}

	// in practice, I ran into some issues where the dumps that I created had duplicate blocks, This map is used to
	// detect and skip any kind of duplicates
	blockHashSet := make(map[ethcommon.Hash]struct{}, 0)

	// insertion into the chain will fail if blocks are numbered non-sequnetially. This is used to throw an error if we
	// encounter blocks out of order. In the future, we should have a flag if we want to use original numbering or if
	// we want to create new numbering
	var lastNumber uint64 = 0
	var i uint64
	for i = 1; i < blocksToRead; i = i + 1 {
		// read a polyblock which is a generic interface that can be marshalled into different formats
		b, err := br.ReadBlock()
		if err != nil {
			return fmt.Errorf("could not read block %d due to error: %w", i, err)
		}

		if _, hasKey := blockHashSet[b.Hash()]; hasKey {
			log.Trace().Str("blockhash", b.Hash().String()).Msg("skipping duplicate block")
			continue
		}
		blockHashSet[b.Hash()] = struct{}{}

		if b.Number().Uint64()-1 != lastNumber {
			return fmt.Errorf("encountered non consecutive block numbers on input. Got %s and expected %d", b.Number().String(), lastNumber+1)
		}
		lastNumber = b.Number().Uint64()

		// convert the generic rpc block into a block for edge. I suppose we'll need to think about other blockchain
		// forging at somet poing, but for now edge & supernets seem to be the real use case
		edgeBlock := PolyblockToEdge(b)

		// The parent hash value will not make sense, so we'll overwrite this when the value from our local parent block.
		edgeBlock.Header.ParentHash = parentBlock.Header.ComputeHash().Hash

		// The Transactions Root should be the same (i think?), but we'll set it
		edgeBlock.Header.TxRoot = edgebuildroot.CalculateTransactionsRoot(edgeBlock.Transactions)

		blockCreator, err := bh.Blockchain.GetConsensus().GetBlockCreator(edgeBlock.Header)
		if err != nil {
			return err
		}

		// This will execute the block and apply the transaction to the state
		txn, err := bh.Executor.ProcessBlock(parentBlock.Header.StateRoot, edgeBlock, blockCreator)
		if err != nil {
			return fmt.Errorf("unable to process block %d %s: %w", i, edgeBlock.Hash().String(), err)
		}

		if err := bh.Blockchain.GetConsensus().PreCommitState(edgeBlock.Header, txn); err != nil {
			return fmt.Errorf("could not pre commit state: %w", err)
		}

		// many of the headers are going to be different, so we'll get all of the headers and recompute the hash
		_, newRoot := txn.Commit()
		edgeBlock.Header.GasUsed = txn.TotalGas()
		edgeBlock.Header.ReceiptsRoot = edgebuildroot.CalculateReceiptsRoot(txn.Receipts())
		edgeBlock.Header.StateRoot = newRoot
		edgeBlock.Header.Hash = edgeBlock.Header.ComputeHash().Hash

		// This is an optional step but helpful to catch some mistakes in implementation
		err = bc.VerifyFinalizedBlock(edgeBlock)
		if err != nil {
			return fmt.Errorf("unable to verify finalized block: %w", err)
		}

		// This might be worth putting behind a flag at somet point, but we need some way to distribute native token
		// from mining. This is a hacky way to do it and right now, I'm not including transaction fees
		minerBalance := txn.GetBalance(edgetypes.BytesToAddress(edgeBlock.Header.Miner))
		minerBalance = minerBalance.Add(minerBalance, blockReward)
		txn.Txn().SetBalance(edgetypes.BytesToAddress(edgeBlock.Header.Miner), minerBalance)

		// after doing the irregular state change, i need to update the block headers again with the new root hash and
		// block hash
		_, newRoot = txn.Commit()
		edgeBlock.Header.StateRoot = newRoot
		edgeBlock.Header.Hash = edgeBlock.Header.ComputeHash().Hash

		// at this point the block should be OK to write to the local database?
		err = bc.WriteBlock(edgeBlock, "polycli")
		if err != nil {
			return fmt.Errorf("unable to write block: %w", err)
		}
		parentBlock = edgeBlock
	}
	return nil
}

// PolyblockToEdge will take the generic PolyBlock interface and convert it into a edge compatible block
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

		if IsEmptyAddress(addr.Bytes()) {
			// The edge code that determines if a contract call is a contract creation checks for a nil address rather
			// than an address that's all zeros
			etx.To = nil
		} else {
			etx.To = &addr
		}

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

// IsEmptyAddress will just check a slice of bytes to check if it's all zeros or not
func IsEmptyAddress(addr []byte) bool {
	for _, v := range addr {
		if v != 0 {
			return false
		}
	}
	return true
}

// GenerateRandomBlock in most cases we can use existing blocks and transactions for forgeries and testing, but at some
// point we might want to generate complete random blocks especially if we want to model state size after 10 - 20 years
// of operation
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
