package signer

import (
	_ "embed"

	"github.com/0xPolygon/polygon-cli/cmd/signer/create"
	importcmd "github.com/0xPolygon/polygon-cli/cmd/signer/import"
	"github.com/0xPolygon/polygon-cli/cmd/signer/list"
	"github.com/0xPolygon/polygon-cli/cmd/signer/sign"
	"github.com/0xPolygon/polygon-cli/flag"
	"github.com/0xPolygon/polygon-cli/signer"
	"github.com/spf13/cobra"
)

//go:embed usage.md
var signerUsage string

var SignerCmd = &cobra.Command{
	Use:   "signer",
	Short: "Utilities for security signing transactions.",
	Long:  signerUsage,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		signer.InputOpts.PrivateKey, err = flag.GetPrivateKey(cmd)
		if err != nil {
			return err
		}
		return nil
	},
	Args: cobra.NoArgs,
}

func init() {
	SignerCmd.AddCommand(sign.SignCmd)
	SignerCmd.AddCommand(create.CreateCmd)
	SignerCmd.AddCommand(list.ListCmd)
	SignerCmd.AddCommand(importcmd.ImportCmd)

	f := SignerCmd.PersistentFlags()
	f.StringVar(&signer.InputOpts.Keystore, "keystore", "", "use keystore in given folder or file")
	f.StringVar(&signer.InputOpts.PrivateKey, flag.PrivateKey, "", "use provided hex encoded private key")
	f.StringVar(&signer.InputOpts.KMS, "kms", "", "AWS or GCP if key is stored in cloud")
	f.StringVar(&signer.InputOpts.KeyID, "key-id", "", "ID of key to be used for signing")
	f.StringVar(&signer.InputOpts.UnsafePassword, "unsafe-password", "", "non-interactively specified password for unlocking keystore")

	f.StringVar(&signer.InputOpts.SignerType, "type", "london", "type of signer to use: latest, cancun, london, eip2930, eip155")
	f.StringVar(&signer.InputOpts.DataFile, "data-file", "", "file name holding data to be signed")

	f.Uint64Var(&signer.InputOpts.ChainID, "chain-id", 0, "chain ID for transactions")

	f.StringVar(&signer.InputOpts.GCPProjectID, "gcp-project-id", "", "GCP project ID to use")
	f.StringVar(&signer.InputOpts.GCPRegion, "gcp-location", "europe-west2", "GCP region to use")
	f.StringVar(&signer.InputOpts.GCPKeyRingID, "gcp-keyring-id", "polycli-keyring", "GCP keyring ID to be used")
	f.StringVar(&signer.InputOpts.GCPImportJob, "gcp-import-job-id", "", "GCP import job ID to use when importing key")
	f.IntVar(&signer.InputOpts.GCPKeyVersion, "gcp-key-version", 1, "GCP crypto key version to use")
}
