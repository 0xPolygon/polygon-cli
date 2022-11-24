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
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	ethcore "github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethparams "github.com/ethereum/go-ethereum/params"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	inputForgeFile      *string
	inputForgeCache     *int
	inputForgeHandles   *int
	inputForgeAncient   *string
	inputForgeNamespace *string
	inputForgeEngine    *string
	inputForgeChainID   *int
)

type (
	fakeBlockReader struct {
		isInitialized bool
		scanner       *bufio.Scanner
	}
)

func (f *fakeBlockReader) Init() {
	f.scanner = bufio.NewScanner(os.Stdin)
	f.isInitialized = true
}

func (f *fakeBlockReader) GetFakeBlock() (*ethtypes.Header, []*ethtypes.Transaction, error) {
	if !f.scanner.Scan() {
		return nil, nil, fmt.Errorf("no more fake blocks left")
	}

	blockText := f.scanner.Bytes()
	var raw json.RawMessage
	err := json.Unmarshal(blockText, &raw)
	if err != nil {
		return nil, nil, err
	}
	// Decode header and transactions.
	var head *ethtypes.Header
	var body rpcBlock
	if err := json.Unmarshal(raw, &head); err != nil {
		return nil, nil, err
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, nil, err
	}

	txs := make([]*ethtypes.Transaction, len(body.Transactions))
	for i, tx := range body.Transactions {
		txs[i] = tx.tx
	}
	return head, txs, nil

}

// forgeCmd represents the forge command
var forgeCmd = &cobra.Command{
	Use:   "forge",
	Short: "Forge a fake state",
	Long: `
This command expects input to be a stream of json encoded block
objects. The goal is to take those blocks and write them into a
database in order to forge a long and complicated history for new
block chain clients and tests.

`,
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := rawdb.NewLevelDBDatabaseWithFreezer(*inputForgeFile, *inputForgeCache, *inputForgeHandles, *inputForgeAncient, *inputForgeNamespace, false)
		if err != nil {
			return err
		}

		fbr := new(fakeBlockReader)
		fbr.Init()
		// Idea here is that different engines could be used
		// to forge the blocks. We'll read the input and
		// decide which engin to instantiate
		switch *inputForgeEngine {
		case "ethhash":
			config := ethparams.ChainConfig{ChainID: big.NewInt(int64(*inputForgeChainID))}
			gen := &ethcore.Genesis{
				Config:     &config,
				Nonce:      0,
				Timestamp:  uint64(time.Unix(0, 0).Unix()),
				ParentHash: ethcommon.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
				ExtraData:  nil,
				GasLimit:   20_000_000,
				GasUsed:    0,
				Difficulty: big.NewInt(0),
				Mixhash:    ethcommon.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
				Coinbase:   ethcommon.HexToAddress("0x0000000000000000000000000000000000000000000000000000000000000000"),
				Alloc:      make(map[ethcommon.Address]ethcore.GenesisAccount, 0),
				BaseFee:    big.NewInt(0),
			}

			var jsonData []byte
			jsonData, err = gen.MarshalJSON()
			if err != nil {
				return err
			}

			err = os.WriteFile(*inputForgeFile+"/genesis.json", jsonData, 0777)
			if err != nil {
				return err
			}

			var gblock *ethtypes.Block
			gblock, err = gen.Commit(db)
			if err != nil {
				return err
			}

			engine := ethash.NewFaker()

			blocks, receipts := ethcore.GenerateChain(&config, gblock, engine, db, 100000, func(i int, gen *ethcore.BlockGen) {
				genBlock(i, gen, fbr)
			})
			_ = blocks
			_ = receipts

		}

		return err
	},
}

func genBlock(blockNum int, bg *ethcore.BlockGen, fbr *fakeBlockReader) {
	log.Trace().Int("blocknumber", blockNum).Msg("Forging")
	_, txs, err := fbr.GetFakeBlock()

	if err != nil {
		log.Error().Err(err).Msg("Unable to get a fake block")
		return
	}

	for _, v := range txs {
		if v == nil {
			continue
		}
		bg.AddTx(v)
	}
}

/*
// This function is unused.
func readBlocks(bc *ethcore.BlockChain, gblock *ethtypes.Block) error {
	scanner := bufio.NewScanner(os.Stdin)
	parentHash := gblock.Hash()
	parentNumber := gblock.Number()
	for scanner.Scan() {
		blockText := scanner.Bytes()
		var raw json.RawMessage
		err := json.Unmarshal(blockText, &raw)
		if err != nil {
			fmt.Println(err.Error())
		}
		b, err := readBlock(raw, parentHash, parentNumber)
		if err != nil {
			fmt.Println(err.Error())
		}

		parentHash = b.Hash()
		parentNumber = b.Number()
		fmt.Printf("Forging block: %d\n", b.Number())
		fmt.Println(b)

		_, err = bc.InsertChain(ethtypes.Blocks{b})
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}
	return nil
}
*/

type rpcBlock struct {
	Hash         ethcommon.Hash   `json:"hash"`
	Transactions []rpcTransaction `json:"transactions"`
	UncleHashes  []ethcommon.Hash `json:"uncles"`
}
type rpcTransaction struct {
	tx *ethtypes.Transaction
	txExtraInfo
}
type txExtraInfo struct {
	BlockNumber *string            `json:"blockNumber,omitempty"`
	BlockHash   *ethcommon.Hash    `json:"blockHash,omitempty"`
	From        *ethcommon.Address `json:"from,omitempty"`
}

/*
// This function is unused
// https://github.com/ethereum/go-ethereum/blob/d901d85377c2c2f05f09f423c7d739c0feecd90a/ethclient/ethclient.go#L110
func readBlock(raw json.RawMessage, parentHash ethcommon.Hash, parentNumber *big.Int) (*ethtypes.Block, error) {
	// Decode header and transactions.
	var head *ethtypes.Header
	var body rpcBlock
	if err := json.Unmarshal(raw, &head); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, err
	}
	// Quick-verify transaction and uncle lists. This mostly helps with debugging the server.
	if head.UncleHash == ethtypes.EmptyUncleHash && len(body.UncleHashes) > 0 {
		return nil, fmt.Errorf("server returned non-empty uncle list but block header indicates no uncles")
	}
	if head.UncleHash != ethtypes.EmptyUncleHash && len(body.UncleHashes) == 0 {
		return nil, fmt.Errorf("server returned empty uncle list but block header indicates uncles")
	}
	if head.TxHash == ethtypes.EmptyRootHash && len(body.Transactions) > 0 {
		return nil, fmt.Errorf("server returned non-empty transaction list but block header indicates no transactions")
	}
	if head.TxHash != ethtypes.EmptyRootHash && len(body.Transactions) == 0 {
		return nil, fmt.Errorf("server returned empty transaction list but block header indicates transactions")
	}
	txs := make([]*ethtypes.Transaction, len(body.Transactions))
	for i, tx := range body.Transactions {
		txs[i] = tx.tx
	}
	head.ParentHash = parentHash
	parentNumber.Add(parentNumber, big.NewInt(1))
	head.Number = parentNumber
	head.Difficulty = big.NewInt(131072)
	head.GasLimit = 20000000
	head.MixDigest = ethcommon.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")

	head.Time = uint64(time.Now().Unix())
	return ethtypes.NewBlockWithHeader(head).WithBody(txs, nil), nil
}
*/

func init() {
	rootCmd.AddCommand(forgeCmd)
	inputForgeFile = forgeCmd.PersistentFlags().StringP("datadir", "f", "./chaindata", "")
	inputForgeCache = forgeCmd.PersistentFlags().IntP("cache", "c", 512, "")
	inputForgeHandles = forgeCmd.PersistentFlags().Int("handles", 5120, "")
	inputForgeAncient = forgeCmd.PersistentFlags().StringP("datadir.ancient", "a", "./chaindata/ancient", "")
	inputForgeNamespace = forgeCmd.PersistentFlags().StringP("namespace", "n", "eth/db/chaindata/", "")
	inputForgeChainID = forgeCmd.PersistentFlags().IntP("chain-id", "i", 1337, "The chain id to use when configuring the blockchain")
	inputForgeEngine = forgeCmd.PersistentFlags().StringP("engine", "e", "ethhash", "The engine to use for the blockchain")
}
