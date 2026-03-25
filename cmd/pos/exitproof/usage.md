# `polycli pos exit-proof`

Generate the ABI-encoded exit payload required to withdraw tokens from Polygon PoS to Ethereum.

## Usage

```bash
polycli pos exit-proof [flags]
```

## Description

When a user burns tokens on the Polygon PoS sidechain, they must later submit a cryptographic proof to the Ethereum `WithdrawManager` contract to finalise the withdrawal. This command constructs that proof by:

1. Fetching the burn transaction receipt from L2 to determine the block number and transaction index.
2. Fetching all receipts for the burn block via `eth_getBlockReceipts` and reconstructing the receipts Merkle Patricia Trie (MPT).
3. Generating an MPT proof (sibling nodes from root to leaf) for the burn receipt.
4. Fetching the checkpoint by its ID from the `RootChain` contract on L1.
5. Fetching the block headers for the checkpoint range and building a binary Merkle proof that the burn block hash is included.
6. ABI-encoding all of the above into the payload expected by `startExitWithBurntTokens(bytes)`.

## Flags

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--l1-rpc-url` | no | `https://eth.llamarpc.com` | L1 (Ethereum) RPC URL |
| `--l2-rpc-url` | no | `https://polygon-rpc.com` | L2 (Polygon PoS) RPC URL |
| `--tx-hash` | yes | — | burn transaction hash on L2 |
| `--root-chain-address` | no | `0x86E4Dc95c7FBdBf52e33D563BbDB00823894C287` | `RootChain` contract address on L1 |
| `--checkpoint-id` | yes | — | checkpoint ID covering the burn block (visible on Polygonscan under the Checkpoint tab) |
| `--log-index` | no | `0` | index of the burn log within the receipt; `0` works for most ERC20 withdrawals — increase if the token emits extra logs before the burn event |

## Output

Writes `0x`-prefixed ABI-encoded payload bytes to stdout. All progress and diagnostic messages go to stderr so the output can be captured directly in a shell variable.

## Example

```bash
# Step 1: burn tokens on L2 (e.g., call ERC20 withdraw on the ChildToken contract)
# burn_tx_hash=0x...

# Step 2: wait for the checkpoint to be submitted (~30 min on mainnet)

# Step 3: generate the exit proof
payload=$(polycli pos exit-proof \
  --l1-rpc-url "${L1_RPC_URL}" \
  --l2-rpc-url "${L2_RPC_URL}" \
  --tx-hash "${burn_tx_hash}" \
  --checkpoint-id "${CHECKPOINT_ID}")

# Step 4: start the exit on L1
cast send \
  --rpc-url "${L1_RPC_URL}" \
  --private-key "${PRIVATE_KEY}" \
  "${WITHDRAW_MANAGER_PROXY_ADDRESS}" \
  "startExitWithBurntTokens(bytes)" \
  "${payload}"

# Step 5: process the exit on L1
cast send \
  --rpc-url "${L1_RPC_URL}" \
  --private-key "${PRIVATE_KEY}" \
  "${WITHDRAW_MANAGER_PROXY_ADDRESS}" \
  "processExits(address)" \
  "${POL_TOKEN_ADDRESS}"
```

## Contract Addresses (Mainnet)

| Contract | Address |
|----------|---------|
| `RootChain` | `0x86E4Dc95c7FBdBf52e33D563BbDB00823894C287` |
| `WithdrawManagerProxy` | `0x2A88696e0fFA76bAA1338F2C74497cC013495922` |
