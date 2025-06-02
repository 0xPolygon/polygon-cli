# Load Testing Guide: Using polycli with Kurtosis

## Overview

We frequently receive questions about TPS (transactions per second) benchmarks and how to stress test various types of transactions. Since `polycli` already includes a suite of common load generation scenarios, this guide aims to help
developers and integration teams perform realistic and reproducible load tests against Kurtosis-based environments.

This document focuses on:

- Practical examples using `polycli` and `Kurtosis`
- Standardized, reproducible workflows
- Reference TPS numbers for baseline performance (e.g., `CDK-Erigon`, `OP-Geth`)

## Goals

- Enable teams to run meaningful load tests on different blockchain stacks
- Provide a consistent, easy-to-follow methodology using `Kurtosis` and `polycli`
- Reduce confusion around setup and execution
- Offer baseline metrics for guidance and comparison

## Load Test Scenarios

The following load scenarios are supported out of the box in `polycli`:

1. **Native Token Transfers using EOAs**
2. **ERC20 Token Transfers**
3. **NFT Mints** (ERC721 and ERC1155)
4. **Uniswap-style Swaps**

Each scenario includes example commands and configuration tips to simulate realistic conditions.

---

## Prerequisites

- [Kurtosis CLI](https://docs.kurtosis.com/) installed
- [polycli](https://github.com/0xpolygon/polygon-cli?tab=readme-ov-file#install) installed and configured
- Funded accounts with private keys available to `polycli`
- Access to performance monitoring tools (optional but recommended, e.g., `Grafana` + `Prometheus`)

## Configuring Kurtosis for Load Testing

Launch a Kurtosis environment with your desired stack (e.g., CDK-Erigon or OP-Geth), for this example we will use the [kurtosis-cdk](https://github.com/0xPolygon/kurtosis-cdk) repository:

````bash
kurtosis run --enclave cdk github.com/0xPolygon/kurtosis-cdk
````

After the whole environment is started, set the env var `$rpc_url` with the RPC endpoint URL that will be used by
`polycli`.

```bash
export $rpc_url=$(kurtosis port print cdk cdk-erigon-rpc-001 rpc)
```

## Running Load Tests with polycli

`polycli` comes with a built-in command called `loadtest`, check more details about this command with:

```bash
polycli loadtest --help
```

Before starting using `polycli` to run load tests, it's important to mention some important flags to the `loadtest` 
command.

  - `--rpc-url`: defines the RPC URL that the test will call when calling the network RPC.
  - `--private-key`: defines the private key of the account used to send the transactions to the network.
  - `--mode`: defines what kind of load test you want to perform, EOA txs, ERC20, etc.
  - `--verbosity`: defines the log level that will be printed to the console: `0 to Silent`, `100 to Panic`, 
  `200 to Fatal`, `300 to Error`, `400 to Warning`, `500 to Info`, `600 to Debug`, `700 to Trace`, the default is `500`.  
  - `--requests`: defines the number of requests that will be sent to the network by each concurrent execution.
  - `--concurrency`: define the number of concurrent executions of the load test. For example, if `--requests` is set to
  `10` 
  and `--concurrency` is set to `2`, then 2 load test executions will start concurrently and each concurrent execution
  will send 10 requests, making 20 requests in total.
  - `--rate-limit`: defines the number of requests that can be sent per second to the network, this limits the requests
  sent across the concurrent executions

Here is a template that you can use to start writing you own load tests using `polycli`:

```bash
polycli loadtest \
  --rpc-url <RPC_URL>
  --private-key <PRIVATE-KEY>
  --mode <MODE> \
  --verbosity <LOG_LEVEL> \
  --concurrency <NUMBER_OF_CONCURRENT_REQUESTS> \
  --requests <NUMBER_OF_REQUESTS> \
  --rate-limit <MAX_REQUESTS_PER_SECOND>
```

As mentioned before, `polycli` test various types of transactions, please take a look for the flag `--mode` for the
different.

>In the steps below you will find examples on how to run `polycli` for specific scenarios.

### Native Token Transfers

To test `Native Token Transfers using EOAs`, you can use the following command

```bash
polycli loadtest --rpc-url $rpc_url --private-key $private_key --mode t
```

### ERC20 Token Transfers

```bash
polycli loadtest --rpc-url $rpc_url --private-key $private_key --mode 2
```

### NFT Mints (ERC721/1155)

```bash
polycli loadtest --rpc-url $rpc_url --private-key $private_key --mode 7
```

### Uniswap-style Swaps

```bash
polycli loadtest --rpc-url $rpc_url --private-key $private_key --mode v3
```

## Multi-account support

The loadtest command provided by `polycli` send all the transactions with the same sender by default, but some times
the test requires multiple accounts to be used in order to avoid the pool queue limits, the `polycli` supports
multi-account by setting the following flags:

sending-address-count

- `sending-address-count`: defines the number of accounts that will be created to send the test txs
- `address-funding-amount`: defines the amount that will be funded to each sending account that will be created
- `pre-fund-sending-addresses`: if set to true, the sending addresses will be funded at the start of the execution of
the load test, funding all the created accounts concurrently, otherwise all sending accounts will be funded only when
used for the first time.
- `keep-funded-amount`: by default, at the final of the load test execution, all the funds funded to the sending
accounts that were created during the execution will be returned to the funding account, which is the account set in the 
`--private-key` flag. In case you want the funded funds to remain in the created accounts, set it to `true`
- `sending-addresses-file`: defines the path to a file containing a list of private keys to be used by the load test.
Instead of creating new accounts, the test will use only the accounts defined in the file with the private keys as the
accounts that will send the test txs

### Proxy support

Generally network RPCs will have some sort of protection to avoid DoS attacks, for cases where you have to by pass the
IP check or any other change to the requests made by the `polycli` load test, you can specify the `--proxy` flag with
the proxy url.

## Example Benchmark Results

TBD

| Scenario             | CDK-Erigon | OP-Geth  | 
|----------------------|------------|----------|
| Native Transfers     | ~000 TPS   | ~000 TPS |
| ERC20 Transfers      | ~000 TPS   | ~000 TPS |
| NFT Mints (ERC721)   | ~000 TPS   | ~000 TPS |
| Uniswap Swaps        | ~000 TPS   | ~000 TPS |

> **Note**: These numbers are based on controlled testnet conditions and may vary depending on hardware, network conditions, and configuration.

## Monitoring and Visualization (Optional)

TBD

---

## Source

For more details, see the [`polycli` repository](https://github.com/0xPolygon/polygon-cli`).

---

_Last updated: 02-JUN-2025_
