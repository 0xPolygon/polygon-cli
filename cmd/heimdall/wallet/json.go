package wallet

import "encoding/json"

// unmarshalJSON is a thin alias used by the store and import paths.
// Centralised so we can swap in a stricter decoder later without
// touching every call site.
func unmarshalJSON(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
