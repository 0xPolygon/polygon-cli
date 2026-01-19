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
	"github.com/0xPolygon/polygon-cli/flag"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/google/tink/go/kwp/subtle"
	"github.com/manifoldco/promptui"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// SignerOpts are the input arguments for these commands
type SignerOpts struct {
	Keystore       string
	PrivateKey     string
	KMS            string
	KeyID          string
	UnsafePassword string
	DataFile       string
	SignerType     string
	ChainID        uint64
	GCPProjectID   string
	GCPRegion      string
	GCPKeyRingID   string
	GCPImportJob   string
	GCPKeyVersion  int
}

var InputSignerOpts = SignerOpts{}

//go:embed usage.md
var signerUsage string

var SignerCmd = &cobra.Command{
	Use:   "signer",
	Short: "Utilities for security signing transactions.",
	Long:  signerUsage,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		InputSignerOpts.PrivateKey, err = flag.GetPrivateKey(cmd)
		if err != nil {
			return err
		}
		return nil
	},
	Args: cobra.NoArgs,
}

func GetTxDataToSign() (*ethtypes.Transaction, error) {
	if InputSignerOpts.DataFile == "" {
		return nil, fmt.Errorf("datafile not specified")
	}
	dataToSign, err := os.ReadFile(InputSignerOpts.DataFile)
	if err != nil {
		return nil, err
	}

	// TODO at some point we should support signing other data types besides transactions
	var txArgs apitypes.SendTxArgs
	if err = json.Unmarshal(dataToSign, &txArgs); err != nil {
		return nil, err
	}
	var tx *ethtypes.Transaction
	tx, err = txArgs.ToTransaction()
	if err != nil {
		log.Error().Err(err).Str("txArgs", txArgs.String()).Msg("unable to convert the arguments to a transaction")
		return nil, err
	}
	return tx, nil

}
func Sign(pk *ecdsa.PrivateKey) error {
	tx, err := GetTxDataToSign()
	if err != nil {
		return err
	}
	signer, err := GetSigner()
	if err != nil {
		return err
	}
	signedTx, err := ethtypes.SignTx(tx, signer, pk)
	if err != nil {
		return err
	}
	return OutputSignedTx(signedTx)
}

func OutputSignedTx(signedTx *ethtypes.Transaction) error {
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
	// https://pkg.go.dev/cloud.google.com/go#hdr-Client_Options
	c, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return err
	}
	defer c.Close()
	parent := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s", InputSignerOpts.GCPProjectID, InputSignerOpts.GCPRegion, InputSignerOpts.GCPKeyRingID)

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

		pubKey, err := getPublicKeyByName(ctx, c, fmt.Sprintf("%s/cryptoKeyVersions/%d", resp.Name, InputSignerOpts.GCPKeyVersion))
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
	parent := fmt.Sprintf("projects/%s/locations/%s", InputSignerOpts.GCPProjectID, InputSignerOpts.GCPRegion)
	id := InputSignerOpts.GCPKeyRingID
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
	parent := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s", InputSignerOpts.GCPProjectID, InputSignerOpts.GCPRegion, InputSignerOpts.GCPKeyRingID)
	id := InputSignerOpts.KeyID

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
	parent := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s", InputSignerOpts.GCPProjectID, InputSignerOpts.GCPRegion, InputSignerOpts.GCPKeyRingID)
	id := InputSignerOpts.GCPImportJob

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
	name := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", InputSignerOpts.GCPProjectID, InputSignerOpts.GCPRegion, InputSignerOpts.GCPKeyRingID, InputSignerOpts.KeyID)
	importJob := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/importJobs/%s", InputSignerOpts.GCPProjectID, InputSignerOpts.GCPRegion, InputSignerOpts.GCPKeyRingID, InputSignerOpts.GCPImportJob)
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
	key, err := crypto.HexToECDSA(InputSignerOpts.PrivateKey)
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
		return nil, fmt.Errorf("unable to marshal private key: %w", err)
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

	importJobName := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/importJobs/%s", InputSignerOpts.GCPProjectID, InputSignerOpts.GCPRegion, InputSignerOpts.GCPKeyRingID, InputSignerOpts.GCPImportJob)

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

func (g *GCPKMS) Sign(ctx context.Context, tx *ethtypes.Transaction) error {
	name := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s/cryptoKeyVersions/%d", InputSignerOpts.GCPProjectID, InputSignerOpts.GCPRegion, InputSignerOpts.GCPKeyRingID, InputSignerOpts.KeyID, InputSignerOpts.GCPKeyVersion)

	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %w", err)
	}
	defer client.Close()

	signer, err := GetSigner()
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
		return fmt.Errorf("asymmetric sign: request corrupted in-transit")
	}
	if result.Name != req.Name {
		return fmt.Errorf("asymmetric sign: request corrupted in-transit")
	}
	if int64(crc32c(result.Signature)) != result.SignatureCrc32C.Value {
		return fmt.Errorf("asymmetric sign: response corrupted in-transit")
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

	// Feels like a hack, but I can't figure out a better way to determine the recovery ID than this since google isn't returning it. More research is required
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
		Stringer("ethAddress", pubKeyAddr).
		Msg("Recovered pub key")

	signedTx, err := tx.WithSignature(signer, ethSig)
	if err != nil {
		return err
	}

	return OutputSignedTx(signedTx)
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
func GetKeystorePassword() (string, error) {
	if InputSignerOpts.UnsafePassword != "" {
		return InputSignerOpts.UnsafePassword, nil
	}
	return PasswordPrompt.Run()
}

func SanityCheck(cmd *cobra.Command, args []string) error {
	// Strip off the 0x if it's included in the private key hex
	InputSignerOpts.PrivateKey = strings.TrimPrefix(InputSignerOpts.PrivateKey, "0x")

	// normalize the format of the kms argument
	InputSignerOpts.KMS = strings.ToUpper(InputSignerOpts.KMS)

	keyStoreMethods := 0
	if InputSignerOpts.KMS != "" {
		keyStoreMethods += 1
	}
	if InputSignerOpts.PrivateKey != "" && cmd.Name() != "import" {
		keyStoreMethods += 1
	}
	if InputSignerOpts.Keystore != "" {
		keyStoreMethods += 1
	}
	if keyStoreMethods > 1 {
		return fmt.Errorf("multiple conflicting keystore sources were specified")
	}
	pwErr := PasswordValidation(InputSignerOpts.UnsafePassword)
	if InputSignerOpts.UnsafePassword != "" && pwErr != nil {
		return pwErr
	}

	if InputSignerOpts.KMS == "GCP" {
		if InputSignerOpts.GCPProjectID == "" {
			return fmt.Errorf("GCP project id must be specified")
		}

		if InputSignerOpts.GCPRegion == "" {
			return fmt.Errorf("location is required")
		}

		if InputSignerOpts.GCPKeyRingID == "" {
			return fmt.Errorf("GCP keyring ID is required")
		}
		if InputSignerOpts.KeyID == "" && cmd.Name() != "list" {
			return fmt.Errorf("key id is required")
		}
	}

	return nil
}

func PasswordValidation(inputPw string) error {
	if len(inputPw) < 6 {
		return fmt.Errorf("password only had %d characters, 8 or more required", len(inputPw))
	}
	return nil
}

var PasswordPrompt = promptui.Prompt{
	Label:    "Password",
	Validate: PasswordValidation,
	Mask:     '*',
}

func GetSigner() (ethtypes.Signer, error) {
	chainID := new(big.Int).SetUint64(InputSignerOpts.ChainID)
	switch InputSignerOpts.SignerType {
	case "latest":
		return ethtypes.LatestSignerForChainID(chainID), nil
	case "cancun":
		return ethtypes.NewCancunSigner(chainID), nil
	case "london":
		return ethtypes.NewLondonSigner(chainID), nil
	case "eip2930":
		return ethtypes.NewEIP2930Signer(chainID), nil
	case "eip155":
		return ethtypes.NewEIP155Signer(chainID), nil
	}
	return nil, fmt.Errorf("signer %s is not recognized", InputSignerOpts.SignerType)
}

func init() {
	f := SignerCmd.PersistentFlags()
	f.StringVar(&InputSignerOpts.Keystore, "keystore", "", "use keystore in given folder or file")
	f.StringVar(&InputSignerOpts.PrivateKey, flag.PrivateKey, "", "use provided hex encoded private key")
	f.StringVar(&InputSignerOpts.KMS, "kms", "", "AWS or GCP if key is stored in cloud")
	f.StringVar(&InputSignerOpts.KeyID, "key-id", "", "ID of key to be used for signing")
	f.StringVar(&InputSignerOpts.UnsafePassword, "unsafe-password", "", "non-interactively specified password for unlocking keystore")

	f.StringVar(&InputSignerOpts.SignerType, "type", "london", "type of signer to use: latest, cancun, london, eip2930, eip155")
	f.StringVar(&InputSignerOpts.DataFile, "data-file", "", "file name holding data to be signed")

	f.Uint64Var(&InputSignerOpts.ChainID, "chain-id", 0, "chain ID for transactions")

	// https://github.com/golang/oauth2/issues/241
	f.StringVar(&InputSignerOpts.GCPProjectID, "gcp-project-id", "", "GCP project ID to use")
	f.StringVar(&InputSignerOpts.GCPRegion, "gcp-location", "europe-west2", "GCP region to use")
	// What is dead may never die https://cloud.google.com/kms/docs/faq#cannot_delete
	f.StringVar(&InputSignerOpts.GCPKeyRingID, "gcp-keyring-id", "polycli-keyring", "GCP keyring ID to be used")
	f.StringVar(&InputSignerOpts.GCPImportJob, "gcp-import-job-id", "", "GCP import job ID to use when importing key")
	f.IntVar(&InputSignerOpts.GCPKeyVersion, "gcp-key-version", 1, "GCP crypto key version to use")
}
