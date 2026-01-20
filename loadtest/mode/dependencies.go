package mode

import (
	"math/rand"

	"github.com/0xPolygon/polygon-cli/bindings/tester"
	"github.com/0xPolygon/polygon-cli/bindings/tokens"
	"github.com/0xPolygon/polygon-cli/bindings/uniswapv3"
	uniswap "github.com/0xPolygon/polygon-cli/loadtest/uniswapv3"
	"github.com/0xPolygon/polygon-cli/rpctypes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

// Dependencies holds shared resources needed by mode implementations.
type Dependencies struct {
	// Ethereum clients
	Client    *ethclient.Client
	RPCClient *ethrpc.Client

	// Contract instances
	LoadTesterContract *tester.LoadTester
	LoadTesterAddress  common.Address
	ERC20Contract      *tokens.ERC20
	ERC20Address       common.Address
	ERC721Contract     *tokens.ERC721
	ERC721Address      common.Address

	// Mode-specific data
	RecallTransactions []rpctypes.PolyTransaction
	IndexedActivity    *IndexedActivity

	// Random source for deterministic randomness
	RandSource *rand.Rand

	// UniswapV3 configuration (populated when uniswapv3 mode is used)
	UniswapV3Config     *uniswap.UniswapV3Config
	UniswapV3PoolConfig *uniswap.PoolConfig
	UniswapV3Pool       *uniswapv3.IUniswapV3Pool
}

// IndexedActivity holds indexed blockchain activity data for RPC testing.
type IndexedActivity struct {
	BlockNumbers    []string
	TransactionIDs  []string
	BlockIDs        []string
	Addresses       []string
	ERC20Addresses  []string
	ERC721Addresses []string
	Contracts       []string
	BlockNumber     uint64
	Transactions    []rpctypes.PolyTransaction
}
