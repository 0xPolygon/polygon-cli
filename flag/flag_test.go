package flag

import (
	"fmt"
	"os"
	"strconv"
	"testing"

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
