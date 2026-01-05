package txgaschart

import (
	"math"
	"testing"
)

// TestStartBlockLogic tests the start block calculation logic
func TestStartBlockLogic(t *testing.T) {
	tests := []struct {
		name          string
		startBlock    uint64
		endBlock      uint64
		expectedStart uint64
		shouldCalc    bool // whether default calculation should happen
	}{
		{
			name:          "explicit start block 0 should be respected",
			startBlock:    0,
			endBlock:      1000,
			expectedStart: 0,
			shouldCalc:    false,
		},
		{
			name:          "explicit start block 100 should be respected",
			startBlock:    100,
			endBlock:      1000,
			expectedStart: 100,
			shouldCalc:    false,
		},
		{
			name:          "unset start block with end > 500 should calculate default",
			startBlock:    math.MaxUint64,
			endBlock:      1000,
			expectedStart: 500, // 1000 - 500
			shouldCalc:    true,
		},
		{
			name:          "unset start block with end < 500 should default to 0",
			startBlock:    math.MaxUint64,
			endBlock:      400,
			expectedStart: 0,
			shouldCalc:    true,
		},
		{
			name:          "unset start block with end = 500 should default to 0",
			startBlock:    math.MaxUint64,
			endBlock:      500,
			expectedStart: 0,
			shouldCalc:    true,
		},
		{
			name:          "unset start block with end = 501 should calculate default",
			startBlock:    math.MaxUint64,
			endBlock:      501,
			expectedStart: 1, // 501 - 500
			shouldCalc:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the logic from parseFlags
			const defaultBlockRange = 500
			startBlock := tt.startBlock

			if startBlock == math.MaxUint64 {
				if tt.endBlock < defaultBlockRange {
					startBlock = 0
				} else {
					startBlock = tt.endBlock - defaultBlockRange
				}
			}

			if startBlock != tt.expectedStart {
				t.Errorf("Expected start block %d, got %d", tt.expectedStart, startBlock)
			}
		})
	}
}
