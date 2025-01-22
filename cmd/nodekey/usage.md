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

# Generate an ED25519 nodekey from a private key (in hex format).
$ polycli nodekey --private-key 2a4ae8c4c250917781d38d95dafbb0abe87ae2c9aea02ed7c7524685358e49c2 | jq
$ polycli nodekey --private-key 0x2a4ae8c4c250917781d38d95dafbb0abe87ae2c9aea02ed7c7524685358e49c2 | jq
{
  "PublicKey": "93e8717f46b146ebfb99159eb13a5d044c191998656c8b79007b16051bb1ff762d09884e43783d898dd47f6220af040206cabbd45c9a26bb278a522c3d538a1f",
  "PrivateKey": "2a4ae8c4c250917781d38d95dafbb0abe87ae2c9aea02ed7c7524685358e49c2",
  "ENR": "enode://93e8717f46b146ebfb99159eb13a5d044c191998656c8b79007b16051bb1ff762d09884e43783d898dd47f6220af040206cabbd45c9a26bb278a522c3d538a1f@0.0.0.0:30303?discport=0"
}
```
