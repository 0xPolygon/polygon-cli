package p2p

import (
	"errors"
	"fmt"
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
}

// TxValidator applies the stateless half of a node's transaction admission
// rules to transactions the sensor is about to rebroadcast.
//
// It deliberately does not check nonce or balance: the sensor holds no state, and
// those checks would need an RPC round trip per sender.
//
// It also does not filter on the head block's base fee. Bor applies no such
// filter when it gossips a transaction it has just accepted (eth/handler.go
// BroadcastTransactions walks the pool's feed unconditionally); the base fee only
// gates its periodic re-broadcast of stuck pending transactions
// (legacypool.identifyStuckTransactions). Forwarding what we just received is
// the first case, not the second. A fractional base fee floor lived here for a
// while and was removed: on ~100k mainnet transactions it produced 97% of all
// rejections and every other check produced none, so it was doing nearly all of
// the "spam" filtering while dropping transactions bor itself would have
// relayed. The sensor's head also lags the chain by a median of 2 blocks and up
// to 7, so the comparison was made against a stale base fee. What it does catch is
// everything a peer can make up for free -- forged signatures, transactions signed
// for another chain, oversized payloads, gas below the intrinsic cost, fee caps
// that can never be mined -- which is the bulk of what a spamming peer sends.
type TxValidator struct {
	signer types.Signer
	opts   *txpool.ValidationOptions

	// minGasPrice is nil when the floor is disabled.
	minGasPrice *big.Int
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

	// go-ethereum gates every timestamp-scheduled fork -- Shanghai, Cancun,
	// Prague -- on the chain being post-merge, which ValidateTransaction infers
	// from a zero header difficulty. Bor keeps a non-zero difficulty forever, so
	// handing it the header as-is pins validation to London and rejects every
	// EIP-7702 transaction as an unsupported type. Validate against a copy with
	// the difficulty zeroed, which also covers a nil difficulty: geth
	// dereferences it unconditionally, and that panic would land on the peer's
	// protocol goroutine.
	postMerge := *head
	postMerge.Difficulty = new(big.Int)
	head = &postMerge

	if err := txpool.ValidateTransaction(tx, head, v.signer, v.opts); err != nil {
		return err
	}

	if v.minGasPrice != nil && tx.GasFeeCapIntCmp(v.minGasPrice) < 0 {
		return fmt.Errorf("%w: gas fee cap %v, minimum needed %v", errBelowMinGasPrice, tx.GasFeeCap(), v.minGasPrice)
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
	default:
		return "other"
	}
}

// broadcastChainConfig builds a chain config for the given chain ID with every
// fork through Prague active from genesis. Validate zeroes the header difficulty
// so the timestamp-scheduled forks in here actually take effect.
//
// The sensor does not know its chain's fork schedule -- it knows a network ID
// and a fork ID hash, neither of which yields one -- so this is an assumption,
// and it is the assumption that a chain the sensor is pointed at is current.
//
// IT IS NOT A PURELY PERMISSIVE ONE. Later forks add transaction types, which
// only ever costs an over-accept, but they also tighten: Shanghai caps init code
// size (EIP-3860) and Prague adds the calldata floor gas cost (EIP-7623). On a
// chain that has not adopted those, a large deployment or a calldata-heavy
// transaction is legal there and dropped here, showing up as
// init_code_too_large or intrinsic_gas.
//
// Prague is still the right floor for the networks this targets -- Polygon since
// Bhilai, Ethereum since Pectra -- and pinning an older fork set has the larger
// failure mode, silently dropping every EIP-7702 transaction as an unsupported
// type. Osaka and later are left off because their caps (MaxTxGas, Amsterdam's
// floor gas rules) would tighten further against chains that have not taken
// them. An operator on a pre-Prague chain should expect those two rejection
// reasons and turn validation off, or the floors down, accordingly.
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
