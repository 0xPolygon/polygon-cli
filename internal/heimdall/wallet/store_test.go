package wallet

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

// TestFindAccountIndexBounds asserts the integer-index form neither
// panics nor wraps on values that exceed a 31-bit index: they fall
// through to the address parser and surface a usage error.
func TestFindAccountIndexBounds(t *testing.T) {
	ks := keystore.NewKeyStore(t.TempDir(), keystore.LightScryptN, keystore.LightScryptP)

	// In-range index against an empty keystore: clean out-of-range error.
	if _, err := FindAccount(ks, "1"); err == nil || !strings.Contains(err.Error(), "out of range") {
		t.Errorf("index 1 on empty keystore: got %v, want out-of-range error", err)
	}
	// Values above 2^31-1 are no longer parsed as indexes (they could
	// wrap int on 32-bit platforms); they must fall through and be
	// rejected as a non-address.
	for _, in := range []string{"2147483648", "4294967295", "18446744073709551615"} {
		if _, err := FindAccount(ks, in); err == nil || !strings.Contains(err.Error(), "neither an address") {
			t.Errorf("FindAccount(%q): got %v, want non-address usage error", in, err)
		}
	}
}
