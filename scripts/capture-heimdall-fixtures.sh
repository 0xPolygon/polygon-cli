#!/usr/bin/env bash
# Capture JSON fixtures from a live Heimdall v2 node for use in unit
# tests under internal/heimdall/client/testdata/.
#
# Idempotent: re-running overwrites existing fixtures with the latest
# response. A per-endpoint success case is the hard requirement; we do
# not attempt to synthesise error cases that require specific broken
# node state.
#
# Usage:
#   scripts/capture-heimdall-fixtures.sh
#
# Environment:
#   HEIMDALL_REST_URL   REST gateway (default http://172.19.0.2:1317)
#   HEIMDALL_RPC_URL    CometBFT RPC (default http://172.19.0.2:26657)
#   FIXTURE_DIR         override output dir (default internal/heimdall/client/testdata)

set -euo pipefail

REST="${HEIMDALL_REST_URL:-http://172.19.0.2:1317}"
RPC="${HEIMDALL_RPC_URL:-http://172.19.0.2:26657}"

repo_root="$(git -C "$(dirname "$0")" rev-parse --show-toplevel)"
OUT="${FIXTURE_DIR:-$repo_root/internal/heimdall/client/testdata}"

mkdir -p "$OUT/rest" "$OUT/rpc"

need() {
  command -v "$1" >/dev/null 2>&1 || { echo "missing dep: $1" >&2; exit 1; }
}
need curl
need jq

fetch_rest() {
  local name="$1" path="$2"
  local out="$OUT/rest/$name.json"
  local tmp
  tmp="$(mktemp)"
  if ! curl -sS --max-time 10 --fail "$REST$path" -o "$tmp"; then
    echo "FAIL $path" >&2
    rm -f "$tmp"
    return 0
  fi
  if ! jq . "$tmp" > "$out"; then
    echo "BAD JSON from $path" >&2
    mv "$tmp" "$out" # keep raw bytes for debugging
    return 0
  fi
  rm -f "$tmp"
  echo "ok rest/$name.json"
}

fetch_rpc() {
  local name="$1" method="$2"
  shift 2
  local params="{}"
  if [ "$#" -gt 0 ]; then
    params="$1"
  fi
  local out="$OUT/rpc/$name.json"
  local payload
  payload="$(jq -nc --arg m "$method" --argjson p "$params" '{jsonrpc:"2.0", id:1, method:$m, params:$p}')"
  local tmp
  tmp="$(mktemp)"
  if ! curl -sS --max-time 10 --fail -H 'Content-Type: application/json' -d "$payload" "$RPC" -o "$tmp"; then
    echo "FAIL rpc $method" >&2
    rm -f "$tmp"
    return 0
  fi
  if ! jq . "$tmp" > "$out"; then
    mv "$tmp" "$out"
    return 0
  fi
  rm -f "$tmp"
  echo "ok rpc/$name.json"
}

# -----------------------------------------------------------------------------
# REST captures — mirror the command taxonomy in HEIMDALLCAST_REQUIREMENTS.md §3.2
# -----------------------------------------------------------------------------

# Checkpoints
fetch_rest checkpoints_params              /checkpoints/params
fetch_rest checkpoints_count               /checkpoints/count
fetch_rest checkpoints_latest              /checkpoints/latest
fetch_rest checkpoints_buffer              /checkpoints/buffer
fetch_rest checkpoints_last_no_ack         /checkpoints/last-no-ack
fetch_rest checkpoints_overview            /checkpoints/overview
fetch_rest checkpoints_list                '/checkpoints/list?pagination.limit=3&pagination.reverse=true'

# Resolve an id we can look up.
CP_COUNT="$(curl -sS --max-time 5 --fail "$REST/checkpoints/count" | jq -r '.ack_count // .count // empty' || true)"
if [ -n "${CP_COUNT:-}" ]; then
  fetch_rest "checkpoints_by_id"           "/checkpoints/${CP_COUNT}"
fi

# Spans
fetch_rest bor_params                      /bor/params
fetch_rest bor_spans_latest                /bor/spans/latest
fetch_rest bor_spans_list                  '/bor/spans/list?pagination.limit=3&pagination.reverse=true'
SPAN_LATEST="$(curl -sS --max-time 5 --fail "$REST/bor/spans/latest" | jq -r '.span.id // empty' || true)"
if [ -n "${SPAN_LATEST:-}" ]; then
  fetch_rest bor_spans_by_id               "/bor/spans/${SPAN_LATEST}"
  fetch_rest bor_spans_seed                "/bor/spans/seed/${SPAN_LATEST}"
fi
fetch_rest bor_producer_votes              /bor/producer-votes
fetch_rest bor_validator_performance_score /bor/validator-performance-score

# Milestones
fetch_rest milestones_params               /milestones/params
fetch_rest milestones_count                /milestones/count
fetch_rest milestones_latest               /milestones/latest
MS_COUNT="$(curl -sS --max-time 5 --fail "$REST/milestones/count" | jq -r '.count // empty' || true)"
if [ -n "${MS_COUNT:-}" ] && [ "${MS_COUNT}" -gt 0 ] 2>/dev/null; then
  fetch_rest milestones_by_number          "/milestones/${MS_COUNT}"
fi

# Validators
fetch_rest stake_validators_set            /stake/validators-set
fetch_rest stake_total_power               /stake/total-power
fetch_rest stake_proposers_current         /stake/proposers/current
fetch_rest stake_proposers_n               /stake/proposers/5
V_ID="$(curl -sS --max-time 5 --fail "$REST/stake/validators-set" | jq -r '.validator_set.validators[0].val_id // empty' || true)"
V_SIGNER="$(curl -sS --max-time 5 --fail "$REST/stake/validators-set" | jq -r '.validator_set.validators[0].signer // empty' || true)"
if [ -n "${V_ID:-}" ]; then
  fetch_rest stake_validator_by_id         "/stake/validator/${V_ID}"
fi
if [ -n "${V_SIGNER:-}" ]; then
  fetch_rest stake_signer                  "/stake/signer/${V_SIGNER}"
  fetch_rest stake_validator_status        "/stake/validator-status/${V_SIGNER}"
fi

# State-sync / clerk
fetch_rest clerk_event_records_count       /clerk/event-records/count
fetch_rest clerk_event_records_list        '/clerk/event-records/list?page=1&limit=3'
fetch_rest clerk_event_records_latest_id   /clerk/event-records/latest-id
SS_COUNT="$(curl -sS --max-time 5 --fail "$REST/clerk/event-records/count" | jq -r '.count // empty' || true)"
if [ -n "${SS_COUNT:-}" ] && [ "${SS_COUNT}" -gt 0 ] 2>/dev/null; then
  fetch_rest clerk_event_record_by_id      "/clerk/event-records/${SS_COUNT}"
fi

# Topup
fetch_rest topup_dividend_account_root     /topup/dividend-account-root
if [ -n "${V_SIGNER:-}" ]; then
  fetch_rest topup_dividend_account        "/topup/dividend-account/${V_SIGNER}"
  fetch_rest topup_account_proof           "/topup/account-proof/${V_SIGNER}"
fi

# Chain manager
fetch_rest chainmanager_params             /chainmanager/params

# Cosmos SDK — accounts / balances (pick a known signer if available)
if [ -n "${V_SIGNER:-}" ]; then
  fetch_rest cosmos_auth_account           "/cosmos/auth/v1beta1/accounts/${V_SIGNER}"
  fetch_rest cosmos_bank_balance_pol       "/cosmos/bank/v1beta1/balances/${V_SIGNER}/by_denom?denom=pol"
fi

# -----------------------------------------------------------------------------
# CometBFT JSON-RPC captures
# -----------------------------------------------------------------------------

fetch_rpc status                status
fetch_rpc abci_info              abci_info
fetch_rpc health                 health
fetch_rpc net_info               net_info
fetch_rpc num_unconfirmed_txs    num_unconfirmed_txs
fetch_rpc unconfirmed_txs        unconfirmed_txs
fetch_rpc consensus_state        consensus_state

HEIGHT="$(curl -sS --max-time 5 --fail -H 'Content-Type: application/json' \
  -d '{"jsonrpc":"2.0","id":1,"method":"status","params":{}}' "$RPC" \
  | jq -r '.result.sync_info.latest_block_height // empty' || true)"

if [ -n "${HEIGHT:-}" ]; then
  fetch_rpc block_latest        block        "$(jq -nc --arg h "$HEIGHT" '{height:$h}')"
  fetch_rpc commit_latest       commit       "$(jq -nc --arg h "$HEIGHT" '{height:$h}')"
  fetch_rpc validators          validators   "$(jq -nc --arg h "$HEIGHT" '{height:$h}')"
fi

# Pick the first tx hash in the latest block, if any.
TX_HASH="$(curl -sS --max-time 5 --fail -H 'Content-Type: application/json' \
  -d '{"jsonrpc":"2.0","id":1,"method":"block","params":{}}' "$RPC" \
  | jq -r '.result.block.data.txs // [] | .[0] // empty' || true)"
if [ -n "${TX_HASH:-}" ]; then
  # The RPC `/tx?hash=...` wants a 0x-prefixed SHA256 of the raw tx.
  TX_HASH_HEX="$(printf '%s' "$TX_HASH" | base64 -d 2>/dev/null | sha256sum | awk '{print "0x"toupper($1)}')"
  if [ -n "${TX_HASH_HEX:-}" ]; then
    fetch_rpc tx                tx           "$(jq -nc --arg h "$TX_HASH_HEX" '{hash:$h}')"
  fi
fi

fetch_rpc tx_search_msgs        tx_search    '{"query":"tx.height>0","prove":false,"page":"1","per_page":"3","order_by":"desc"}'

echo "captured into $OUT"
