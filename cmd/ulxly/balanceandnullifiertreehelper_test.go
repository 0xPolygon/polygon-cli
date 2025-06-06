package ulxly_test

import (
	"encoding/json"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/0xPolygon/polygon-cli/cmd/ulxly"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/testvectors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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
		totalTokenBalance, ok := big.NewInt(0).SetString(transition.UpdateLeaf.Value.String(), 0)
		require.Equal(t, true, ok)

		assert.Equal(t, transition.UpdateLeaf.Path, BoolArrayToString(token.ToBits()))

		root, err := balancer.UpdateBalanceTree(token, totalTokenBalance)
		require.NoError(t, err)
		assert.Equal(t, transition.NewRoot, root)
	}
}

func TestBalanceTree2(t *testing.T) {
	balancer, err := ulxly.NewBalanceTree()
	require.NoError(t, err)

	token := ulxly.TokenInfo{
		OriginNetwork:      big.NewInt(0),
		OriginTokenAddress: common.HexToAddress("0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"),
	}
	totalTokenBalance, ok := big.NewInt(0).SetString("100000000000000000000", 0)
	require.Equal(t, true, ok)

	root, err := balancer.UpdateBalanceTree(token, totalTokenBalance)
	require.NoError(t, err)

	token2 := ulxly.TokenInfo{
		OriginNetwork:      big.NewInt(0),
		OriginTokenAddress: common.HexToAddress("0xa23fd6e51aad88f6f4ce6ab8827279cfffb92300"),
	}
	totalToken2Balance, ok := big.NewInt(0).SetString("10000000000000000000", 0)
	require.Equal(t, true, ok)
	root2, err := balancer.UpdateBalanceTree(token2, totalToken2Balance)
	require.NoError(t, err)

	totalToken2Balance = big.NewInt(0)
	root3, err := balancer.UpdateBalanceTree(token2, totalToken2Balance)
	require.NoError(t, err)

	totalToken2Balance = big.NewInt(0)
	root4, err := balancer.UpdateBalanceTree(token, totalToken2Balance)
	require.NoError(t, err)

	t.Log("balancer root: ", root.String())
	t.Log("balancer root2: ", root2.String())
	t.Log("balancer root3: ", root3.String())
	t.Log("balancer root4: ", root4.String())
	assert.Equal(t, "0xb89931f7384aeddb5c136a679d54464007e2d828d4741bec626ff92aeb4b12d4", root4.String())
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

func TestTokenString(t *testing.T) {
	token := ulxly.TokenInfo{
		OriginNetwork:      big.NewInt(0),
		OriginTokenAddress: common.HexToAddress("0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"),
	}
	newToken, err := ulxly.TokenInfoStringToStruct(token.String())
	require.NoError(t, err)
	assert.Equal(t, token.OriginNetwork.String(), newToken.OriginNetwork.String())
	assert.Equal(t, token.OriginTokenAddress.String(), newToken.OriginTokenAddress.String())
	assert.Equal(t, token, newToken)
	assert.Equal(t, token.String(), newToken.String())

	token2 := ulxly.TokenInfo{
		OriginNetwork:      big.NewInt(19),
		OriginTokenAddress: common.HexToAddress("0xf00fd6e51aad88f6f4ce6ab8827279cfffb92006"),
	}
	newToken2, err := ulxly.TokenInfoStringToStruct(token2.String())
	require.NoError(t, err)
	assert.Equal(t, token2.OriginNetwork.String(), newToken2.OriginNetwork.String())
	assert.Equal(t, token2.OriginTokenAddress.String(), newToken2.OriginTokenAddress.String())
	assert.Equal(t, token2, newToken2)
	assert.Equal(t, token2.String(), newToken2.String())
}

func TestCheckHash(t *testing.T) {
	balanceTreeRoot := common.HexToHash("0xb89931f7384aeddb5c136a679d54464007e2d828d4741bec626ff92aeb4b12d4")
	nullifierTreeRoot := common.HexToHash("0xe2c3ed4052eeb1d60514b4c38ece8d73a27f37fa5b36dcbf338e70de95798caa")
	ler_counter := uint32(0)

	initPessimisticRoot := crypto.Keccak256Hash(balanceTreeRoot.Bytes(), nullifierTreeRoot.Bytes(), ulxly.Uint32ToBytesLittleEndian(ler_counter))
	assert.Equal(t, "0xc89c9c0f2ebd19afa9e5910097c43e56fb4aff3a06ddee8d7c9bae09bc769184", initPessimisticRoot.String())
}
