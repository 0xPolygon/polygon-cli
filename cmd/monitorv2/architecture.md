# MonitorV2 Architecture

## Overview

We want to have a powerful tool like [atop](https://www.atoptool.nl/)
for monitoring the performance of EVM based blockchains. Here are the
features:

- Live monitoring of the tip of the chain showing recent blocks
- A block view that shows details of a particular block
- A raw block view that shows the pretty printed json view
- A transaction view that shows the details of a particular
  transaction
- A generic "search" feature that takes as input a particular block
  number, block hash, or transaction hash and will jump directly to
  that view (if it's found)
- A high level overview of the chain
  - Chain ID
  - Gas Prices
- A deep info view that prints a snapshot of all of the information
  that we have. This could be very useful for rollups where there are
  specific system contracts or addresses or extra RPCs that can
  provide additional context
- Event and function signature decoding powered by 4byte.directory
- Multiple renderers - The default is a TUI, but the data can also be
  rendered as a JSON stream
- Reorg detection - If a block hash changes, we can rewind and update
  the store. The depth to check for reorgs is configurable

## Design Principles

- Non-blocking design. I.e. when something is loading, the UI remains
  responsive
- Configurable - Many aspects of the tool can be configured

## Architecture Components

- RPC - these would be basic abstractions over the RPC requests and
  responses. The idea is to generalize a little over top of the
  standard RPC values
- Store - is an interface that abstracts away the details of how the
  blocks are stored locally (on-disk or in memory) and also the
  details of how the RPC requests are made.
  - We'll have a command line option (--store [memory, sqlite, etc])
    that would allow a store to be configured
- Indexer - is responsible for making requests to the store to
  populate the data
  - The approach here will be to use simple interfaces. E.g. we'll use
    polling rather than trying to use websockets and subscriptions
- Renderer - reads data from the store and renders it as a TUI or a
  JSON Stream or something else
- TUI - will be based of of tview because it's battle tested and
  trusted by important tools like k9s

## Data Flow

- The source of truth is the remote RPC
- The remote RPC is read via go-ethereum's RPC client
- The store will make the RPC requests and update it's own data
- The indexer will trigger the store to retrive data
- The renderer will read data from ths store and print

## Key Differences from Monitor V1

The main difference between V1 and V2 monitor is that the V2 monitor
is more decomposed and meant to be easier to maintain and upgrade.

