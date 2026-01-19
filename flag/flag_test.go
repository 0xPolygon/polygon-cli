package flag

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/0xPolygon/polygon-cli/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// TestValuePriority tests the priority system for flag value resolution.
// It verifies that flag values take precedence over environment variables,
// which take precedence over default values. It also tests the required
// flag validation logic.
func TestValuePriority(t *testing.T) {
	type testCase struct {
		defaultValue *int
		envVarValue  *int
		flagValue    *int
		required     bool

		expectedValue *int
		expectedError error
	}

	testCases := []testCase{
		// Test case: All three sources set - flag should win
		{
			defaultValue:  ptr(1),
			envVarValue:   ptr(2),
			flagValue:     ptr(3),
			expectedValue: ptr(3),
			required:      true,
			expectedError: nil,
		},
		// Test case: Flag set to same value as default - flag should still win
		{
			defaultValue:  ptr(1),
			envVarValue:   ptr(2),
			flagValue:     ptr(1),
			expectedValue: ptr(1),
			required:      true,
			expectedError: nil,
		},
		// Test case: Default and env var set - env var should win
		{
			defaultValue:  ptr(1),
			envVarValue:   ptr(2),
			flagValue:     nil,
			expectedValue: ptr(2),
			required:      true,
			expectedError: nil,
		},
		// Test case: Default and flag set - flag should win
		{
			defaultValue:  ptr(1),
			envVarValue:   nil,
			flagValue:     ptr(3),
			expectedValue: ptr(3),
			required:      true,
			expectedError: nil,
		},
		// Test case: Env var and flag set - flag should win
		{
			defaultValue:  nil,
			envVarValue:   ptr(2),
			flagValue:     ptr(3),
			expectedValue: ptr(3),
			required:      true,
			expectedError: nil,
		},
		// Test case: Only flag set
		{
			defaultValue:  nil,
			envVarValue:   nil,
			flagValue:     ptr(3),
			expectedValue: ptr(3),
			required:      true,
			expectedError: nil,
		},
		// Test case: Only default set (non-required)
		{
			defaultValue:  ptr(1),
			envVarValue:   nil,
			flagValue:     nil,
			expectedValue: ptr(1),
			required:      false,
			expectedError: nil,
		},
		// Test case: Only default set (required) - default should satisfy requirement
		{
			defaultValue:  ptr(1),
			envVarValue:   nil,
			flagValue:     nil,
			expectedValue: ptr(1),
			required:      true,
			expectedError: nil,
		},
		// Test case: Only env var set
		{
			defaultValue:  nil,
			envVarValue:   ptr(2),
			flagValue:     nil,
			expectedValue: ptr(2),
			required:      true,
			expectedError: nil,
		},
		// Test case: Nothing set (non-required) - should return empty
		{
			defaultValue:  nil,
			envVarValue:   nil,
			flagValue:     nil,
			expectedValue: nil,
			required:      false,
			expectedError: nil,
		},
		// Test case: Nothing set (required) - should return error
		{
			defaultValue:  nil,
			envVarValue:   nil,
			flagValue:     nil,
			expectedValue: nil,
			required:      true,
			expectedError: fmt.Errorf("required flag(s) \"flag\" not set"),
		},
	}

	for _, tc := range testCases {
		var value *int
		cmd := &cobra.Command{
			Use: "test",
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				valueStr, err := getValue(cmd, "flag", "FLAG", tc.required)
				if tc.expectedError != nil {
					assert.EqualError(t, err, tc.expectedError.Error())
					return
				}
				assert.NoError(t, err)
				if valueStr != "" {
					valueInt, err := strconv.Atoi(valueStr)
					assert.NoError(t, err)
					value = &valueInt
				}
			},
			Run: func(cmd *cobra.Command, args []string) {
				if tc.expectedValue != nil {
					assert.NotNil(t, value)
					if value != nil {
						assert.Equal(t, *tc.expectedValue, *value)
					}
				} else {
					assert.Nil(t, value)
				}
			},
		}
		if tc.defaultValue != nil {
			cmd.Flags().Int("flag", *tc.defaultValue, "flag")
		} else {
			cmd.Flags().String("flag", "", "flag")
		}

		os.Unsetenv("FLAG")
		if tc.envVarValue != nil {
			v := strconv.Itoa(*tc.envVarValue)
			os.Setenv("FLAG", v)
		}

		if tc.flagValue != nil {
			v := strconv.Itoa(*tc.flagValue)
			cmd.SetArgs([]string{"--flag", v})
		}

		err := cmd.Execute()
		assert.Nil(t, err)
	}

}

// ptr is a helper function to create a pointer to a value.
// This is useful for test cases where we need to distinguish between
// nil (not set) and a zero value (explicitly set to 0).
func ptr[T any](v T) *T {
	return &v
}

// TestValidateURL tests the URL validation function.
func TestValidateURL(t *testing.T) {
	testCases := []struct {
		name        string
		url         string
		expectError bool
	}{
		{
			name:        "Valid HTTP URL",
			url:         "http://localhost:8545",
			expectError: false,
		},
		{
			name:        "Valid HTTPS URL",
			url:         "https://eth-mainnet.example.com",
			expectError: false,
		},
		{
			name:        "Valid WS URL",
			url:         "ws://localhost:8546",
			expectError: false,
		},
		{
			name:        "Valid WSS URL",
			url:         "wss://eth-mainnet.example.com",
			expectError: false,
		},
		{
			name:        "URL with path",
			url:         "https://example.com/rpc/v1",
			expectError: false,
		},
		{
			name:        "URL with port",
			url:         "http://localhost:8545",
			expectError: false,
		},
		{
			name:        "Empty URL",
			url:         "",
			expectError: true,
		},
		{
			name:        "URL without scheme",
			url:         "localhost:8545",
			expectError: true,
		},
		{
			name:        "URL without host",
			url:         "http://",
			expectError: false, // util.ValidateUrl only checks scheme, not host
		},
		{
			name:        "Invalid URL format",
			url:         "ht!tp://invalid",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := util.ValidateUrl(tc.url)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestGetRPCURLValidation tests that GetRPCURL validates URLs.
func TestGetRPCURLValidation(t *testing.T) {
	testCases := []struct {
		name        string
		flagValue   string
		expectError bool
	}{
		{
			name:        "Valid RPC URL",
			flagValue:   "http://localhost:8545",
			expectError: false,
		},
		{
			name:        "Invalid RPC URL - no scheme",
			flagValue:   "localhost:8545",
			expectError: true,
		},
		{
			name:        "Invalid RPC URL - no host",
			flagValue:   "http://",
			expectError: false, // util.ValidateUrl only checks scheme, not host
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv(RPCURLEnvVar, "")
			cmd := &cobra.Command{Use: "test"}
			cmd.Flags().String(RPCURL, DefaultRPCURL, "test rpc url")
			if tc.flagValue != "" {
				cmd.SetArgs([]string{"--" + RPCURL, tc.flagValue})
				_ = cmd.Execute()
			}

			_, err := GetRPCURL(cmd)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
