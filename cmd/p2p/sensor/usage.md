Running the sensor will do peer discovery and continue to watch for blocks and
transactions from those peers. This is useful for observing the network for
forks and reorgs without the need to run the entire full node infrastructure.

The sensor can persist data to various backends including Google Cloud Datastore
or JSON output. If no nodes.json file exists at the specified path, it will be
created automatically.

The bootnodes may change, so refer to the [Polygon Knowledge Layer][bootnodes]
if the sensor is not discovering peers.

## Metrics

The sensor exposes Prometheus metrics at `http://localhost:2112/metrics`
(configurable via `--prom-port`). For a complete list of available metrics, see
[polycli_p2p_sensor_metrics.md](polycli_p2p_sensor_metrics.md).

## Examples

### Mainnet

To run a Polygon Mainnet sensor, copy the `genesis.json` from [here][mainnet-genesis].

```bash
polycli p2p sensor nodes.json \
  --bootnodes "enode://b8f1cc9c5d4403703fbf377116469667d2b1823c0daf16b7250aa576bacf399e42c3930ccfcb02c5df6879565a2b8931335565f0e8d3f8e72385ecf4a4bf160a@3.36.224.80:30303,enode://8729e0c825f3d9cad382555f3e46dcff21af323e89025a0e6312df541f4a9e73abfa562d64906f5e59c51fe6f0501b3e61b07979606c56329c020ed739910759@54.194.245.5:30303" \
  --network-id 137 \
  --sensor-id "sensor" \
  --write-blocks=true \
  --write-block-events=true \
  --write-txs=true \
  --write-tx-events=true \
  --genesis-hash "0xa9c28ce2141b56c474f1dc504bee9b01eb1bd7d1a507580d5519d4437a97de1b" \
  --fork-id "22d523b2" \
  --rpc "https://polygon-rpc.com" \
  --discovery-dns "enrtree://AKUEZKN7PSKVNR65FZDHECMKOJQSGPARGTPPBI7WS2VUL4EGR6XPC@pos.polygon-peers.io" \
  --pprof \
  --verbosity 700 \
  --pretty-logs=true \
  --database "json"
```

### Amoy

To run a Polygon Amoy sensor, copy the `genesis.json` from [here][amoy-genesis].

```bash
polycli p2p sensor amoy-nodes.json \
  --bootnodes "enode://b8f1cc9c5d4403703fbf377116469667d2b1823c0daf16b7250aa576bacf399e42c3930ccfcb02c5df6879565a2b8931335565f0e8d3f8e72385ecf4a4bf160a@3.36.224.80:30303,enode://8729e0c825f3d9cad382555f3e46dcff21af323e89025a0e6312df541f4a9e73abfa562d64906f5e59c51fe6f0501b3e61b07979606c56329c020ed739910759@54.194.245.5:30303" \
  --network-id 80002 \
  --sensor-id "sensor-amoy" \
  --write-blocks=true \
  --write-block-events=true \
  --write-txs=true \
  --write-tx-events=true \
  --genesis-hash "0x7202b2b53c5a0836e773e319d18922cc756dd67432f9a1f65352b61f4406c697" \
  --fork-id "8b7e4175" \
  --rpc "https://rpc-amoy.polygon.technology" \
  --discovery-dns "enrtree://AKUEZKN7PSKVNR65FZDHECMKOJQSGPARGTPPBI7WS2VUL4EGR6XPC@amoy.polygon-peers.io" \
  --pprof \
  --verbosity 700 \
  --pretty-logs=true \
  --database "json"
```

[mainnet-genesis]: https://github.com/0xPolygon/bor/blob/master/builder/files/genesis-mainnet-v1.json
[amoy-genesis]: https://github.com/0xPolygon/bor/blob/master/builder/files/genesis-amoy.json
[bootnodes]: https://docs.polygon.technology/pos/reference/seed-and-bootnodes/
