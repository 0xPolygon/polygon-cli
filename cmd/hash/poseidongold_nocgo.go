//go:build !cgo

package hash

// poseidongoldAvailable reports whether this binary was built with the
// CGO-only poseidongold implementation. vectorized-poseidon-gold has no
// pure-Go fallback, so CGO_ENABLED=0 builds (e.g. cross-compiled darwin
// binaries) ship without it.
const poseidongoldAvailable = false

func poseidongoldHash(_ [8]uint64, _ [4]uint64) ([4]uint64, error) {
	return [4]uint64{}, errPoseidongoldUnavailable
}
