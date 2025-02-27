# `polycli fix-nonce-gap`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Send txs to fix the nonce gap for a specific account

```bash
polycli fix-nonce-gap [flags]
```

## Usage

This command will check the account current nonce against the max nonce found in the pool. In case of a nonce gap is found, txs will be sent to fill those gaps.

To fix a nonce gap, we can use a command like this:

```bash
polycli fix-nonce-gap \
    --rpc-url https://sepolia.drpc.org
    --private-key 0x32430699cd4f46ab2422f1df4ad6546811be20c9725544e99253a887e971f92b
```

In case the RPC doesn't provide the `txpool_content` endpoint, the flag `--max-nonce` can be set to define the max nonce. The command will generate TXs from the current nonce up to the max nonce set.

```bash
polycli fix-nonce-gap \
    --rpc-url https://sepolia.drpc.org
    --private-key 0x32430699cd4f46ab2422f1df4ad6546811be20c9725544e99253a887e971f92b
    --max-nonce
```

By default, the command will skip TXs found in the pool, for example, let's assume the current nonce is 10, there is a TX with nonce 15 and 20 in the pool. When sending TXs to fill the gaps, the TXs 15 and 20 will be skipped. IN case you want to force these TXs to be replaced, you must provide the flag `--replace`.

```bash
polycli fix-nonce-gap \
    --rpc-url https://sepolia.drpc.org
    --private-key 0x32430699cd4f46ab2422f1df4ad6546811be20c9725544e99253a887e971f92b
    --replace
```
## Flags

```bash
  -h, --help                 help for fix-nonce-gap
      --max-nonce uint       when set, the max nonce will be this value instead of trying to get it from the pool
      --private-key string   the private key to be used when sending the txs to fix the nonce gap
      --replace              replace the existing txs in the pool
  -r, --rpc-url string       The RPC endpoint url (default "http://localhost:8545")
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
