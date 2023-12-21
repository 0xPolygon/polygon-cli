package signer

import (
	"bytes"
	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	"context"
	"crypto/ecdsa"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	accounts2 "github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/manifoldco/promptui"
	"github.com/maticnetwork/polygon-cli/gethkeystore"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"hash/crc32"
	"math/big"
	"os"
	"strings"
	"time"
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
	gcpProjectID   *string
	gcpRegion      *string
	gcpKeyRingID   *string
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
		if *inputSignerOpts.kms == "GCP" {
			tx, err := getTxDataToSign()
			if err != nil {
				return err
			}
			foo := GCPKMS{}
			return foo.Sign(cmd.Context(), tx)
		}
		return fmt.Errorf("not implemented")
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
		if *inputSignerOpts.kms == "GCP" {
			foo := GCPKMS{}
			err := foo.CreateKeyRing(cmd.Context())
			if err != nil {
				return err
			}
			err = foo.CreateKey(cmd.Context())
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func getTxDataToSign() (*types.Transaction, error) {
	dataToSign, err := os.ReadFile(*inputSignerOpts.dataFile)
	if err != nil {
		return nil, err
	}
	// var tx types.DynamicFeeTx
	// var tx types.LegacyTx
	var tx apitypes.SendTxArgs
	err = json.Unmarshal(dataToSign, &tx)
	if err != nil {
		// TODO in the future it might make sense to sign arbitrary data?
		return nil, err
	}
	// tx.ChainID = new(big.Int).SetUint64(*inputSignerOpts.chainID)
	return tx.ToTransaction(), nil

}
func sign(pk *ecdsa.PrivateKey) error {
	tx, err := getTxDataToSign()
	if err != nil {
		return err
	}
	signer, err := getSigner()
	if err != nil {
		return err
	}
	signedTx, err := types.SignTx(tx, signer, pk)
	if err != nil {
		return err
	}
	return outputSignedTx(signedTx)
}

func outputSignedTx(signedTx *types.Transaction) error {
	rawTx, err := signedTx.MarshalBinary()
	if err != nil {
		return err
	}
	rawHexString := hex.EncodeToString(rawTx)
	out := make(map[string]any, 0)
	out["signedTx"] = signedTx
	out["rawSignedTx"] = rawHexString
	outJSON, err := json.Marshal(out)
	if err != nil {
		return err
	}
	fmt.Println(string(outJSON))
	return nil
}

type GCPKMS struct{}

func (g *GCPKMS) CreateKeyRing(ctx context.Context) error {
	parent := fmt.Sprintf("projects/%s/locations/%s", *inputSignerOpts.gcpProjectID, *inputSignerOpts.gcpRegion)
	id := *inputSignerOpts.gcpKeyRingID
	log.Info().Str("parent", parent).Str("id", id).Msg("Creating keyring")
	// Create the client.
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %w", err)
	}
	defer client.Close()

	result, err := client.GetKeyRing(ctx, &kmspb.GetKeyRingRequest{Name: fmt.Sprintf("%s/keyRings/%s", parent, id)})
	if err != nil {
		nf := strings.Contains(err.Error(), "not found")
		if !nf {
			return err
		}
	}
	if err == nil {
		log.Info().Str("name", result.Name).Msg("key ring already exists")
		return nil
	}
	log.Info().Str("id", id).Msg("key ring not found - creating")

	// Build the request.
	req := &kmspb.CreateKeyRingRequest{
		Parent:    parent,
		KeyRingId: id,
	}

	// Call the API.
	result, err = client.CreateKeyRing(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create key ring: %w", err)
	}
	log.Info().Str("name", result.Name).Msg("Created key ring")
	return nil
}

func (g *GCPKMS) CreateKey(ctx context.Context) error {
	parent := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s", *inputSignerOpts.gcpProjectID, *inputSignerOpts.gcpRegion, *inputSignerOpts.gcpKeyRingID)
	id := *inputSignerOpts.keyID

	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %w", err)
	}
	defer client.Close()

	// Build the request.
	req := &kmspb.CreateCryptoKeyRequest{
		Parent:      parent,
		CryptoKeyId: id,
		CryptoKey: &kmspb.CryptoKey{
			Purpose: kmspb.CryptoKey_ASYMMETRIC_SIGN,
			VersionTemplate: &kmspb.CryptoKeyVersionTemplate{
				Algorithm:       kmspb.CryptoKeyVersion_EC_SIGN_SECP256K1_SHA256,
				ProtectionLevel: kmspb.ProtectionLevel_HSM,
			},

			// Optional: customize how long key versions should be kept before destroying.
			DestroyScheduledDuration: durationpb.New(24 * time.Hour),
		},
	}

	// Call the API.
	result, err := client.CreateCryptoKey(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Info().Str("parent", parent).Str("id", id).Msg("key already exists")
			return nil
		}
		return fmt.Errorf("failed to create key: %w", err)
	}
	log.Info().Str("name", result.Name).Msg("created key")
	return nil

}

type publicKeyInfo struct {
	Raw       asn1.RawContent
	Algorithm pkix.AlgorithmIdentifier
	PublicKey asn1.BitString
}

func (g *GCPKMS) Sign(ctx context.Context, tx *types.Transaction) error {
	// TODO we might need to set a version as a parameter
	name := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s/cryptoKeyVersions/%d", *inputSignerOpts.gcpProjectID, *inputSignerOpts.gcpRegion, *inputSignerOpts.gcpKeyRingID, *inputSignerOpts.keyID, 1)

	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %w", err)
	}
	defer client.Close()

	signer, err := getSigner()
	if err != nil {
		return err
	}
	digest := signer.Hash(tx)

	// Optional but recommended: Compute digest's CRC32C.
	crc32c := func(data []byte) uint32 {
		t := crc32.MakeTable(crc32.Castagnoli)
		return crc32.Checksum(data, t)

	}
	digestCRC32C := crc32c(digest.Bytes())

	req := &kmspb.AsymmetricSignRequest{
		Name: name,
		Digest: &kmspb.Digest{
			Digest: &kmspb.Digest_Sha256{
				Sha256: digest.Bytes(),
			},
		},
		DigestCrc32C: wrapperspb.Int64(int64(digestCRC32C)),
	}

	// Call the API.
	result, err := client.AsymmetricSign(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to sign digest: %w", err)
	}

	// Optional, but recommended: perform integrity verification on result.
	// For more details on ensuring E2E in-transit integrity to and from Cloud KMS visit:
	// https://cloud.google.com/kms/docs/data-integrity-guidelines
	if result.VerifiedDigestCrc32C == false {
		return fmt.Errorf("AsymmetricSign: request corrupted in-transit")
	}
	if result.Name != req.Name {
		return fmt.Errorf("AsymmetricSign: request corrupted in-transit")
	}
	if int64(crc32c(result.Signature)) != result.SignatureCrc32C.Value {
		return fmt.Errorf("AsymmetricSign: response corrupted in-transit")
	}

	pubKeyResponse, err := client.GetPublicKey(ctx, &kmspb.GetPublicKeyRequest{Name: name})
	if err != nil {
		return fmt.Errorf("failed to get public key: %w", err)
	}

	block, _ := pem.Decode([]byte(pubKeyResponse.Pem))
	var gcpPubKey publicKeyInfo
	if _, err := asn1.Unmarshal(block.Bytes, &gcpPubKey); err != nil {
		return err
	}

	// Verify Elliptic Curve signature.
	var parsedSig struct{ R, S *big.Int }
	if _, err = asn1.Unmarshal(result.Signature, &parsedSig); err != nil {
		return fmt.Errorf("asn1.Unmarshal: %w", err)
	}
	ethSig := make([]byte, 0)

	ethSig = append(ethSig, bigIntTo32Bytes(parsedSig.R)...)
	ethSig = append(ethSig, bigIntTo32Bytes(parsedSig.S)...)
	ethSig = append(ethSig, 0)

	// Feels like a hack, but I cna't figure out a better way to determine the recovery ID than this since google isn't returning it. More research is required
	pubKey, err := crypto.Ecrecover(digest.Bytes(), ethSig)
	if err != nil || !bytes.Equal(pubKey, gcpPubKey.PublicKey.Bytes) {
		ethSig[64] = 1
	}
	pubKey, err = crypto.Ecrecover(digest.Bytes(), ethSig)
	if err != nil || !bytes.Equal(pubKey, gcpPubKey.PublicKey.Bytes) {
		return fmt.Errorf("unable to determine recovery identifier value: %w", err)
	}
	pubKeyAddr := common.BytesToAddress(crypto.Keccak256(pubKey[1:])[12:])
	log.Info().
		Str("hexSignature", hex.EncodeToString(result.Signature)).
		Str("ethSignature", hex.EncodeToString(ethSig)).
		Msg("Got signature")

	log.Info().
		Str("recoveredPub", hex.EncodeToString(pubKey)).
		Str("gcpPub", hex.EncodeToString(gcpPubKey.PublicKey.Bytes)).
		Str("ethAddress", pubKeyAddr.String()).
		Msg("Recovered pub key")

	signedTx, err := tx.WithSignature(signer, ethSig)
	if err != nil {
		return err
	}

	return outputSignedTx(signedTx)
}

func bigIntTo32Bytes(num *big.Int) []byte {
	// Convert big.Int to a 32-byte array
	b := num.Bytes()
	if len(b) < 32 {
		// Left-pad with zeros if needed
		b = append(make([]byte, 32-len(b)), b...)
	}
	return b
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

	if *inputSignerOpts.kms == "GCP" {
		if *inputSignerOpts.gcpProjectID == "" {
			return fmt.Errorf("a GCP project id must be specified")
		}

		if *inputSignerOpts.gcpRegion == "" {
			return fmt.Errorf("a location is required")
		}

		if *inputSignerOpts.gcpKeyRingID == "" {
			return fmt.Errorf("a GCP Keyring ID is needed")
		}
		if *inputSignerOpts.keyID == "" {
			return fmt.Errorf("a key id is required")
		}
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

	inputSignerOpts.signerType = SignerCmd.PersistentFlags().String("type", "london", "The type of signer to use: latest, cancun, london, eip2930, eip155")
	inputSignerOpts.dataFile = SignerCmd.PersistentFlags().String("data-file", "", "File name holding data to be signed")

	inputSignerOpts.chainID = SignerCmd.PersistentFlags().Uint64("chain-id", 0, "The chain id for the transactions.")

	// https://github.com/golang/oauth2/issues/241
	inputSignerOpts.gcpProjectID = SignerCmd.PersistentFlags().String("gcp-project-id", "", "The GCP Project ID to use")
	inputSignerOpts.gcpRegion = SignerCmd.PersistentFlags().String("gcp-location", "europe-west2", "The GCP Region to use")
	// What is dead may never die https://cloud.google.com/kms/docs/faq#cannot_delete
	inputSignerOpts.gcpKeyRingID = SignerCmd.PersistentFlags().String("gcp-keyring-id", "polycli-keyring", "The GCP Keyring ID to be used")

	SignerCmd.AddCommand(SignCmd)
	SignerCmd.AddCommand(CreateCmd)
}
