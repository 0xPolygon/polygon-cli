# `polycli fork`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Take a forked block and walk up the chain to do analysis.

```bash
polycli fork blockhash url [flags]
```

## Usage

Occasionally, we'll want to analyze the details of a side chain to understand in detail who was proposing the blocks, what was the difficultly, and just generally get better understanding of block propagation.

```bash
$ polycli fork 0x053d84d5215684c8ae810a4729f7c9b54d65a80b128a27aeddcd7dc295a0cebd https://polygon-rpc.com
```

In order to use this, you'll need to have a blockhash of a block that was part of a fork / side chain. Once you have that, you can run `fork` against a node to get the details of the fork and the canonical chain.

## Flags

```bash
  -h, --help   help for fork
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
