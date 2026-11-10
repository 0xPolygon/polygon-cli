package nodelist

import "testing"

func TestLimitMustBePositive(t *testing.T) {
	for _, tc := range []struct {
		limit   int
		wantErr bool
	}{{-5, true}, {0, true}, {1, false}, {100, false}} {
		inputNodeListParams.Limit = tc.limit
		err := NodeListCmd.PreRunE(NodeListCmd, []string{"nodes.json"})
		if (err != nil) != tc.wantErr {
			t.Errorf("--limit %d: err = %v, want error %v", tc.limit, err, tc.wantErr)
		}
	}
}
