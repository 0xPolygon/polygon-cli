package ulxly_test

import (
	"encoding/json"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/0xPolygon/polygon-cli/cmd/ulxly"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/testvectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBalanceTree(t *testing.T) {
	balancer, err := ulxly.NewBalanceTree()
	require.NoError(t, err)

	data, err := os.ReadFile("testvectors/balancetree.json")
	require.NoError(t, err)

	var testVectors vectors.TestVector[vectors.BalanceLeaf]
	err = json.Unmarshal(data, &testVectors)
	require.NoError(t, err)

	for _, transition := range testVectors.Transitions {
		token := ulxly.TokenInfo{
			OriginNetwork:      big.NewInt(0).SetUint64(uint64(transition.UpdateLeaf.Key.OriginNetwork)),
			OriginTokenAddress: transition.UpdateLeaf.Key.OriginTokenAddress,
		}
		totalTokenBalance := big.NewInt(1)
		totalTokenBalance, ok := big.NewInt(0).SetString(transition.UpdateLeaf.Value.String(), 0)
		require.Equal(t, true, ok)

		assert.Equal(t, transition.UpdateLeaf.Path, BoolArrayToString(token.ToBits()))

		root, err := balancer.UpdateBalanceTree(token, totalTokenBalance)
		require.NoError(t, err)
		assert.Equal(t, transition.NewRoot, root)
	}
}

func TestNullifierTree(t *testing.T) {
	nullifier, err := ulxly.NewNullifierTree()
	require.NoError(t, err)

	data, err := os.ReadFile("testvectors/nullifiertree.json")
	require.NoError(t, err)

	var testVectors vectors.TestVector[vectors.NullifierLeaf]
	err = json.Unmarshal(data, &testVectors)
	require.NoError(t, err)

	for _, transition := range testVectors.Transitions {
		n := ulxly.NullifierKey{
			NetworkID: transition.UpdateLeaf.Key.NetworkID,
			Index:     transition.UpdateLeaf.Key.Index,
		}
		assert.Equal(t, transition.UpdateLeaf.Path, BoolArrayToString(n.ToBits()))
		root, err := nullifier.UpdateNullifierTree(n)
		require.NoError(t, err)
		assert.Equal(t, transition.NewRoot, root)
	}
}

func BoolArrayToString(bits []bool) string {
	var b strings.Builder
	for _, bit := range bits {
		if bit {
			b.WriteByte('1')
		} else {
			b.WriteByte('0')
		}
	}
	return b.String()
}
