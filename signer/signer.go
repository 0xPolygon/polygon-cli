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

// Opts are the input arguments for signer commands.
type Opts struct {
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

var InputOpts = Opts{}

func GetTxDataToSign() (*ethtypes.Transaction, error) {
	if InputOpts.DataFile == "" {
		return nil, fmt.Errorf("datafile not specified")
	}
	dataToSign, err := os.ReadFile(InputOpts.DataFile)
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
	c, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return err
	}
	defer c.Close()
	parent := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s", InputOpts.GCPProjectID, InputOpts.GCPRegion, InputOpts.GCPKeyRingID)

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

		pubKey, err := getPublicKeyByName(ctx, c, fmt.Sprintf("%s/cryptoKeyVersions/%d", resp.Name, InputOpts.GCPKeyVersion))
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
	parent := fmt.Sprintf("projects/%s/locations/%s", InputOpts.GCPProjectID, InputOpts.GCPRegion)
	id := InputOpts.GCPKeyRingID
	log.Info().Str("parent", parent).Str("id", id).Msg("Creating keyring")
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

	req := &kmspb.CreateKeyRingRequest{
		Parent:    parent,
		KeyRingId: id,
	}

	result, err = client.CreateKeyRing(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create key ring: %w", err)
	}
	log.Info().Str("name", result.Name).Msg("Created key ring")
	return nil
}

func (g *GCPKMS) CreateKey(ctx context.Context) error {
	parent := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s", InputOpts.GCPProjectID, InputOpts.GCPRegion, InputOpts.GCPKeyRingID)
	id := InputOpts.KeyID

	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %w", err)
	}
	defer client.Close()

	req := &kmspb.CreateCryptoKeyRequest{
		Parent:      parent,
		CryptoKeyId: id,
		CryptoKey: &kmspb.CryptoKey{
			Purpose: kmspb.CryptoKey_ASYMMETRIC_SIGN,
			VersionTemplate: &kmspb.CryptoKeyVersionTemplate{
				Algorithm:       kmspb.CryptoKeyVersion_EC_SIGN_SECP256K1_SHA256,
				ProtectionLevel: kmspb.ProtectionLevel_HSM,
			},
			DestroyScheduledDuration: durationpb.New(24 * time.Hour),
		},
	}

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
	parent := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s", InputOpts.GCPProjectID, InputOpts.GCPRegion, InputOpts.GCPKeyRingID)
	id := InputOpts.GCPImportJob

	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %w", err)
	}
	defer client.Close()

	req := &kmspb.CreateImportJobRequest{
		Parent:      parent,
		ImportJobId: id,
		ImportJob: &kmspb.ImportJob{
			ProtectionLevel: kmspb.ProtectionLevel_HSM,
			ImportMethod:    kmspb.ImportJob_RSA_OAEP_3072_SHA1_AES_256,
		},
	}

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
	name := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", InputOpts.GCPProjectID, InputOpts.GCPRegion, InputOpts.GCPKeyRingID, InputOpts.KeyID)
	importJob := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/importJobs/%s", InputOpts.GCPProjectID, InputOpts.GCPRegion, InputOpts.GCPKeyRingID, InputOpts.GCPImportJob)
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
	key, err := crypto.HexToECDSA(InputOpts.PrivateKey)
	if err != nil {
		return nil, err
	}

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
		Version:       1,
		PrivateKey:    key.D.FillBytes(privateKey),
		NamedCurveOID: nil,
		PublicKey:     asn1.BitString{Bytes: elliptic.Marshal(key.Curve, key.X, key.Y)}, //nolint:staticcheck
	})
	if err != nil {
		return nil, fmt.Errorf("unable to marshal private key: %w", err)
	}
	keyBytes, err := asn1.Marshal(privKey)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal full private key")
	}

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

	importJobName := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/importJobs/%s", InputOpts.GCPProjectID, InputOpts.GCPRegion, InputOpts.GCPKeyRingID, InputOpts.GCPImportJob)

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

	wrappedWrappingKey, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, pub, kwpKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to wrap KWP key: %w", err)
	}

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
}

func (g *GCPKMS) Sign(ctx context.Context, tx *ethtypes.Transaction) error {
	name := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s/cryptoKeyVersions/%d", InputOpts.GCPProjectID, InputOpts.GCPRegion, InputOpts.GCPKeyRingID, InputOpts.KeyID, InputOpts.GCPKeyVersion)

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

	result, err := client.AsymmetricSign(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to sign digest: %w", err)
	}

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

	var parsedSig struct{ R, S *big.Int }
	if _, err = asn1.Unmarshal(result.Signature, &parsedSig); err != nil {
		return fmt.Errorf("asn1.Unmarshal: %w", err)
	}
	ethSig := make([]byte, 0)

	ethSig = append(ethSig, bigIntTo32Bytes(parsedSig.R)...)
	ethSig = append(ethSig, bigIntTo32Bytes(parsedSig.S)...)
	ethSig = append(ethSig, 0)

	pubKey, err := crypto.Ecrecover(digest.Bytes(), ethSig)
	if err != nil || !bytes.Equal(pubKey, gcpPubKey.PublicKey.Bytes) {
		ethSig[64] = 1
	}
	pubKey, err = crypto.Ecrecover(digest.Bytes(), ethSig)
	if err != nil || !bytes.Equal(pubKey, gcpPubKey.PublicKey.Bytes) {
		return fmt.Errorf("unable to determine recovery identifier value: %w", err)
	}
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
	b := num.Bytes()
	if len(b) < 32 {
		b = append(make([]byte, 32-len(b)), b...)
	}
	return b
}

func GetKeystorePassword() (string, error) {
	if InputOpts.UnsafePassword != "" {
		return InputOpts.UnsafePassword, nil
	}
	return PasswordPrompt.Run()
}

func SanityCheck(cmd *cobra.Command, args []string) error {
	InputOpts.PrivateKey = strings.TrimPrefix(InputOpts.PrivateKey, "0x")
	InputOpts.KMS = strings.ToUpper(InputOpts.KMS)

	keyStoreMethods := 0
	if InputOpts.KMS != "" {
		keyStoreMethods += 1
	}
	if InputOpts.PrivateKey != "" && cmd.Name() != "import" {
		keyStoreMethods += 1
	}
	if InputOpts.Keystore != "" {
		keyStoreMethods += 1
	}
	if keyStoreMethods > 1 {
		return fmt.Errorf("multiple conflicting keystore sources were specified")
	}
	pwErr := PasswordValidation(InputOpts.UnsafePassword)
	if InputOpts.UnsafePassword != "" && pwErr != nil {
		return pwErr
	}

	if InputOpts.KMS == "GCP" {
		if InputOpts.GCPProjectID == "" {
			return fmt.Errorf("GCP project id must be specified")
		}

		if InputOpts.GCPRegion == "" {
			return fmt.Errorf("location is required")
		}

		if InputOpts.GCPKeyRingID == "" {
			return fmt.Errorf("GCP keyring ID is required")
		}
		if InputOpts.KeyID == "" && cmd.Name() != "list" {
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
	chainID := new(big.Int).SetUint64(InputOpts.ChainID)
	switch InputOpts.SignerType {
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
	return nil, fmt.Errorf("signer %s is not recognized", InputOpts.SignerType)
}
