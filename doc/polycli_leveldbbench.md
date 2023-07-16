# `polycli leveldbbench`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Perform a level db benchmark

```bash
polycli leveldbbench [flags]
```

## Usage

This command is meant to give us a sense of the system level
performance for leveldb:

```bash
go run main.go leveldbbench --degree-of-parallelism 2 | jq '.' > result.json
```

In many cases, we'll want to emulate the performance characteristics
of `bor` or `geth`. This is the basic IO pattern when `bor` is in sync:

```text
Process Name = bor
     Kbytes              : count     distribution
         0 -> 1          : 0        |                                        |
         2 -> 3          : 0        |                                        |
         4 -> 7          : 10239    |****************                        |
         8 -> 15         : 25370    |****************************************|
        16 -> 31         : 7082     |***********                             |
        32 -> 63         : 1241     |*                                       |
        64 -> 127        : 58       |                                        |
       128 -> 255        : 11       |                                        |
```

This is the IO pattern when `bor` is getting in sync.

```text
Process Name = bor
     Kbytes              : count     distribution
         0 -> 1          : 0        |                                        |
         2 -> 3          : 0        |                                        |
         4 -> 7          : 23089    |*************                           |
         8 -> 15         : 70350    |****************************************|
        16 -> 31         : 11790    |******                                  |
        32 -> 63         : 1193     |                                        |
        64 -> 127        : 204      |                                        |
       128 -> 255        : 271      |                                        |
       256 -> 511        : 1381     |                                        |
```

This gives us a sense of the relative size of the IOPs. We'd also want
to get a sense of the read/write ratio. This is some sample data from
bor while syncing:

```text
12:48:08 loadavg: 5.86 6.22 7.13 16/451 56297

READS  WRITES R_Kb     W_Kb     PATH
307558 1277   4339783  30488    /var/lib/bor/data/bor/chaindata/

12:48:38 loadavg: 6.46 6.32 7.14 3/452 56298

READS  WRITES R_Kb     W_Kb     PATH
309904 946    4399349  26051    /var/lib/bor/data/bor/chaindata/

```

During the same period of time this is what the IO looks like from a
node that's in sync.

```text
12:48:05 loadavg: 1.55 1.85 2.03 18/416 88371

READS  WRITES R_Kb     W_Kb     PATH
124530 488    1437436  12165    /var/lib/bor/data/bor/chaindata/

12:48:35 loadavg: 4.14 2.44 2.22 1/416 88371

READS  WRITES R_Kb     W_Kb     PATH
81282  215    823530   4610     /var/lib/bor/data/bor/chaindata/

```

If we want to simulate `bor` behavior, we can leverage this data to
configure the leveldb benchmark tool.


| Syncing | Reads   | Writes | Read (kb) | Write (kb) | RW Ratio | kb/r | kb/w |
|---------|---------|--------|-----------|------------|----------|------|------|
| TRUE    | 307,558 |  1,277 | 4,339,783 | 30,488     |      241 | 14.1 | 23.9 |
| TRUE    | 309,904 |    946 | 7,399,349 | 26,051     |      328 | 23.9 | 27.5 |
| FALSE   | 124,530 |    488 | 1,437,436 | 12,165     |      255 | 11.5 | 24.9 |
| FALSE   | 51,282  |    215 | 823,530   | 4,610      |      239 | 16.1 | 21.4 |

The number of IOps while syncing is a lot higher. The only other
obvious difference is that the IOp size is a bit larger while syncing
as well.

- Syncing
  - Read Write Ratio - 275:1 
  - Small IOp - 10kb
  - Large IOp - 256kb
  - Small Large Ratio - 10:1
- Synced
  - Read Write Ratio - 250:1
  - Small IOp - 10kb
  - Larg IOp - 32kb
  - Small Large Ratio - 10:1

## Flags

```bash
      --cache-size int                the number of megabytes to use as our internal cache size (default 512)
      --degree-of-parallelism uint8   The number of concurrent iops we'll perform (default 1)
      --dont-fill-read-cache          if false, then random reads will be cached
      --handles int                   defines the capacity of the open files caching. Use -1 for zero, this has same effect as specifying NoCacher to OpenFilesCacher. (default 500)
  -h, --help                          help for leveldbbench
      --key-size uint                 The byte length of the keys that we'll use (default 8)
      --nil-read-opts                 if true we'll use nil read opt (this is what geth/bor does)
      --no-merge-write                allows disabling write merge
      --overwrite-count uint          the number of times to overwrite the data (default 5)
      --read-limit uint               the number of reads will attempt to complete in a given test (default 10000000)
      --read-strict                   if true the rand reads will be made in strict mode
      --sequential-reads              if true we'll perform reads sequentially
      --sequential-writes             if true we'll perform writes in somewhat sequential manner
      --size-kb-distribution string   the size distribution to use while testing (default "4-7:23089,8-15:70350,16-31:11790,32-63:1193,64-127:204,128-255:271,256-511:1381")
      --sync-writes                   sync each write
      --write-limit uint              The number of entries to write in the db (default 1000000)
      --write-zero                    if true, we'll write 0s rather than random data
```

The command also inherits flags from parent commands.

```bash
      --config string   config file (default is $HOME/.polygon-cli.yaml)
      --pretty-logs     Should logs be in pretty format or JSON (default true)
  -v, --verbosity int   0 - Silent
                        100 Fatal
                        200 Error
                        300 Warning
                        400 Info
                        500 Debug
                        600 Trace (default 400)
```

## See also

- [polycli](polycli.md) - A Swiss Army knife of blockchain tools.
