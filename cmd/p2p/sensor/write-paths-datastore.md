# Sensor Datastore write paths

How a devp2p message becomes entities in GCP Datastore, for
`--database=datastore` (`p2p/database/datastore.go`). Handlers are in
`p2p/protocol.go`.

See also: [Datastore schema](/cmd/p2p/sensor/schema-datastore.md),
[ClickHouse write paths](/cmd/p2p/sensor/write-paths-clickhouse.md).

The contrast with the ClickHouse backend is the whole point of this page: almost
every path here is a **read-modify-write** inside a `RunInTransaction` retry loop
(`MaxAttempts = 5`), and four separate paths contend on the *same* `blocks` entity.
That is what keeps the earliest-seen timestamps correct, and it is the shape the
ClickHouse schema was designed to stop emulating.

Writes are dispatched through `runAsync`, a semaphore sized by
`--max-db-concurrency` (default 10000). `Close` acquires every slot to guarantee no
write is in flight before closing the client.

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
        M7(["peer ticker, 2s"])
    end

    subgraph api["Database interface"]
        A1[WriteBlockEvents]
        A1b[WriteBlockHashFirstSeen]
        A2[WriteBlock]
        A3[WriteBlockHeaders]
        A4[WriteBlockBody]
        A5[WriteTransactions]
        A6[WriteTransactionEvents]
        A7[WritePeers]
    end

    subgraph rmw["Read-modify-write: RunInTransaction, up to 5 attempts"]
        R1["writeBlock
        Get then conditional Put"]
        R2["writeBlockHeader
        Get then conditional Put"]
        R3["writeBlockBody
        Get then conditional Put"]
        R4["writeBlockHashFirstSeen
        Get then conditional Put"]
    end

    subgraph dedup["Read-then-write"]
        D1["writeTransactions
        GetMulti then PutMulti
        skips earlier TimeFirstSeen"]
    end

    subgraph blind["Blind writes"]
        B1["writeEvents
        PutMulti"]
        B2["writeEvent
        Put"]
        B3["PutMulti"]
    end

    subgraph kinds["Datastore kinds"]
        K1[(blocks)]
        K2[(block_events)]
        K3[(transactions)]
        K4[(transaction_events)]
        K5[(peers)]
    end

    M1 --> A1
    M1 --> A1b
    M2 --> A2
    M3 --> A3
    M4 --> A4
    M5 --> A5
    M6 --> A6
    M7 --> A7

    A1 --> B1
    A1b --> R4
    A2 --> B2
    A2 --> R1
    A3 --> R2
    A4 --> R3
    A5 --> D1
    A6 --> B1
    A7 --> B3

    R1 -.->|"nested, inside the txn"| D1
    R1 -.->|"each uncle, nested"| R2
    R3 -.->|"nested, inside the txn"| D1
    R3 -.->|"each uncle, nested"| R2

    R1 --> K1
    R2 --> K1
    R3 --> K1
    R4 --> K1
    B1 --> K2
    B1 --> K4
    B2 --> K2
    D1 --> K3
    B3 --> K5
```

## The dotted edges

They are the part that does not show up in a method list, and they are worth
tracing. `writeBlock` and `writeBlockBody` call `writeTransactions` (itself a
`GetMulti` + `PutMulti`) and `writeBlockHeader` (itself a whole transaction) *from
inside their own transaction closure*. On contention the closure re-runs, so those
nested round trips re-run with it, up to five times.

## Per-method detail

| Method | Path | Shape |
| --- | --- | --- |
| `WriteBlock` | `writeEvent` + `writeBlock` | Blind event `Put`, then RMW on the block; conditionally nests transactions and uncle headers |
| `WriteBlockHeaders` | `writeBlockHeader` | RMW per header; skips the write if an equal-or-earlier `TimeFirstSeen` exists |
| `WriteBlockBody` | `writeBlockBody` | RMW; fills `Transactions` / `Uncles` key lists only if still nil |
| `WriteBlockEvents` | `writeEvents` | Batched `PutMulti`, no read. Announced heights are discarded |
| `WriteBlockHashFirstSeen` | `writeBlockHashFirstSeen` | RMW on `blocks` purely to keep the earliest hash-announce time |
| `WriteTransactions` | `writeTransactions` | `GetMulti` to skip txs already stored with an earlier or equal `TimeFirstSeen`, then `PutMulti` |
| `WriteTransactionEvents` | `writeEvents` | Batched `PutMulti`, no read |
| `WritePeers` | `PutMulti` | One entity per connected peer, every 2s |
| `HasBlock` | `Get` | Point lookup by block key |
| `NodeList` | query | Scans `block_events` ordered by `-Time`, collecting distinct `PeerId` |

Note that `WriteBlockHeaders` and `WriteBlockBody` deliberately write **no** events:
headers and bodies only arrive because the sensor asked for them, so the sighting is
recorded when the hash announcement comes in instead.

## Why this differs from ClickHouse

| | Datastore | ClickHouse |
| --- | --- | --- |
| Write shape | read-modify-write per entity | append-only, batched |
| Contention | 4 paths on the same `blocks` entity | none; writers never read |
| "Keep the earliest sighting" | `tx.Get` then conditional `tx.Put` | derived at read time from the sighting stream |
| `WriteBlockHashFirstSeen` | its own transaction on `blocks` | no-op |
| Backpressure | semaphore, blocks the caller | fixed buffers, drop-on-full |
| Peer id in events | enode URL (`peer.URLv4()`) | devp2p node id |
| Block ↔ tx link | key list on the block entity | `block_txs` table |
| Uncles | written as full `blocks` entities | hashes on `block_bodies` |
| Expiry | `TTL` field + a manual delete job | `TTL` clauses, whole-partition drops |
| Block numbers | strings, lexicographic ranges | `UInt64` |
| `NodeList` | scan `block_events` by `-Time` | `peers_current` rollup |
