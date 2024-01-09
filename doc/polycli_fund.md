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
      --addresses strings         Comma-separated list of wallet addresses to fund
      --contract-address string   The address of a pre-deployed Funder contract
  -a, --eth-amount float          The amount of ether to send to each wallet (default 0.05)
  -f, --file string               The output JSON file path for storing the addresses and private keys of funded wallets (default "wallets.json")
      --hd-derivation             Derive wallets to fund from the private key in a deterministic way (default true)
  -h, --help                      help for fund
  -n, --number uint               The number of wallets to fund (default 10)
      --private-key string        The hex encoded private key that we'll use to send transactions (default "0x42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa")
  -r, --rpc-url string            The RPC endpoint url (default "http://localhost:8545")
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
