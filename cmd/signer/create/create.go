package create

import (
	_ "embed"
	"encoding/hex"
	"fmt"

	"github.com/0xPolygon/polygon-cli/signer"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

//go:embed usage.md
var usage string

var CreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new key.",
	Long:    usage,
	Args:    cobra.NoArgs,
	PreRunE: signer.SanityCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := signer.InputOpts
		if opts.Keystore == "" && opts.KMS == "" {
			log.Info().Msg("Generating new private hex key and writing to stdout")
			pk, err := crypto.GenerateKey()
			if err != nil {
				return err
			}
			k := hex.EncodeToString(crypto.FromECDSA(pk))
			fmt.Println(k)
			return nil
		}
		if opts.Keystore != "" {
			ks := keystore.NewKeyStore(opts.Keystore, keystore.StandardScryptN, keystore.StandardScryptP)
			pk, err := crypto.GenerateKey()
			if err != nil {
				return err
			}
			password, err := signer.GetKeystorePassword()
			if err != nil {
				return err
			}
			acc, err := ks.ImportECDSA(pk, password)
			if err != nil {
				return err
			}
			log.Info().Str("address", acc.Address.String()).Msg("imported new account")
			return nil
		}
		if opts.KMS == "GCP" {
			gcpKMS := signer.GCPKMS{}
			err := gcpKMS.CreateKeyRing(cmd.Context())
			if err != nil {
				return err
			}
			err = gcpKMS.CreateKey(cmd.Context())
			if err != nil {
				return err
			}
		}
		return nil
	},
}
