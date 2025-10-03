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
