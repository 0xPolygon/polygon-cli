# `polycli ecrecover`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Recovers and returns the public key of the signature.

```bash
polycli ecrecover [flags]
```

## Usage

Recover public key from block

```bash
polycli ecrecover -r https://polygon-mumbai-bor.publicnode.com -b 45200775
> Recovering signer from block #45200775
> 0x5082F249cDb2f2c1eE035E4f423c46EA2daB3ab1

polycli ecrecover -r https://polygon-rpc.com
> Recovering signer from block #52888893
> 0xeEDBa2484aAF940f37cd3CD21a5D7C4A7DAfbfC0

polycli ecrecover -f block-52888893.json
> Recovering signer from block #52888893
> 0xeEDBa2484aAF940f37cd3CD21a5D7C4A7DAfbfC0

cat block-52888893.json | go run main.go ecrecover
> Recovering signer from block #52888893
> 0xeEDBa2484aAF940f37cd3CD21a5D7C4A7DAfbfC0
```

JSON Data passed in follows object definition [here](https://www.quicknode.com/docs/ethereum/eth_getBlockByNumber)

## Flags

```bash
  -b, --block-number uint   block number to check the extra data for (default: latest)
  -f, --file string         path to a file containing block information in JSON format
  -h, --help                help for ecrecover
  -r, --rpc-url string      the RPC endpoint URL
  -t, --tx string           transaction data in hex format
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
