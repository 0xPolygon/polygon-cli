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
