// Package client provides thin HTTP clients for Heimdall's REST gateway
// and CometBFT JSON-RPC endpoint.
//
// Neither client decodes response bodies. They return raw bytes and
// leave unmarshalling to the caller so subcommands can choose between
// typed decode (for rendered output) and direct JSON passthrough
// (for --json / --field).
package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Transport abstracts how a request is carried out so that --curl can
// short-circuit execution and emit the equivalent command instead.
type Transport interface {
	// Do performs the request and returns the response body, status
	// code, and error. Implementations must honour req.Context().
	Do(req *http.Request) ([]byte, int, error)
}

// HTTPTransport is the default Transport: it forwards to an
// *http.Client.
type HTTPTransport struct {
	Client *http.Client
}

// Do implements Transport.
func (t *HTTPTransport) Do(req *http.Request) ([]byte, int, error) {
	resp, err := t.Client.Do(req)
	if err != nil {
		return nil, 0, &NetworkError{Err: err}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, &NetworkError{Err: err}
	}
	return body, resp.StatusCode, nil
}

// CurlTransport replaces the HTTP call with an equivalent `curl`
// command dumped to Out. Last printed curl is also captured for tests.
type CurlTransport struct {
	Out     io.Writer
	Headers map[string]string
	Last    string
}

// Do implements Transport by rendering the request as a curl
// one-liner and writing it to t.Out; the HTTP body is never sent.
// Returns an empty body with status 0 so callers know no real
// response is available.
func (t *CurlTransport) Do(req *http.Request) ([]byte, int, error) {
	cmd, err := BuildCurl(req, t.Headers)
	if err != nil {
		return nil, 0, err
	}
	t.Last = cmd
	if t.Out != nil {
		fmt.Fprintln(t.Out, cmd)
	}
	return nil, 0, nil
}

// BuildCurl renders a request as an equivalent curl invocation,
// including any headers merged from extra. Caller-side secrets are
// the caller's problem; nothing is redacted.
func BuildCurl(req *http.Request, extra map[string]string) (string, error) {
	var b strings.Builder
	b.WriteString("curl -sS")
	if req.Method != http.MethodGet {
		fmt.Fprintf(&b, " -X %s", req.Method)
	}
	for k, vs := range req.Header {
		for _, v := range vs {
			fmt.Fprintf(&b, " -H %q", k+": "+v)
		}
	}
	for k, v := range extra {
		fmt.Fprintf(&b, " -H %q", k+": "+v)
	}
	if req.Body != nil {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return "", fmt.Errorf("reading request body for curl: %w", err)
		}
		// Restore so the caller can still send it if they want.
		req.Body = io.NopCloser(bytes.NewReader(body))
		fmt.Fprintf(&b, " -d %q", string(body))
	}
	fmt.Fprintf(&b, " %q", req.URL.String())
	return b.String(), nil
}

// RESTClient wraps net/http for Heimdall's REST gateway.
type RESTClient struct {
	BaseURL   string
	Headers   map[string]string
	Transport Transport
}

// NewRESTClient returns a RESTClient configured from the resolved
// config.
func NewRESTClient(base string, timeout time.Duration, headers map[string]string, insecure bool) *RESTClient {
	tlsCfg := &tls.Config{}
	if insecure {
		tlsCfg.InsecureSkipVerify = true
	}
	httpClient := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: tlsCfg,
		},
	}
	return &RESTClient{
		BaseURL:   strings.TrimRight(base, "/"),
		Headers:   cloneHeaders(headers),
		Transport: &HTTPTransport{Client: httpClient},
	}
}

// Get issues a GET against path (starting with '/') and returns the raw
// response body plus the HTTP status code. 2xx responses return
// (body, status, nil). 4xx/5xx responses return (body, status,
// *HTTPError).
func (c *RESTClient) Get(ctx context.Context, path string, query url.Values) ([]byte, int, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u := c.BaseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("building GET %s: %w", u, err)
	}
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}
	body, status, err := c.Transport.Do(req)
	if err != nil {
		return body, status, err
	}
	if status >= 400 {
		return body, status, &HTTPError{Method: http.MethodGet, URL: u, StatusCode: status, Body: body}
	}
	return body, status, nil
}

// Post issues a POST with the given body and Content-Type.
func (c *RESTClient) Post(ctx context.Context, path string, contentType string, body []byte) ([]byte, int, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u := c.BaseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return nil, 0, fmt.Errorf("building POST %s: %w", u, err)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}
	resp, status, err := c.Transport.Do(req)
	if err != nil {
		return resp, status, err
	}
	if status >= 400 {
		return resp, status, &HTTPError{Method: http.MethodPost, URL: u, StatusCode: status, Body: resp}
	}
	return resp, status, nil
}

func cloneHeaders(h map[string]string) map[string]string {
	if len(h) == 0 {
		return nil
	}
	out := make(map[string]string, len(h))
	for k, v := range h {
		out[k] = v
	}
	return out
}
