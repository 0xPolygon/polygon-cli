package render

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestRenderKV(t *testing.T) {
	var buf bytes.Buffer
	input := map[string]any{
		"id":         "38835",
		"proposer":   "0x6dc2dd54f24979ec26212794c71afefed722280c",
		"startBlock": "36942195",
	}
	if err := RenderKV(&buf, input, Options{}); err != nil {
		t.Fatalf("RenderKV: %v", err)
	}
	got := buf.String()
	for _, want := range []string{"id          38835", "proposer    0x6dc2dd54", "startBlock  36942195"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
}

func TestRenderJSONNormalizesBytes(t *testing.T) {
	var buf bytes.Buffer
	// base64("abcd") = "YWJjZA==" -> 0x61626364
	input := map[string]any{
		"root_hash": "YWJjZA==",
		"count":     "7",
	}
	if err := RenderJSON(&buf, input, Options{}); err != nil {
		t.Fatalf("RenderJSON: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, `"root_hash": "0x61626364"`) {
		t.Errorf("byte field not normalized: %s", got)
	}
	if !strings.Contains(got, `"count": "7"`) {
		t.Errorf("uint64 string not preserved: %s", got)
	}
}

func TestRenderJSONRawPreservesBase64(t *testing.T) {
	var buf bytes.Buffer
	input := map[string]any{"root_hash": "YWJjZA=="}
	if err := RenderJSON(&buf, input, Options{Raw: true}); err != nil {
		t.Fatalf("RenderJSON: %v", err)
	}
	if !strings.Contains(buf.String(), `"YWJjZA=="`) {
		t.Errorf("expected raw base64 preserved, got %s", buf.String())
	}
}

func TestFieldPluckingBareOutput(t *testing.T) {
	var buf bytes.Buffer
	input := map[string]any{
		"id":       "38835",
		"proposer": "0x6dc2dd54",
	}
	if err := RenderKV(&buf, input, Options{Fields: []string{"proposer"}}); err != nil {
		t.Fatalf("RenderKV: %v", err)
	}
	got := buf.String()
	if got != "0x6dc2dd54\n" {
		t.Errorf("bare field output = %q, want %q", got, "0x6dc2dd54\n")
	}
}

func TestFieldPluckingMultipleFields(t *testing.T) {
	var buf bytes.Buffer
	input := map[string]any{
		"id":       "38835",
		"proposer": "0x6dc2dd54",
		"other":    "nope",
	}
	if err := RenderKV(&buf, input, Options{Fields: []string{"id", "proposer"}}); err != nil {
		t.Fatalf("RenderKV: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "38835") || !strings.Contains(out, "0x6dc2dd54") {
		t.Errorf("multi-field output missing values: %q", out)
	}
	if strings.Contains(out, "nope") {
		t.Errorf("multi-field output leaked non-selected field: %q", out)
	}
}

func TestFieldPluckNestedPath(t *testing.T) {
	var buf bytes.Buffer
	input := map[string]any{
		"block": map[string]any{
			"header": map[string]any{
				"height": "100",
			},
		},
	}
	if err := RenderKV(&buf, input, Options{Fields: []string{"block.header.height"}}); err != nil {
		t.Fatalf("RenderKV: %v", err)
	}
	if got := strings.TrimRight(buf.String(), "\n"); got != "100" {
		t.Errorf("nested pluck = %q, want 100", got)
	}
}

func TestRenderTable(t *testing.T) {
	var buf bytes.Buffer
	recs := []map[string]any{
		{"id": "1", "power": "100"},
		{"id": "2", "power": "50"},
	}
	if err := RenderTable(&buf, recs, Options{}); err != nil {
		t.Fatalf("RenderTable: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "id") || !strings.Contains(got, "power") {
		t.Errorf("table missing headers: %s", got)
	}
	if !strings.Contains(got, "100") || !strings.Contains(got, "50") {
		t.Errorf("table missing rows: %s", got)
	}
}

func TestWriteHintColorHandling(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteHint(&buf, HintBufferEmpty, Options{Color: "never"}); err != nil {
		t.Fatalf("WriteHint: %v", err)
	}
	if strings.Contains(buf.String(), "\x1b[") {
		t.Errorf("ANSI leaked with color=never: %q", buf.String())
	}

	buf.Reset()
	if err := WriteHint(&buf, HintBufferEmpty, Options{Color: "always"}); err != nil {
		t.Fatalf("WriteHint: %v", err)
	}
	if !strings.Contains(buf.String(), "\x1b[2m") {
		t.Errorf("ANSI missing with color=always: %q", buf.String())
	}
}

func TestDetectHintsIsOld(t *testing.T) {
	m := map[string]any{"is_old": false, "signer": "0xabc"}
	hints := DetectHints(m)
	if len(hints) != 1 || hints[0].Key != HintIsOldRenamed.Key {
		t.Errorf("expected is-old hint, got %+v", hints)
	}
}

func TestDetectHintsBufferEmpty(t *testing.T) {
	m := map[string]any{"proposer": "0x0000000000000000000000000000000000000000"}
	hints := DetectHints(m)
	if len(hints) != 1 || hints[0].Key != HintBufferEmpty.Key {
		t.Errorf("expected buffer-empty hint, got %+v", hints)
	}
}

func TestAnnotateUnixSeconds(t *testing.T) {
	now := time.Unix(1776640801, 0).UTC().Add(2*time.Hour + 4*time.Minute)
	got := annotateAt("1776640801", now)
	if !strings.HasPrefix(got, "1776640801  (2026-04-19 23:20:01 UTC,") {
		t.Errorf("AnnotateUnixSeconds = %q", got)
	}
	if !strings.Contains(got, "2h 4m ago") {
		t.Errorf("expected 2h 4m ago, got %q", got)
	}

	if got := annotateAt("not-a-number", now); got != "not-a-number" {
		t.Errorf("non-integer should pass through, got %q", got)
	}
}

func TestWatchExitsOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var out, errOut bytes.Buffer

	fn := func(ctx context.Context) (string, error) {
		return "hello", nil
	}
	done := make(chan error, 1)
	go func() {
		done <- Watch(ctx, &out, &errOut, 50*time.Millisecond, fn)
	}()
	// Let it produce the first snapshot, then cancel.
	time.Sleep(10 * time.Millisecond)
	start := time.Now()
	cancel()
	select {
	case <-done:
		if elapsed := time.Since(start); elapsed > 100*time.Millisecond {
			t.Errorf("watch took %s to cancel, want <100ms", elapsed)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("watch did not exit within 500ms of cancel")
	}
	if !strings.Contains(out.String(), "hello") {
		t.Errorf("initial snapshot missing: %q", out.String())
	}
}

func TestWatchDetectsChanges(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var out, errOut bytes.Buffer

	calls := 0
	fn := func(ctx context.Context) (string, error) {
		calls++
		if calls <= 1 {
			return "v1", nil
		}
		return "v2", nil
	}
	done := make(chan error, 1)
	go func() { done <- Watch(ctx, &out, &errOut, 10*time.Millisecond, fn) }()
	time.Sleep(60 * time.Millisecond)
	cancel()
	<-done

	got := out.String()
	if !strings.Contains(got, "v1") {
		t.Errorf("missing v1 snapshot: %q", got)
	}
	if !strings.Contains(got, "+ v2") || !strings.Contains(got, "- v1") {
		t.Errorf("expected diff +v2 / -v1 in:\n%s", got)
	}
}

// Smoke-check that JSON fields survive through pluckFields even when
// the plucked value is a nested object.
func TestPluckNestedObjectSurvives(t *testing.T) {
	input := map[string]any{
		"info": map[string]any{
			"version": "1.2.3",
			"chain":   "amoy",
		},
	}
	var buf bytes.Buffer
	if err := RenderJSON(&buf, input, Options{Fields: []string{"info"}}); err != nil {
		t.Fatalf("RenderJSON: %v", err)
	}
	// When only one field is plucked, bare output is emitted; for an
	// object, bare output is a JSON encoding.
	var out any
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("plucked nested object not valid JSON: %q", buf.String())
	}
	m, ok := out.(map[string]any)
	if !ok || m["version"] != "1.2.3" {
		t.Errorf("unexpected shape: %v", out)
	}
}
