# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

polygon-cli is a Swiss Army knife of blockchain tools for building, testing, and running blockchain applications. It's a collection of utilities primarily focused on Polygon/Ethereum ecosystems.

## Architecture Overview

The project follows a command-based architecture using Cobra framework:

- **Main Entry**: `main.go` simply calls `cmd.Execute()`
- **Command Structure**: Each command is organized in its own package under `cmd/` (e.g., `cmd/loadtest/`, `cmd/monitor/`)
- **Bindings**: Go bindings for smart contracts are in `bindings/` (generated from Solidity contracts in `contracts/`)
- **Utilities**: Common utilities are in `util/` package

Key architectural patterns:
- Commands are self-contained with their own usage documentation (`.md` files)
- Heavy use of code generation for documentation, protobuf, contract bindings, and RPC types
- Docker-based generation workflows to ensure consistency

## Common Development Commands

### Building and Installation
```bash
# Build the binary
make build

# Install to ~/go/bin/
make install

# Cross-compile for different platforms
make cross        # With CGO
make simplecross  # Without CGO
```

### Testing
```bash
# Run all tests with coverage
make test

# Run specific test
go test -v ./cmd/loadtest/...

# Run load test against local node
make geth          # Start local geth node
make geth-loadtest # Fund account and run load test
```

### Code Quality
```bash
# Run all linters (includes tidy, vet, golangci-lint)
make lint

# Individual linter commands
make tidy          # Clean up go.mod
make fmt           # Format code
make vet           # Run go vet and shadow
make golangci-lint # Run golangci-lint
```

### Code Generation
```bash
# Generate everything (docs, proto, bindings, etc.)
make gen

# Individual generation commands
make gen-doc              # Generate CLI documentation
make gen-proto            # Generate protobuf stubs
make gen-go-bindings      # Generate contract bindings
make gen-load-test-modes  # Generate loadtest mode strings
make gen-json-rpc-types   # Generate JSON RPC types
```

### Contract Development
```bash
# Work with smart contracts
cd contracts/
make build            # Build contracts with Foundry
make gen-go-bindings  # Generate Go bindings
```

## Adding New Features

When adding a new command:
1. Create a new package under `cmd/your-command/`
2. Add the command to `cmd/root.go` in the `NewPolycliCommand()` function
3. Create a usage documentation file (e.g., `yourCommandUsage.md`)
4. Run `make gen-doc` to update the main documentation
5. If adding a new loadtest mode, run `make gen-load-test-modes` after using stringer

## CI/CD Considerations

The CI pipeline (`/.github/workflows/ci.yml`) runs:
- Linting (golangci-lint, shadow)
- Tests
- Generation checks (ensures all generated files are up-to-date)
- Load tests against both geth and anvil

Always run `make gen` before committing if you've changed anything that affects code generation.

## Key Dependencies

- Go 1.24+ required
- Foundry (for smart contract compilation)
- Docker (for generation tasks)
- Additional tools: jq, bc, protoc (for development)

## Environment Configuration

The tool supports configuration via:
- CLI flags (highest priority)
- Environment variables
- Config file (`~/.polygon-cli.yaml`)
- Viper is used for configuration management

## Logging

- Use zerolog for structured, performant logging throughout the project

## Development Guidelines
- Use conventional commit messages

## Development Memories
- Use `make build` to build polycli

## Code Quality Checklist

**CRITICAL**: Before writing any code, systematically check these categories to avoid rework:

### 1. Security
- **HTML/Template Injection**: Always use `html.EscapeString()` for any data interpolated into HTML, even if currently from trusted sources
- **Input Validation**: Validate all user inputs at boundaries (flags, API inputs)
- **SQL Injection**: Use parameterized queries, never string concatenation
- **Command Injection**: Never pass user input directly to shell commands
- **Question to ask**: "What data is untrusted? Where does it flow? Is it escaped/validated at every output point?"

### 2. Resource Management & Performance
- **Goroutine Lifecycle**: Every goroutine must have a clear termination condition via context cancellation
- **Timer Cleanup**: Use `time.NewTimer()` + `defer timer.Stop()`, never `time.After()` in select statements (causes goroutine leaks)
- **Channel Buffers**: Use small fixed buffers (e.g., `concurrency*2`), never proportional to total dataset size
- **Memory Allocation**: Consider behavior with 10x, 100x, 1000x expected input
- **Question to ask**: "How does every goroutine, timer, and channel clean up on cancellation? What's the memory footprint at scale?"

### 3. Context Propagation
- **Never create root contexts**: Always thread `context.Context` through call chains; never use `context.Background()` in the middle of operations
- **Cancellation Flow**: Context should flow through every I/O operation, long-running task, and goroutine
- **Timeout Management**: Create child contexts with `context.WithTimeout(parentCtx, duration)`, not `context.WithTimeout(context.Background(), duration)`
- **Question to ask**: "Does context flow through all long-running operations? Will Ctrl+C immediately stop everything?"

### 4. Data Integrity & Determinism
- **Completeness**: Data collection operations must fetch ALL requested data or fail entirely - never produce partial results
- **Retry Logic**: Failed operations should retry (with backoff) before failing
- **Idempotency**: Same input parameters should produce identical output every time
- **Validation**: Verify expected vs actual data counts; fail loudly if mismatched
- **Question to ask**: "If I run this twice with the same parameters, will I get identical results? What makes this non-deterministic?"

### 5. Error Handling
- **Error Wrapping**: Use `fmt.Errorf("context: %w", err)` to wrap errors with context
- **Single-line Messages**: Put context before `%w` in single line: `fmt.Errorf("failed after %d attempts: %w", n, err)`
- **Failure Modes**: Consider and handle all failure paths explicitly
- **Logging Levels**: Use appropriate levels (Error for failures, Warn for retries, Info for progress)
- **Question to ask**: "What can fail? How is each failure mode handled? Are errors properly wrapped?"

### 6. Concurrency Patterns
- **Channel Closing**: Close channels in the correct goroutine (usually the sender); use atomic counters to coordinate
- **Worker Pools**: Use `sync.WaitGroup` to wait for workers; protect shared state with mutexes or channels
- **Race Conditions**: Run with `-race` flag during testing
- **Goroutine Leaks**: Ensure every goroutine can exit on context cancellation
- **Question to ask**: "Who closes each channel? Can any goroutine block forever? Does this have race conditions?"

### 7. Testing & Validation
- **Test Coverage**: Write tests for edge cases, not just happy paths
- **Error Injection**: Test retry logic, failure modes, and error paths
- **Resource Limits**: Test with large inputs to verify scalability
- **Cancellation**: Test that context cancellation stops operations immediately
- **Question to ask**: "What edge cases exist? How do I test failure modes?"

### Common Patterns to Apply by Default

```go
// ✅ DO: Timer cleanup
timer := time.NewTimer(500 * time.Millisecond)
defer timer.Stop()
select {
case <-timer.C:
case <-ctx.Done():
    return ctx.Err()
}

// ❌ DON'T: Timer leak
select {
case <-time.After(500 * time.Millisecond): // Leaks if ctx cancels first
case <-ctx.Done():
}

// ✅ DO: HTML escaping
html := fmt.Sprintf(`<div>%s</div>`, html.EscapeString(userInput))

// ❌ DON'T: HTML injection risk
html := fmt.Sprintf(`<div>%s</div>`, userInput)

// ✅ DO: Context propagation
func outputPDF(ctx context.Context, ...) error {
    timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()
    ...
}

// ❌ DON'T: Context.Background in call chain
func outputPDF(...) error {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    ...
}

// ✅ DO: Deterministic data collection with retries
for attempt := 1; attempt <= maxRetries; attempt++ {
    if data, err := fetch(); err == nil {
        return data
    }
}
return fmt.Errorf("failed after %d attempts", maxRetries)

// ❌ DON'T: Skip failures (non-deterministic)
if data, err := fetch(); err != nil {
    log.Warn("skipping failed item")
    continue
}

// ✅ DO: Fixed channel buffer
ch := make(chan T, concurrency*2)

// ❌ DON'T: Buffer proportional to input size
ch := make(chan T, totalItems) // Can allocate GB of memory

// ✅ DO: Goroutine with cancellation
go func() {
    for {
        select {
        case <-ctx.Done():
            return
        case item := <-inputChan:
            process(item)
        }
    }
}()

// ❌ DON'T: Goroutine without cancellation path
go func() {
    for item := range inputChan { // Blocks forever if ctx cancels
        process(item)
    }
}()
```

## Code Style

### Cobra Flags
- Flag names: lowercase with hyphens (kebab-case), e.g., `--output-file`
- Usage strings: lowercase, no ending punctuation, e.g., `"path to output file"`
- Remove unnecessary leading articles and filler words (e.g., "the", "a", "an") from usage strings
- Use `PersistentFlags()` only when flags need to be inherited by subcommands; otherwise use `Flags()`
- When defining multiple flags, use `f := cmd.Flags()` to avoid repetition
- Prefer `Var()` flag methods (e.g., `StringVar`, `IntVar`, `BoolVar`) over non-Var methods (e.g., `String`, `Int`, `Bool`) to bind directly to variables:
  ```go
  f := cmd.Flags()
  f.StringVar(&myVar, "name", "", "description")
  f.IntVar(&count, "count", 0, "description")
  ```
- Flag variables should be non-pointer types unless there's a specific need for pointers (e.g., distinguishing unset from zero value):
  ```go
  var myVar string  // preferred
  var count int     // preferred
  f.StringVar(&myVar, "name", "", "description")
  f.IntVar(&count, "count", 0, "description")
  ```

### Cobra Command Arguments
- Prefer to use Cobra built-in validators (`cobra.NoArgs`, `cobra.ExactArgs(n)`, `cobra.MinimumNArgs(n)`, `cobra.MaximumNArgs(n)`, `cobra.ArbitraryArgs`) instead of custom `Args: func(cmd *cobra.Command, args []string) error` functions, and move argument parsing/validation logic to `PreRunE` hook

### Cobra Commands
- Command `Short` descriptions: sentence case with ending period, e.g., `"Generate a node list to seed a node."`
- Command `Long` descriptions: consider using embedded usage.md file via `//go:embed usage.md` pattern; when using inline strings, use sentence case with ending period for complete sentences
- Command `Short` should be brief (~50 characters or less), appears in help menus and command lists
- Command `Long` provides detailed explanation, can be empty if `Short` is sufficient
