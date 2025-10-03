# `polycli p2p query`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Query block header(s) from node and prints the output.

```bash
polycli p2p query [enode/enr] [flags]
```

## Usage

Query header of single block or range of blocks given a single enode/enr.
	
This command will initially establish a handshake and exchange status message
from the peer. Then, it will query the node for block(s) given the start block
and the amount of blocks to query and print the results.
## Flags

```bash
  -a, --amount uint        amount of blocks to query (default 1)
  -h, --help               help for query
  -s, --start-block uint   block number to start querying from
```

The command also inherits flags from parent commands.

```bash
      --config string   config file (default is $HOME/.polygon-cli.yaml)
      --pretty-logs     output logs in pretty format instead of JSON (default true)
  -v, --verbosity int   0 - silent
                        100 panic
                        200 fatal
                        300 error
                        400 warning
                        500 info
                        600 debug
                        700 trace (default 500)
```

## See also

- [polycli p2p](polycli_p2p.md) - Set of commands related to devp2p.
