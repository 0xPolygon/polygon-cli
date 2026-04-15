package flag

import (
	"fmt"
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

// GasValue is a custom flag type for gas values with unit support.
// It implements the pflag.Value interface to enable parsing values like "100gwei" or "1ether".
type GasValue struct {
	Val *uint64
}

// unitMultipliers maps unit names to their multipliers in wei.
var unitMultipliers = map[string]*big.Int{
	"wei":   big.NewInt(1),
	"gwei":  big.NewInt(1e9),
	"ether": big.NewInt(1e18),
	"eth":   big.NewInt(1e18),
}

// gasPattern matches a number (possibly decimal) followed by an optional unit.
var gasPattern = regexp.MustCompile(`^([0-9]+(?:\.[0-9]+)?)\s*([a-zA-Z]*)$`)

// String returns the string representation of the gas value in wei.
func (g *GasValue) String() string {
	if g.Val == nil {
		return "0"
	}
	return strconv.FormatUint(*g.Val, 10)
}

// Set parses a gas value string with optional unit (e.g., "100gwei", "1000000000").
func (g *GasValue) Set(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return fmt.Errorf("empty gas value")
	}

	matches := gasPattern.FindStringSubmatch(s)
	if matches == nil {
		return fmt.Errorf("invalid gas value format: %q", s)
	}

	numStr := matches[1]
	unit := strings.ToLower(matches[2])

	// If no unit specified, treat as wei
	if unit == "" {
		unit = "wei"
	}

	multiplier, ok := unitMultipliers[unit]
	if !ok {
		return fmt.Errorf("unknown gas unit: %q (supported: wei, gwei, ether, eth)", unit)
	}

	result, err := parseWithUnit(numStr, multiplier)
	if err != nil {
		return err
	}

	*g.Val = result
	return nil
}

// Type returns the type string for this flag value.
func (g *GasValue) Type() string {
	return "gas"
}

// parseWithUnit parses a number string and multiplies it by the given unit multiplier.
// Handles decimal values like "0.1gwei" correctly.
func parseWithUnit(numStr string, multiplier *big.Int) (uint64, error) {
	// Handle decimal numbers
	if strings.Contains(numStr, ".") {
		parts := strings.Split(numStr, ".")
		if len(parts) != 2 {
			return 0, fmt.Errorf("invalid number format: %q", numStr)
		}

		// Parse the integer part
		intPart := parts[0]
		if intPart == "" {
			intPart = "0"
		}
		intVal, ok := new(big.Int).SetString(intPart, 10)
		if !ok {
			return 0, fmt.Errorf("invalid integer part: %q", intPart)
		}

		// Parse the decimal part
		decPart := parts[1]
		if decPart == "" {
			// No decimal digits, just use integer part
			result := new(big.Int).Mul(intVal, multiplier)
			return bigIntToUint64(result)
		}

		// Calculate the fractional value
		// e.g., "0.1" with multiplier 1e9 -> 0.1 * 1e9 = 1e8
		decLen := len(decPart)
		decVal, ok := new(big.Int).SetString(decPart, 10)
		if !ok {
			return 0, fmt.Errorf("invalid decimal part: %q", decPart)
		}

		// Scale the decimal: decVal * multiplier / 10^decLen
		divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decLen)), nil)
		scaledDec := new(big.Int).Mul(decVal, multiplier)
		scaledDec.Div(scaledDec, divisor)

		// Total = intPart * multiplier + scaledDec
		result := new(big.Int).Mul(intVal, multiplier)
		result.Add(result, scaledDec)

		return bigIntToUint64(result)
	}

	// Integer value
	val, ok := new(big.Int).SetString(numStr, 10)
	if !ok {
		return 0, fmt.Errorf("invalid number: %q", numStr)
	}

	result := new(big.Int).Mul(val, multiplier)
	return bigIntToUint64(result)
}

// bigIntToUint64 safely converts a big.Int to uint64, returning an error on overflow.
func bigIntToUint64(val *big.Int) (uint64, error) {
	if val.Sign() < 0 {
		return 0, fmt.Errorf("negative gas value not allowed")
	}
	if !val.IsUint64() {
		return 0, fmt.Errorf("gas value %s exceeds uint64 max (%d)", val.String(), uint64(math.MaxUint64))
	}
	return val.Uint64(), nil
}

// ParseGasUnit parses a gas value string with optional unit and returns the value in wei.
// This is a convenience function for use outside of flag parsing.
func ParseGasUnit(s string) (uint64, error) {
	var result uint64
	gv := &GasValue{Val: &result}
	if err := gv.Set(s); err != nil {
		return 0, err
	}
	return result, nil
}
