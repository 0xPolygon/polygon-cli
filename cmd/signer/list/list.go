package list

import (
	_ "embed"
	"fmt"

	"github.com/0xPolygon/polygon-cli/cmd/signer"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

//go:embed usage.md
var usage string

var ListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List the keys in the keyring / keystore.",
	Long:    usage,
	Args:    cobra.NoArgs,
	PreRunE: signer.SanityCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := signer.InputSignerOpts
		if opts.Keystore != "" {
			ks := keystore.NewKeyStore(opts.Keystore, keystore.StandardScryptN, keystore.StandardScryptP)
			accounts := ks.Accounts()
			for idx, a := range accounts {
				log.Info().Str("account", a.Address.String()).Int("index", idx).Msg("Account")
			}
			return nil
		}
		if opts.KMS == "GCP" {
			gcpKMS := signer.GCPKMS{}
			return gcpKMS.ListKeyRingKeys(cmd.Context())
		}
		return fmt.Errorf("unable to list accounts")
	},
}
