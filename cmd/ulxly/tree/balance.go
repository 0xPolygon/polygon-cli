package tree

import (
	"context"
	"crypto/tls"
	_ "embed"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
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

//go:embed computeBalanceTreeUsage.md
var computeBalanceTreeUsage string

var balanceTreeOptions = &BalanceTreeOptions{}

var BalanceTreeCmd = &cobra.Command{
	Use:   "compute-balance-tree",
	Short: "Compute the balance tree given the deposits.",
	Long:  computeBalanceTreeUsage,
	RunE: func(cmd *cobra.Command, args []string) error {
		return balanceTree()
	},
	SilenceUsage: true,
}

func init() {
	f := BalanceTreeCmd.Flags()
	f.StringVar(&balanceTreeOptions.L2ClaimsFile, ArgL2ClaimsFileName, "", "ndjson file with l2 claim events data")
	f.StringVar(&balanceTreeOptions.L2DepositsFile, ArgL2DepositsFileName, "", "ndjson file with l2 deposit events data")
	f.StringVar(&balanceTreeOptions.BridgeAddress, ArgBridgeAddress, "", "bridge address")
	f.StringVarP(&balanceTreeOptions.RpcURL, ArgRPCURL, "r", "", "RPC URL")
	f.Uint32Var(&balanceTreeOptions.L2NetworkID, ArgL2NetworkID, 0, "L2 network ID")
	f.BoolVar(&balanceTreeOptions.Insecure, ArgInsecure, false, "skip TLS certificate verification")
}

func balanceTree() error {
	l2NetworkID := balanceTreeOptions.L2NetworkID
	bridgeAddress := common.HexToAddress(balanceTreeOptions.BridgeAddress)

	var client *ethclient.Client
	var err error

	if balanceTreeOptions.Insecure {
		client, err = createInsecureEthClient(balanceTreeOptions.RpcURL)
	} else {
		client, err = ethclient.DialContext(context.Background(), balanceTreeOptions.RpcURL)
	}

	if err != nil {
		return err
	}
	defer client.Close()
	l2RawClaimsData, l2RawDepositsData, err := getBalanceTreeData(balanceTreeOptions)
	if err != nil {
		return err
	}
	root, balances, err := computeBalanceTree(client, bridgeAddress, l2RawClaimsData, l2NetworkID, l2RawDepositsData)
	if err != nil {
		return err
	}
	type BalanceEntry struct {
		OriginNetwork      uint32         `json:"originNetwork"`
		OriginTokenAddress common.Address `json:"originTokenAddress"`
		TotalSupply        string         `json:"totalSupply"`
	}

	var balanceEntries []BalanceEntry
	for tokenKey, balance := range balances {
		if balance.Cmp(big.NewInt(0)) == 0 {
			continue
		}

		var token TokenInfo
		token, err = TokenInfoStringToStruct(tokenKey)
		if err != nil {
			return err
		}

		if token.OriginNetwork.Uint64() == uint64(l2NetworkID) {
			continue
		}

		balanceEntries = append(balanceEntries, BalanceEntry{
			OriginNetwork:      uint32(token.OriginNetwork.Uint64()),
			OriginTokenAddress: token.OriginTokenAddress,
			TotalSupply:        balance.String(),
		})
	}

	// Create the response structure
	response := struct {
		Root     string         `json:"root"`
		Balances []BalanceEntry `json:"balances"`
	}{
		Root:     root.String(),
		Balances: balanceEntries,
	}

	// Marshal to JSON with proper formatting
	jsonOutput, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonOutput))
	return nil
}

// createInsecureEthClient creates an Ethereum client with TLS verification disabled
func createInsecureEthClient(rpcURL string) (*ethclient.Client, error) {
	// WARNING: This disables TLS certificate verification
	log.Warn().Msg("WARNING: TLS certificate verification is disabled. This is unsafe for production use.")

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	rpcClient, err := ethrpc.DialOptions(context.Background(), rpcURL, ethrpc.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}

	return ethclient.NewClient(rpcClient), nil
}
