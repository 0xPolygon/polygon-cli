# Gas Manager

The Gas Manager package provides sophisticated control over transaction gas usage and pricing during load testing. It acts as a throttling and pricing mechanism that helps simulate realistic blockchain transaction patterns, including oscillating gas limits and dynamic gas pricing strategies.

## Overview

The Gas Manager package consists of three main components that work together:

1. **Gas Vault**: A budget-based gas limit controller
2. **Gas Provider**: Supplies gas budget to the vault based on configurable patterns
3. **Gas Pricer**: Determines transaction gas prices using different strategies

## Architecture

### Gas Vault

The `GasVault` acts as a semaphore/throttle mechanism that stores a gas budget and allows controlled spending:

- **`AddGas(uint64)`**: Adds gas to the available budget
- **`SpendOrWaitAvailableBudget(uint64)`**: Attempts to spend gas; blocks if insufficient budget is available
- **`GetAvailableBudget()`**: Returns current available gas budget

The vault prevents overspending by blocking transaction sending until sufficient budget is replenished by the gas provider.

### Gas Provider

The `GasProvider` interface defines how gas budget is supplied to the vault. The package includes an `OscillatingGasProvider` implementation that:

- Monitors blockchain for new blocks
- Uses configurable wave patterns to determine gas budget per block
- Automatically adds gas to the vault when new blocks are detected

### Gas Pricer

The `GasPricer` determines gas prices for transactions using pluggable strategies:

**Available Strategies:**

1. **Estimated** (default): Uses network-suggested gas price via `eth_gasPrice`
2. **Fixed**: Always returns a constant gas price
3. **Dynamic**: Cycles through a list of predefined gas prices with random variation

## Wave Patterns

The Gas Manager supports multiple oscillation wave patterns to simulate varying network conditions:

### Flat Wave

Maintains constant gas limit per block:
y = Target

### Sine Wave

Smooth oscillation following a sinusoidal pattern:
y = Amplitude × sin(2π/Period × x) + Target

### Square Wave

Alternates between high and low values:
y = Target ± Amplitude (alternates at half-period intervals)

### Triangle Wave

Linear increase and decrease:
y increases from (Target - Amplitude) to (Target + Amplitude)
  then decreases back linearly over the period

### Sawtooth Wave

Linear increase with sharp drop:
y increases linearly from (Target - Amplitude) to (Target + Amplitude)
  then resets sharply

## Usage in Loadtest

The `polycli loadtest` command integrates the Gas Manager through several flags:

### Gas Limit Control (Wave Configuration)

```bash
# Configure the oscillation wave pattern
--gas-manager-oscillation-wave string   # Wave type: flat, sine, square, triangle, sawtooth (default: "flat")
--gas-manager-target uint64             # Target gas limit baseline (default: 30000000)
--gas-manager-period uint64             # Period in blocks for wave oscillation (default: 1)
--gas-manager-amplitude uint64          # Amplitude of oscillation (default: 0)

Gas Price Control (Pricing Strategy)

# Select and configure pricing strategy
--gas-manager-price-strategy string                # Strategy: estimated, fixed, dynamic (default: "estimated")
--gas-manager-fixed-gas-price-wei uint64           # Fixed price in wei (default: 300000000)
--gas-manager-dynamic-gas-prices-wei string        # Comma-separated prices for dynamic strategy
                                                   # Use 0 for network-suggested price
                                                   # (default: "0,1000000,0,10000000,0,100000000")
--gas-manager-dynamic-gas-prices-variation float64 # Variation ±percentage for dynamic prices (default: 0.3)

Examples

Example 1: Constant Load with Fixed Gas Price

Simulate steady transaction flow with predictable gas price:

polycli loadtest \
  --rpc-url http://localhost:8545 \
  --gas-manager-oscillation-wave flat \
  --gas-manager-target 30000000 \
  --gas-manager-price-strategy fixed \
  --gas-manager-fixed-gas-price-wei 50000000000

Result: Sends up to 30M gas per block at constant 50 Gwei

Example 2: Sine Wave with Estimated Pricing

Simulate gradually increasing/decreasing load following network gas prices:

polycli loadtest \
  --rpc-url http://localhost:8545 \
  --gas-manager-oscillation-wave sine \
  --gas-manager-target 20000000 \
  --gas-manager-amplitude 10000000 \
  --gas-manager-period 100 \
  --gas-manager-price-strategy estimated

Result: Gas limit oscillates between 10M and 30M over 100 blocks using network gas prices

Example 3: Square Wave Traffic Pattern

Simulate bursty traffic with alternating high/low load:

polycli loadtest \
  --rpc-url http://localhost:8545 \
  --gas-manager-oscillation-wave square \
  --gas-manager-target 25000000 \
  --gas-manager-amplitude 15000000 \
  --gas-manager-period 20 \
  --gas-manager-price-strategy estimated

Result: Alternates between 10M gas (low) and 40M gas (high) every 10 blocks

Example 4: Dynamic Gas Prices with Variation

Simulate diverse user behavior with varying gas prices:

polycli loadtest \
  --rpc-url http://localhost:8545 \
  --gas-manager-oscillation-wave flat \
  --gas-manager-target 30000000 \
  --gas-manager-price-strategy dynamic \
  --gas-manager-dynamic-gas-prices-wei "1000000000,5000000000,10000000000,0" \
  --gas-manager-dynamic-gas-prices-variation 0.2

Result: Cycles through prices: 1 Gwei, 5 Gwei, 10 Gwei, and network-suggested, each with ±20% random variation

Example 5: Stress Test with Sawtooth Pattern

Gradually increase load then reset, simulating growing congestion:

polycli loadtest \
  --rpc-url http://localhost:8545 \
  --gas-manager-oscillation-wave sawtooth \
  --gas-manager-target 20000000 \
  --gas-manager-amplitude 15000000 \
  --gas-manager-period 50 \
  --gas-manager-price-strategy dynamic \
  --gas-manager-dynamic-gas-prices-wei "0,2000000000,0,5000000000,0,10000000000"

Result: Gas limit ramps from 5M to 35M over 50 blocks, then resets; alternates between network price and fixed prices

Visualization with tx-gas-chart

You can visualize the gas patterns generated by your loadtest using the tx-gas-chart command:

# Run loadtest with sine wave pattern
polycli loadtest \
  --rpc-url http://localhost:8545 \
  --gas-manager-oscillation-wave sine \
  --gas-manager-period 100 \
  --gas-manager-amplitude 10000000 \
  --target-address 0xYourAddress

# Generate chart to visualize results
polycli tx-gas-chart \
  --rpc-url http://localhost:8545 \
  --target-address 0xYourAddress \
  --output sine_wave_result.png

How It Works

Initialization Flow

1. Setup Gas Vault: Creates vault with zero initial budget
2. Configure Wave: Instantiates selected wave pattern with period, amplitude, and target
3. Create Gas Provider: Wraps wave pattern in OscillatingGasProvider
4. Start Provider:
  - Adds initial gas budget based on wave's starting Y value
  - Begins watching for new blocks
5. Setup Gas Pricer: Creates pricer with selected strategy

Transaction Flow

When a transaction is ready to send:

1. Gas Limit Decision:
  - Loadtest calls gasVault.SpendOrWaitAvailableBudget(gasLimit)
  - If vault has sufficient budget: deducts amount and continues
  - If vault lacks budget: blocks until provider adds more gas
2. Gas Price Decision:
  - Loadtest calls gasPricer.GetGasPrice()
  - Returns price based on selected strategy:
      - Estimated: Returns nil (loadtest queries network)
    - Fixed: Returns configured constant
    - Dynamic: Returns next price in sequence with variation
3. Block Progression:
  - When new block is detected, provider's onNewHeader callback fires
  - Wave advances to next position via MoveNext()
  - Provider adds gas equal to wave's new Y value to vault

Throttling Behavior

The Gas Manager creates realistic transaction throttling:

- High wave values → More gas budget → More concurrent transactions
- Low wave values → Less gas budget → Throttled transaction sending
- Zero budget → Complete blocking until next block

This simulates network congestion and varying block space availability.

Implementation Notes

- Gas Vault uses mutex for thread-safe budget management
- Gas Provider watches blocks with 1-second polling interval
- Wave patterns pre-compute all points during initialization
- Dynamic gas prices use atomic operations for thread-safe index management
- Zero gas price in dynamic strategy indicates "use network-suggested price"
- Variation in dynamic strategy applies random multiplier: price × (1 ± variation%)

Configuration Tips

For realistic network simulation:
- Use sine or triangle waves with period matching expected congestion cycles
- Use estimated pricing to follow real market conditions
- Set target to average block gas limit

For stress testing:
- Use square waves for sudden load changes
- Use sawtooth for gradually increasing stress
- Use dynamic pricing with high variation to test diverse scenarios

For consistent benchmarking:
- Use flat wave with zero amplitude
- Use fixed pricing for reproducible results
- Set target to desired constant load

See Also

- ../loadtestUsage.md
- ../txgaschart/usage.md
- Wave visualization examples: cmd/txgaschart/examples/