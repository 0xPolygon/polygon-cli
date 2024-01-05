package signer

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	_ "embed"
	"encoding/asn1"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"hash/crc32"
	"math/big"
	"os"
	"strings"
	"time"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	accounts2 "github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/google/tink/go/kwp/subtle"
	"github.com/manifoldco/promptui"
	"github.com/maticnetwork/polygon-cli/gethkeystore"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// signerOpts are the input arguments for these commands
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
	gcpImportJob   *string
	gcpKeyVersion  *int
}

var inputSignerOpts = signerOpts{}

//go:embed usage.md
var signerUsage string

//go:embed signCmdUsage.md
var signCmdUsage string

//go:embed createCmdUsage.md
var createCmdUsage string

//go:embed listCmdUsage.md
var listCmdUsage string

//go:embed importCmdUsage.md
var importCmdUsage string

var SignerCmd = &cobra.Command{
	Use:   "signer",
	Short: "Utilities for security signing transactions",
	Long:  signerUsage,
	Args:  cobra.NoArgs,
}

var SignCmd = &cobra.Command{
	Use:     "sign",
	Short:   "Sign tx data",
	Long:    signCmdUsage,
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
			gcpKMS := GCPKMS{}
			return gcpKMS.Sign(cmd.Context(), tx)
		}
		return fmt.Errorf("not implemented")
	},
}

var CreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new key",
	Long:    createCmdUsage,
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
			gcpKMS := GCPKMS{}
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

var ListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List the keys in the keyring / keystore",
	Long:    listCmdUsage,
	Args:    cobra.NoArgs,
	PreRunE: sanityCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		if *inputSignerOpts.keystore != "" {
			ks := keystore.NewKeyStore(*inputSignerOpts.keystore, keystore.StandardScryptN, keystore.StandardScryptP)
			accounts := ks.Accounts()
			for idx, a := range accounts {
				log.Info().Str("account", a.Address.String()).Int("index", idx).Msg("Account")
			}
			return nil
		}
		if *inputSignerOpts.kms == "GCP" {
			gcpKMS := GCPKMS{}
			return gcpKMS.ListKeyRingKeys(cmd.Context())
		}
		return fmt.Errorf("unable to list accounts")
	},
}

var ImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a private key into the keyring / keystore",
	Long:  importCmdUsage,
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := sanityCheck(cmd, args); err != nil {
			return err
		}
		if err := cmd.MarkFlagRequired("private-key"); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if *inputSignerOpts.keystore != "" {
			ks := keystore.NewKeyStore(*inputSignerOpts.keystore, keystore.StandardScryptN, keystore.StandardScryptP)
			pk, err := crypto.HexToECDSA(*inputSignerOpts.privateKey)
			if err != nil {
				return err
			}
			pass, err := getKeystorePassword()
			if err != nil {
				return err
			}
			_, err = ks.ImportECDSA(pk, pass)
			return err
		}
		if *inputSignerOpts.kms == "GCP" {
			gcpKMS := GCPKMS{}
			if err := gcpKMS.CreateImportJob(cmd.Context()); err != nil {
				return err
			}
			return gcpKMS.ImportKey(cmd.Context())
		}
		return fmt.Errorf("unable to import key")
	},
}

func getTxDataToSign() (*types.Transaction, error) {
	if *inputSignerOpts.dataFile == "" {
		return nil, fmt.Errorf("no datafile was specified to sign")
	}
	dataToSign, err := os.ReadFile(*inputSignerOpts.dataFile)
	if err != nil {
		return nil, err
	}

	// TODO at some point we should support signing other data types besides transactions
	var tx apitypes.SendTxArgs
	err = json.Unmarshal(dataToSign, &tx)
	if err != nil {
		return nil, err
	}
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

func (g *GCPKMS) ListKeyRingKeys(ctx context.Context) error {
	// This snippet has been automatically generated and should be regarded as a code template only.
	// It will require modifications to work:
	// - It may require correct/in-range values for request initialization.
	// - It may require specifying regional endpoints when creating the service client as shown in:
	//   https://pkg.go.dev/cloud.google.com/go#hdr-Client_Options
	c, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return err
	}
	defer c.Close()
	parent := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s", *inputSignerOpts.gcpProjectID, *inputSignerOpts.gcpRegion, *inputSignerOpts.gcpKeyRingID)

	req := &kmspb.ListCryptoKeysRequest{
		Parent: parent,
	}
	it := c.ListCryptoKeys(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		pubKey, err := getPublicKeyByName(ctx, c, fmt.Sprintf("%s/cryptoKeyVersions/%d", resp.Name, *inputSignerOpts.gcpKeyVersion))
		if err != nil {
			log.Error().Err(err).Str("name", resp.Name).Msg("key not found")
			continue
		}
		ethAddress := gcpPubKeyToEthAddress(pubKey)

		log.Info().Str("CryptoKeyBackend", resp.CryptoKeyBackend).
			Str("DestroyScheduledDuration", resp.DestroyScheduledDuration.String()).
			Str("CreateTime", resp.CreateTime.String()).
			Str("Purpose", resp.Purpose.String()).
			Str("ProtectionLevel", resp.VersionTemplate.ProtectionLevel.String()).
			Str("Algorithm", resp.VersionTemplate.Algorithm.String()).
			Str("ETHAddress", ethAddress.String()).
			Str("Name", resp.Name).Msg("got key")

	}
	return nil
}
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

func (g *GCPKMS) CreateImportJob(ctx context.Context) error {
	// parent := "projects/PROJECT_ID/locations/global/keyRings/my-key-ring"
	// id := "my-import-job"
	parent := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s", *inputSignerOpts.gcpProjectID, *inputSignerOpts.gcpRegion, *inputSignerOpts.gcpKeyRingID)
	id := *inputSignerOpts.gcpImportJob

	// Create the client.
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %w", err)
	}
	defer client.Close()

	// Build the request.
	req := &kmspb.CreateImportJobRequest{
		Parent:      parent,
		ImportJobId: id,
		ImportJob: &kmspb.ImportJob{
			// See allowed values and their descriptions at
			// https://cloud.google.com/kms/docs/algorithms#protection_levels
			ProtectionLevel: kmspb.ProtectionLevel_HSM,
			// See allowed values and their descriptions at
			// https://cloud.google.com/kms/docs/key-wrapping#import_methods
			ImportMethod: kmspb.ImportJob_RSA_OAEP_3072_SHA1_AES_256,
		},
	}

	// Call the API.
	result, err := client.CreateImportJob(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Info().Str("name", parent).Msg("import job already exists")
			return nil
		}
		return fmt.Errorf("failed to create import job: %w", err)
	}
	log.Info().Str("name", result.Name).Msg("created import job")

	return nil
}

func (g *GCPKMS) ImportKey(ctx context.Context) error {
	name := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", *inputSignerOpts.gcpProjectID, *inputSignerOpts.gcpRegion, *inputSignerOpts.gcpKeyRingID, *inputSignerOpts.keyID)
	importJob := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/importJobs/%s", *inputSignerOpts.gcpProjectID, *inputSignerOpts.gcpRegion, *inputSignerOpts.gcpKeyRingID, *inputSignerOpts.gcpImportJob)
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %w", err)
	}
	defer client.Close()

	wrappedKey, err := wrapKeyForGCPKMS(ctx, client)
	if err != nil {
		return err
	}
	req := &kmspb.ImportCryptoKeyVersionRequest{
		Parent:     name,
		Algorithm:  kmspb.CryptoKeyVersion_EC_SIGN_SECP256K1_SHA256,
		WrappedKey: wrappedKey,
		ImportJob:  importJob,
	}

	result, err := client.ImportCryptoKeyVersion(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Info().Str("name", name).Msg("key already exists")
			return nil
		}
		return fmt.Errorf("failed to import key: %w", err)
	}
	log.Info().Str("name", result.Name).Msg("imported key")
	return nil

}

func wrapKeyForGCPKMS(ctx context.Context, client *kms.KeyManagementClient) ([]byte, error) {
	// Generate a ECDSA keypair, and format the private key as PKCS #8 DER.
	key, err := crypto.HexToECDSA(*inputSignerOpts.privateKey)
	if err != nil {
		return nil, err
	}
	// These are a lot of hacks because the default x509 library doesn't seem to support the secp256k1 curve
	// START HACKS
	// keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to format private key: %w", err)
	// }

	// https://docs.rs/k256/latest/src/k256/lib.rs.html#116
	oidNamedCurveP256K1 := asn1.ObjectIdentifier{1, 3, 132, 0, 10}
	oidBytes, err := asn1.Marshal(oidNamedCurveP256K1)
	if err != nil {
		return nil, fmt.Errorf("x509: failed to marshal curve OID: %w", err)
	}
	var privKey pkcs8
	oidPublicKeyECDSA := asn1.ObjectIdentifier{1, 2, 840, 10045, 2, 1}
	privKey.Algo = pkix.AlgorithmIdentifier{
		Algorithm: oidPublicKeyECDSA,
		Parameters: asn1.RawValue{
			FullBytes: oidBytes,
		},
	}
	privateKey := make([]byte, (key.Curve.Params().N.BitLen()+7)/8)
	privKey.PrivateKey, err = asn1.Marshal(ecPrivateKey{
		Version:       1, // This is not the GCP Cryptokey version!
		PrivateKey:    key.D.FillBytes(privateKey),
		NamedCurveOID: nil,
		// It looks like elliptic.Marshal is deprecated, but it's still being used in the core library as of go 1.21.5, so I don't want to switch to ecdh especially since it's not obvious how to do so
		// https://cs.opensource.google/go/go/+/refs/tags/go1.21.5:src/crypto/x509/x509.go;l=106
		PublicKey: asn1.BitString{Bytes: elliptic.Marshal(key.Curve, key.X, key.Y)}, //nolint:staticcheck
	})
	if err != nil {
		return nil, fmt.Errorf("unable to marshal private key %w", err)
	}
	keyBytes, err := asn1.Marshal(privKey)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal full private key")
	}
	// END HACKS

	// Generate a temporary 32-byte key for AES-KWP and wrap the key material.
	kwpKey := make([]byte, 32)
	if _, err = rand.Read(kwpKey); err != nil {
		return nil, fmt.Errorf("failed to generate AES-KWP key: %w", err)
	}
	kwp, err := subtle.NewKWP(kwpKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create KWP cipher: %w", err)
	}
	wrappedTarget, err := kwp.Wrap(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to wrap target key with KWP: %w", err)
	}

	importJobName := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/importJobs/%s", *inputSignerOpts.gcpProjectID, *inputSignerOpts.gcpRegion, *inputSignerOpts.gcpKeyRingID, *inputSignerOpts.gcpImportJob)

	// Retrieve the public key from the import job.
	importJob, err := client.GetImportJob(ctx, &kmspb.GetImportJobRequest{
		Name: importJobName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve import job: %w", err)
	}
	pubBlock, _ := pem.Decode([]byte(importJob.PublicKey.Pem))
	pubAny, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse import job public key: %w", err)
	}
	pub, ok := pubAny.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("unexpected public key type %T, want *rsa.PublicKey", pubAny)
	}

	// Wrap the KWP key using the import job key.
	wrappedWrappingKey, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, pub, kwpKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to wrap KWP key: %w", err)
	}

	// Concatenate the wrapped KWP key and the wrapped target key.
	combined := append(wrappedWrappingKey, wrappedTarget...)
	return combined, nil

}

type ecPrivateKey struct {
	Version       int
	PrivateKey    []byte
	NamedCurveOID asn1.ObjectIdentifier `asn1:"optional,explicit,tag:0"`
	PublicKey     asn1.BitString        `asn1:"optional,explicit,tag:1"`
}

type publicKeyInfo struct {
	Raw       asn1.RawContent
	Algorithm pkix.AlgorithmIdentifier
	PublicKey asn1.BitString
}
type pkcs8 struct {
	Version    int
	Algo       pkix.AlgorithmIdentifier
	PrivateKey []byte
	// optional attributes omitted.
}

func (g *GCPKMS) Sign(ctx context.Context, tx *types.Transaction) error {
	name := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s/cryptoKeyVersions/%d", *inputSignerOpts.gcpProjectID, *inputSignerOpts.gcpRegion, *inputSignerOpts.gcpKeyRingID, *inputSignerOpts.keyID, *inputSignerOpts.gcpKeyVersion)

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
	if !result.VerifiedDigestCrc32C {
		return fmt.Errorf("AsymmetricSign: request corrupted in-transit")
	}
	if result.Name != req.Name {
		return fmt.Errorf("AsymmetricSign: request corrupted in-transit")
	}
	if int64(crc32c(result.Signature)) != result.SignatureCrc32C.Value {
		return fmt.Errorf("AsymmetricSign: response corrupted in-transit")
	}

	gcpPubKey, err := getPublicKeyByName(ctx, client, name)
	if err != nil {
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
	// pubKeyAddr := common.BytesToAddress(crypto.Keccak256(pubKey[1:])[12:])
	pubKeyAddr := gcpPubKeyToEthAddress(gcpPubKey)
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

func gcpPubKeyToEthAddress(gcpPubKey *publicKeyInfo) common.Address {
	pubKeyAddr := common.BytesToAddress(crypto.Keccak256(gcpPubKey.PublicKey.Bytes[1:])[12:])
	return pubKeyAddr

}
func getPublicKeyByName(ctx context.Context, client *kms.KeyManagementClient, name string) (*publicKeyInfo, error) {
	pubKeyResponse, err := client.GetPublicKey(ctx, &kmspb.GetPublicKeyRequest{Name: name})
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode([]byte(pubKeyResponse.Pem))
	var gcpPubKey publicKeyInfo
	if _, err = asn1.Unmarshal(block.Bytes, &gcpPubKey); err != nil {
		return nil, err
	}
	return &gcpPubKey, nil
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
	// Strip off the 0x if it's included in the private key hex
	*inputSignerOpts.privateKey = strings.TrimPrefix(*inputSignerOpts.privateKey, "0x")

	// normalize the format of the kms argument
	*inputSignerOpts.kms = strings.ToUpper(*inputSignerOpts.kms)

	keyStoreMethods := 0
	if *inputSignerOpts.kms != "" {
		keyStoreMethods += 1
	}
	if *inputSignerOpts.privateKey != "" && cmd.Name() != "import" {
		keyStoreMethods += 1
	}
	if *inputSignerOpts.keystore != "" {
		keyStoreMethods += 1
	}
	if keyStoreMethods > 1 {
		return fmt.Errorf("Multiple conflicting keystore sources were specified")
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
		if *inputSignerOpts.keyID == "" && cmd.Name() != "list" {
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
	SignerCmd.AddCommand(SignCmd)
	SignerCmd.AddCommand(CreateCmd)
	SignerCmd.AddCommand(ListCmd)
	SignerCmd.AddCommand(ImportCmd)

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
	inputSignerOpts.gcpImportJob = SignerCmd.PersistentFlags().String("gcp-import-job-id", "", "The GCP Import Job ID to use when importing a key")
	inputSignerOpts.gcpKeyVersion = SignerCmd.PersistentFlags().Int("gcp-key-version", 1, "The GCP crypto key version to use")
}
