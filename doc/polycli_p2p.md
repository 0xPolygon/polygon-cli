# `polycli p2p`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Set of commands related to devp2p.

## Usage

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

```bash
polycli p2p sensor nodes.json \
  --bootnodes "enode://b8f1cc9c5d4403703fbf377116469667d2b1823c0daf16b7250aa576bacf399e42c3930ccfcb02c5df6879565a2b8931335565f0e8d3f8e72385ecf4a4bf160a@3.36.224.80:30303,enode://8729e0c825f3d9cad382555f3e46dcff21af323e89025a0e6312df541f4a9e73abfa562d64906f5e59c51fe6f0501b3e61b07979606c56329c020ed739910759@54.194.245.5:30303" \
  --network-id 137 \
  --sensor-id sensor \
  --rpc "https://polygon-rpc.com"
```

#### Mumbai

```bash
polycli p2p sensor mumbai-nodes.json \
  --bootnodes "enode://bdcd4786a616a853b8a041f53496d853c68d99d54ff305615cd91c03cd56895e0a7f6e9f35dbf89131044e2114a9a782b792b5661e3aff07faf125a98606a071@43.200.206.40:30303,enode://209aaf7ed549cf4a5700fd833da25413f80a1248bd3aa7fe2a87203e3f7b236dd729579e5c8df61c97bf508281bae4969d6de76a7393bcbd04a0af70270333b3@54.216.248.9:30303" \
  --network-id 80001 \
  --sensor-id sensor \
  --rpc "https://polygon-mumbai-bor.publicnode.com" \
  --genesis-hash 0x7b66506a9ebdbf30d32b43c5f15a3b1216269a1ec3a75aa3182b86176a2b1ca7 \
  --fork-id 0c015a91
```

#### Amoy

```bash
polycli p2p sensor amoy-nodes.json \
  --bootnodes "enode://bce861be777e91b0a5a49d58a51e14f32f201b4c6c2d1fbea6c7a1f14756cbb3f931f3188d6b65de8b07b53ff28d03b6e366d09e56360d2124a9fc5a15a0913d@54.217.171.196:30303,enode://4a3dc0081a346d26a73d79dd88216a9402d2292318e2db9947dbc97ea9c4afb2498dc519c0af04420dc13a238c279062da0320181e7c1461216ce4513bfd40bf@13.251.184.185:30303" \
  --network-id 80002 \
  --sensor-id sensor \
  --rpc "https://rpc-amoy.polygon.technology" \
  --genesis-hash 0x7202b2b53c5a0836e773e319d18922cc756dd67432f9a1f65352b61f4406c697 \
  --fork-id b4f6ec4f
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

[bootnodes]: https://wiki.polygon.technology/docs/pos/operate/node/full-node-binaries/#configure-bor-seeds-mainnet

## Flags

```bash
  -h, --help   help for p2p
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
- [polycli p2p crawl](polycli_p2p_crawl.md) - Crawl a network on the devp2p layer and generate a nodes JSON file.

- [polycli p2p nodelist](polycli_p2p_nodelist.md) - Generate a node list to seed a node

- [polycli p2p ping](polycli_p2p_ping.md) - Ping node(s) and return the output.

- [polycli p2p query](polycli_p2p_query.md) - Query block header(s) from node and prints the output.

- [polycli p2p sensor](polycli_p2p_sensor.md) - Start a devp2p sensor that discovers other peers and will receive blocks and transactions.

