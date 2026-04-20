// Package msgs wires the per-Msg subcommands shared by
// `polycli heimdall mktx`, `polycli heimdall send`, and
// `polycli heimdall estimate`. The package owns the shared flag bag
// (TxOpts), the msg-subcommand registry, and the common
// build/sign/broadcast/simulate pipeline.
//
// For W3 the only implemented Msg is `withdraw` (MsgWithdrawFeeTx).
// Additional Msg types land alongside their subcommands in W4 by
// calling RegisterFactory in an init or Register call — see
// msgs/registry.go for the contract.
package msgs

import (
	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
)

// Mode is the action performed by an umbrella command: build only,
// build + broadcast, or build + simulate. Each Msg subcommand
// inspects the mode to tune output (e.g. skip broadcasting in
// `mktx`) and to skip fetching account info when inputs are
// sufficient without it.
type Mode int

const (
	// ModeMkTx builds a TxRaw and prints it. Never broadcasts.
	ModeMkTx Mode = iota
	// ModeSend builds, signs, broadcasts, and waits for inclusion.
	ModeSend
	// ModeEstimate builds, signs, and calls /cosmos/tx/v1beta1/simulate.
	ModeEstimate
)

// String returns the umbrella command name for a mode. Used in
// usage/error strings.
func (m Mode) String() string {
	switch m {
	case ModeMkTx:
		return "mktx"
	case ModeSend:
		return "send"
	case ModeEstimate:
		return "estimate"
	default:
		return "unknown"
	}
}

// TxOpts is the shared flag bag every msg subcommand receives. The
// fields are populated by cobra via RegisterFlags; a single TxOpts is
// allocated per (mode, msg) pair because cobra flag variables must be
// addressable and persist across the command's lifetime.
//
// TxOpts intentionally carries only the "how to build/sign/broadcast"
// knobs — per-message fields (e.g. `--amount` on withdraw) live on
// the msg subcommand itself.
type TxOpts struct {
	// Shared flags injected from the parent heimdall command. We keep
	// the pointer so each msg subcommand can resolve the network
	// config at run time.
	Global *config.Flags

	// Wallet: how to obtain the signing key.
	From           string
	KeystoreDir    string
	KeystoreFile   string
	Account        string
	Password       string
	PasswordFile   string
	PrivateKey     string
	Mnemonic       string
	MnemonicIndex  uint32
	DerivationPath string

	// Gas / fee.
	Gas           uint64
	GasAdjustment float64
	GasPrice      float64
	Fee           string
	Memo          string

	// Account overrides. Non-zero values skip the auto-fetch via
	// /cosmos/auth/v1beta1/accounts.
	AccountNumber uint64
	Sequence      uint64

	// Sign / broadcast.
	SignMode      string
	DryRun        bool
	Async         bool
	Confirmations uint64
	Force         bool

	// Output.
	JSONOut bool
}

// RegisterFlags attaches the shared tx flags to cmd. Call exactly
// once per msg-subcommand instance; both mktx/send/estimate share the
// same flag surface but each owns its own TxOpts so cobra can parse
// distinct invocations correctly.
//
// Flag names follow the cast-style dash-separated convention from the
// heimdall CLAUDE.md. No leading articles, lowercase usage strings,
// no ending punctuation. Global/network flags (--rest-url, --rpc-url,
// --chain-id, --timeout, --json) are inherited from the heimdall
// parent command via its PersistentFlags.
func RegisterFlags(cmd *cobra.Command, opts *TxOpts, mode Mode) {
	f := cmd.Flags()
	// Wallet.
	f.StringVar(&opts.From, "from", "", "signer address (20-byte hex)")
	f.StringVar(&opts.KeystoreDir, "keystore-dir", "", "keystore directory (overrides ETH_KEYSTORE)")
	f.StringVar(&opts.KeystoreFile, "keystore-file", "", "explicit keystore JSON file path")
	f.StringVar(&opts.Account, "account", "", "address or index into keystore (overrides --from for key lookup)")
	f.StringVar(&opts.Password, "password", "", "keystore password (mutually exclusive with --password-file)")
	f.StringVar(&opts.PasswordFile, "password-file", "", "path to file containing keystore password")
	f.StringVar(&opts.PrivateKey, "private-key", "", "hex-encoded secp256k1 private key (unsafe outside local dev)")
	f.StringVar(&opts.Mnemonic, "mnemonic", "", "BIP-39 mnemonic used to derive the signing key")
	f.Uint32Var(&opts.MnemonicIndex, "mnemonic-index", 0, "address index when deriving from --mnemonic")
	f.StringVar(&opts.DerivationPath, "derivation-path", "", "BIP-32 derivation path (default m/44'/60'/0'/0/<index>)")

	// Gas / fee.
	f.Uint64Var(&opts.Gas, "gas", 0, "gas limit (0 means estimate via simulation)")
	f.Float64Var(&opts.GasAdjustment, "gas-adjustment", 1.3, "multiplier applied to simulated gas to pick final gas limit")
	f.Float64Var(&opts.GasPrice, "gas-price", 0, "fee price per gas unit in the default denom")
	f.StringVar(&opts.Fee, "fee", "", "explicit fee coin amount, e.g. 10000pol (overrides --gas-price)")
	f.StringVar(&opts.Memo, "memo", "", "optional tx memo")

	// Account overrides.
	f.Uint64Var(&opts.AccountNumber, "account-number", 0, "override fetched account number")
	f.Uint64Var(&opts.Sequence, "sequence", 0, "override fetched sequence")

	// Sign / broadcast.
	f.StringVar(&opts.SignMode, "sign-mode", "direct", "signing mode (direct|amino-json)")
	f.BoolVar(&opts.DryRun, "dry-run", false, "build the tx but do not broadcast")
	f.BoolVar(&opts.Async, "async", false, "use BROADCAST_MODE_ASYNC and skip inclusion polling")
	f.Uint64Var(&opts.Confirmations, "confirmations", 0, "after inclusion, wait for N additional blocks")
	f.BoolVar(&opts.Force, "force", false, "bypass safety guards for L1-mirroring message types")

	// Output.
	f.BoolVar(&opts.JSONOut, "json", false, "emit JSON instead of key/value output")

	// Mode-specific tweaks.
	switch mode {
	case ModeMkTx:
		// `mktx` is a pure build; hide broadcast-only flags so `--help`
		// doesn't advertise things that do nothing.
		_ = cmd.Flags().MarkHidden("async")
		_ = cmd.Flags().MarkHidden("confirmations")
		_ = cmd.Flags().MarkHidden("dry-run")
	case ModeEstimate:
		// `estimate` never broadcasts either.
		_ = cmd.Flags().MarkHidden("async")
		_ = cmd.Flags().MarkHidden("confirmations")
		_ = cmd.Flags().MarkHidden("dry-run")
	}
}
