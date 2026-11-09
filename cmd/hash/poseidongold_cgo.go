//go:build cgo

package hash

import "github.com/gateway-fm/vectorized-poseidon-gold/src/vectorizedposeidongold"

// poseidongoldAvailable reports whether this binary was built with the
// CGO-only poseidongold implementation.
const poseidongoldAvailable = true

func poseidongoldHash(in [8]uint64, capacity [4]uint64) ([4]uint64, error) {
	return vectorizedposeidongold.Hash(in, capacity)
}
