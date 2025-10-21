# Gas Limiter

## Goal

Control the flow of transactions by limiting the amount of gas it can use. The gas limiter package works basically as a semaphore/throttle mechanism to help the application to know if a determined amount of gas can be spent given the rules defined by the gas provider.

## Components

- Gas Vault: Stores gas budget and allow external components to request it
- Gas Provider: Provides gas budget to a Gas Vault, the provider is the responsible to determine when the budget will be provided, this means it can be periodically, it can watch for events like new blocks on chain, etc.

## How it works

- The application creates a Gas Vault
- Next, the application
  - Creates a Gas provider and set it to provide gas budget to the previously created Gas Vault
  - Start the Gas provider
- When the application needs to do an operation that require gas, it requests the gas to the Gas Vault via `SpendOrWaitAvailableBudget`
  - If the vault has enough budget, it will "spend" the amount from the vault and allow the application to continue.
  - If the vault doesn't have enough budget, it will wait until the budget is provided by the gas provider and the application will hang.

