The `nodekey` command is still in progress, but the idea is to have a
simple command for generating a node key. Most clients will generate
this on the fly, but if we want to store the key pair during an
automated provisioning process, it's helpful to have the output be
structured

```shell
# this will generate a secp256k1 key for devp2p protocol
polycli nodekey

# generate a networking keypair for avail
polycli nodekey --protocol libp2p

# generate a networking keypair for edge
polycli nodekey --protocol libp2p --key-type secp256k1 --marshal-protobuf
```
