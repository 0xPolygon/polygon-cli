package ulxly_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/0xPolygon/polygon-cli/cmd/ulxly"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestBalanceTree(t *testing.T) {
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

	fmt.Println("balancer root: ", root.String())
	fmt.Println("balancer root2: ", root2.String())
	fmt.Println("balancer root3: ", root3.String())
}

func TestNullifierTree(t *testing.T) {
	nullifier, err := ulxly.NewNullifierTree()
	require.NoError(t, err)

	n := ulxly.NullifierKey{
		NetworkID: 0,
		Index:     0,
	}
	root := nullifier.UpdateNullifierTree(n)
	require.NoError(t, err)

	n2 := ulxly.NullifierKey{
		NetworkID: 0,
		Index:     2,
	}
	root2 := nullifier.UpdateNullifierTree(n2)
	require.NoError(t, err)

	n3 := ulxly.NullifierKey{
		NetworkID: 0,
		Index:     1,
	}
	root3 := nullifier.UpdateNullifierTree(n3)
	require.NoError(t, err)

	fmt.Println("nullifier root: ", root.String())
	fmt.Println("nullifier root2: ", root2.String())
	fmt.Println("nullifier root3: ", root3.String())
}
