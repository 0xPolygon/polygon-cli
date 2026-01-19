package importcmd

import (
	_ "embed"
	"fmt"

	"github.com/0xPolygon/polygon-cli/cmd/signer"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

//go:embed usage.md
var usage string

var ImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a private key into the keyring / keystore.",
	Long:  usage,
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := signer.SanityCheck(cmd, args); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := signer.InputSignerOpts
		if opts.Keystore != "" {
			ks := keystore.NewKeyStore(opts.Keystore, keystore.StandardScryptN, keystore.StandardScryptP)
			pk, err := crypto.HexToECDSA(opts.PrivateKey)
			if err != nil {
				return err
			}
			pass, err := signer.GetKeystorePassword()
			if err != nil {
				return err
			}
			_, err = ks.ImportECDSA(pk, pass)
			return err
		}
		if opts.KMS == "GCP" {
			gcpKMS := signer.GCPKMS{}
			if err := gcpKMS.CreateImportJob(cmd.Context()); err != nil {
				return err
			}
			return gcpKMS.ImportKey(cmd.Context())
		}
		return fmt.Errorf("unable to import key")
	},
}
