package p2p

import (
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
)

const (
	// maxBroadcastTxSize mirrors geth's txMaxSize (4 slots of 32KB): the largest
	// transaction a node will admit to its pool, so anything larger is dead weight
	// no matter how well formed it is.
	maxBroadcastTxSize = 4 * 32 * 1024

	// baseFeeRatioScale is the fixed-point denominator used to apply the
	// fractional base fee floor with integer math.
	baseFeeRatioScale = 1000

	// maxBaseFeeRatio bounds --broadcast-min-basefee-ratio. Anything above a few
	// multiples of the base fee would drop every transaction on the network, so
	// the cap exists to catch a misconfigured flag rather than to be useful.
	maxBaseFeeRatio = 100.0
)

// errBelowMinGasPrice is returned when a transaction's fee cap is under the
// operator-configured floor. It is distinct from txpool.ErrTxGasPriceTooLow
// (which the pool raises for the tip floor) so the two show up as separate
// rejection reasons in metrics.
var errBelowMinGasPrice = errors.New("gas fee cap below configured minimum")

// TxValidatorOptions configures a TxValidator.
type TxValidatorOptions struct {
	// ChainID is the chain the sensor is attached to. It binds the signer, so a
	// transaction signed for a different chain fails sender recovery.
	ChainID uint64

	// MinGasPrice is a floor on the transaction's gas fee cap. Nil or zero
	// disables the check.
	MinGasPrice *big.Int

	// MinTip is a floor on the transaction's gas tip cap, equivalent to a
	// node's --txpool.pricelimit. Nil or zero disables the check.
	MinTip *big.Int

	// MinBaseFeeRatio rejects transactions whose gas fee cap is below this
	// fraction of the head block's base fee. 1.0 drops everything that cannot be
	// included at the current base fee; 0 disables the check.
	MinBaseFeeRatio float64
}

// TxValidator applies the stateless half of a node's transaction admission
// rules to transactions the sensor is about to rebroadcast.
//
// It deliberately does not check nonce or balance: the sensor holds no state, and
// those checks would need an RPC round trip per sender. What it does catch is
// everything a peer can make up for free -- forged signatures, transactions signed
// for another chain, oversized payloads, gas below the intrinsic cost, fee caps
// that can never be mined -- which is the bulk of what a spamming peer sends.
type TxValidator struct {
	signer types.Signer
	opts   *txpool.ValidationOptions

	// minGasPrice is nil when the floor is disabled.
	minGasPrice *big.Int
	// baseFeeRatio is the MinBaseFeeRatio scaled by baseFeeRatioScale, or nil
	// when the check is disabled.
	baseFeeRatio *big.Int
}

// NewTxValidator creates a transaction validator for the given chain.
func NewTxValidator(opts TxValidatorOptions) (*TxValidator, error) {
	if opts.ChainID == 0 {
		return nil, errors.New("chain ID must be non-zero")
	}
	if opts.MinGasPrice != nil && opts.MinGasPrice.Sign() < 0 {
		return nil, fmt.Errorf("minimum gas price cannot be negative: %v", opts.MinGasPrice)
	}
	if opts.MinTip != nil && opts.MinTip.Sign() < 0 {
		return nil, fmt.Errorf("minimum tip cannot be negative: %v", opts.MinTip)
	}
	if math.IsNaN(opts.MinBaseFeeRatio) || opts.MinBaseFeeRatio < 0 || opts.MinBaseFeeRatio > maxBaseFeeRatio {
		return nil, fmt.Errorf("base fee ratio must be between 0 and %v: %v", maxBaseFeeRatio, opts.MinBaseFeeRatio)
	}

	config := broadcastChainConfig(opts.ChainID)

	minTip := new(big.Int)
	if opts.MinTip != nil {
		minTip.Set(opts.MinTip)
	}

	v := &TxValidator{
		signer: types.LatestSigner(config),
		opts: &txpool.ValidationOptions{
			Config: config,
			// Blob transactions are excluded on purpose: they only travel with
			// their sidecars, which the sensor never receives, so forwarding one
			// would propagate a transaction no peer can use.
			Accept: 1<<types.LegacyTxType |
				1<<types.AccessListTxType |
				1<<types.DynamicFeeTxType |
				1<<types.SetCodeTxType,
			MaxSize: maxBroadcastTxSize,
			MinTip:  minTip,
		},
	}

	if opts.MinGasPrice != nil && opts.MinGasPrice.Sign() > 0 {
		v.minGasPrice = new(big.Int).Set(opts.MinGasPrice)
	}
	if opts.MinBaseFeeRatio > 0 {
		v.baseFeeRatio = big.NewInt(int64(opts.MinBaseFeeRatio * baseFeeRatioScale))
	}

	return v, nil
}

// Validate reports whether a transaction is worth forwarding to peers, given
// the current head block. A nil return means the transaction should be
// broadcast; any error means it should be dropped from the broadcast path.
func (v *TxValidator) Validate(tx *types.Transaction, head *types.Header) error {
	if tx == nil {
		return errors.New("nil transaction")
	}
	if head == nil {
		return errors.New("nil head header")
	}

	// ValidateTransaction reads head.Difficulty to decide whether the chain is
	// post-merge. Headers reconstructed from RPC or RLP can leave it nil, and a
	// nil dereference there would kill the peer's protocol goroutine, so fill in
	// a copy rather than trusting the caller.
	if head.Difficulty == nil {
		clone := *head
		clone.Difficulty = new(big.Int)
		head = &clone
	}

	if err := txpool.ValidateTransaction(tx, head, v.signer, v.opts); err != nil {
		return err
	}

	if v.minGasPrice != nil && tx.GasFeeCapIntCmp(v.minGasPrice) < 0 {
		return fmt.Errorf("%w: gas fee cap %v, minimum needed %v", errBelowMinGasPrice, tx.GasFeeCap(), v.minGasPrice)
	}

	if v.baseFeeRatio != nil && head.BaseFee != nil {
		floor := new(big.Int).Mul(head.BaseFee, v.baseFeeRatio)
		floor.Div(floor, big.NewInt(baseFeeRatioScale))
		if tx.GasFeeCapIntCmp(floor) < 0 {
			return fmt.Errorf("%w: gas fee cap %v, floor %v", core.ErrFeeCapTooLow, tx.GasFeeCap(), floor)
		}
	}

	return nil
}

// RejectReason maps a validation error to a short, bounded label for metrics.
// Unrecognized errors collapse to "other" so a peer cannot inflate Prometheus
// label cardinality by crafting transactions.
func RejectReason(err error) string {
	switch {
	case err == nil:
		return ""
	case errors.Is(err, core.ErrTxTypeNotSupported):
		return "unsupported_type"
	case errors.Is(err, txpool.ErrOversizedData):
		return "oversized"
	case errors.Is(err, vm.ErrMaxInitCodeSizeExceeded):
		return "init_code_too_large"
	case errors.Is(err, txpool.ErrNegativeValue):
		return "negative_value"
	case errors.Is(err, txpool.ErrGasLimit), errors.Is(err, core.ErrGasLimitTooHigh):
		return "gas_limit"
	case errors.Is(err, core.ErrFeeCapVeryHigh), errors.Is(err, core.ErrTipVeryHigh):
		return "fee_overflow"
	case errors.Is(err, core.ErrTipAboveFeeCap):
		return "tip_above_fee_cap"
	case errors.Is(err, txpool.ErrInvalidSender):
		return "invalid_signature"
	case errors.Is(err, core.ErrNonceMax):
		return "nonce_max"
	case errors.Is(err, core.ErrIntrinsicGas), errors.Is(err, core.ErrFloorDataGas):
		return "intrinsic_gas"
	case errors.Is(err, core.ErrInsufficientFunds):
		return "value_overflow"
	case errors.Is(err, txpool.ErrTxGasPriceTooLow):
		return "tip_too_low"
	case errors.Is(err, errBelowMinGasPrice):
		return "below_min_gas_price"
	case errors.Is(err, core.ErrFeeCapTooLow):
		return "fee_cap_below_basefee"
	default:
		return "other"
	}
}

// broadcastChainConfig builds a permissive chain config for the given chain ID:
// every fork through Prague is active from genesis.
//
// The sensor does not know its chain's fork schedule (it only knows the network
// ID and a fork ID hash), and this config is used solely to decide what to
// forward. Being generous about which transaction types are legal is the safe
// direction: the cost of accepting a type the chain has not enabled yet is one
// rebroadcast, whereas pinning an older fork set would silently stop forwarding
// legitimate traffic -- every EIP-7702 transaction, for instance. Osaka and
// later are left off so their tighter caps do not reject transactions the chain
// still accepts.
func broadcastChainConfig(chainID uint64) *params.ChainConfig {
	zero := uint64(0)
	return &params.ChainConfig{
		ChainID:             new(big.Int).SetUint64(chainID),
		HomesteadBlock:      big.NewInt(0),
		EIP150Block:         big.NewInt(0),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
		IstanbulBlock:       big.NewInt(0),
		MuirGlacierBlock:    big.NewInt(0),
		BerlinBlock:         big.NewInt(0),
		LondonBlock:         big.NewInt(0),
		ShanghaiTime:        &zero,
		CancunTime:          &zero,
		PragueTime:          &zero,
		BlobScheduleConfig: &params.BlobScheduleConfig{
			Cancun: params.DefaultCancunBlobConfig,
			Prague: params.DefaultPragueBlobConfig,
		},
	}
}
