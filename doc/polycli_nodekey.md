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

Generate an [ED25519](https://en.wikipedia.org/wiki/Curve25519) nodekey from a private key (in hex format).

```bash
polycli nodekey --private-key 2a4ae8c4c250917781d38d95dafbb0abe87ae2c9aea02ed7c7524685358e49c2 | jq
```

```json
{
  "PublicKey": "93e8717f46b146ebfb99159eb13a5d044c191998656c8b79007b16051bb1ff762d09884e43783d898dd47f6220af040206cabbd45c9a26bb278a522c3d538a1f",
  "PrivateKey": "2a4ae8c4c250917781d38d95dafbb0abe87ae2c9aea02ed7c7524685358e49c2",
  "ENR": "enode://93e8717f46b146ebfb99159eb13a5d044c191998656c8b79007b16051bb1ff762d09884e43783d898dd47f6220af040206cabbd45c9a26bb278a522c3d538a1f@0.0.0.0:30303?discport=0"
}
```

Generate an [Secp256k1](https://en.bitcoin.it/wiki/Secp256k1) nodekey from a private key (in hex format).

```bash
polycli nodekey --private-key 2a4ae8c4c250917781d38d95dafbb0abe87ae2c9aea02ed7c7524685358e49c2 --key-type secp256k1 | jq
```

```json
{
  "address": "99AA9FC116C1E5E741E9EC18BD1FD232130A5C44",
  "pub_key": {
    "type": "comet/PubKeySecp256k1Uncompressed",
    "value": "BBNYN0nMJsgo0Fp3kVW85PRGBNe7Gdz1XBFuTWQ7D8FnKRb2JYO3i3FK2UiA5+gTSxYu1K66KdYjQYP1mOkH09g="
  },
  "priv_key": {
    "type": "comet/PrivKeySecp256k1Uncompressed",
    "value": "OP72E0D7GEi/4VySpolVudLW7uPJm+6PWEtFKJmvp1M="
  }
}
```

## Flags

```bash
  -f, --file string          a file with the private nodekey (in hex format)
  -h, --help                 help for nodekey
  -i, --ip string            the IP to be associated with this address (default "0.0.0.0")
      --key-type string      ed25519|secp256k1|ecdsa|rsa (default "ed25519")
  -m, --marshal-protobuf     marshal libp2p key to protobuf format instead of raw
      --private-key string   use the provided private key (in hex format)
      --protocol string      devp2p|libp2p|pex|seed-libp2p (default "devp2p")
  -S, --seed uint            a numeric seed value (default 271828)
  -s, --sign                 sign the node record
  -t, --tcp int              the TCP port to be associated with this address (default 30303)
  -u, --udp int              the UDP port to be associated with this address
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
