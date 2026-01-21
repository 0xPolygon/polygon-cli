// Package combined provides the compute-balance-nullifier-tree command.
package combined

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/0xPolygon/polygon-cli/bindings/ulxly"
	ulxlycommon "github.com/0xPolygon/polygon-cli/cmd/ulxly/common"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly/tree"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	ArgL2ClaimsFileName   = "l2-claims-file"
	ArgL2DepositsFileName = "l2-deposits-file"
	ArgBridgeAddress      = "bridge-address"
	ArgRPCURL             = "rpc-url"
	ArgL2NetworkID        = "l2-network-id"
	ArgInsecure           = "insecure"
)

//go:embed usage.md
var usage string

var combinedBalanceTreeOptions = &tree.BalanceTreeOptions{}

var Cmd = &cobra.Command{
	Use:   "compute-balance-nullifier-tree",
	Short: "Compute the balance tree and the nullifier tree given the deposits and claims.",
	Long:  usage,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nullifierAndBalanceTree()
	},
	SilenceUsage: true,
}

func init() {
	f := Cmd.Flags()
	f.StringVar(&combinedBalanceTreeOptions.L2ClaimsFile, ArgL2ClaimsFileName, "", "ndjson file with l2 claim events data")
	f.StringVar(&combinedBalanceTreeOptions.L2DepositsFile, ArgL2DepositsFileName, "", "ndjson file with l2 deposit events data")
	f.StringVar(&combinedBalanceTreeOptions.BridgeAddress, ArgBridgeAddress, "", "bridge address")
	f.StringVarP(&combinedBalanceTreeOptions.RpcURL, ArgRPCURL, "r", "", "RPC URL")
	f.Uint32Var(&combinedBalanceTreeOptions.L2NetworkID, ArgL2NetworkID, 0, "L2 network ID")
	f.BoolVar(&combinedBalanceTreeOptions.Insecure, ArgInsecure, false, "skip TLS certificate verification")
}

func nullifierAndBalanceTree() error {
	l2NetworkID := combinedBalanceTreeOptions.L2NetworkID
	bridgeAddress := common.HexToAddress(combinedBalanceTreeOptions.BridgeAddress)

	var client *ethclient.Client
	var err error

	if combinedBalanceTreeOptions.Insecure {
		client, err = ulxlycommon.CreateInsecureEthClient(combinedBalanceTreeOptions.RpcURL)
	} else {
		client, err = ethclient.DialContext(context.Background(), combinedBalanceTreeOptions.RpcURL)
	}

	if err != nil {
		return err
	}
	defer client.Close()
	l2RawClaimsData, l2RawDepositsData, err := tree.GetBalanceTreeData(combinedBalanceTreeOptions)
	if err != nil {
		return err
	}
	bridgeV2, err := ulxly.NewUlxly(bridgeAddress, client)
	if err != nil {
		return err
	}
	ler_count, err := bridgeV2.LastUpdatedDepositCount(&bind.CallOpts{Pending: false})
	if err != nil {
		return err
	}
	log.Info().Msgf("Last LER count: %d", ler_count)
	balanceTreeRoot, _, err := tree.ComputeBalanceTree(client, bridgeAddress, l2RawClaimsData, l2NetworkID, l2RawDepositsData)
	if err != nil {
		return err
	}
	nullifierTreeRoot, err := tree.ComputeNullifierTree(l2RawClaimsData)
	if err != nil {
		return err
	}
	initPessimisticRoot := crypto.Keccak256Hash(balanceTreeRoot.Bytes(), nullifierTreeRoot.Bytes(), tree.Uint32ToBytesLittleEndian(ler_count))
	fmt.Printf(`
	{
		"balanceTreeRoot": "%s",
		"nullifierTreeRoot": "%s",
		"initPessimisticRoot": "%s"
	}
	`, balanceTreeRoot.String(), nullifierTreeRoot.String(), initPessimisticRoot.String())
	return nil
}
