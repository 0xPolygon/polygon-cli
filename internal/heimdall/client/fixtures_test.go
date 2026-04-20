package client

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestFixturesAreValidJSON verifies that every file committed under
// testdata/{rest,rpc} parses as JSON and carries the envelope keys
// the command layer depends on.
func TestFixturesAreValidJSON(t *testing.T) {
	root := "testdata"
	subdirs := []string{"rest", "rpc"}
	for _, sub := range subdirs {
		dir := filepath.Join(root, sub)
		entries, err := os.ReadDir(dir)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			t.Fatalf("reading %s: %v", dir, err)
		}
		for _, e := range entries {
			if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
				continue
			}
			p := filepath.Join(dir, e.Name())
			raw, err := os.ReadFile(p)
			if err != nil {
				t.Errorf("%s: %v", p, err)
				continue
			}
			var v any
			if err := json.Unmarshal(raw, &v); err != nil {
				t.Errorf("%s: not valid JSON: %v", p, err)
				continue
			}
			// RPC fixtures must be a JSON-RPC 2.0 response envelope.
			if sub == "rpc" {
				m, ok := v.(map[string]any)
				if !ok {
					t.Errorf("%s: rpc fixture is not an object", p)
					continue
				}
				if _, hasResult := m["result"]; !hasResult {
					if _, hasErr := m["error"]; !hasErr {
						t.Errorf("%s: rpc fixture has neither result nor error", p)
					}
				}
			}
		}
	}
}
