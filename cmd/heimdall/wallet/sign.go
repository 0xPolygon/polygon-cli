package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// newSignCmd builds `wallet sign <MESSAGE>`: sign a message with a
// key from the keystore (or with a private key supplied via
// --private-key). Default is EIP-191 personal_sign; --raw signs a
// 32-byte hash directly.
func newSignCmd() *cobra.Command {
	var (
		shared     keystoreSharedFlags
		addrFlag   string
		privateKey string
		raw        bool
	)
	cmd := &cobra.Command{
		Use:   "sign <message-or-hash>",
		Short: "Sign a message with a keystore key.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			msg := args[0]
			priv, err := resolveSigningKey(&shared, privateKey, addrFlag)
			if err != nil {
				return err
			}
			var sig []byte
			if raw {
				// For --raw the argument is always hex (0x-optional).
				hashBytes, err := parseHex(msg, "hash")
				if err != nil {
					return err
				}
				if len(hashBytes) != 32 {
					return &client.UsageError{Msg: fmt.Sprintf("--raw input must decode to 32 bytes, got %d", len(hashBytes))}
				}
				sig, err = signRawHash(priv, hashBytes)
				if err != nil {
					return err
				}
			} else {
				// EIP-191 personal_sign. If the argument is hex, sign
				// the raw decoded bytes to match cast's behaviour.
				payload := []byte(msg)
				if decoded, err := parseHex(msg, "message"); err == nil {
					payload = decoded
				}
				sig, err = signPersonal(priv, payload)
				if err != nil {
					return err
				}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "0x%s\n", hex.EncodeToString(sig))
			return nil
		},
	}
	bindKeystoreFlags(cmd, &shared)
	f := cmd.Flags()
	f.StringVar(&addrFlag, "address", "", "address of the keystore key to sign with")
	f.StringVar(&privateKey, "private-key", "", "hex-encoded private key (skips the keystore)")
	f.BoolVar(&raw, "raw", false, "sign the 32-byte hash directly (no EIP-191 framing)")
	rejectHardwareFlags(cmd)
	return cmd
}

// resolveSigningKey returns the ECDSA private key that a `sign`
// invocation should use. Precedence: --private-key > keystore
// address + password. An explicit --keystore-file resolves by
// reading the address from that file.
func resolveSigningKey(shared *keystoreSharedFlags, privHex, addrFlag string) (*ecdsa.PrivateKey, error) {
	if privHex != "" {
		return parsePrivateKeyHex(privHex)
	}
	identifier := addrFlag
	if identifier == "" {
		identifier = shared.KeystoreFile
	}
	if identifier == "" {
		return nil, &client.UsageError{Msg: "one of --address, --keystore-file, or --private-key is required"}
	}
	dir, err := resolveKeystoreDir(shared.KeystoreDir)
	if err != nil {
		return nil, err
	}
	ks := newKeyStore(dir)
	acc, err := findAccount(ks, identifier)
	if err != nil {
		return nil, err
	}
	password, err := readPassword(shared, os.Stdin, false, "keystore password")
	if err != nil {
		return nil, err
	}
	return decryptKeystoreAccount(acc, password)
}

// parseHex decodes a 0x-prefixed (or bare) hex string into bytes.
// Returns a usage error label tailored to the caller when decoding
// fails. An empty input is treated as a decode failure rather than
// producing an empty byte slice so `wallet sign ""` still errors.
func parseHex(input, label string) ([]byte, error) {
	s := strings.TrimSpace(input)
	s = strings.TrimPrefix(strings.TrimPrefix(s, "0x"), "0X")
	if s == "" {
		return nil, &client.UsageError{Msg: fmt.Sprintf("%s is empty", label)}
	}
	if len(s)%2 != 0 {
		return nil, &client.UsageError{Msg: fmt.Sprintf("%s must have an even number of hex chars, got %d", label, len(s))}
	}
	raw, err := hex.DecodeString(s)
	if err != nil {
		return nil, &client.UsageError{Msg: fmt.Sprintf("decoding %s: %v", label, err)}
	}
	return raw, nil
}
