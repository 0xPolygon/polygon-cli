package cmdutil

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/client"
)

// GRPCErrorBody is the standard gRPC-gateway error envelope returned on
// 4xx/5xx from Heimdall REST. Only `code` and `message` are used here.
type GRPCErrorBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// GRPCCodeUnavailable is the L1-unreachable code surfaced by the
// endpoints that fan out to L1 on the server side (clerk/topup
// sequence + is-old, stake is-old-tx, topup account-proof) when the
// Heimdall node lacks `eth_rpc_url`.
const GRPCCodeUnavailable = 13

// IsL1Unreachable inspects a REST body / error pair and returns true if
// the response looks like "gRPC code 13 because L1 RPC isn't configured
// on this Heimdall". The body may come either from a successful 2xx
// response that still carries a gRPC-error envelope, or from an
// HTTPError (4xx/5xx). The transport layer surfaces "dial tcp" /
// "connection refused" when the REST gateway itself can't reach L1.
func IsL1Unreachable(body []byte, err error) bool {
	var hErr *client.HTTPError
	if errors.As(err, &hErr) && len(hErr.Body) > 0 {
		body = hErr.Body
	}
	if len(body) == 0 {
		return errLooksUnreachable(err)
	}
	var g GRPCErrorBody
	if jerr := json.Unmarshal(body, &g); jerr == nil && g.Code == GRPCCodeUnavailable {
		return true
	}
	return errLooksUnreachable(err)
}

func errLooksUnreachable(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "dial tcp")
}
