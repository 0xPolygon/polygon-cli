package flag

import (
	"fmt"
	"math/big"
)

// BigIntValue is a custom flag type for big.Int values.
// It implements the pflag.Value interface to enable using *big.Int with Cobra flags.
type BigIntValue struct {
	Val *big.Int
}

// String returns the decimal string representation of the big.Int value.
func (b *BigIntValue) String() string {
	return b.Val.String()
}

// Set parses a decimal string and sets the big.Int value.
func (b *BigIntValue) Set(s string) error {
	if _, ok := b.Val.SetString(s, 10); !ok {
		return fmt.Errorf("invalid big integer: %q", s)
	}
	return nil
}

// Type returns the type string for this flag value.
func (b *BigIntValue) Type() string {
	return "big.Int"
}
