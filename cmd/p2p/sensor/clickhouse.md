# Sensor ClickHouse data model

The tables the `--database=clickhouse` backend writes, and how a devp2p message
becomes rows in them. The writer is `p2p/database/clickhouse.go`; the handlers that
drive it are in `p2p/protocol.go`.

The DDL itself is not in this repo — it lives in
[sensor-network-tools `clickhouse/schema.sql`](https://github.com/0xPolygon/sensor-network-tools/blob/main/clickhouse/schema.sql)
(canonical) and is mirrored into `polygon-infrastructure`, which applies it to the
ClickHouse VM.

See also: [Datastore data model](/cmd/p2p/sensor/datastore.md).

Block numbers collide across networks, so **each network gets its own database**
(`mainnet`, `amoy`) and there is no chain column. The DSN selects it:
`--clickhouse-dsn clickhouse://user:pass@host:9000/mainnet`.

## Schema

### The split

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
    }

    block_bodies {
        String hash PK "keyed by hash alone"
        UInt32 tx_count
        UInt16 uncle_count
        Array uncles
    }

    block_txs {
        String block_hash PK
        UInt32 tx_index PK
        String tx_hash FK
        Date seen_date "version column, TTL clock"
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
        Date seen_date "version column, TTL clock"
    }

    block_events {
        UInt64 block_number PK "denormalised, leads sort key"
        String block_hash PK
        LowCardinality sensor_id PK
        LowCardinality node_id FK "devp2p node id"
        LowCardinality source "hash_announce new_block header header_backfill body"
        DateTime64 seen_at PK "ms precision"
        UInt256 total_difficulty "announced value, 0 otherwise"
    }

    tx_events {
        String tx_hash PK
        LowCardinality sensor_id
        LowCardinality node_id FK
        LowCardinality source "hash_announce full_tx"
        DateTime64 seen_at PK
    }

    peers {
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
    peers ||--o{ block_events : "node_id"
    peers ||--o{ tx_events : "node_id"
```

`blocks.parent_hash` points at another `blocks.hash`, forming the chain that reorg
and fork analysis walks. It is annotated on the column rather than drawn as a
relationship, because a self-reference renders as an easily-missed loop (or nothing
at all) depending on the mermaid version.

| Table                   | Engine                            | Sort key                                         | Retention |
| ----------------------- | --------------------------------- | ------------------------------------------------ | --------- |
| `blocks`                | `ReplacingMergeTree` (no version) | `(number, hash)`                                 | forever   |
| `block_bodies`          | `ReplacingMergeTree`              | `(hash)`                                         | forever   |
| `block_txs`             | `ReplacingMergeTree`              | `(block_hash, tx_index)`                         | forever   |
| `transactions`          | `ReplacingMergeTree`              | `(hash)`                                         | 14d       |
| `block_events`          | `MergeTree`                       | `(block_number, block_hash, sensor_id, seen_at)` | 14d       |
| `tx_events`             | `MergeTree`                       | `(tx_hash, seen_at)`                             | 14d       |
| `peers`                 | `MergeTree`                       | `(sensor_id, node_id, seen_at)`                  | 14d       |
| `block_events_first`    | `AggregatingMergeTree`            | `(block_number, block_hash, sensor_id, source)`  | 14d       |
| `tx_events_first`       | `AggregatingMergeTree`            | `(tx_hash)`                                      | 14d       |
| `block_forks`           | `AggregatingMergeTree`            | `(number)`                                       | forever   |
| `peers_current`         | `ReplacingMergeTree(last_seen)`   | `(sensor_id, node_id)`                           | 14d       |
| `reorg_detections`      | `MergeTree`                       | `(start_block, depth, detected_at)`              | forever   |
| `block_latency_metrics` | `MergeTree`                       | `(scope, hours_analyzed, timestamp)`             | forever   |

**Retention is 14 days or forever, never anything in between.** Observations and
anything derived from them expire at 14 days; the content-addressed facts and the
analysis-job tables are kept. So the growth question is only about the `forever`
group — `block_txs` dominates it at roughly 92 GiB/year on mainnet.

Things the diagram cannot carry:

- **`blocks` has no version column.** All rows for a hash are identical, so there
  is nothing to order them by. `block_time` is header-derived, which means a hash's
  duplicates always land in the same partition and the dedup key can actually
  collapse — a dedup key that spans partitions never merges.
- **No encoded block size is stored.** It is only knowable from a full `NewBlock`,
  not from a body delivered on its own, so it is not a function of the hash. Storing
  it let two sensors write differing rows for one block, and the engine picked
  arbitrarily — measured, the `0` won and the real size was discarded.
- **`block_bodies` and `block_txs` are keyed by hash alone, with no `number`.** A
  body can arrive before, or without, its header (the sensor requests the two
  separately and they race), so the height is not reliably known on that path.
  Carrying `number` would mean writing `0` when unknown, reintroducing the exact
  partial-row problem the split removes. The height is one join away.
- **`peers` → events is a join on `node_id`, not a foreign key.** It
  works only because both sides record the devp2p node id. Get this wrong and the
  join silently returns nothing.
- **`logs_bloom` and `extra_data` are raw bytes in a `String` column, not hex.**
  Readers that re-run ecrecover depend on it; hex-decoding `extra_data` corrupts it
  and silently breaks every signer-derived metric.
- **The post-Shanghai/Cancun header fields are deliberately absent.** `mix_digest`
  and `nonce` are all-zero on Bor but stored, because clique's `encodeSigHeader`
  includes them unconditionally and ecrecover needs them (as it needs `base_fee`).
  `withdrawals_root`, `blob_gas_used`, `excess_blob_gas` and `parent_beacon_root` are
  not stored: clique _panics_ if any is non-nil, so they can never take part in the
  seal hash, and they are absent from mainnet and amoy headers. Add one back with
  `ALTER TABLE ADD COLUMN` if that ever changes.
- **A fact table's partition key must be a function of its dedup key**, because
  `ReplacingMergeTree` merges only within a partition. `blocks` satisfies this via
  header-derived `block_time`. `transactions` and `block_txs` originally partitioned
  on `seen_date`, which is ingest-derived, and so stranded duplicates in separate
  partitions permanently: 10.07% of `transactions` rows survived a full
  `OPTIMIZE FINAL`, 115k hashes spanning partitions. A transaction has no intrinsic
  timestamp — unlike a block, whose header supplies one — so `seen_date` could not
  be made key-derived, and both tables now bucket on their own key
  (`cityHash64(hash) % 16`).
- **`seen_date` is the version column on both tables.** It is the one column there
  that is not a function of the key, so without a version the surviving row's value
  was arbitrary. As a version it resolves to the latest sighting, which also means a
  re-announced pending transaction's 14 days restart from when it was last seen.
- **Fact rows and provenance events are gated separately, per write path.** A
  `blocks` row needs `--write-blocks`; an event needs either block-event flag. Mixing
  them is the defect that recurred three times — `new_block`/`header`/`body` behind
  `--write-block-events` alone, then `full_tx` behind `--write-tx-events` alone, then
  `header`/`header_backfill` behind `--write-blocks` because `WriteBlockHeaders`
  returned early before reaching the event check. Headers are requested regardless of
  `--write-blocks`, so that last one produced the events and discarded them.
- **`ttl_only_drop_parts` requires a partition no coarser than the TTL.** It
  suppresses row-level expiry and drops a part only once every row in it has
  expired, so a partition spanning longer than the TTL pins expired rows. Both
  `*_first` rollups partitioned monthly against a 14-day TTL: once background
  merges combined a month's inserts into one part, a row from the 1st survived
  until the 31st's row expired — 44 days, sawtoothing with the calendar. Every
  table using the setting now partitions daily; the two that cannot
  (`peers_current`, hash-bucketed `transactions`) do not use it.
- **Dedup happens on merge, so duplicates are transient, not absent.** Correct
  partitioning makes them converge; it does not stop an unmerged part from holding
  two rows for a key. A reader that must not double-count still needs `FINAL` or
  `LIMIT 1 BY` — what it no longer needs is to compensate for inflation that never
  converges.

### Derived layer

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
        PS[(peers)]
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
        per block x sensor x source - 14d")]
        TSF[("tx_events_first
        per tx, fleet-wide - 14d")]
        BF[("block_forks
        per height - forever")]
        PC[("peers_current
        per sensor x peer - 14d")]
    end

    subgraph views["Read surface"]
        VBL[v_block_latency]
        VBP[v_block_provenance]
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

    BSF -->|"propagation sources only"| VBL
    BL --> VBL
    BSF -->|"every source"| VBP
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
  target, so unlike the fact tables its rows are _not_ identical. Joining it raw
  fans out over unmerged snapshots — measured at 3x inflation on a freshly loaded
  database (74 rows covering 27 distinct peers). `v_peers` collapses them with
  `argMax`.
- **Timestamp names carry their scope.** `sensor_first_seen` / `sensor_last_seen` are
  one sensor's earliest and latest; plain `first_seen` / `last_seen` are across every
  sensor; `first_seen_latency_ms` is how far behind the earliest sensor a given sensor
  was. `v_block_latency` exposes the first two side by side, so a per-sensor value is
  never mistaken for a fleet-wide one. `tx_events_first` keeps a plain `first_seen`
  because its grain is per transaction, which is already across sensors.
- **`v_block_latency.latency_ms` is `NULL` when the header has not been seen.** A
  hash announcement is routinely recorded before its header, and the join would
  otherwise supply `block_time = epoch` and yield a ~1.8e12 ms "latency" that
  destroys any percentile over the column. Filter on `has_header`. Note also that
  `latency_ms` is legitimately _negative_ for many Bor blocks: the header timestamp
  is the proposer's slot time, which can be ahead of actual propagation.

## Write paths

This backend is **pure append**. No `Write*` method reads a row back in order to
modify it; the only read on the path is `HasBlock`, used to decide whether to
backfill a missing parent.

Every write enqueues onto a per-table `rowBatcher` that flushes on a 1s tick or a
size threshold, retrying a failed batch up to 3 times. `add` is **non-blocking with
drop-on-full**, so a slow or unreachable database can never stall the sensor's hot
path — it drops rows and logs the count instead.

```mermaid
flowchart LR
    subgraph msgs["devp2p / timers"]
        M1[NewBlockHashes]
        M2[NewBlock]
        M3[BlockHeaders]
        M4[BlockBodies]
        M5["Transactions and
        PooledTransactions"]
        M6[NewPooledTransactionHashes]
        M7(["peer snapshot ticker
        --peer-snapshot-interval, 30s"])
    end

    subgraph handlers["p2p/protocol.go"]
        H1[handleNewBlockHashes]
        H2[handleNewBlock]
        H3[handleBlockHeaders]
        H4[handleBlockBodies]
        H5[processTransactions]
        H6[handleNewPooledTransactionHashes]
        H7[getParentBlock]
    end

    subgraph api["Database interface"]
        A1["WriteBlockEvents
        carries block heights"]
        A1b["WriteBlockHashFirstSeen
        NO-OP"]
        A2[WriteBlock]
        A3[WriteBlockHeaders]
        A4[WriteBlockBody]
        A5[WriteTransactions]
        A6[WriteTransactionEvents]
        A7[WritePeers]
        A8[HasBlock]
    end

    subgraph tables["ClickHouse, via rowBatcher"]
        T1[(blocks)]
        T2[(block_bodies)]
        T3[(block_txs)]
        T4[(block_events)]
        T5[(transactions)]
        T6[(tx_events)]
        T7[(peers)]
    end

    M1 --> H1
    M2 --> H2
    M3 --> H3
    M4 --> H4
    M5 --> H5
    M6 --> H6
    M7 --> A7

    H1 --> A1
    H1 --> A1b
    H2 --> A2
    H3 --> A3
    H3 --> H7
    H4 --> A4
    H5 --> A5
    H6 --> A6
    H7 --> A8

    A1 -->|"source=hash_announce"| T4
    A2 --> T1
    A2 -->|"tx_count, uncles, size"| T2
    A2 --> T3
    A2 -->|"source=new_block
    + total_difficulty"| T4
    A2 --> T5
    A3 --> T1
    A3 -->|"source=header or
    header_backfill"| T4
    A4 --> T2
    A4 --> T3
    A4 --> T5
    A5 --> T5
    A5 -->|"source=full_tx"| T6
    A6 -->|"source=hash_announce"| T6
    A7 --> T7
    A8 -.->|"point read, bloom index"| T1
```

### Per-method detail

| Method                    | Writes                                                                | Notes                                                 |
| ------------------------- | --------------------------------------------------------------------- | ----------------------------------------------------- |
| `WriteBlock`              | `blocks`, `block_bodies`, `block_txs`, `block_events`, `transactions` | The only path with the whole block                    |
| `WriteBlockHeaders`       | `blocks`, `block_events`                                              | `isParent` picks `header_backfill`; gated separately  |
| `WriteBlockBody`          | `block_bodies`, `block_txs`, `transactions`                           | No event: bodies arrive only when requested           |
| `WriteBlockEvents`        | `block_events`                                                        | Takes `[]BlockAnnouncement`, so heights reach the row |
| `WriteBlockHashFirstSeen` | nothing                                                               | Derived instead, see below                            |
| `WriteTransactions`       | `transactions`, `tx_events`                                           | Event under either tx-event flag                      |
| `WriteTransactionEvents`  | `tx_events`                                                           |                                                       |
| `WritePeers`              | `peers`                                                               | Own ticker, not the 2s metrics tick                   |
| `HasBlock`                | —                                                                     | Reads `blocks` by hash, once per new header           |

### Five things worth reading off this

**`WriteBlockHashFirstSeen` is a no-op.** Earliest first-seen is derived from the
event stream by the `block_events_first` materialized view, so there is nothing to
stamp on the block. The interface method stays because the Datastore backend does
need it.

Because that rollup is keyed by `source`, it replaces both Datastore column pairs at
once — `TimeFirstSeenHash`/`SensorFirstSeenHash` is `source = 'hash_announce'` and
`TimeFirstSeen`/`SensorFirstSeen` is `source = 'header'`. It expires at 14 days like
its source, so it is a pre-aggregation for query speed rather than a retention play;
`v_block_provenance` presents it. **Any latency read must restrict to the propagation
sources** (`hash_announce`, `new_block`), which `v_block_latency` does: the
sensor-requested sources carry no peer and their timestamps say when it chose to
fetch, not how fast the block arrived.

**`blocks` is written by two paths, and that is safe.** Both `WriteBlock` and
`WriteBlockHeaders` write a _complete_ header row — every column is a pure function
of the header — so ordering between them cannot matter. The fields a header cannot
carry (`tx_count`, `uncle_count`, `uncles`) go to `block_bodies`,
which the header path never touches.

This is the defect the schema redesign fixed. Previously the header path wrote
`tx_count = 0` onto the block row, and because the engine kept the row with the
newest version column, a header arriving _after_ the full block — routine during
parent backfill — permanently replaced the real counts with zeros. There is now a
regression test for exactly that ordering
(`TestClickHouseHeaderDoesNotClobberBody`).

**`total_difficulty` is written twice, to two different kinds of table.** _Which
peer announced which value_ is an observation, so it rides on the `block_events`
row (header and hash-announce events write `0` there). _The value itself_ is a
property of the block, so it also goes to `block_total_difficulty` — hash-keyed and
kept forever, written only by the `NewBlock` path and never as `0`.

It needs its own table rather than a `blocks` column for the same reason
`block_bodies` does: the header path does not know it and would write `0`, letting a
header clobber a real value. It used to be read out of `block_events_first`, which
worked only while that rollup outlived the raw stream; once retention was
normalised both were 14 days and `v_blocks.total_difficulty` returned `0` for every
older block — the same value that means "no peer ever announced it to us". Now
absence carries that meaning and the column is `Nullable`, so there is no sentinel.
`TestClickHouseTotalDifficultySurvivesEventExpiry` covers it.

**Block heights have to reach `WriteBlockEvents`.** `block_events` leads its sort
key with `block_number`, which is what lets a reader fetch a whole range in one
query instead of one point lookup per hash. `NewBlockHashes` carries the numbers, so
the interface takes `[]database.BlockAnnouncement` rather than bare hashes —
and since `p2p.NewBlockHashesPacket` is defined as a slice of exactly that type, a
decoded packet is handed to the backend with no copy or conversion.

**Peer identity is the devp2p node id**, not the enode URL, on both event
streams — so events join to `peers` / `peers_current`. The Datastore
backend uses `peer.URLv4()` here, which is why peers and events cannot be joined on
that backend.

### Write volume

Batch sizes are per table (`chBlockBatch` and friends): 5,000 blocks, 5,000 bodies,
20,000 block-txs, 50,000 block events, 20,000 transactions, 50,000 tx events,
2,000 peers.

`tx_events` is the volume driver, and only its `hash_announce` rows are: every peer
that announces a hash produces one, so the count scales with `--max-peers`. The
`--write-tx-events` / `--write-first-tx-event` pair is what bounds it — whether the
table receives every announcement or only first events. `full_tx` rows are recorded
under either flag, because a delivered body is roughly 2 rows per transaction per
sensor rather than a per-peer stream. Blocks work the same way, via
`recordsBlockEvents` / `recordsTxEvents`.

Peer snapshots are their own cadence. The 2s tick in `sensor.go` still drives the
Prometheus gauge and the local peer file, but persisting up to `--max-peers` rows
every 2s is a large amount of near-duplicate data to answer one question ("who is
connected now"), so the database write runs on `--peer-snapshot-interval` (30s) and
reads go through `peers_current`.
