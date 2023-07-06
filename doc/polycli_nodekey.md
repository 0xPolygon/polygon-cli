# `polycli nodekey`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Generate node keys for different blockchain clients and protocols.

```bash
polycli nodekey [flags]
```

## Usage

Generate node keys for different blockchain clients and protocols.
## Flags

```bash
  -f, --file string        A file with the private nodekey in hex format
  -h, --help               help for nodekey
  -i, --ip string          The IP to be associated with this address (default "0.0.0.0")
      --key-type string    ed25519|secp256k1|ecdsa|rsa (default "ed25519")
  -m, --marshal-protobuf   If true the libp2p key will be marshaled to protobuf format rather than raw
      --protocol string    devp2p|libp2p|pex|seed-libp2p (default "devp2p")
  -S, --seed uint          A numeric seed value (default 271828)
  -s, --sign               Should the node record be signed?
  -t, --tcp int            The tcp Port to be associated with this address (default 30303)
  -u, --udp int            The udp Port to be associated with this address
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
