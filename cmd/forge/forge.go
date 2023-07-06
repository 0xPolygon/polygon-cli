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
package forge

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	_ "embed"

	edgeblockchain "github.com/0xPolygon/polygon-edge/blockchain"
	edgechain "github.com/0xPolygon/polygon-edge/chain"
	edgeconsensus "github.com/0xPolygon/polygon-edge/consensus"
	edgedummy "github.com/0xPolygon/polygon-edge/consensus/dummy"
	edgepolybft "github.com/0xPolygon/polygon-edge/consensus/polybft"
	edgecontracts "github.com/0xPolygon/polygon-edge/contracts"
	edgecrypto "github.com/0xPolygon/polygon-edge/crypto"
	edgeserver "github.com/0xPolygon/polygon-edge/server"
	edgestate "github.com/0xPolygon/polygon-edge/state"
	edgeitrie "github.com/0xPolygon/polygon-edge/state/immutable-trie"
	edgeallowlist "github.com/0xPolygon/polygon-edge/state/runtime/allowlist"
	edgetxpool "github.com/0xPolygon/polygon-edge/txpool"
	edgetypes "github.com/0xPolygon/polygon-edge/types"
	edgebuildroot "github.com/0xPolygon/polygon-edge/types/buildroot"
	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/hashicorp/go-hclog"
	"github.com/maticnetwork/polygon-cli/rpctypes"
	"golang.org/x/exp/slices"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type (
	forgeParams struct {
		Client                string
		DataDir               string
		GenesisFile           string
		Verifier              string
		Mode                  string
		Count                 uint64
		BlocksFile            string
		BaseBlockReward       string
		ReceiptsFile          string
		IncludeTxFees         bool
		ShouldReadFirstBlock  bool
		ShouldVerifyBlocks    bool
		ShouldRewriteTxNonces bool
		HasConsecutiveBlocks  bool
		ShouldProcessBlocks   bool

		GenesisData []byte
	}
)

var (
	//go:embed usage.md
	usage        string
	inputForge   forgeParams
	BlockReadEOF = errors.New("no more blocks to read")
)

// forgeCmd represents the forge command
var ForgeCmd = &cobra.Command{
	Use:   "forge",
	Short: "Forge dumped blocks on top of a genesis file.",
	Long:  usage,

	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("forge called")
		blockchain, err := NewEdgeBlockchain()
		if err != nil {
			return err
		}

		blockReader, err := OpenBlockReader(inputForge.BlocksFile, inputForge.Mode)
		if err != nil {
			return err
		}

		receiptReader, err := OpenReceiptReader(inputForge.ReceiptsFile, inputForge.Mode)
		if inputForge.IncludeTxFees && err != nil {
			return err
		}

		err = readAllBlocksToChain(blockchain, blockReader, receiptReader)

		return err
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if inputForge.Client != "edge" {
			return fmt.Errorf("the client %s is not supported. Only Edge is supported", inputForge.Client)
		}
		if !slices.Contains([]string{"json", "proto"}, inputForge.Mode) {
			return fmt.Errorf("output format must one of [json, proto]")
		}
		f, err := os.Open(inputForge.GenesisFile)
		if err != nil {
			return fmt.Errorf("unable to open genesis file: %w", err)
		}
		genesisData, err := io.ReadAll(f)
		if err != nil {
			return fmt.Errorf("unable to read genesis file data: %w", err)
		}
		inputForge.GenesisData = genesisData

		return nil
	},
}

func init() {
	ForgeCmd.PersistentFlags().StringVarP(&inputForge.Client, "client", "c", "edge", "Specify which blockchain client should be use to forge the data")
	ForgeCmd.PersistentFlags().StringVarP(&inputForge.DataDir, "data-dir", "d", "./forged-data", "Specify a folder to be used to store the chain data")
	ForgeCmd.PersistentFlags().StringVarP(&inputForge.GenesisFile, "genesis", "g", "genesis.json", "Specify a file to be used for genesis configuration")
	ForgeCmd.PersistentFlags().StringVarP(&inputForge.Verifier, "verifier", "V", "dummy", "Specify a consensus engine to use for forging")
	ForgeCmd.PersistentFlags().StringVarP(&inputForge.Mode, "mode", "m", "json", "The forge mode indicates how we should get the transactions for our blocks [json, proto]")
	ForgeCmd.PersistentFlags().Uint64VarP(&inputForge.Count, "count", "C", 100, "The number of blocks to try to forge")
	ForgeCmd.PersistentFlags().StringVarP(&inputForge.BlocksFile, "blocks", "b", "", "A file of encoded blocks; the format of this file should match the mode")
	ForgeCmd.PersistentFlags().StringVarP(&inputForge.BaseBlockReward, "base-block-reward", "B", "2_000_000_000_000_000_000", "The amount rewarded for mining blocks")
	ForgeCmd.PersistentFlags().StringVarP(&inputForge.ReceiptsFile, "receipts", "r", "", "A file of encoded receipts; the format of this file should match the mode")
	ForgeCmd.PersistentFlags().BoolVarP(&inputForge.IncludeTxFees, "tx-fees", "t", false, "if the transaction fees should be included when computing block rewards")
	ForgeCmd.PersistentFlags().BoolVarP(&inputForge.ShouldReadFirstBlock, "read-first-block", "R", false, "whether to read the first block, leave false if first block is genesis")
	ForgeCmd.PersistentFlags().BoolVar(&inputForge.ShouldVerifyBlocks, "verify-blocks", true, "whether to verify blocks, set false if forging nonconsecutive blocks")
	ForgeCmd.PersistentFlags().BoolVar(&inputForge.ShouldRewriteTxNonces, "rewrite-tx-nonces", false, "whether to rewrite transaction nonces, set true if forging nonconsecutive blocks")
	ForgeCmd.PersistentFlags().BoolVar(&inputForge.HasConsecutiveBlocks, "consecutive-blocks", true, "whether the blocks file has consecutive blocks")
	ForgeCmd.PersistentFlags().BoolVarP(&inputForge.ShouldProcessBlocks, "process-blocks", "p", true, "whether the transactions in blocks should be processed applied to the state")

	if err := cobra.MarkFlagRequired(ForgeCmd.PersistentFlags(), "blocks"); err != nil {
		log.Error().Err(err).Msg("Unable to mark blocks flag as required")
	}
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
	if err := json.Unmarshal(inputForge.GenesisData, &chainConfig); err != nil {
		return nil, fmt.Errorf("unable to parse genesis data: %w", err)
	}
	logger := hclog.Default()

	stateStorage, err := edgeitrie.NewLevelDBStorage(filepath.Join(inputForge.DataDir, "trie"), logger)
	if err != nil {
		return nil, fmt.Errorf("unable to open leveldb storage: %w", err)
	}
	state := edgeitrie.NewState(stateStorage)

	executor := edgestate.NewExecutor(chainConfig.Params, state, logger)

	// custom write genesis hook per consensus engine
	genesisCreationFactory := map[edgeserver.ConsensusType]edgeserver.GenesisFactoryHook{
		edgeserver.PolyBFTConsensus: edgepolybft.GenesisPostHookFactory,
	}

	engineName := chainConfig.Params.GetEngine()
	if factory, exists := genesisCreationFactory[edgeserver.ConsensusType(engineName)]; exists {
		executor.GenesisPostHook = factory(&chainConfig, engineName)
	}

	// apply allow list genesis data
	if chainConfig.Params.ContractDeployerAllowList != nil {
		edgeallowlist.ApplyGenesisAllocs(chainConfig.Genesis, edgecontracts.AllowListContractsAddr,
			chainConfig.Params.ContractDeployerAllowList)
	}

	initialStateRoot := edgetypes.ZeroHash

	if edgeserver.ConsensusType(engineName) == edgeserver.PolyBFTConsensus {
		polyBFTConfig, configErr := edgepolybft.GetPolyBFTConfig(&chainConfig)
		if configErr != nil {
			return nil, configErr
		}

		if polyBFTConfig.InitialTrieRoot != edgetypes.ZeroHash {
			checkedInitialTrieRoot, hashErr := edgeitrie.HashChecker(polyBFTConfig.InitialTrieRoot.Bytes(), stateStorage)
			if hashErr != nil {
				return nil, fmt.Errorf("error on state root verification %w", hashErr)
			}

			if checkedInitialTrieRoot != polyBFTConfig.InitialTrieRoot {
				return nil, errors.New("invalid initial state root")
			}

			logger.Info("Initial state root checked and correct")

			initialStateRoot = polyBFTConfig.InitialTrieRoot
		}
	}

	genesisRoot, err := executor.WriteGenesis(chainConfig.Genesis.Alloc, initialStateRoot)
	if err != nil {
		return nil, err
	}

	chainConfig.Genesis.StateRoot = genesisRoot
	signer := edgecrypto.NewEIP155Signer(edgechain.AllForksEnabled.At(0), uint64(chainConfig.Params.ChainID))
	bc, err := edgeblockchain.NewBlockchain(logger, inputForge.DataDir, &chainConfig, nil, executor, signer)
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
		&edgetxpool.Config{
			MaxSlots:            1000,
			PriceLimit:          1000,
			MaxAccountEnqueued:  1000,
			DeploymentWhitelist: nil,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create new edge tx pool: %w", err)
	}
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

func readAllBlocksToChain(bh *edgeBlockchainHandle, blockReader BlockReader, receiptReader ReceiptReader) error {
	bc := bh.Blockchain
	blocksToRead := inputForge.Count
	genesisBlock, _ := bc.GetBlockByHash(bc.Genesis(), true)

	inputForge.BaseBlockReward = strings.ReplaceAll(strings.TrimSpace(inputForge.BaseBlockReward), "_", "")
	baseBlockReward := new(big.Int)
	base := 10
	if strings.HasPrefix(inputForge.BaseBlockReward, "0x") {
		base = 16
	}
	baseBlockReward.SetString(inputForge.BaseBlockReward, base)

	parentBlock := genesisBlock

	// Usually, the first block is the genesis block so it will skip it based on the flag.
	var i uint64 = 0
	if !inputForge.ShouldReadFirstBlock {
		_, err := blockReader.ReadBlock()
		if err != nil {
			return fmt.Errorf("could not read off the genesis block from input: %w", err)
		}
		i++
	}

	// in practice, I ran into some issues where the dumps that I created had duplicate blocks, This map is used to
	// detect and skip any kind of duplicates
	blockHashSet := make(map[ethcommon.Hash]struct{}, 0)

	// insertion into the chain will fail if blocks are numbered non-sequentially. This is used to throw an error if we
	// encounter blocks out of order. In the future, we should have a flag if we want to use original numbering or if
	// we want to create new numbering
	var lastNumber uint64 = 0
	var receipt *rpctypes.RawTxReceipt
	for ; i < blocksToRead; i++ {
		// read a polyblock which is a generic interface that can be marshalled into different formats
		block, err := blockReader.ReadBlock()
		if err != nil {
			return fmt.Errorf("could not read block %d due to error: %w", i, err)
		}

		if _, hasKey := blockHashSet[block.Hash()]; hasKey {
			log.Trace().Str("blockhash", block.Hash().String()).Msg("Skipping duplicate block")
			continue
		}
		blockHashSet[block.Hash()] = struct{}{}

		// There are instances where we can import nonconsecutive blocks, skip this
		// error on those instances.
		if inputForge.HasConsecutiveBlocks && block.Number().Uint64()-1 != lastNumber {
			return fmt.Errorf("encountered non consecutive block numbers on input. Got %s and expected %d", block.Number().String(), lastNumber+1)
		}
		lastNumber = block.Number().Uint64()

		// convert the generic rpc block into a block for edge. I suppose we'll need to think about other blockchain
		// forging at some point, but for now edge & supernets seem to be the real use case
		edgeBlock := PolyBlockToEdge(block)

		// The transactions nonces need to be rewritten or else there will be an error.
		if inputForge.ShouldRewriteTxNonces {
			for nonce, tx := range edgeBlock.Transactions {
				tx.Nonce = uint64(nonce)
				log.Debug().Int64("old nonce", int64(tx.Nonce)).Int64("new nonce", int64(nonce)).Str("tx hash", tx.Hash.String()).Msg("Rewrote tx nonce")
			}
		}

		// The parent hash value will not make sense, so we'll overwrite this when the value from our local parent block.
		edgeBlock.Header.ParentHash = parentBlock.Header.ComputeHash().Hash

		// The Transactions Root should be the same (i think?), but we'll set it
		edgeBlock.Header.TxRoot = edgebuildroot.CalculateTransactionsRoot(edgeBlock.Transactions)

		blockCreator, err := bh.Blockchain.GetConsensus().GetBlockCreator(edgeBlock.Header)
		if err != nil {
			return err
		}

		var txn *edgestate.Transition
		if inputForge.ShouldProcessBlocks {
			// This will execute the block and apply the transaction to the state.
			txn, err = bh.Executor.ProcessBlock(parentBlock.Header.StateRoot, edgeBlock, blockCreator)
		} else {
			txn, err = bh.Executor.BeginTxn(parentBlock.Header.StateRoot, edgeBlock.Header, blockCreator)
		}
		if err != nil {
			return fmt.Errorf("unable to process block %d with hash %s at index %d: %w", block.Number().Int64(), edgeBlock.Hash().String(), i, err)
		}

		if err = bh.Blockchain.GetConsensus().PreCommitState(edgeBlock.Header, txn); err != nil {
			return fmt.Errorf("could not pre commit state: %w", err)
		}

		// many of the headers are going to be different, so we'll get all of the headers and recompute the hash
		_, newRoot := txn.Commit()
		edgeBlock.Header.GasUsed = txn.TotalGas()
		edgeBlock.Header.ReceiptsRoot = edgebuildroot.CalculateReceiptsRoot(txn.Receipts())
		edgeBlock.Header.StateRoot = newRoot
		edgeBlock.Header.Hash = edgeBlock.Header.ComputeHash().Hash

		// This is an optional step but helpful to catch some mistakes in implementation.
		if inputForge.ShouldVerifyBlocks {
			_, err = bc.VerifyFinalizedBlock(edgeBlock)
			if err != nil {
				return fmt.Errorf("unable to verify finalized block: %w", err)
			}
		}

		// This might be worth putting behind a flag at some point, but we need some way to distribute native token
		// from mining. This is a hacky way to do it and right now.
		minerBalance := txn.GetBalance(edgetypes.BytesToAddress(edgeBlock.Header.Miner))
		minerTips := big.NewInt(0)
		burnedFee := big.NewInt(0)

		if inputForge.IncludeTxFees {
			totalGasUsed := big.NewInt(0)

			for i := 0; i < len(block.Transactions()); i++ {
				receipt, err = receiptReader.ReadReceipt()
				if err != nil {
					return fmt.Errorf("unable to read receipt to compute transaction fees: %w", err)
				}

				if receipt.BlockNumber.ToBigInt().Cmp(block.Number()) == -1 {
					// There are some receipts that exists which are not in the block
					// transactions. Skip the receipts where receiptBlockNumber is less
					// than blockNumber.
					i -= 1
					continue
				}

				if receipt.BlockNumber.ToBigInt().Cmp(block.Number()) != 0 {
					return fmt.Errorf("receipt block number mismatch, block numbers: %v, %v ", receipt.BlockNumber.ToBigInt(), block.Number())
				}

				totalGasUsed.Add(totalGasUsed, receipt.GasUsed.ToBigInt())
				totalFee := big.NewInt(0).Mul(totalGasUsed, receipt.EffectiveGasPrice.ToBigInt())
				minerTips.Add(minerTips, totalFee)
			}

			burnedFee = burnedFee.Mul(big.NewInt(int64(block.GasUsed())), block.BaseFee())
		}

		blockReward := big.NewInt(0).Add(baseBlockReward, big.NewInt(0).Sub(minerTips, burnedFee))
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

// PolyBlockToEdge will take the generic PolyBlock interface and convert it into an Edge compatible block.
func PolyBlockToEdge(polyBlock rpctypes.PolyBlock) *edgetypes.Block {
	h := new(edgetypes.Header)
	h.ParentHash = edgetypes.Hash(polyBlock.ParentHash())
	h.Sha3Uncles = edgetypes.Hash(polyBlock.UncleHash())
	h.Miner = polyBlock.Miner().Bytes()
	h.StateRoot = edgetypes.Hash(polyBlock.Root())
	h.TxRoot = edgetypes.Hash(polyBlock.TxHash())
	h.ReceiptsRoot = edgetypes.Hash(polyBlock.ReceiptsRoot())
	lb := polyBlock.LogsBloom()
	l := edgetypes.Bloom{}
	copy(l[:], lb)
	h.LogsBloom = l
	h.Difficulty = polyBlock.Difficulty().Uint64()
	h.Number = polyBlock.Number().Uint64()
	h.GasLimit = polyBlock.GasLimit()
	h.GasUsed = polyBlock.GasUsed()
	h.Timestamp = polyBlock.Time()
	h.ExtraData = polyBlock.Extra()
	h.MixHash = edgetypes.Hash{}
	var nonce [8]byte
	binary.LittleEndian.PutUint64(nonce[:], polyBlock.Nonce())
	h.Nonce = edgetypes.Nonce(nonce)
	h.Hash = edgetypes.Hash(polyBlock.Hash())

	txs := polyBlock.Transactions()
	etxs := make([]*edgetypes.Transaction, 0)
	for _, tx := range txs {
		etx := edgetypes.Transaction{}
		etx.Nonce = tx.Nonce()
		etx.GasPrice = tx.GasPrice()
		etx.Gas = tx.Gas()
		addr := edgetypes.Address(tx.To())

		if IsEmptyAddress(addr.Bytes()) {
			// The edge code that determines if a contract call is a contract creation
			// checks for a nil address rather than an address that's all zeros.
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

// IsEmptyAddress will just check a slice of bytes to check if it's all zeros or not.
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
