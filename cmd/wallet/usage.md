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
