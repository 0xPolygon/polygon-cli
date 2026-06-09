package config

//go:generate stringer -type=Mode

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"

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
	RPCURL     string
	ChainID    uint64
	Proxy      string
	RPCHeaders string
	Headers    map[string]string

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
	PrivateTxs         bool
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
	StoreDataSize        uint64
	RecallLength         uint64
	BlockBatchSize       uint64
	ContractAddress      string
	ContractCallData     string
	ContractCallDataFile string
	ContractCallPayable  bool
	BlobFeeCap           uint64

	// Account pool options
	SendingAccountsCount      uint64
	AccountFundingAmount      *big.Int
	PreFundSendingAccounts    bool
	RefundRemainingFunds      bool
	SendingAccountsFile       string
	CheckBalanceBeforeFunding bool
	DumpSendingAccountsFile   string
	AccountsPerFundingTx      uint64
	SequentialNonceFetch      bool
	StopOnInsufficientFunds   bool

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
	Enabled bool

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

	if c.PrivateTxs {
		if err := c.validatePrivateTxsModes(); err != nil {
			return err
		}
	}

	if err := c.resolveContractCallData(); err != nil {
		return err
	}

	return nil
}

// resolveContractCallData loads --calldata-file into ContractCallData when set.
// This lets contract-call payloads exceed the OS single-argument limit
// (MAX_ARG_STRLEN, 128 KiB on Linux), which caps the inline --calldata flag at
// roughly 64 KB of calldata. The file holds the same hex encoding as --calldata
// (function signature + encoded arguments, optional 0x prefix). All whitespace
// is stripped, so line-wrapped hex dumps (e.g. from `xxd`/`od`) are accepted as
// readily as a single line. The two flags are mutually exclusive. The hex is
// validated here so an unusable payload fails fast, before any transactions are
// sent.
func (c *Config) resolveContractCallData() error {
	if c.ContractCallDataFile == "" {
		return nil
	}
	if c.ContractCallData != "" {
		return errors.New("--calldata and --calldata-file are mutually exclusive; specify only one")
	}

	data, err := os.ReadFile(c.ContractCallDataFile)
	if err != nil {
		return fmt.Errorf("unable to read --calldata-file %q: %w", c.ContractCallDataFile, err)
	}

	// Strip all whitespace (incl. interior newlines from wrapped hex dumps),
	// not just the surrounding bytes. strings.Fields splits on any unicode
	// whitespace run; joining yields a single contiguous hex string.
	calldata := strings.Join(strings.Fields(string(data)), "")
	// Empty-check the hex body AFTER removing a leading 0x, so a file holding
	// only "0x" (which decodes to zero bytes without error) is rejected too.
	hexBody := strings.TrimPrefix(calldata, "0x")
	if hexBody == "" {
		return fmt.Errorf("--calldata-file %q is empty", c.ContractCallDataFile)
	}
	if _, err := hex.DecodeString(hexBody); err != nil {
		return fmt.Errorf("--calldata-file %q does not contain valid hex calldata: %w", c.ContractCallDataFile, err)
	}

	c.ContractCallData = calldata
	return nil
}

// validatePrivateTxsModes checks that all specified modes support --private-txs.
func (c *Config) validatePrivateTxsModes() error {
	supported := map[string]bool{
		"t": true, "transaction": true,
		"b": true, "blob": true,
		"cc": true, "contract-call": true,
		"R": true, "recall": true,
	}

	for _, mode := range c.Modes {
		if !supported[mode] {
			return fmt.Errorf("--private-txs is not supported for mode %q; supported modes: transaction, blob, contract-call, recall", mode)
		}
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

// ParseRPCHeaders parses a comma-separated string of key:value pairs into a map.
// Format: "key1:value1,key2:value2"
// Values may contain colons (e.g., "Authorization:Bearer token").
func ParseRPCHeaders(s string) (map[string]string, error) {
	if s == "" {
		return nil, nil
	}

	headers := make(map[string]string)
	for pair := range strings.SplitSeq(s, ",") {
		parts := strings.SplitN(pair, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid header format %q, expected key:value", pair)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" {
			return nil, fmt.Errorf("empty header key in %q", pair)
		}

		headers[key] = value
	}

	return headers, nil
}
