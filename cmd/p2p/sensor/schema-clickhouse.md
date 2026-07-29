# Sensor ClickHouse schema

The tables the `--database=clickhouse` backend writes
(`p2p/database/clickhouse.go`). The DDL itself is not in this repo — it lives in
`clickhouse_schema.sql` in
[sensor-network-tools](https://github.com/0xPolygon/sensor-network-tools)
(canonical) and `polygon-infrastructure` (applied to the ClickHouse VM).

See also: [ClickHouse write paths](/cmd/p2p/sensor/write-paths-clickhouse.md),
[Datastore schema](/cmd/p2p/sensor/schema-datastore.md).

Block numbers collide across networks, so **each network gets its own database**
(`mainnet`, `amoy`) and there is no chain column. The DSN selects it:
`--clickhouse-dsn clickhouse://user:pass@host:9000/mainnet`.

## The split

Every kind of data is one of two things.

**Fact tables** are keyed by a hash, and every column is a pure function of that
hash. Two sensors that see the same block therefore produce byte-identical rows,
so deduplication is a storage optimisation rather than a correctness requirement
and readers never need `FINAL`.

**Observation tables** are append-only event streams — one row per (thing,
sensor, peer, time). Everything situational lives here and nowhere else: which
sensor, which peer, when, how the block was learned, and the announced total
difficulty.

The consequence worth internalising: a partially-known block cannot corrupt a
fully-known one. A header carries no transaction count, so instead of writing
`tx_count = 0` onto the block row, body facts live in their own hash-keyed table
that the header path never touches.

```mermaid
erDiagram
    blocks {
        UInt64 number PK
        String hash PK
        String parent_hash FK "-> blocks.hash, the parent block"
        DateTime block_time "header consensus timestamp"
        LowCardinality signer "precomputed ecrecover"
        LowCardinality coinbase
        UInt64 difficulty
        UInt64 gas_used
        UInt64 gas_limit
        UInt256 base_fee
        String uncle_hash
        String state_root
        String tx_root
        String receipt_root
        String logs_bloom "RAW BYTES not hex"
        String extra_data "RAW BYTES not hex"
        String mix_digest
        UInt64 nonce
        String withdrawals_root
        UInt64 blob_gas_used
        UInt64 excess_blob_gas
        String parent_beacon_root
    }

    block_bodies {
        String hash PK "keyed by hash alone"
        UInt32 tx_count
        UInt16 uncle_count
        Array uncles
        UInt32 size_bytes
    }

    block_txs {
        String block_hash PK
        UInt32 tx_index PK
        String tx_hash FK
        Date seen_date "partition and TTL key"
    }

    transactions {
        String hash PK
        String from_address
        String to_address "empty for contract creation"
        UInt256 value
        UInt64 gas
        UInt256 gas_price
        UInt256 gas_fee_cap
        UInt256 gas_tip_cap
        UInt64 nonce
        UInt8 tx_type
        UInt64 chain_id
        String input_selector "first 4 bytes, no calldata"
        UInt32 input_size
        UInt16 access_list_size
        UInt8 blob_count
        UInt8 auth_list_size
        Date seen_date "partition and TTL key"
    }

    block_events {
        UInt64 block_number PK "denormalised, leads sort key"
        String block_hash PK
        LowCardinality sensor_id PK
        LowCardinality node_id FK "devp2p node id"
        LowCardinality source "hash_announce new_block header header_backfill"
        DateTime64 seen_at PK "ms precision"
        UInt256 total_difficulty "from the announcement"
    }

    tx_events {
        String tx_hash PK
        LowCardinality sensor_id
        LowCardinality node_id FK
        LowCardinality source "hash_announce full_tx"
        DateTime64 seen_at PK
    }

    peer_snapshots {
        LowCardinality sensor_id PK
        String node_id PK
        String name "client version"
        String url "enode URL"
        Array caps
        DateTime64 seen_at PK
    }

    blocks         ||--o| block_bodies : "hash (when a body is seen)"
    blocks         ||--o{ block_txs : "block_hash"
    block_txs      }o--|| transactions : "tx_hash"
    blocks         ||--o{ block_events : "number + hash"
    transactions   ||--o{ tx_events : "hash"
    peer_snapshots ||--o{ block_events : "node_id"
    peer_snapshots ||--o{ tx_events : "node_id"
```

`blocks.parent_hash` points at another `blocks.hash`, forming the chain that reorg
and fork analysis walks. It is annotated on the column rather than drawn as a
relationship, because a self-reference renders as an easily-missed loop (or nothing
at all) depending on the mermaid version.

| Table | Engine | Sort key | Retention |
| --- | --- | --- | --- |
| `blocks` | `ReplacingMergeTree` (no version) | `(number, hash)` | forever |
| `block_bodies` | `ReplacingMergeTree` | `(hash)` | forever |
| `block_txs` | `ReplacingMergeTree` | `(block_hash, tx_index)` | forever |
| `transactions` | `ReplacingMergeTree` | `(hash)` | 14d |
| `block_events` | `MergeTree` | `(block_number, block_hash, sensor_id, seen_at)` | 14d |
| `tx_events` | `MergeTree` | `(tx_hash, seen_at)` | 14d |
| `peer_snapshots` | `MergeTree` | `(sensor_id, node_id, seen_at)` | 3d |

Things the diagram cannot carry:

- **`blocks` has no version column.** All rows for a hash are identical, so there
  is nothing to order them by. `block_time` is header-derived, which means a hash's
  duplicates always land in the same partition and the dedup key can actually
  collapse — a dedup key that spans partitions never merges.
- **`block_bodies` and `block_txs` are keyed by hash alone, with no `number`.** A
  body can arrive before, or without, its header (the sensor requests the two
  separately and they race), so the height is not reliably known on that path.
  Carrying `number` would mean writing `0` when unknown, reintroducing the exact
  partial-row problem the split removes. The height is one join away.
- **`peer_snapshots` → events is a join on `node_id`, not a foreign key.** It
  works only because both sides record the devp2p node id. Get this wrong and the
  join silently returns nothing.
- **`logs_bloom` and `extra_data` are raw bytes in a `String` column, not hex.**
  Readers that re-run ecrecover depend on it; hex-decoding `extra_data` corrupts it
  and silently breaks every signer-derived metric.
- **`seen_date` on `transactions` and `block_txs` is an ingest fact**, not a
  consensus one. On `transactions` it makes expiry a whole-partition drop; on
  `block_txs`, which is kept forever, it only partitions (monthly, so partitions
  don't accumulate). A row first seen either side of a partition boundary is written
  twice, but being content-addressed the copies are identical, so it costs bytes
  rather than correctness.

## Derived layer

Rollups are maintained by materialized views on insert. Every rollup column is a
`SimpleAggregateFunction`, which merges exactly and needs no `-State`/`-Merge`
combinators — readers apply the plain function under a `GROUP BY`.

The `v_*` views are the intended read surface, so latency arithmetic and fork
detection are defined once rather than in each consumer.

```mermaid
flowchart LR
    subgraph sensor["Written by the sensor"]
        BS[(block_events)]
        TS[(tx_events)]
        PS[(peer_snapshots)]
        BL[(blocks)]
        BB[(block_bodies)]
        TX[(transactions)]
    end

    subgraph jobs["Written by the analysis jobs"]
        RD[(reorg_detections)]
        BLM[(block_latency_metrics)]
    end

    subgraph rollups["Rollups: AggregatingMergeTree fed by MVs"]
        BSF[("block_events_first
        per block x sensor - 400d")]
        TSF[("tx_events_first
        per tx, fleet-wide - 90d")]
        BF[("block_forks
        per height - forever")]
        PC[("peers_current
        per sensor x peer - 30d")]
    end

    subgraph views["Read surface"]
        VBL[v_block_latency]
        VB[v_blocks]
        VF[v_forks]
        VSC[v_sensor_coverage]
        VP[v_peers]
        VTP[v_tx_propagation]
        VR[v_reorgs]
    end

    BS -->|MV| BSF
    TS -->|MV| TSF
    BL -->|MV| BF
    PS -->|MV| PC

    BSF --> VBL
    BL --> VBL
    BSF --> VSC
    BF --> VF
    PC --> VP
    TSF --> VTP
    TX --> VTP
    BL --> VB
    BB --> VB
    RD --> VR
```

`block_latency_metrics` is a report table, one row per job run per scope (`all` vs
`validated`). Its sort key `(scope, hours_analyzed, timestamp)` turns the job's
previous-run lookup into a reverse primary-index seek; on Datastore the same read
had to over-fetch 1000 rows and linearly scan for the matching pair.

Two traps in the derived layer:

- **`v_peers` is not optional convenience.** `peers_current` is a genuine upsert
  target, so unlike the fact tables its rows are *not* identical. Joining it raw
  fans out over unmerged snapshots — measured at 3x inflation on a freshly loaded
  database (74 rows covering 27 distinct peers). `v_peers` collapses them with
  `argMax`.
- **`v_block_latency.latency_ms` is `NULL` when the header has not been seen.** A
  hash announcement is routinely recorded before its header, and the join would
  otherwise supply `block_time = epoch` and yield a ~1.8e12 ms "latency" that
  destroys any percentile over the column. Filter on `has_header`. Note also that
  `latency_ms` is legitimately *negative* for many Bor blocks: the header timestamp
  is the proposer's slot time, which can be ahead of actual propagation.
