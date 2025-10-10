# `polycli fund`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Bulk fund crypto wallets automatically.

```bash
polycli fund [flags]
```

## Usage

Bulk fund crypto wallets automatically.

```bash
# Fund wallets specified by the user.
$ polycli fund --addresses=0x5eD3BE7a1cDafd558F88a673345889dC75837aA2,0x1Ec6efdBd371D6444779eAE7B7e16907e0c8eC27
3:58PM INF Starting bulk funding wallets
3:58PM INF Using addresses provided by the user
3:58PM INF Wallet(s) funded! 💸
3:58PM INF Total execution time: 1.020693583s

# Fund 20 random wallets using a pre-deployed contract address.
$ polycli fund --number=20 --contract-address=0xf5a73e7cfcc83b7e8ce2e17eb44f050e8071ee60
3:58PM INF Starting bulk funding wallets
3:58PM INF Deriving wallets from the default mnemonic
3:58PM INF Wallet(s) derived count=20
3:58PM INF Wallet(s) funded! 💸
3:58PM INF Total execution time: 396.814917ms

# Fund 20 random wallets.
$ polycli fund --number 20 --hd-derivation=false
3:58PM INF Starting bulk funding wallets
3:58PM INF Generating random wallets
3:58PM INF Wallet(s) generated count=20
3:58PM INF Wallets' address(es) and private key(s) saved to file fileName=wallets.json
3:58PM INF Wallet(s) funded! 💸
3:58PM INF Total execution time: 1.027506s

# Fund wallets from a key file (one private key in hex per line).
$ polycli fund --key-file=keys.txt
3:58PM INF Starting bulk funding wallets
3:58PM INF Wallet(s) derived from key file count=3
3:58PM INF Wallet(s) funded! 💸
3:58PM INF Total execution time: 1.2s
```

Extract from `wallets.json`.

```json
[
  {
    "Address": "0xc1A44B1e37EE1fca4C6Fd5562c730d5b8525e4C6",
    "PrivateKey": "c1a8f737fd9f78aee361bfd856f9b2e99f853a5fe5efa2131fb030acdcee762b"
  },
  {
    "Address": "0x5D8121cf716B70d3e345adB58157752304eED5C3",
    "PrivateKey": "fbc57de542cef10fdcdf99e5578ffb5508992e9a8623ea4a39ab957d77e9b849"
  },
  ...
]
```

Check the balances of the wallets.

```bash
$ cast balance 0xc1A44B1e37EE1fca4C6Fd5562c730d5b8525e4C6
50000000000000000

$ cast balance 0x5D8121cf716B70d3e345adB58157752304eED5C3
50000000000000000
...
```

## Flags

```bash
      --addresses strings         comma-separated list of wallet addresses to fund
      --contract-address string   address of pre-deployed Funder contract
      --eth-amount big.Int        amount of wei to send to each wallet (default 50000000000000000)
  -f, --file string               output JSON file path for storing addresses and private keys of funded wallets (default "wallets.json")
      --hd-derivation             derive wallets to fund from private key in deterministic way (default true)
  -h, --help                      help for fund
      --key-file string           file containing accounts private keys, one per line
  -n, --number uint               number of wallets to fund (default 10)
      --private-key string        hex encoded private key to use for sending transactions (default "0x42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa")
  -r, --rpc-url string            RPC endpoint URL (default "http://localhost:8545")
      --seed string               seed string for deterministic wallet generation (e.g., 'ephemeral_test')
```

The command also inherits flags from parent commands.

```bash
      --config string      config file (default is $HOME/.polygon-cli.yaml)
      --pretty-logs        output logs in pretty format instead of JSON (default true)
  -v, --verbosity string   log level (string or int):
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

- [polycli](polycli.md) - A Swiss Army knife of blockchain tools.
