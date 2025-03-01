package flag_loader

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

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
		{
			defaultValue:  ptr(1),
			envVarValue:   ptr(2),
			flagValue:     ptr(3),
			expectedValue: ptr(3),
			required:      true,
			expectedError: nil,
		},
		{
			defaultValue:  ptr(1),
			envVarValue:   ptr(2),
			flagValue:     ptr(1),
			expectedValue: ptr(1),
			required:      true,
			expectedError: nil,
		},
		{
			defaultValue:  ptr(1),
			envVarValue:   ptr(2),
			flagValue:     nil,
			expectedValue: ptr(2),
			required:      true,
			expectedError: nil,
		},
		{
			defaultValue:  ptr(1),
			envVarValue:   nil,
			flagValue:     ptr(3),
			expectedValue: ptr(3),
			required:      true,
			expectedError: nil,
		},
		{
			defaultValue:  nil,
			envVarValue:   ptr(2),
			flagValue:     ptr(3),
			expectedValue: ptr(3),
			required:      true,
			expectedError: nil,
		},
		{
			defaultValue:  nil,
			envVarValue:   nil,
			flagValue:     ptr(3),
			expectedValue: ptr(3),
			required:      true,
			expectedError: nil,
		},
		{
			defaultValue:  ptr(1),
			envVarValue:   nil,
			flagValue:     nil,
			expectedValue: ptr(1),
			required:      false,
			expectedError: nil,
		},
		{
			defaultValue:  nil,
			envVarValue:   ptr(2),
			flagValue:     nil,
			expectedValue: ptr(2),
			required:      true,
			expectedError: nil,
		},
		{
			defaultValue:  nil,
			envVarValue:   nil,
			flagValue:     nil,
			expectedValue: ptr(0),
			required:      false,
			expectedError: nil,
		},
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
				valueStr, err := getFlagValue(cmd, "flag", "FLAG", tc.required)
				if tc.expectedError != nil {
					assert.EqualError(t, err, tc.expectedError.Error())
					return
				}
				assert.NoError(t, err)
				valueInt, err := strconv.Atoi(*valueStr)
				assert.NoError(t, err)
				value = &valueInt
			},
			Run: func(cmd *cobra.Command, args []string) {
				if tc.expectedValue != nil {
					assert.Equal(t, *tc.expectedValue, *value)
				} else {
					assert.Nil(t, value)
				}
			},
		}
		if tc.defaultValue != nil {
			cmd.Flags().Int("flag", *tc.defaultValue, "flag")
		} else {
			cmd.Flags().Int("flag", 0, "flag")
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

func ptr[T any](v T) *T {
	return &v
}
