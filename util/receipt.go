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

// WaitPreconf waits for a preconf status check with a specified timeout.
// Uses exponential backoff with jitter, similar to WaitReceiptWithTimeout.
func WaitPreconf(ctx context.Context, client *ethclient.Client, txHash common.Hash, timeout time.Duration) (bool, error) {
	if timeout == 0 {
		timeout = 1 * time.Minute
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	const initialDelayMs = 100
	timer := time.NewTimer(0)
	defer timer.Stop()
	// Drain initial timer since we want to check immediately on first iteration
	<-timer.C

	for attempt := uint(0); ; attempt++ {
		var res any
		err := client.Client().CallContext(timeoutCtx, &res, "eth_checkPreconfStatus", txHash.Hex())
		if err == nil {
			return res.(bool), nil
		}

		// Calculate delay with exponential backoff and jitter
		baseDelay := time.Duration(initialDelayMs) * time.Millisecond
		exponentialDelay := baseDelay * time.Duration(1<<attempt)

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

		timer.Reset(totalDelay)
		select {
		case <-timeoutCtx.Done():
			return false, timeoutCtx.Err()
		case <-timer.C:
			// Continue
		}
	}
}
