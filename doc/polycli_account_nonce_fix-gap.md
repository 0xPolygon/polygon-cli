# `polycli account nonce fix-gap`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Send txs to fix the nonce gap for a specific account

```bash
polycli account nonce fix-gap [flags]
```

## Usage

This command will check the account current nonce against the max nonce found in the pool. In case of a nonce gap is found, txs will be sent to fill those gaps.

To fix a nonce gap, we can use a command like this:

```bash
polycli account nonce fix-gap \
    --rpc-url https://sepolia.drpc.org
    --private-key 0x32430699cd4f46ab2422f1df4ad6546811be20c9725544e99253a887e971f92b
```

## Flags

```bash
  -h, --help                 help for fix-gap
      --private-key string   the private key to be used when sending the txs to fix the nonce gap
```

The command also inherits flags from parent commands.

```bash
      --config string    config file (default is $HOME/.polygon-cli.yaml)
      --pretty-logs      Should logs be in pretty format or JSON (default true)
  -r, --rpc-url string   The RPC endpoint url (default "http://localhost:8545")
  -v, --verbosity int    0 - Silent
                         100 Panic
                         200 Fatal
                         300 Error
                         400 Warning
                         500 Info
                         600 Debug
                         700 Trace (default 500)
```

## See also

- [polycli account nonce](polycli_account_nonce.md) - 
