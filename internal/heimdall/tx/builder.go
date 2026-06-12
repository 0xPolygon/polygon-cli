package tx

import (
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
	hproto "github.com/0xPolygon/polygon-cli/internal/heimdall/proto"
)

// Msg is a Heimdall / Cosmos SDK message that can be packed into a
// google.protobuf.Any and serialized as part of a TxBody.
//
// Implementations are tiny — one per concrete Msg type — and live with
// the command that builds them (W3/W4 wave). The first implementation,
// for MsgWithdrawFeeTx, ships with this package to exercise the
// builder end-to-end.
type Msg interface {
	// TypeURL is the Any.type_url for the message, e.g.
	// "/heimdallv2.topup.MsgWithdrawFeeTx".
	TypeURL() string
	// Marshal returns the proto-serialized message (the payload that
	// goes inside Any.value).
	Marshal() ([]byte, error)
	// AminoName is the amino.name option on the message, used for
	// SIGN_MODE_LEGACY_AMINO_JSON. Return "" if amino-JSON is not
	// supported for this message; the builder will refuse to sign.
	AminoName() string
	// AminoJSON returns the message as a JSON-serializable Go value
	// (typically map[string]any) for the legacy amino-JSON sign doc.
	// Implementations should emit uint64 / int64 fields as decimal
	// strings to match cosmos-sdk behavior.
	AminoJSON() (any, error)
}

// WithdrawFeeMsg is the starter Msg for MsgWithdrawFeeTx. Additional
// Msg types land alongside their subcommands in W3/W4.
type WithdrawFeeMsg struct {
	Proposer string
	Amount   string
}

// TypeURL implements Msg.
func (m *WithdrawFeeMsg) TypeURL() string { return hproto.MsgWithdrawFeeTxTypeURL }

// Marshal implements Msg.
func (m *WithdrawFeeMsg) Marshal() ([]byte, error) {
	if m.Proposer == "" {
		return nil, fmt.Errorf("WithdrawFeeMsg: proposer is required")
	}
	if m.Amount == "" {
		return nil, fmt.Errorf("WithdrawFeeMsg: amount is required")
	}
	p := &hproto.MsgWithdrawFeeTx{Proposer: m.Proposer, Amount: m.Amount}
	return p.Marshal(), nil
}

// AminoName implements Msg.
func (m *WithdrawFeeMsg) AminoName() string { return "heimdallv2/topup/MsgWithdrawFeeTx" }

// AminoJSON implements Msg.
func (m *WithdrawFeeMsg) AminoJSON() (any, error) {
	return map[string]any{
		"proposer": m.Proposer,
		"amount":   m.Amount,
	}, nil
}

// Builder constructs a TxRaw. Call the With* setters, then Sign (which
// produces the raw bytes) or SignAndEncode (which also base64/hex
// encodes them). Builders are single-use; reuse after Sign is
// undefined.
type Builder struct {
	msgs          []Msg
	memo          string
	timeoutHeight uint64
	fee           hproto.Fee
	gasLimit      uint64
	signMode      SignMode
	chainID       string
	accountNumber uint64
	sequence      uint64
	// accountFetched is true once ResolveAccount has populated
	// accountNumber and sequence. Direct callers who set both fields by
	// hand (via WithAccountNumber / WithSequence) can skip the
	// auto-fetch.
	accountFetched bool
	// Signing key. Populated by Sign; the builder never stores the key
	// longer than necessary.
}

// NewBuilder returns a fresh Builder with direct sign mode and a
// zero-value fee. The caller must populate at least WithChainID and
// one Msg before Sign.
func NewBuilder() *Builder {
	return &Builder{signMode: SignModeDirect}
}

// AddMsg appends a Msg to the builder's TxBody. Repeatable; order is
// preserved and determines the order of signer_infos.
func (b *Builder) AddMsg(m Msg) *Builder {
	b.msgs = append(b.msgs, m)
	return b
}

// WithChainID sets the chain id used for SIGN_MODE_DIRECT replay
// protection.
func (b *Builder) WithChainID(id string) *Builder { b.chainID = id; return b }

// WithMemo sets the tx memo. Default empty.
func (b *Builder) WithMemo(memo string) *Builder { b.memo = memo; return b }

// WithTimeoutHeight sets the absolute block height after which the tx
// is invalid. Default 0 (no timeout).
func (b *Builder) WithTimeoutHeight(h uint64) *Builder { b.timeoutHeight = h; return b }

// WithFee sets the fee amount coins. Gas limit is set separately via
// WithGasLimit.
func (b *Builder) WithFee(coins ...hproto.Coin) *Builder {
	b.fee.Amount = append(b.fee.Amount[:0], coins...)
	return b
}

// WithGasLimit sets the TxBody gas limit.
func (b *Builder) WithGasLimit(g uint64) *Builder { b.gasLimit = g; b.fee.GasLimit = g; return b }

// WithFeePayer sets the optional fee payer address on Fee.
func (b *Builder) WithFeePayer(addr string) *Builder { b.fee.Payer = addr; return b }

// WithFeeGranter sets the optional fee granter address on Fee.
func (b *Builder) WithFeeGranter(addr string) *Builder { b.fee.Granter = addr; return b }

// WithSignMode sets the signing mode (direct or amino-json).
func (b *Builder) WithSignMode(m SignMode) *Builder { b.signMode = m; return b }

// WithAccountNumber overrides the auto-fetched account number.
func (b *Builder) WithAccountNumber(n uint64) *Builder {
	b.accountNumber = n
	b.accountFetched = true
	return b
}

// WithSequence overrides the auto-fetched sequence.
func (b *Builder) WithSequence(s uint64) *Builder { b.sequence = s; return b }

// Account is the subset of /cosmos/auth/v1beta1/accounts/{addr} that
// the builder reads. Exposed for tests / advanced callers who want to
// fetch once and feed into multiple builders.
type Account struct {
	Address       string
	AccountNumber uint64
	Sequence      uint64
}

// AccountFetcher is the narrow interface the builder uses to look up
// signer accounts. The production implementation wraps
// client.RESTClient; tests inject a fake.
type AccountFetcher interface {
	FetchAccount(ctx context.Context, addr string) (*Account, error)
}

// RESTAccountFetcher is the default AccountFetcher. It calls
// /cosmos/auth/v1beta1/accounts/{addr} and parses the BaseAccount
// response.
type RESTAccountFetcher struct {
	Client *client.RESTClient
}

// FetchAccount implements AccountFetcher.
func (f *RESTAccountFetcher) FetchAccount(ctx context.Context, addr string) (*Account, error) {
	body, status, err := f.Client.Get(ctx, "/cosmos/auth/v1beta1/accounts/"+url.PathEscape(addr), nil)
	if err != nil {
		return nil, fmt.Errorf("fetching account %s: %w", addr, err)
	}
	if status == 0 && body == nil {
		return nil, fmt.Errorf("--curl transport: account fetch cannot be simulated; pass --account-number/--sequence explicitly")
	}
	var parsed struct {
		Account struct {
			Address       string `json:"address"`
			AccountNumber string `json:"account_number"`
			Sequence      string `json:"sequence"`
		} `json:"account"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("decoding account %s: %w", addr, err)
	}
	a := &Account{Address: parsed.Account.Address}
	if parsed.Account.AccountNumber != "" {
		n, err := strconv.ParseUint(parsed.Account.AccountNumber, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parsing account_number %q: %w", parsed.Account.AccountNumber, err)
		}
		a.AccountNumber = n
	}
	if parsed.Account.Sequence != "" {
		n, err := strconv.ParseUint(parsed.Account.Sequence, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parsing sequence %q: %w", parsed.Account.Sequence, err)
		}
		a.Sequence = n
	}
	return a, nil
}

// ResolveAccount populates the builder's accountNumber and sequence by
// asking fetcher for the given signer address. Callers that set both
// values via WithAccountNumber and WithSequence can skip this step.
func (b *Builder) ResolveAccount(ctx context.Context, fetcher AccountFetcher, addr string) error {
	if fetcher == nil {
		return fmt.Errorf("ResolveAccount: fetcher is nil")
	}
	acc, err := fetcher.FetchAccount(ctx, addr)
	if err != nil {
		return err
	}
	b.accountNumber = acc.AccountNumber
	// Don't overwrite an explicitly-set sequence: operators sometimes
	// pre-compute sequences to send a batch of txs.
	if b.sequence == 0 {
		b.sequence = acc.Sequence
	}
	b.accountFetched = true
	return nil
}

// Sign builds the TxRaw bytes using priv as the single signing key.
// chainID and accountNumber must be set (via WithChainID and either
// WithAccountNumber or ResolveAccount). Returns the canonical TxRaw
// proto-encoded bytes, ready for /cosmos/tx/v1beta1/txs.
func (b *Builder) Sign(priv *ecdsa.PrivateKey) ([]byte, error) {
	if priv == nil {
		return nil, fmt.Errorf("Sign: private key is nil")
	}
	if len(b.msgs) == 0 {
		return nil, fmt.Errorf("Sign: no messages")
	}
	if b.chainID == "" {
		return nil, fmt.Errorf("Sign: chain id is required")
	}
	if b.gasLimit == 0 {
		return nil, fmt.Errorf("Sign: gas limit is required (set via WithGasLimit)")
	}

	// Build TxBody.
	anys := make([]*hproto.Any, 0, len(b.msgs))
	for _, m := range b.msgs {
		val, err := m.Marshal()
		if err != nil {
			return nil, fmt.Errorf("marshalling msg %s: %w", m.TypeURL(), err)
		}
		anys = append(anys, &hproto.Any{TypeURL: m.TypeURL(), Value: val})
	}
	body := &hproto.TxBody{
		Messages:      anys,
		Memo:          b.memo,
		TimeoutHeight: b.timeoutHeight,
	}
	bodyBytes := body.Marshal()

	// Build AuthInfo with a single signer whose pubkey is derived
	// from priv.
	pubKey := uncompressedPubkey(priv)
	authInfo := &hproto.AuthInfo{
		SignerInfos: []*hproto.SignerInfo{{
			PublicKey: hproto.PubKeyAny(pubKey),
			ModeInfo:  &hproto.ModeInfo{Single: &hproto.ModeInfoSingle{Mode: int32(b.signMode)}},
			Sequence:  b.sequence,
		}},
		Fee: &hproto.Fee{
			Amount:   b.fee.Amount,
			GasLimit: b.gasLimit,
			Payer:    b.fee.Payer,
			Granter:  b.fee.Granter,
		},
	}
	authInfoBytes := authInfo.Marshal()

	// Compute the sign digest for the chosen mode.
	var digest []byte
	switch b.signMode {
	case SignModeDirect:
		_, digest = signBytesDirect(bodyBytes, authInfoBytes, b.chainID, b.accountNumber)
	case SignModeAminoJSON:
		if name := b.msgs[0].AminoName(); name == "" {
			return nil, fmt.Errorf("amino-json unsupported for message type %s", b.msgs[0].TypeURL())
		}
		_, d, err := signBytesAminoJSON(b, b.accountNumber)
		if err != nil {
			return nil, err
		}
		digest = d
	default:
		return nil, fmt.Errorf("unsupported sign mode %d", b.signMode)
	}

	sig, err := signDigest(priv, digest)
	if err != nil {
		return nil, err
	}

	raw := &hproto.TxRaw{
		BodyBytes:     bodyBytes,
		AuthInfoBytes: authInfoBytes,
		Signatures:    [][]byte{sig},
	}
	return raw.Marshal(), nil
}

// SignAndEncode signs and returns the TxRaw bytes in both base64 and
// hex encodings. Callers typically print one or both.
func (b *Builder) SignAndEncode(priv *ecdsa.PrivateKey) (raw []byte, b64, hex string, err error) {
	raw, err = b.Sign(priv)
	if err != nil {
		return nil, "", "", err
	}
	return raw, base64.StdEncoding.EncodeToString(raw), encodeHex(raw), nil
}

// encodeHex prints raw as lower-case 0x-prefixed hex, the form
// operators paste into `polycli heimdall publish` or debug tooling.
func encodeHex(raw []byte) string {
	const alphabet = "0123456789abcdef"
	out := strings.Builder{}
	out.Grow(2 + 2*len(raw))
	out.WriteString("0x")
	for _, c := range raw {
		out.WriteByte(alphabet[c>>4])
		out.WriteByte(alphabet[c&0x0f])
	}
	return out.String()
}
