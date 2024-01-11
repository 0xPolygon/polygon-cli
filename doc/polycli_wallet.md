# `polycli wallet`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Create or inspect BIP39(ish) wallets.

```bash
polycli wallet [create|inspect] [flags]
```

## Usage

This command is meant to simplify the operations of creating
wallets. This command can take a seed phrase and spit out child
accounts or generate new accmounts along with a seed phrase. It can
generate portable wallets to be used across ETH, BTC, PoS, Substrate,
etc.

In the example, we're generating a wallet with a few flags that are
used to configure how many wallets are generated and how the seed
phrase is used to generate the wallets.

```bash
$ polycli wallet create --raw-entropy --root-only --words 15 --language english
```

In addition to generating wallets with new mnemonics, you can use a
known mnemonic to generate wallets. **Caution** entering your seed
phrase in the command line should only be done for test
mnemonics. Never do this with a real seed phrase.

The example below is a test vector from Substrate.

[BIP-0039](https://github.com/paritytech/substrate-bip39/blob/eef2f86337d2dab075806c12948e8a098aa59d59/src/lib.rs#L74) where the expected seed is `44e9d125f037ac1d51f0a7d3649689d422c2af8b1ec8e00d71db4d7bf6d127e33f50c3d5c84fa3e5399c72d6cbbbbc4a49bf76f76d952f479d74655a2ef2d453`

```bash
$ polycli wallet inspect --raw-entropy --root-only --language english --password "Substrate" --mnemonic "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
```

This command also leverages the BIP-0032 library for hierarchically derived wallets.

```bash
$ polycli wallet create --path "m/44'/0'/0'" --addresses 5
```

## Flags

```bash
      --addresses uint         The number of addresses to generate (default 10)
  -h, --help                   help for wallet
      --iterations uint        Number of pbkdf2 iterations to perform (default 2048)
      --language string        Which language to use [ChineseSimplified, ChineseTraditional, Czech, English, French, Italian, Japanese, Korean, Spanish] (default "english")
      --mnemonic string        A mnemonic phrase used to generate entropy
      --mnemonic-file string   A mneomonic phrase written in a file used to generate entropy
      --password string        Password used along with the mnemonic
      --password-file string   Password stored in a file used along with the mnemonic
      --path string            What would you like the derivation path to be (default "m/44'/60'/0'")
      --raw-entropy            substrate and polkda dot don't follow strict bip39 and use raw entropy
      --root-only              don't produce HD accounts. Just produce a single wallet
      --words int              The number of words to use in the mnemonic (default 24)
```

The command also inherits flags from parent commands.

```bash
      --config string   config file (default is $HOME/.polygon-cli.yaml)
      --pretty-logs     Should logs be in pretty format or JSON (default true)
  -v, --verbosity int   0 - Silent
                        100 Panic
                        200 Fatal
                        300 Error
                        400 Warning
                        500 Info
                        600 Debug
                        700 Trace (default 500)
```

## See also

- [polycli](polycli.md) - A Swiss Army knife of blockchain tools.
