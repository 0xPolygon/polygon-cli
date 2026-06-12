package client

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
)

// HTTPError is returned when the Heimdall node responds with a non-2xx
// status. Body is the raw response body; StatusCode is the HTTP code.
type HTTPError struct {
	Method     string
	URL        string
	StatusCode int
	Body       []byte
}

// Error implements error. Keeps the formatting terse but tries to
// surface enough context to debug from a log line alone.
func (e *HTTPError) Error() string {
	snippet := string(e.Body)
	if len(snippet) > 256 {
		snippet = snippet[:256] + "..."
	}
	return fmt.Sprintf("%s %s: HTTP %d: %s", e.Method, e.URL, e.StatusCode, snippet)
}

// NotFound reports whether the error represents a 404.
func (e *HTTPError) NotFound() bool { return e.StatusCode == 404 }

// NetworkError marks transport-level failures (DNS, TCP reset, TLS
// handshake, etc.).
type NetworkError struct {
	Err error
}

func (e *NetworkError) Error() string { return "network error: " + e.Err.Error() }
func (e *NetworkError) Unwrap() error { return e.Err }

// UsageError marks a caller-side mistake (bad flag, bad argument).
type UsageError struct {
	Msg string
}

func (e *UsageError) Error() string { return "usage error: " + e.Msg }

// ExitCode maps an error to a cast-style exit code.
// 0 -> success, 1 -> node error / not-found, 2 -> network error,
// 3 -> usage error, 4 -> signing error. The signing path lives in
// the tx builder and wraps errors as *SignError (declared elsewhere
// in W2).
func ExitCode(err error) int {
	if err == nil {
		return config.ExitOK
	}
	var uErr *UsageError
	if errors.As(err, &uErr) {
		return config.ExitUsageErr
	}
	var nErr *NetworkError
	if errors.As(err, &nErr) {
		return config.ExitNetErr
	}
	var hErr *HTTPError
	if errors.As(err, &hErr) {
		return config.ExitNodeErr
	}
	// Context cancellation and deadline exceeded are treated as
	// network errors: the caller didn't get useful data from the node.
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return config.ExitNetErr
	}
	var netErr net.Error
	if errors.As(err, &netErr) {
		return config.ExitNetErr
	}
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		return config.ExitNetErr
	}
	return config.ExitNodeErr
}

