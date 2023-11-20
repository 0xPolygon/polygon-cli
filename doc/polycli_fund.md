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
      --concurrency int            Concurrency level for speeding up funding wallets (default 2)
      --funding-wallet-pk string   Corresponding private key for funding wallet address, ensure you remove leading 0x
  -h, --help                       help for fund
      --output-file string         Specify the output JSON file name (default "funded_wallets.json")
      --rpc-url string             The RPC endpoint url (default "http://localhost:8545")
      --wallet-count int           Number of wallets to fund (default 2)
      --wallet-funding-amt float   Amount to fund each wallet with (default 0.05)
      --wallet-funding-gas uint    Gas for each wallet funding transaction (default 100000)
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
