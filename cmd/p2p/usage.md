### Ping

Pinging a peer is useful to determine information about the peer and retrieving
the `Hello` and `Status` messages. By default, it will listen to the peer after
the status exchange for blocks and transactions. To disable this behavior, set
the `--listen` flag.

```bash
polycli p2p ping <enode/enr or nodes.json file>
```

### Sensor

Running the sensor will do peer discovery and continue to watch for blocks and
transactions from those peers. This is useful for observing the network for
forks and reorgs without the need to run the entire full node infrastructure.

The bootnodes may change, so refer to the [Wiki][bootnodes] if the sensor is not
discovering peers.

#### Mainnet

To run a Polygon Mainnet sensor, copy the `genesis.json` from
[here][mainnet-genesis].

```bash
polycli p2p sensor nodes.json \
  --bootnodes "enode://b8f1cc9c5d4403703fbf377116469667d2b1823c0daf16b7250aa576bacf399e42c3930ccfcb02c5df6879565a2b8931335565f0e8d3f8e72385ecf4a4bf160a@3.36.224.80:30303,enode://8729e0c825f3d9cad382555f3e46dcff21af323e89025a0e6312df541f4a9e73abfa562d64906f5e59c51fe6f0501b3e61b07979606c56329c020ed739910759@54.194.245.5:30303" \
  --network-id 137 \
  --sensor-id sensor \
  --rpc "https://polygon-rpc.com"
```

#### Mumbai

To run a Polygon Mumbai sensor, copy the `genesis.json` from
[here][mumbai-genesis].

```bash
polycli p2p sensor nodes.json \
  --bootnodes "enode://bdcd4786a616a853b8a041f53496d853c68d99d54ff305615cd91c03cd56895e0a7f6e9f35dbf89131044e2114a9a782b792b5661e3aff07faf125a98606a071@43.200.206.40:30303,enode://209aaf7ed549cf4a5700fd833da25413f80a1248bd3aa7fe2a87203e3f7b236dd729579e5c8df61c97bf508281bae4969d6de76a7393bcbd04a0af70270333b3@54.216.248.9:30303" \
  --network-id 80001 \
  --sensor-id sensor \
  --rpc "https://polygon-mumbai-bor.publicnode.com" \
  --genesis-hash 0x7b66506a9ebdbf30d32b43c5f15a3b1216269a1ec3a75aa3182b86176a2b1ca7 \
  --fork-id 0c015a91
```

### Crawl

To crawl the network for nodes and write the output json to a file. This will
not engage in block or transaction propagation, but it can give a good indicator
of network size, and the output json can be used to quick start other nodes.

```bash
polycli p2p crawl nodes.json \
  --bootnodes "enode://0cb82b395094ee4a2915e9714894627de9ed8498fb881cec6db7c65e8b9a5bd7f2f25cc84e71e89d0947e51c76e85d0847de848c7782b13c0255247a6758178c@44.232.55.71:30303,enode://88116f4295f5a31538ae409e4d44ad40d22e44ee9342869e7d68bdec55b0f83c1530355ce8b41fbec0928a7d75a5745d528450d30aec92066ab6ba1ee351d710@159.203.9.164:30303,enode://4be7248c3a12c5f95d4ef5fff37f7c44ad1072fdb59701b2e5987c5f3846ef448ce7eabc941c5575b13db0fb016552c1fa5cca0dda1a8008cf6d63874c0f3eb7@3.93.224.197:30303,enode://32dd20eaf75513cf84ffc9940972ab17a62e88ea753b0780ea5eca9f40f9254064dacb99508337043d944c2a41b561a17deaad45c53ea0be02663e55e6a302b2@3.212.183.151:30303" \
  --network-id 137
```

[mainnet-genesis]: https://github.com/maticnetwork/bor/blob/master/builder/files/genesis-mainnet-v1.json
[mumbai-genesis]: https://github.com/maticnetwork/bor/blob/master/builder/files/genesis-testnet-v4.json
[bootnodes]: https://wiki.polygon.technology/docs/pos/operate/node/full-node-binaries/#configure-bor-seeds-mainnet
