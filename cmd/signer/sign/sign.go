package sign

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/0xPolygon/polygon-cli/cmd/signer"
	"github.com/0xPolygon/polygon-cli/gethkeystore"
	accounts2 "github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

//go:embed usage.md
var usage string

var SignCmd = &cobra.Command{
	Use:     "sign",
	Short:   "Sign tx data.",
	Long:    usage,
	Args:    cobra.NoArgs,
	PreRunE: signer.SanityCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := signer.InputSignerOpts
		if opts.Keystore == "" && opts.PrivateKey == "" && opts.KMS == "" {
			return fmt.Errorf("no valid keystore was specified")
		}

		if opts.Keystore != "" {
			ks := keystore.NewKeyStore(opts.Keystore, keystore.StandardScryptN, keystore.StandardScryptP)
			accounts := ks.Accounts()
			var accountToUnlock *accounts2.Account
			for _, a := range accounts {
				if a.Address.String() == opts.KeyID {
					accountToUnlock = &a
					break
				}
			}
			if accountToUnlock == nil {
				accountStrings := ""
				for _, a := range accounts {
					accountStrings += a.Address.String() + " "
				}
				return fmt.Errorf("account with address %s not found in list [%s]", opts.KeyID, accountStrings)
			}
			password, err := signer.GetKeystorePassword()
			if err != nil {
				return err
			}

			err = ks.Unlock(*accountToUnlock, password)
			if err != nil {
				return err
			}

			log.Info().Str("path", accountToUnlock.URL.Path).Msg("Unlocked account")
			encryptedKey, err := os.ReadFile(accountToUnlock.URL.Path)
			if err != nil {
				return err
			}
			privKey, err := gethkeystore.DecryptKeystoreFile(encryptedKey, password)
			if err != nil {
				return err
			}
			return signer.Sign(privKey)
		}

		if opts.PrivateKey != "" {
			pk, err := crypto.HexToECDSA(opts.PrivateKey)
			if err != nil {
				return err
			}
			return signer.Sign(pk)
		}
		if opts.KMS == "GCP" {
			tx, err := signer.GetTxDataToSign()
			if err != nil {
				return err
			}
			gcpKMS := signer.GCPKMS{}
			return gcpKMS.Sign(cmd.Context(), tx)
		}
		return fmt.Errorf("not implemented")
	},
}
