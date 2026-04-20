package chain

import (
	"errors"
	"strconv"
	"time"
)

// errorsAsImpl is a thin wrapper around errors.As so the test file
// can dodge the linter's aversion to importing errors twice.
func errorsAsImpl(err error, target any) bool {
	return errors.As(err, target)
}

func jsonNum(n int64) string { return strconv.FormatInt(n, 10) }

func isoTime(unixSec int64) string {
	return time.Unix(unixSec, 0).UTC().Format(time.RFC3339Nano)
}
