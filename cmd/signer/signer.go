package signer

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	accounts2 "github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/manifoldco/promptui"
	"github.com/maticnetwork/polygon-cli/gethkeystore"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"math/big"
	"os"
)

type signerOpts struct {
	keystore       *string
	privateKey     *string
	kms            *string
	keyID          *string
	unsafePassword *string
	dataFile       *string
	signerType     *string
	chainID        *uint64
}

type signedTx struct {
}

var inputSignerOpts = signerOpts{}

var SignerCmd = &cobra.Command{
	Use:   "signer",
	Short: "Utilities for security signing transactions",
	Long:  "TODO",
	Args:  cobra.NoArgs,
}

var SignCmd = &cobra.Command{
	Use:     "sign",
	Short:   "Sign tx data",
	Long:    "TODO",
	Args:    cobra.NoArgs,
	PreRunE: sanityCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		if *inputSignerOpts.keystore == "" && *inputSignerOpts.privateKey == "" && *inputSignerOpts.kms == "" {
			return fmt.Errorf("no valid keystore was specified")
		}

		if *inputSignerOpts.keystore != "" {
			ks := keystore.NewKeyStore(*inputSignerOpts.keystore, keystore.StandardScryptN, keystore.StandardScryptP)
			accounts := ks.Accounts()
			var accountToUnlock *accounts2.Account
			for _, a := range accounts {
				if a.Address.String() == *inputSignerOpts.keyID {
					accountToUnlock = &a
					break
				}
			}
			if accountToUnlock == nil {
				accountStrings := ""
				for _, a := range accounts {
					accountStrings += a.Address.String() + " "
				}
				return fmt.Errorf("the account with address <%s> could not be found in list [%s]", *inputSignerOpts.keyID, accountStrings)
			}
			password, err := getKeystorePassword()
			if err != nil {
				return err
			}

			err = ks.Unlock(*accountToUnlock, password)
			if err != nil {
				return err
			}
			// chainID := new(big.Int).SetUint64(*inputSignerOpts.chainID)

			// ks.SignTx(*accountToUnlock, &tx, chainID)
			log.Info().Str("path", accountToUnlock.URL.Path).Msg("Unlocked account")
			encryptedKey, err := os.ReadFile(accountToUnlock.URL.Path)
			if err != nil {
				return err
			}
			privKey, err := gethkeystore.DecryptKeystoreFile(encryptedKey, password)
			if err != nil {
				return err
			}
			return sign(privKey)
		}

		if *inputSignerOpts.privateKey != "" {
			pk, err := crypto.HexToECDSA(*inputSignerOpts.privateKey)
			if err != nil {
				return err
			}
			return sign(pk)

		}
		return nil
	},
}

var CreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new key",
	Long:    "TODO",
	Args:    cobra.NoArgs,
	PreRunE: sanityCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		if *inputSignerOpts.keystore == "" && *inputSignerOpts.kms == "" {
			log.Info().Msg("Generating new private hex key and writing to stdout")
			pk, err := crypto.GenerateKey()
			if err != nil {
				return err
			}
			k := hex.EncodeToString(crypto.FromECDSA(pk))
			fmt.Println(k)
			return nil
		}
		if *inputSignerOpts.keystore != "" {
			ks := keystore.NewKeyStore(*inputSignerOpts.keystore, keystore.StandardScryptN, keystore.StandardScryptP)
			pk, err := crypto.GenerateKey()
			if err != nil {
				return err
			}
			password, err := getKeystorePassword()
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
		return nil
	},
}

func sign(pk *ecdsa.PrivateKey) error {
	dataToSign, err := os.ReadFile(*inputSignerOpts.dataFile)
	if err != nil {
		return err
	}
	var tx types.LegacyTx
	err = json.Unmarshal(dataToSign, &tx)
	if err != nil {
		// TODO in the future it might make sense to sign arbitrary data?
		return err
	}
	signer, err := getSigner()
	if err != nil {
		return err
	}
	signedTx, err := types.SignNewTx(pk, signer, &tx)
	if err != nil {
		return err
	}
	rlpData := make([]byte, 0)
	buf := bytes.NewBuffer(rlpData)
	err = signedTx.EncodeRLP(buf)
	if err != nil {
		return err
	}
	rawHexString := hex.EncodeToString(buf.Bytes())
	out := make(map[string]any, 0)
	out["txData"] = tx
	out["signedTx"] = signedTx
	out["rawSignedTx"] = rawHexString
	outJSON, err := json.Marshal(out)
	if err != nil {
		return err
	}
	fmt.Println(string(outJSON))
	return nil
}

func getKeystorePassword() (string, error) {
	if *inputSignerOpts.unsafePassword != "" {
		return *inputSignerOpts.unsafePassword, nil
	}
	return passwordPrompt.Run()
}

func sanityCheck(cmd *cobra.Command, args []string) error {
	keyStoreMethods := 0
	if *inputSignerOpts.kms != "" {
		keyStoreMethods += 1
	}
	if *inputSignerOpts.privateKey != "" {
		keyStoreMethods += 1
	}
	if *inputSignerOpts.keystore != "" {
		keyStoreMethods += 1
	}
	if keyStoreMethods > 1 {
		return fmt.Errorf("Multiple conflicting keystore mults were specified")
	}
	pwErr := passwordValidation(*inputSignerOpts.unsafePassword)
	if *inputSignerOpts.unsafePassword != "" && pwErr != nil {
		return pwErr
	}
	return nil
}

func passwordValidation(inputPw string) error {
	if len(inputPw) < 6 {
		return fmt.Errorf("Password only had %d character. 8 or more required", len(inputPw))
	}
	return nil
}

var passwordPrompt = promptui.Prompt{
	Label:    "Password",
	Validate: passwordValidation,
	Mask:     '*',
}

func getSigner() (types.Signer, error) {
	chainID := new(big.Int).SetUint64(*inputSignerOpts.chainID)
	switch *inputSignerOpts.signerType {
	case "latest":
		return types.LatestSignerForChainID(chainID), nil
	case "cancun":
		return types.NewCancunSigner(chainID), nil
	case "london":
		return types.NewLondonSigner(chainID), nil
	case "eip2930":
		return types.NewEIP2930Signer(chainID), nil
	case "eip155":
		return types.NewEIP155Signer(chainID), nil
	}
	return nil, fmt.Errorf("signer %s is not recognized", *inputSignerOpts.signerType)
}

func init() {
	inputSignerOpts.keystore = SignerCmd.PersistentFlags().String("keystore", "", "Use the keystore in the given folder or file")
	inputSignerOpts.privateKey = SignerCmd.PersistentFlags().String("private-key", "", "Use the provided hex encoded private key")
	inputSignerOpts.kms = SignerCmd.PersistentFlags().String("kms", "", "AWS or GCP if the key is stored in the cloud")
	inputSignerOpts.keyID = SignerCmd.PersistentFlags().String("key-id", "", "The id of the key to be used for signing")
	inputSignerOpts.unsafePassword = SignerCmd.PersistentFlags().String("unsafe-password", "", "A non-interactively specified password for unlocking the keystore")

	inputSignerOpts.signerType = SignerCmd.PersistentFlags().String("type", "latest", "The type of signer to use: latest, cancun, london, eip2930, eip155")
	inputSignerOpts.dataFile = SignerCmd.PersistentFlags().String("data-file", "", "File name holding data to be signed")

	inputSignerOpts.chainID = SignerCmd.PersistentFlags().Uint64("chain-id", 0, "The chain id for the transactions.")

	SignerCmd.AddCommand(SignCmd)
	SignerCmd.AddCommand(CreateCmd)
}
