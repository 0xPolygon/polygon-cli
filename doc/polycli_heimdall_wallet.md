# `polycli heimdall wallet`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Manage keystores, keys, and message signatures.

## Usage

Local key and keystore management, compatible with Foundry's `cast wallet`.

All subcommands are offline. Keystores are written in the go-ethereum
v3 JSON format, which is byte-for-byte compatible with Foundry. Any
existing `cast wallet` keystores under `~/.foundry/keystores/` are
picked up automatically.

The keystore directory is chosen in the following order, highest
priority first:

1. `--keystore-dir` flag.
2. `ETH_KEYSTORE` environment variable.
3. `~/.foundry/keystores/` if it already exists.
4. `~/.polycli/keystores/` (default; created on demand).

```bash
# Generate a new random key and write to the keystore.
polycli heimdall wallet new

# Generate a new BIP-39 mnemonic and print the first address.
polycli heimdall wallet new-mnemonic --print-only

# Inspect an existing key.
polycli heimdall wallet address 0x1234...
polycli heimdall wallet address --private-key 0xabc...
polycli heimdall wallet address --mnemonic "abandon abandon ... about"

# Derive a range of addresses from a mnemonic.
polycli heimdall wallet derive --mnemonic "abandon abandon ... about" --count 5

# Sign a message (EIP-191 personal_sign) and verify it.
polycli heimdall wallet sign "hello" --address 0x1234...
polycli heimdall wallet verify 0x1234... "hello" 0x<signature>

# Import a private key, a keystore file, or a mnemonic.
polycli heimdall wallet import --private-key 0xabc...
polycli heimdall wallet import --source-keystore-file path/to/UTC--...
polycli heimdall wallet import --mnemonic "abandon ... about"

# List / remove / change password.
polycli heimdall wallet list
polycli heimdall wallet remove 0x1234... --yes
polycli heimdall wallet change-password 0x1234...

# Emit the public key in both compressed and uncompressed form.
polycli heimdall wallet public-key 0x1234...

# Plaintext key export (guarded by friction flag).
polycli heimdall wallet private-key 0x1234... --i-understand-the-risks
polycli heimdall wallet decrypt-keystore path/to/UTC--... --i-understand-the-risks
```

Hardware wallets (`--ledger`, `--trezor`), `vanity`, and `sign-auth`
are intentionally out of scope. Use `cast wallet` directly for those.

## Flags

```bash
  -h, --help   help for wallet
```

The command also inherits flags from parent commands.

```bash
      --amoy                     shortcut for --network amoy (default)
      --chain-id string          chain id used for signing
      --color string             color mode (auto|always|never) (default "auto")
      --config string            config file (default is $HOME/.polygon-cli.yaml)
      --curl                     print the equivalent curl command instead of executing
      --denom string             fee denom
      --heimdall-config string   path to heimdall config TOML (default ~/.polycli/heimdall.toml)
  -k, --insecure                 accept invalid TLS certs
      --json                     emit JSON instead of key/value
      --mainnet                  shortcut for --network mainnet
  -N, --network string           named network preset (amoy|mainnet)
      --no-color                 disable color output
      --pretty-logs              output logs in pretty format instead of JSON (default true)
      --raw                      preserve raw bytes (no 0x-hex normalization)
  -r, --rest-url string          heimdall REST gateway URL
      --rpc-headers string       extra request headers, comma-separated key=value pairs
  -R, --rpc-url string           cometBFT RPC URL
      --timeout int              HTTP timeout in seconds
  -v, --verbosity string         log level (string or int):
                                   0   - silent
                                   100 - panic
                                   200 - fatal
                                   300 - error
                                   400 - warn
                                   500 - info (default)
                                   600 - debug
                                   700 - trace (default "info")
```

## See also

- [polycli heimdall](polycli_heimdall.md) - Query and interact with a Heimdall v2 node.
- [polycli heimdall wallet address](polycli_heimdall_wallet_address.md) - Show the address for a key or keystore file.

- [polycli heimdall wallet change-password](polycli_heimdall_wallet_change-password.md) - Change a keystore entry's password.

- [polycli heimdall wallet decrypt-keystore](polycli_heimdall_wallet_decrypt-keystore.md) - Decrypt a keystore file to its plaintext private key.

- [polycli heimdall wallet derive](polycli_heimdall_wallet_derive.md) - Derive addresses from a BIP-39 mnemonic.

- [polycli heimdall wallet import](polycli_heimdall_wallet_import.md) - Import an existing key into the keystore.

- [polycli heimdall wallet list](polycli_heimdall_wallet_list.md) - List keys in the keystore.

- [polycli heimdall wallet new](polycli_heimdall_wallet_new.md) - Generate a new key in the keystore.

- [polycli heimdall wallet new-mnemonic](polycli_heimdall_wallet_new-mnemonic.md) - Generate a new BIP-39 mnemonic and derive a key.

- [polycli heimdall wallet private-key](polycli_heimdall_wallet_private-key.md) - Print the plaintext private key for a keystore entry.

- [polycli heimdall wallet public-key](polycli_heimdall_wallet_public-key.md) - Print the secp256k1 public key for a key.

- [polycli heimdall wallet remove](polycli_heimdall_wallet_remove.md) - Remove a key from the keystore.

- [polycli heimdall wallet sign](polycli_heimdall_wallet_sign.md) - Sign a message with a keystore key.

- [polycli heimdall wallet sign-auth](polycli_heimdall_wallet_sign-auth.md) - Not supported — use `cast wallet sign-auth`.

- [polycli heimdall wallet vanity](polycli_heimdall_wallet_vanity.md) - Not supported — use `cast wallet vanity`.

- [polycli heimdall wallet verify](polycli_heimdall_wallet_verify.md) - Verify a signature against an address.

