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

The `nodekey` command is still in progress, but the idea is to have a
simple command for generating a node key.

Most clients will generate this on the fly, but if we want to store
the key pair during an automated provisioning process, it's helpful to
have the output be structured.

```bash
# This will generate a secp256k1 key for devp2p protocol.
$ polycli nodekey

# Generate a networking keypair for libp2p.
$ polycli nodekey --protocol libp2p

# Generate a networking keypair for edge.
$ polycli nodekey --protocol libp2p --key-type secp256k1 --marshal-protobuf
```

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
