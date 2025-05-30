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
  -c, --concurrency uint      Number of txs to send concurrently. Default is one request at a time. (default 1)
      --file string           Provide a filename with transactions to publish
  -h, --help                  help for publish
      --job-queue-size uint   Number of jobs we can put in the job queue for workers to process. (default 100)
      --rate-limit uint       Rate limit in txs per second. Default is no rate limit.
      --rpc-url string        The RPC URL of the network (default "http://localhost:8545")
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
