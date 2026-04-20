package wallet

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

// ResolveKeystoreDir returns the keystore directory to use per the
// precedence rule:
//
//  1. override (from --keystore-dir flag), if non-empty
//  2. ETH_KEYSTORE environment variable
//  3. ~/.foundry/keystores, if it already exists (honour existing cast
//     users without migration)
//  4. ~/.polycli/keystores (the default)
//
// createDefault controls whether step 4 creates the fallback directory
// on demand. Keystore-management commands (`wallet new`, `wallet
// import`) pass true so the default is materialised the first time an
// operator uses it. Signing commands (`mktx`, `send`, `estimate`) pass
// false: they shouldn't silently create a keystore dir just because an
// operator typo'd an address, and they should surface a clear "account
// not found" error instead.
//
// The returned path is absolute and logged at debug so operators can
// see why a given path was chosen.
func ResolveKeystoreDir(override string, createDefault bool) (string, error) {
	switch {
	case override != "":
		abs, err := filepath.Abs(override)
		if err != nil {
			return "", fmt.Errorf("resolving --keystore-dir %q: %w", override, err)
		}
		log.Debug().Str("source", "flag").Str("path", abs).Msg("heimdall keystore dir")
		return abs, nil
	case os.Getenv("ETH_KEYSTORE") != "":
		abs, err := filepath.Abs(os.Getenv("ETH_KEYSTORE"))
		if err != nil {
			return "", fmt.Errorf("resolving ETH_KEYSTORE %q: %w", os.Getenv("ETH_KEYSTORE"), err)
		}
		log.Debug().Str("source", "env").Str("path", abs).Msg("heimdall keystore dir")
		return abs, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolving home directory: %w", err)
	}
	foundry := filepath.Join(home, ".foundry", "keystores")
	if st, err := os.Stat(foundry); err == nil && st.IsDir() {
		log.Debug().Str("source", "foundry").Str("path", foundry).Msg("heimdall keystore dir")
		return foundry, nil
	}
	polycli := filepath.Join(home, ".polycli", "keystores")
	if createDefault {
		if err := os.MkdirAll(polycli, 0o700); err != nil {
			return "", fmt.Errorf("creating %s: %w", polycli, err)
		}
	}
	log.Debug().Str("source", "default").Str("path", polycli).Msg("heimdall keystore dir")
	return polycli, nil
}
