package client

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestRESTClientGetSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test") != "present" {
			t.Errorf("X-Test header missing, got %q", r.Header.Get("X-Test"))
		}
		if r.URL.Path != "/checkpoints/count" {
			t.Errorf("path = %q, want /checkpoints/count", r.URL.Path)
		}
		if got := r.URL.Query().Get("pagination.limit"); got != "10" {
			t.Errorf("pagination.limit = %q, want 10", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"count":"42"}`)
	}))
	defer srv.Close()

	c := NewRESTClient(srv.URL, 5*time.Second, map[string]string{"X-Test": "present"}, false)
	body, status, err := c.Get(context.Background(), "/checkpoints/count", url.Values{"pagination.limit": {"10"}})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if status != 200 {
		t.Errorf("status = %d, want 200", status)
	}
	if !strings.Contains(string(body), `"count":"42"`) {
		t.Errorf("body = %q", body)
	}
}

func TestRESTClientGet404(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, `{"code":5,"message":"not found"}`)
	}))
	defer srv.Close()

	c := NewRESTClient(srv.URL, 5*time.Second, nil, false)
	_, _, err := c.Get(context.Background(), "/checkpoints/999999", nil)
	if err == nil {
		t.Fatal("expected error on 404, got nil")
	}
	var hErr *HTTPError
	if !errors.As(err, &hErr) {
		t.Fatalf("error type = %T, want *HTTPError", err)
	}
	if !hErr.NotFound() {
		t.Errorf("NotFound() = false, want true")
	}
	if ExitCode(err) != 1 {
		t.Errorf("ExitCode = %d, want 1", ExitCode(err))
	}
}

func TestRESTClientGetTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
		case <-time.After(2 * time.Second):
		}
	}))
	defer srv.Close()

	c := NewRESTClient(srv.URL, 50*time.Millisecond, nil, false)
	start := time.Now()
	_, _, err := c.Get(context.Background(), "/anything", nil)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if elapsed := time.Since(start); elapsed > 500*time.Millisecond {
		t.Errorf("timeout took %s, want <500ms", elapsed)
	}
	if ExitCode(err) != 2 {
		t.Errorf("ExitCode for timeout = %d, want 2", ExitCode(err))
	}
}

func TestRESTClientContextCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer srv.Close()

	c := NewRESTClient(srv.URL, 5*time.Second, nil, false)
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() {
		_, _, err := c.Get(ctx, "/anything", nil)
		done <- err
	}()

	time.Sleep(10 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("expected cancellation error, got nil")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("cancellation not honoured within 500ms")
	}
}

func TestRESTClientCurlTransport(t *testing.T) {
	var buf bytes.Buffer
	curl := &CurlTransport{Out: &buf, Headers: map[string]string{"X-Extra": "yep"}}
	c := &RESTClient{BaseURL: "http://example.test", Transport: curl}

	_, _, err := c.Get(context.Background(), "/checkpoints/count", url.Values{"x": {"1"}})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "curl -sS") {
		t.Errorf("missing curl prefix: %q", out)
	}
	if !strings.Contains(out, "/checkpoints/count?x=1") {
		t.Errorf("missing path+query: %q", out)
	}
	if !strings.Contains(out, "X-Extra: yep") {
		t.Errorf("missing extra header: %q", out)
	}
}

func TestExitCodeMappings(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{"nil", nil, 0},
		{"usage", &UsageError{Msg: "bad flag"}, 3},
		{"http", &HTTPError{StatusCode: 404}, 1},
		{"network", &NetworkError{Err: errors.New("dns")}, 2},
		{"ctx-cancel", context.Canceled, 2},
		{"ctx-deadline", context.DeadlineExceeded, 2},
	}
	for _, tc := range tests {
		if got := ExitCode(tc.err); got != tc.want {
			t.Errorf("%s: ExitCode = %d, want %d", tc.name, got, tc.want)
		}
	}
}
