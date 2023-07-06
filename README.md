We run a lot of different blockchain technologies. Different tools often
have inconsistent tooling and this makes automation and operations
painful. The goal of this codebase is to standardize some of our
commonly needed tools and provide interfaces and formats.

# Summary

- [Install](#install)
- [Features](doc/polycli.md)
- [Testing with Geth](#testing-with-geth)
- [Reference](#reference)

# Install

Requirements:

- [Go](https://go.dev/)
- make
- jq
- bc
- protoc (Only required for `make generate`)

To install, clone this repo and run:

```bash
$ make install
```

By default, the commands will be in `$HOME/go/bin/`, so for ease, we
recommend adding that path to your shell's startup file by adding the
following line:

```bash
export PATH="$HOME/go/bin:$PATH"
```

# Testing with Geth

While working on some of the Polygon CLI tools, we'll run geth in dev
mode in order to make sure the various functions work properly. First,
we'll startup geth.

```shell
# Geth
./build/bin/geth --dev --dev.period 2 --http --http.addr localhost --http.port 8545 --http.api admin,debug,web3,eth,txpool,personal,miner,net --verbosity 5 --rpc.gascap 50000000  --rpc.txfeecap 0 --miner.gaslimit  10 --miner.gasprice 1 --gpo.blocks 1 --gpo.percentile 1 --gpo.maxprice 10 --gpo.ignoreprice 2 --dev.gaslimit 50000000
```

In the logs, we'll see a line that says IPC endpoint opened:

```example
INFO [08-14|16:09:31.451] Starting peer-to-peer node               instance=Geth/v1.10.21-stable-67109427/darwin-arm64/go1.18.1
WARN [08-14|16:09:31.451] P2P server will be useless, neither dialing nor listening
DEBUG[08-14|16:09:31.452] IPCs registered                          namespaces=admin,debug,web3,eth,txpool,personal,clique,miner,net,engine
INFO [08-14|16:09:31.452] IPC endpoint opened                      url=/var/folders/zs/k8swqskj1t79cgnjh6yt0fqm0000gn/T/geth.ipc
INFO [08-14|16:09:31.452] Generated ephemeral JWT secret           secret=0xdfa5c30e07ef1041d15a2dbf0865386305330128b792d4a461cddb9bf38e416e
```

I'll usually then use that line to attach

```shell
./build/bin/geth attach /var/folders/zs/k8swqskj1t79cgnjh6yt0fqm0000gn/T/geth.ipc
```

After attaching to geth, we can fund the default load testing account
with some currency.

```shell
eth.coinbase==eth.accounts[0]
eth.sendTransaction({from: eth.coinbase, to: "0x85da99c8a7c2c95964c8efd687e95e632fc533d6", value: web3.toWei(5000, "ether")})
```

Then we can generate some load to make sure that there are some blocks
with transactions being created. `1337` is the chain id that's used in
local geth.

```shell
polycli loadtest --verbosity 700 --chain-id 1337 --concurrency 1 --requests 1000 --rate-limit 5 --mode c http://127.0.0.1:8545
```

Then we can monitor the chain:

```bash
polycli monitor http://127.0.0.1:8545
```

# Reference

Sending some value to the default load testing account

Listening for re-orgs

```shell
socat - UNIX-CONNECT:/var/folders/zs/k8swqskj1t79cgnjh6yt0fqm0000gn/T/geth.ipc
{"id": 1, "method": "eth_subscribe", "params": ["newHeads"]}
```

Useful RPCs when testing

```shell
curl -v -H 'Content-Type: application/json' -d '{"jsonrpc":"2.0", "id": 1, "method": "net_version", "params": []}' https://polygon-rpc.com
curl -v -H 'Content-Type: application/json' -d '{"jsonrpc":"2.0", "id": 1, "method": "eth_blockNumber", "params": []}' https://polygon-rpc.com
curl -v -H 'Content-Type: application/json' -d '{"jsonrpc":"2.0", "id": 1, "method": "eth_getBlockByNumber", "params": ["0x1DE8531", true]}' https://polygon-rpc.com
curl -v -H 'Content-Type: application/json' -d '{"jsonrpc":"2.0", "id": 1, "method": "clique_getSigner", "params": ["0x1DE8531", true]}' https://polygon-rpc.com
curl -v -H 'Content-Type: application/json' -d '{"jsonrpc":"2.0", "id": 1, "method": "eth_getBalance", "params": ["0x85da99c8a7c2c95964c8efd687e95e632fc533d6", "latest"]}' https://polygon-rpc.com
curl -v -H 'Content-Type: application/json' -d '{"jsonrpc":"2.0", "id": 1, "method": "eth_getCode", "params": ["0x79954f948079ee9ef1d15eff3e07ceaef7cdf3b4", "latest"]}' https://polygon-rpc.com


curl -v -H 'Content-Type: application/json' -d '{"jsonrpc":"2.0", "id": 1, "method": "txpool_inspect", "params": []}' http://localhost:8545
curl -v -H 'Content-Type: application/json' -d '{"jsonrpc":"2.0", "id": 1, "method": "txpool_status", "params": []}' http://localhost:8545
curl -v -H 'Content-Type: application/json' -d '{"jsonrpc":"2.0", "id": 1, "method": "eth_gasPrice", "params": []}' http://localhost:8545
curl -v -H 'Content-Type: application/json' -d '{"jsonrpc":"2.0", "id": 1, "method": "admin_peers", "params": []}' http://localhost:8545
```

```shell
websocat ws://34.208.176.205:9944
{"jsonrpc":"2.0", "id": 1, "method": "chain_subscribeNewHead", "params": []}
```
