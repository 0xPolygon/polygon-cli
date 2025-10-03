# `polycli publish`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Publish transactions to the network with high-throughput

```bash
polycli publish [flags]
```

## Usage

This command publish transactions with high-throughput.

The command accepts a list of rlp hex encoded transactions that can be provided via a file, 
command line or stdin.

Internally it uses a worker pool strategy that can be dimensioned via flag, so it can be adjusted 
for optimal performance depending on the hardware available.

Since this command focus on high-throughput, please ensure the RPC will not rate-limit the requests.

Below are some example of how to use it

File: to use a file, set the file path using the --file flag
```bash
polycli publish --rpc-url https://sepolia.drpc.org --file /home/tclemos/txs
```

Command Line: to use command line args, set as many args you need when calling the command
```bash
polycli publish --rpc-url https://sepolia.drpc.org 0x000...001 0x000...002 0x000...003 0x000...004 ...
```

Stdin: to use std int, run the command without file or 0x args and then type one tx rlp per line
```bash
polycli publish --rpc-url https://sepolia.drpc.org

```

## Flags

```bash
  -c, --concurrency uint      number of txs to send concurrently (default: one at a time) (default 1)
      --file string           provide a filename with transactions to publish
  -h, --help                  help for publish
      --job-queue-size uint   number of jobs we can put in the job queue for workers to process (default 100)
      --rate-limit uint       rate limit in txs per second (default: no limit)
      --rpc-url string        RPC URL of network (default "http://localhost:8545")
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
