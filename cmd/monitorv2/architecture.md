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
- TUI support for mouse clicking
- Sortable columns for adjusting the view based on the currently
  indexed blocks

## Design Principles

- Non-blocking design. I.e. when something is loading, the UI remains
  responsive
- Configurable - Many aspects of the tool can be configured
- Unified data access - Single interface for all chain-related data
- Intelligent caching - Different TTL strategies for different data types
- Capability awareness - Graceful handling of unsupported RPC methods

## Architecture Components

### ChainStore (Previously Store)

The **ChainStore** is a unified interface that abstracts away the details of how 
blockchain data is accessed and cached. It consolidates both block data and 
chain metadata into a single coherent interface.

Key features:
- **Unified Interface**: Replaces the previous BlockStore with comprehensive chain data access
- **Intelligent Caching**: TTL-based caching with different strategies:
  - Static data (ChainID): Cached indefinitely
  - Semi-static (Safe/Finalized blocks): 5 minutes TTL
  - Frequent (Gas price, Fee history): 30 seconds TTL
  - Very frequent (Pending/Queued txs): 5 seconds TTL
  - Block-aligned (Base fee): Cached per block
- **Capability Detection**: Automatically tests RPC methods and gracefully handles unsupported endpoints
- **Configurable TTL**: Different cache expiration strategies for different data types

Store implementations:
- **PassthroughStore**: Direct RPC passthrough with intelligent caching (current implementation)
- **Future stores**: SQLite, memory-based, hybrid approaches

### Indexer

The **Indexer** is responsible for fetching blockchain data and coordinating between
the ChainStore and Renderers. It provides a clean abstraction layer and manages
the flow of data through the system.

Key responsibilities:
- **Block fetching**: Polls for new blocks and publishes them to renderers
- **Gap detection**: Identifies and handles missing blocks
- **Parallel processing**: Concurrent block fetching for improved performance
- **Delegation**: Provides unified access to all ChainStore methods
- **Channel-based communication**: Non-blocking data flow to renderers

### Renderer

The **Renderer** interface supports multiple output formats:

#### TviewRenderer (TUI)
- **Page-based architecture**: Home, Block Detail, Transaction Detail, Info, Help pages
- **Comprehensive status pane**: Real-time chain information including:
  - Current timestamp (full date-time for screenshots)
  - RPC endpoint URL (for network identification)
  - Chain ID, gas prices, pending transactions
  - Safe/finalized block numbers, base fees
- **Enhanced blocks table**: 9-column display with:
  - Block number, relative time, truncated hash
  - Transaction count, human-readable size
  - Gas usage with percentage and formatted numbers
  - State root hash
- **Keyboard shortcuts**: Intuitive navigation with h/i/q/Esc keys
- **Real-time updates**: 10-second refresh cycle for chain info

#### JSONRenderer
- Structured JSON output for automation and scripting

### Data Types and Formatting

The system includes sophisticated data formatting for human readability:
- **Relative timestamps**: "5m ago", "2h ago", "3d ago"
- **Byte sizes**: "1.5KB", "2.3MB" with proper units
- **Number formatting**: Comma-separated thousands (e.g., "13,402,300")
- **Gas percentages**: Utilization display (e.g., "44.7%")
- **Hash truncation**: Smart abbreviation with ellipsis

## Current Implementation Status

### Completed Features
- Unified ChainStore architecture replacing BlockStore
- Intelligent caching system with TTL strategies
- RPC capability detection and graceful error handling
- Page-based TUI with keyboard navigation
- Comprehensive 9-column blocks table
- Real-time status pane with chain information
- Multiple renderer support (JSON, TUI)
- Human-readable data formatting utilities

### In Progress
- Chain info update channels in indexer
- Block detail page implementation
- Transaction detail page implementation

### Future Features
- Search functionality (block/tx lookup)
- Raw JSON block view
- Mouse support for TUI
- Reorg detection and handling
- Event/function signature decoding via 4byte.directory
- Deep info view for rollup-specific data
- Additional store implementations (SQLite, memory)

## Data Flow

```
┌─────────┐    ┌──────────────┐    ┌─────────┐    ┌──────────┐
│   RPC   │ -> │  ChainStore  │ -> │ Indexer │ -> │ Renderer │
│ Network │    │   (Cached)   │    │         │    │   (TUI)  │
└─────────┘    └──────────────┘    └─────────┘    └──────────┘
```

1. **RPC Network**: Source of truth for all blockchain data
2. **ChainStore**: Unified data access with intelligent caching and capability detection
3. **Indexer**: Coordinates data flow and provides abstraction layer
4. **Renderer**: Presents data in user-friendly format (TUI, JSON, etc.)

### Caching Strategy

The ChainStore implements a sophisticated caching strategy based on data characteristics:

- **Static data** (ChainID): Never expires, fetched once
- **Semi-static data** (Safe/Finalized blocks): 5-minute TTL
- **Frequent data** (Gas prices, Fee history): 30-second TTL  
- **Very frequent data** (Pending transactions): 5-second TTL
- **Block-aligned data** (Base fees): Cached per block number

## Key Differences from Monitor V1

The main differences between V1 and V2 monitor:

1. **Unified Data Access**: Single ChainStore interface for all blockchain data
2. **Intelligent Caching**: TTL-based strategies reduce RPC load significantly
3. **Capability Awareness**: Graceful handling of different RPC endpoint capabilities
4. **Enhanced UI**: Rich 9-column blocks table with comprehensive status information
5. **Better Architecture**: Clear separation between data access, coordination, and presentation
6. **Future-Ready**: Extensible design for additional features and store implementations

The V2 monitor is more decomposed, maintainable, and provides significantly richer
information display while being more efficient with RPC usage through intelligent caching.
