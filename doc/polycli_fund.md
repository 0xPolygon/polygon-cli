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

```bash
$ polycli fund \
  --rpc-url="https://rootchain-devnetsub.zkevmdev.net"  \
  --funding-wallet-pk="REPLACE" \
  --wallet-count=5 \
  --wallet-funding-amt=0.00015 \
  --wallet-funding-gas=50000 \
  --concurrency=5 \
  --output-file="/opt/funded_wallets.json"
```
## Flags

```bash
  -a, --amount float         The amount of eth to send to each wallet (default 0.05)
  -c, --concurrency uint     The concurrency level for speeding up funding wallets (default 2)
  -f, --file string          The output JSON file path for storing the addresses and private keys of funded wallets (default "wallets.json")
  -g, --gas uint             The cost of funding a wallet (default 21000)
  -h, --help                 help for fund
      --private-key string   The hex encoded private key that we'll use to send transactions (default "0x42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa")
  -r, --rpc-url string       The RPC endpoint url (default "http://localhost:8545")
  -w, --wallets uint         The number of wallets to fund (default 2)
```

The command also inherits flags from parent commands.

```bash
      --config string   config file (default is $HOME/.polygon-cli.yaml)
      --pretty-logs     Should logs be in pretty format or JSON (default true)
  -v, --verbosity int   0 - Silent
                        100 Fatal
                        200 Error
                        300 Warning
                        400 Info
                        500 Debug
                        600 Trace (default 400)
```

## See also

- [polycli](polycli.md) - A Swiss Army knife of blockchain tools.
