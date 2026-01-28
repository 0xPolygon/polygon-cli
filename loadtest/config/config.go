package config

//go:generate stringer -type=Mode

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/0xPolygon/polygon-cli/loadtest/uniswapv3"
	"github.com/ethereum/go-ethereum/common"
)

// Mode represents the type of load test to perform.
type Mode int

// Mode constants.
// If you add a new constant, run: make gen-load-test-modes
const (
	ModeERC20 Mode = iota
	ModeERC721
	ModeBlob
	ModeContractCall
	ModeDeploy
	ModeIncrement
	ModeRandom
	ModeRecall
	ModeRPC
	ModeStore
	ModeTransaction
	ModeUniswapV3
)

// Config holds all load test parameters.
type Config struct {
	// Network connection
	RPCURL  string
	ChainID uint64
	Proxy   string

	// Test parameters
	Requests    int64
	Concurrency int64
	BatchSize   uint64
	TimeLimit   int64
	Seed        int64

	// Transaction options
	PrivateKey         string
	ToAddress          string
	EthAmountInWei     uint64
	RandomRecipients   bool
	LegacyTxMode       bool
	FireAndForget      bool
	CheckForPreconf    bool
	PreconfStatsFile   string
	WaitForReceipt     bool
	ReceiptRetryMax    uint
	ReceiptRetryDelay  uint // initial delay in milliseconds
	OutputRawTxOnly    bool
	StartNonce         uint64
	GasPriceMultiplier float64

	// Gas options
	ForceGasLimit         uint64
	ForceGasPrice         uint64
	ForcePriorityGasPrice uint64
	MaxBaseFeeWei         uint64

	// Rate limiting
	RateLimit                  float64
	AdaptiveRateLimit          bool
	AdaptiveTargetSize         uint64
	AdaptiveRateLimitIncrement uint64
	AdaptiveCycleDuration      uint64
	AdaptiveBackoffFactor      float64

	// Mode configuration
	Modes []string

	// Call-only options
	EthCallOnly            bool
	EthCallOnlyLatestBlock bool

	// Contract addresses
	LoadTestContractAddress string
	ERC20Address            string
	ERC721Address           string

	// Mode-specific options
	StoreDataSize       uint64
	RecallLength        uint64
	BlockBatchSize      uint64
	ContractAddress     string
	ContractCallData    string
	ContractCallPayable bool
	BlobFeeCap          uint64

	// Account pool options
	SendingAccountsCount      uint64
	AccountFundingAmount      *big.Int
	PreFundSendingAccounts    bool
	RefundRemainingFunds      bool
	SendingAccountsFile       string
	CheckBalanceBeforeFunding bool
	DumpSendingAccountsFile   string
	AccountsPerFundingTx      uint64

	// Summary output
	ShouldProduceSummary bool
	SummaryOutputMode    string

	// UniswapV3-specific config (set by uniswapv3 subcommand)
	UniswapV3 *UniswapV3Config

	// Gas manager config (optional, for gas oscillation features)
	GasManager *GasManagerConfig

	// Computed fields (populated during initialization)
	CurrentGasPrice       *big.Int
	CurrentGasTipCap      *big.Int
	CurrentNonce          *uint64
	ECDSAPrivateKey       *ecdsa.PrivateKey
	FromETHAddress        *common.Address
	ToETHAddress          *common.Address
	ContractETHAddress    *common.Address
	SendAmount            *big.Int
	ChainSupportBaseFee   bool
	ParsedModes           []Mode
	MultiMode             bool
	BigGasPriceMultiplier *big.Float
}

// GasManagerConfig holds gas manager configuration for oscillation waves and pricing strategies.
type GasManagerConfig struct {
	// Oscillation wave options
	OscillationWave string // flat, sine, square, triangle, sawtooth
	Target          uint64 // target gas limit baseline
	Period          uint64 // period in blocks
	Amplitude       uint64 // amplitude of oscillation

	// Pricing strategy options
	PriceStrategy             string  // estimated, fixed, dynamic
	FixedGasPriceWei          uint64  // for fixed strategy
	DynamicGasPricesWei       string  // comma-separated prices for dynamic strategy
	DynamicGasPricesVariation float64 // ±percentage variation for dynamic
}

// UniswapV3Config holds UniswapV3-specific configuration.
type UniswapV3Config struct {
	// Pre-deployed contract addresses (as hex strings).
	FactoryV3                          string
	Multicall                          string
	ProxyAdmin                         string
	TickLens                           string
	NFTDescriptorLib                   string
	NonfungibleTokenPositionDescriptor string
	TransparentUpgradeableProxy        string
	NonfungiblePositionManager         string
	Migrator                           string
	Staker                             string
	QuoterV2                           string
	SwapRouter                         string
	WETH9                              string
	PoolToken0                         string
	PoolToken1                         string

	// Pool and swap parameters.
	PoolFees        float64
	SwapAmountInput uint64
}

// Validate validates the Config and returns an error if any validation fails.
func (c *Config) Validate() error {
	if c.AdaptiveBackoffFactor <= 0.0 {
		return fmt.Errorf("the backoff factor needs to be non-zero positive. Given: %f", c.AdaptiveBackoffFactor)
	}

	if c.WaitForReceipt && c.ReceiptRetryMax <= 1 {
		return errors.New("when waiting for a receipt, use a max retry greater than 1")
	}

	if c.EthCallOnly {
		if c.PreFundSendingAccounts || c.SendingAccountsFile != "" || c.SendingAccountsCount > 0 {
			return errors.New("pre-funding accounts with call only mode doesn't make sense")
		}
		if c.WaitForReceipt {
			return errors.New("waiting for receipts doesn't make sense with call only mode")
		}
	}

	if c.GasPriceMultiplier == 0 {
		return errors.New("gas price multiplier should be non-zero")
	}

	return nil
}

// Validate validates the UniswapV3Config and returns an error if any validation fails.
func (c *UniswapV3Config) Validate() error {
	switch fees := c.PoolFees; fees {
	case float64(uniswapv3.StableTier), float64(uniswapv3.StandardTier), float64(uniswapv3.ExoticTier):
		// Valid fee tier.
	default:
		return fmt.Errorf("UniswapV3 only supports a few pool tiers which are stable: %f%%, standard: %f%%, and exotic: %f%%",
			float64(uniswapv3.StableTier), float64(uniswapv3.StandardTier), float64(uniswapv3.ExoticTier))
	}

	if c.SwapAmountInput == 0 {
		return errors.New("swap amount input has to be greater than zero")
	}

	if (c.PoolToken0 != "") != (c.PoolToken1 != "") {
		return errors.New("both pool tokens must be empty or specified. Specifying only one token is not allowed")
	}

	return nil
}

// Validate validates the GasManagerConfig and returns an error if any validation fails.
func (c *GasManagerConfig) Validate() error {
	switch c.OscillationWave {
	case "flat", "sine", "square", "triangle", "sawtooth":
		// Valid wave type.
	default:
		return fmt.Errorf("invalid oscillation wave type: %s", c.OscillationWave)
	}

	switch c.PriceStrategy {
	case "estimated", "fixed", "dynamic":
		// Valid strategy.
	default:
		return fmt.Errorf("invalid price strategy: %s", c.PriceStrategy)
	}

	return nil
}
