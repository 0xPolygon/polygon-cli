package chainstore

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
)

// CapabilityManager manages RPC method capability detection
type CapabilityManager struct {
	mu           sync.RWMutex
	capabilities map[string]bool
	lastChecked  time.Time
	ttl          time.Duration
	client       *rpc.Client
}

// NewCapabilityManager creates a new capability manager
func NewCapabilityManager(client *rpc.Client, ttl time.Duration) *CapabilityManager {
	return &CapabilityManager{
		capabilities: make(map[string]bool),
		ttl:          ttl,
		client:       client,
	}
}

// IsMethodSupported checks if a method is supported
func (cm *CapabilityManager) IsMethodSupported(method string) bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// If capabilities are stale, return conservative false
	if time.Since(cm.lastChecked) > cm.ttl {
		return false
	}

	supported, exists := cm.capabilities[method]
	return exists && supported
}

// GetSupportedMethods returns all supported methods
func (cm *CapabilityManager) GetSupportedMethods() []string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var methods []string
	for method, supported := range cm.capabilities {
		if supported {
			methods = append(methods, method)
		}
	}
	return methods
}

// RefreshCapabilities tests and caches method capabilities
func (cm *CapabilityManager) RefreshCapabilities(ctx context.Context) error {
	log.Debug().Msg("Refreshing RPC capabilities")

	// Methods to test for capability
	methodsToTest := []string{
		"eth_chainId",
		"eth_gasPrice",
		"eth_feeHistory",
		"eth_maxPriorityFeePerGas",
		"eth_getBlockByNumber",
		"txpool_status",
		"txpool_inspect",
		"txpool_content",
		"engine_forkchoiceUpdatedV1",
		"engine_forkchoiceUpdatedV2",
		"engine_forkchoiceUpdatedV3",
		"debug_getRawBlock",
		"debug_getRawHeader",
	}

	cm.mu.Lock()
	defer cm.mu.Unlock()

	for _, method := range methodsToTest {
		supported := cm.testMethodCapability(ctx, method)
		cm.capabilities[method] = supported

		log.Debug().
			Str("method", method).
			Bool("supported", supported).
			Msg("Tested RPC method capability")
	}

	cm.lastChecked = time.Now()

	supportedCount := 0
	for _, supported := range cm.capabilities {
		if supported {
			supportedCount++
		}
	}

	log.Info().
		Int("totalMethods", len(methodsToTest)).
		Int("supportedMethods", supportedCount).
		Msg("RPC capability detection completed")

	return nil
}

// testMethodCapability tests if a specific method is supported
func (cm *CapabilityManager) testMethodCapability(ctx context.Context, method string) bool {
	// Create a test context with short timeout
	testCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Test method with minimal parameters
	switch method {
	case "eth_chainId":
		var result string
		err := cm.client.CallContext(testCtx, &result, method)
		return err == nil

	case "eth_gasPrice":
		var result string
		err := cm.client.CallContext(testCtx, &result, method)
		return err == nil

	case "eth_feeHistory":
		var result interface{}
		// Test with minimal parameters: 1 block, latest block, no percentiles
		err := cm.client.CallContext(testCtx, &result, method, "0x1", "latest", nil)
		return err == nil

	case "eth_maxPriorityFeePerGas":
		var result string
		err := cm.client.CallContext(testCtx, &result, method)
		return err == nil

	case "eth_getBlockByNumber":
		var result interface{}
		err := cm.client.CallContext(testCtx, &result, method, "latest", false)
		return err == nil

	case "txpool_status":
		var result interface{}
		err := cm.client.CallContext(testCtx, &result, method)
		return err == nil

	case "txpool_inspect":
		var result interface{}
		err := cm.client.CallContext(testCtx, &result, method)
		return err == nil

	case "txpool_content":
		var result interface{}
		err := cm.client.CallContext(testCtx, &result, method)
		return err == nil

	case "engine_forkchoiceUpdatedV1", "engine_forkchoiceUpdatedV2", "engine_forkchoiceUpdatedV3":
		// These are consensus layer methods, likely not supported by most nodes
		// Test with empty/invalid parameters to see if method exists
		var result interface{}
		err := cm.client.CallContext(testCtx, &result, method, nil, nil)
		// Even if it fails due to invalid params, if the method exists we'll get a different error
		// than "method not found"
		return err != nil && !isMethodNotFoundError(err)

	case "debug_getRawBlock":
		var result interface{}
		err := cm.client.CallContext(testCtx, &result, method, "latest")
		return err == nil

	case "debug_getRawHeader":
		var result interface{}
		err := cm.client.CallContext(testCtx, &result, method, "latest")
		return err == nil

	default:
		// For unknown methods, try a generic call
		var result interface{}
		err := cm.client.CallContext(testCtx, &result, method)
		return err == nil
	}
}

// isMethodNotFoundError checks if the error indicates the method is not found
func isMethodNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	errorStr := strings.ToLower(err.Error())
	// Common patterns for method not found errors
	methodNotFoundPatterns := []string{
		"method not found",
		"method not supported",
		"unknown method",
		"the method does not exist",
		"not supported",
	}

	for _, pattern := range methodNotFoundPatterns {
		if strings.Contains(errorStr, pattern) {
			return true
		}
	}

	return false
}
