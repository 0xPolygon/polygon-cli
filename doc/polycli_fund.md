# `polycli fund`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Bulk fund many crypto wallets automatically.

```bash
polycli fund [flags]
```

## Usage

```bash
$ polycli fund \
  --wallet-count=5 \
  --funding-wallet-pk="REPLACE" \
  --chain-id=100 \
  --concurrency=5 \
  --rpc-url="https://rootchain-devnetsub.zkevmdev.net"  \
  --wallet-funding-amt=0.00015 \
  --wallet-funding-gas=50000 \
  --output-file="/opt/generated_keys.json"
  --verbosity=true
```
## Flags

```bash
      --chain-id int               The chain id for the transactions.
      --concurrency int            Concurrency level for speeding up funding wallets (default 5)
      --funding-wallet-pk string   Corresponding private key for funding wallet address, ensure you remove leading 0x
  -h, --help                       help for fund
      --output-file string         Specify the output CSV file name (default "wallets.csv")
      --rpc-url string             The RPC endpoint url (default "http://localhost:8545")
      --verbosity                  Global verbosity flag (true/false) (default true)
      --wallet-count int           Number of wallets to fund (default 2)
      --wallet-funding-amt float   Amount to fund each wallet with (default 0.05)
      --wallet-funding-gas uint    Gas for each wallet funding transaction (default 50000)
```

The command also inherits flags from parent commands.

```bash
      --config string   config file (default is $HOME/.polygon-cli.yaml)
      --pretty-logs     Should logs be in pretty format or JSON (default true)
```

## See also

- [polycli](polycli.md) - A Swiss Army knife of blockchain tools.
