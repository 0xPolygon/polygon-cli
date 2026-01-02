package util

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// WaitReceipt waits for a transaction receipt with default parameters.
func WaitReceipt(ctx context.Context, client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	return internalWaitReceipt(ctx, client, txHash, 0, 0, 0)
}

// WaitReceiptWithRetries waits for a transaction receipt with retries and exponential backoff.
func WaitReceiptWithRetries(ctx context.Context, client *ethclient.Client, txHash common.Hash, maxRetries uint, initialDelayMs uint) (*types.Receipt, error) {
	return internalWaitReceipt(ctx, client, txHash, maxRetries, initialDelayMs, 0)
}

// WaitReceiptWithTimeout waits for a transaction receipt with a specified timeout.
func WaitReceiptWithTimeout(ctx context.Context, client *ethclient.Client, txHash common.Hash, timeout time.Duration) (*types.Receipt, error) {
	return internalWaitReceipt(ctx, client, txHash, 0, 0, timeout)
}

// internalWaitReceipt waits for a transaction receipt with retries, exponential backoff, and a timeout.
func internalWaitReceipt(ctx context.Context, client *ethclient.Client, txHash common.Hash, maxRetries uint, initialDelayMs uint, timeout time.Duration) (*types.Receipt, error) {
	// Set defaults for zero values
	effectiveTimeout := timeout
	if effectiveTimeout == 0 {
		effectiveTimeout = 1 * time.Minute // Default: 1 minute
	}

	// Create context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, effectiveTimeout)
	defer cancel()

	effectiveInitialDelayMs := initialDelayMs
	if effectiveInitialDelayMs == 0 {
		effectiveInitialDelayMs = 100 // Default: 100ms
	} else if effectiveInitialDelayMs < 10 {
		effectiveInitialDelayMs = 10 // Minimum 10ms
	}

	for attempt := uint(0); ; attempt++ {
		receipt, err := client.TransactionReceipt(timeoutCtx, txHash)
		if err == nil && receipt != nil {
			return receipt, nil
		}

		// If maxRetries > 0 and we've reached the limit, exit
		// Note: effectiveMaxRetries is always > 0 due to default above
		if maxRetries > 0 && attempt >= maxRetries-1 {
			return nil, fmt.Errorf("failed to get receipt after %d attempts: %w", maxRetries, err)
		}

		// Calculate delay
		baseDelay := time.Duration(effectiveInitialDelayMs) * time.Millisecond
		exponentialDelay := baseDelay * time.Duration(1<<attempt)

		// Add cap to prevent extremely long delays
		maxDelay := 30 * time.Second
		if exponentialDelay > maxDelay {
			exponentialDelay = maxDelay
		}

		maxJitter := exponentialDelay / 2
		if maxJitter <= 0 {
			maxJitter = 1 * time.Millisecond
		}
		jitter := time.Duration(rand.Int63n(int64(maxJitter)))
		totalDelay := exponentialDelay + jitter

		select {
		case <-timeoutCtx.Done():
			return nil, timeoutCtx.Err()
		case <-time.After(totalDelay):
			// Continue
		}
	}
}

func WaitReceiptNew(ctx context.Context, client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	var maxRetries uint = 10
	var timeout time.Duration = time.Minute
	var delay time.Duration = 2 * time.Second

	// Create context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for attempt := uint(0); ; attempt++ {
		receipt, err := client.TransactionReceipt(timeoutCtx, txHash)
		if err == nil && receipt != nil {
			return receipt, nil
		}

		// If maxRetries > 0 and we've reached the limit, exit
		// Note: effectiveMaxRetries is always > 0 due to default above
		if maxRetries > 0 && attempt >= maxRetries-1 {
			return nil, fmt.Errorf("failed to get receipt after %d attempts: %w", maxRetries, err)
		}

		select {
		case <-timeoutCtx.Done():
			return nil, timeoutCtx.Err()
		case <-time.After(delay):
			// Continue
		}
	}
}

func WaitPreconf(ctx context.Context, client *ethclient.Client, txHash common.Hash) (bool, error) {
	var maxRetries uint = 15
	var timeout time.Duration = time.Minute
	var delay time.Duration = 1 * time.Second

	// Create context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for attempt := uint(0); ; attempt++ {
		var res interface{}
		err := client.Client().CallContext(ctx, &res, "eth_checkPreconfStatus", txHash.Hex())
		if err == nil {
			return res.(bool), nil
		}

		// If maxRetries > 0 and we've reached the limit, exit
		// Note: effectiveMaxRetries is always > 0 due to default above
		if maxRetries > 0 && attempt >= maxRetries-1 {
			return false, fmt.Errorf("failed to get receipt after %d attempts: %w", maxRetries, err)
		}

		select {
		case <-timeoutCtx.Done():
			return false, timeoutCtx.Err()
		case <-time.After(delay):
			// Continue
		}
	}
}
