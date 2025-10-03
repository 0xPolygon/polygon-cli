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

To run a Polygon Mainnet sensor, copy the `genesis.json` from
[here][mainnet-genesis].

```bash
polycli p2p sensor nodes.json \
  --bootnodes "enode://b8f1cc9c5d4403703fbf377116469667d2b1823c0daf16b7250aa576bacf399e42c3930ccfcb02c5df6879565a2b8931335565f0e8d3f8e72385ecf4a4bf160a@3.36.224.80:30303,enode://8729e0c825f3d9cad382555f3e46dcff21af323e89025a0e6312df541f4a9e73abfa562d64906f5e59c51fe6f0501b3e61b07979606c56329c020ed739910759@54.194.245.5:30303" \
  --network-id 137 \
  --sensor-id sensor \
  --rpc "https://polygon-rpc.com"
```

#### Amoy

To run a Polygon Amoy sensor, copy the `genesis.json` from
[here][amoy-genesis].

```bash
polycli p2p sensor nodes.json \
  --bootnodes "enode://0ef8758cafc0063405f3f31fe22f2a3b566aa871bd7cd405e35954ec8aa7237c21e1ccc1f65f1b6099ab36db029362bc2fecf001a771b3d9803bbf1968508cef@35.197.249.21:30303,enode://c9c8c18cde48b41d46ced0c564496aef721a9b58f8724025a0b1f3f26f1b826f31786f890f8f8781e18b16dbb3c7bff805c7304d1273ac11630ed25a3f0dc41c@34.89.39.114:30303" \
  --network-id 80002 \
  --sensor-id "sensor" \
  --genesis-hash "0x7202b2b53c5a0836e773e319d18922cc756dd67432f9a1f65352b61f4406c697" \
  --fork-id "8b7e4175" 
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

[mainnet-genesis]: https://github.com/0xPolygon/bor/blob/master/builder/files/genesis-mainnet-v1.json
[amoy-genesis]: https://github.com/0xPolygon/bor/blob/master/builder/files/genesis-amoy.json
[bootnodes]: https://wiki.polygon.technology/docs/pos/operate/node/full-node-binaries/#configure-bor-seeds-mainnet

## Flags

```bash
  -h, --help   help for p2p
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
- [polycli p2p crawl](polycli_p2p_crawl.md) - Crawl a network on the devp2p layer and generate a nodes JSON file.

- [polycli p2p nodelist](polycli_p2p_nodelist.md) - Generate a node list to seed a node.

- [polycli p2p ping](polycli_p2p_ping.md) - Ping node(s) and return the output.

- [polycli p2p query](polycli_p2p_query.md) - Query block header(s) from node and prints the output.

- [polycli p2p sensor](polycli_p2p_sensor.md) - Start a devp2p sensor that discovers other peers and will receive blocks and transactions.

