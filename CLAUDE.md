# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

This is an implementation of **Hyperlane for the Cosmos SDK**, enabling seamless interchain communication following the Hyperlane spec. It allows Cosmos SDK chains to communicate with other blockchains using Hyperlane without relying on CosmWasm.

### Key Modules

- **x/core**: Implements fundamental Hyperlane protocol functionalities for dispatching and processing messages
  - `01_interchain_security`: Implements Interchain Security Modules (ISMs) including No-Op, Merkle Root Multi-Sig, Message ID Multi-Sig, and Routing ISM
  - `02_post_dispatch`: Implements post-dispatch hooks including IGP (Interchain Gas Paymaster), Merkle Tree, and No-Op hooks
- **x/warp**: Extends core functionality for token creation and cross-chain transfers, supporting both collateral and synthetic tokens

## Common Development Commands

### Building and Testing
```bash
# Build all (requires Docker)
make all

# Build simulation app
make build-simapp

# Run tests
make test

# Run specific module tests
go test -cover -mod=readonly ./x/... ./util/...

# Run with race detection
go test -race -coverprofile=coverage.txt ./...
```

### Code Quality
```bash
# Format code
make format

# Run linter
make lint

# Or directly
go run mvdan.cc/gofumpt -l -w .
go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2 run --exclude-dirs scripts --timeout=10m
```

### Protocol Buffers
```bash
# Generate all protobuf files
make proto-all

# Individual commands
make proto-gen     # Generate protobuf files
make proto-format  # Format proto files
make proto-lint    # Lint proto files
```

### Quick Start Test Chain
```bash
make build-simapp
cd build
./hypd init-sample-chain --home test
./hypd start --home test
```

## Architecture

### Message Flow
1. **Dispatch**: Messages are dispatched through a mailbox with required and default hooks
2. **Post-Dispatch Hooks**: IGP for gas payment, Merkle Tree for proof generation
3. **Cross-Chain Delivery**: Messages are relayed to destination chains
4. **ISM Verification**: Destination chain verifies messages using configured ISMs
5. **Process**: Valid messages are processed by the recipient

### Key Components

**Mailbox** (x/core):
- Manages message dispatch and receipt
- Configurable default and required hooks
- Domain-based routing

**ISMs** (x/core/01_interchain_security):
- No-Op ISM: Accepts all messages (testing only)
- Merkle Root Multi-Sig ISM: Requires merkle proof and validator signatures
- Message ID Multi-Sig ISM: Requires direct message signatures
- Routing ISM: Routes to different ISMs based on origin domain

**Post-Dispatch Hooks** (x/core/02_post_dispatch):
- IGP: Handles interchain gas payments
- Merkle Tree Hook: Generates merkle proofs for messages
- No-Op Hook: Pass-through for testing

**Warp Module** (x/warp):
- Collateral Tokens: Lock tokens on origin, mint synthetic on destination
- Synthetic Tokens: Burn synthetic on origin, unlock collateral on destination

### Code Organization

Keeper files follow this pattern:
- `msg_server_*.go`: Entry points for message handling, should be clean and delegate to logic files
- `query_server_*.go`: Entry points for queries, use pagination utilities from `./util`
- `logic_*.go`: Business logic implementation
- `hook_*.go` / `ism_*.go`: Specific implementations

## Testing Approach

- Unit tests use Ginkgo/Gomega for BDD-style testing
- Integration tests in `tests/integration/` use a simulated app
- All tests follow AAA pattern (Arrange, Act, Assert)
- Mock implementations available in `tests/integration/mock.go`

## Important Notes

- Go version: 1.22.11
- Uses Cosmos SDK v0.50.12
- Protobuf generation uses Docker image `ghcr.io/cosmos/proto-builder:0.15.3`
- All modules use dependency injection for integration
- Domain IDs follow Hyperlane's domain registry (see https://docs.hyperlane.xyz/docs/reference/domains)