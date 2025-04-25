<!--

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking CLI commands and REST routes used by end-users.
"API Breaking" for breaking exported APIs used by developers building on SDK.
"State Machine Breaking" for any changes that result in a different AppState given same genesisState and txList.

-->

# CHANGELOG

An '!' indicates a state machine breaking change.

## [v1.0.0-rc0](https://github.com/dymensionxyz/hyperlane-cosmos/releases/tag/v1.0.0-rc0) - 2025-04-15

### Features

- ! [#98](https://github.com/dymensionxyz/hyperlane-cosmos/pull/98) Add renounce ownership functionalities and validate new owner.
- ! [#100](https://github.com/dymensionxyz/hyperlane-cosmos/pull/100) Add support for the Routing ISM.

### Improvements

- [#104](https://github.com/dymensionxyz/hyperlane-cosmos/pull/104) Add more events and align existing ones.
- [#99](https://github.com/dymensionxyz/hyperlane-cosmos/pull/99) Add registration queries for routers, extract create mailbox logic.
- [#102](https://github.com/dymensionxyz/hyperlane-cosmos/pull/102) Add Query Commands for CLI.

## [v1.0.0-beta0](https://github.com/dymensionxyz/hyperlane-cosmos/releases/tag/v1.0.0-beta0) - 2025-03-27

**Initial Release of the Hyperlane Cosmos SDK Module** 🚀

This module integrates the **Hyperlane messaging protocol**
([Hyperlane Docs](https://docs.hyperlane.xyz/)), enabling seamless interchain
communication. It also provides full support for **token bridges**,
secured by **multi-signature Interchain Security Modules**.

### **Key Features**

- **Mailbox Functionality** – Send and receive messages securely across chains.
- **Warp Routes (Token Bridging)**
  - **Collateral Tokens** – Native asset bridging.
  - **Synthetic Tokens** – Wrapped asset representation.
  - **TokenRouter** – Manages token flows across chains.
  - **GasRouter**
- **Interchain Security Modules (ISMs)**
  - **Merkle-Root-Multisig-ISM** – Secure verification using Merkle roots.
  - **MessageId-Multisig-ISM** – Ensures integrity with message ID-based validation.
- **Post Dispatch Hooks**
  - **Merkle Tree Hook** – Supports Merkle-based verification for Multisig ISMs.
  - **InterchainGasPaymaster** – Provides gas prices for destination chains and interchain gas payments.

The module includes a comprehensive test suite and a preconfigured minimal
Cosmos SDK app.
