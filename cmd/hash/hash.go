package hash

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	_ "embed"

	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/md4"       //nolint:staticcheck
	"golang.org/x/crypto/ripemd160" //nolint:staticcheck
	"golang.org/x/crypto/sha3"

	"github.com/gateway-fm/vectorized-poseidon-gold/src/vectorizedposeidongold"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/rs/zerolog/log"
)

var (
	//go:embed usage.md
	usage                  string
	supportedHashFunctions = []string{
		"md4",
		"md5",
		"sha1",
		"sha224",
		"sha256",
		"sha384",
		"sha512",
		"ripemd160",
		"sha3_224",
		"sha3_256",
		"sha3_384",
		"sha3_512",
		"sha512_224",
		"sha512_256",
		"blake2s_256",
		"blake2b_256",
		"blake2b_384",
		"blake2b_512",
		"keccak256",
		"keccak512",
		"poseidon",
		"poseidongold",
	}
	inputFileName string
)

// hashCmd represents the hash command
var HashCmd = &cobra.Command{
	Use:   fmt.Sprintf("hash [%s]", strings.Join(supportedHashFunctions, "|")),
	Short: "Provide common crypto hashing functions.",
	Long:  usage,
	Run: func(cmd *cobra.Command, args []string) {
		data, err := getInputData(cmd, args)
		if err != nil {
			cmd.PrintErrf("There was an error reading input for hashing: %s", err.Error())
			return
		}
		h, err := getHash(args[0])
		if err != nil {
			cmd.PrintErrf("There was an error creating the hash function: %s", err.Error())
			return
		}
		h.Write(data)
		hashOut := h.Sum(nil)
		cmd.Println(hex.EncodeToString(hashOut))

	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("expected 1 argument to specify hash function. got %d", len(args))
		}
		if slices.Contains(supportedHashFunctions, args[0]) {
			return nil
		}

		return fmt.Errorf("the name %s is not recognized. Please use one of the following: %s", args[0], strings.Join(supportedHashFunctions, ","))
	},
}

func init() {
	f := HashCmd.Flags()
	f.StringVar(&inputFileName, "file", "", "Provide a filename to read and hash")
}

func getHash(name string) (hash.Hash, error) {
	switch name {
	case "md4":
		return md4.New(), nil
	case "md5":
		return md5.New(), nil
	case "sha1":
		return sha1.New(), nil
	case "sha224":
		return sha256.New224(), nil
	case "sha256":
		return sha256.New(), nil
	case "sha384":
		return sha512.New384(), nil
	case "sha512":
		return sha512.New(), nil
	case "ripemd160":
		return ripemd160.New(), nil
	case "sha3_224":
		return sha3.New224(), nil
	case "sha3_256":
		return sha3.New256(), nil
	case "sha3_384":
		return sha3.New384(), nil
	case "sha3_512":
		return sha3.New512(), nil
	case "sha512_224":
		return sha512.New512_224(), nil
	case "sha512_256":
		return sha512.New512_256(), nil
	case "blake2s_256":
		return blake2s.New256(nil)
	case "blake2b_256":
		return blake2b.New256(nil)
	case "blake2b_384":
		return blake2b.New384(nil)
	case "blake2b_512":
		return blake2b.New512(nil)
	case "keccak256":
		return sha3.NewLegacyKeccak256(), nil
	case "keccak512":
		return sha3.NewLegacyKeccak512(), nil
	case "poseidon":
		return poseidon.New(16)
	case "poseidongold":
		pg := new(poseidongoldWrapper)
		return pg, nil
	}
	var h hash.Hash
	return h, fmt.Errorf("unable to create a hash function for %s", name)
}

func getInputData(cmd *cobra.Command, args []string) ([]byte, error) {
	// first check and see if we have an input file
	if inputFileName != "" {
		// If we get here, we're going to assume the user
		// wants to hash a file and we're not going to look
		// for other input sources
		return os.ReadFile(inputFileName)
	}

	// This is a little tricky. If a user provides multiple args that aren't quoted, it could be confusing
	if len(args) > 1 {
		concat := strings.Join(args[1:], " ")
		return []byte(concat), nil
	}

	return io.ReadAll(os.Stdin)
}

type poseidongoldWrapper struct {
	b *bytes.Buffer
}

func (p *poseidongoldWrapper) Write(toWrite []byte) (n int, err error) {
	if p.b == nil {
		buf := make([]byte, 0)
		p.b = bytes.NewBuffer(buf)
	}
	n, err = p.b.Write(toWrite)
	if err != nil {
		log.Error().Err(err).Msg("Unable to write to poseidon buffer")
		return
	}
	return
}

func (p *poseidongoldWrapper) Sum(b []byte) []byte {
	cap := [4]uint64{}
	buf := make([]byte, 64)

	var err error
	for {
		n, _ := p.b.Read(buf)
		log.Info().Bytes("current-buffer", buf).Msg("summing")
		if n > 64 {
			panic("What?? that shouldn't have happened")
		}
		hashInput, _ := bytesToUints(buf)
		cap, err = vectorizedposeidongold.Hash(hashInput, cap)
		if err != nil {
			log.Error().Err(err).Msg("Unable to hash input")
			// Exit since we don't want to return a known bad hash value here
			os.Exit(1)
		}
		if n < 64 {
			log.Info().Int("n", n).Msg("done")
			break
		}

	}
	return Uint64ArrayToBytes(cap)
}
func (p *poseidongoldWrapper) Reset() {
	p.b.Reset()
}
func (p *poseidongoldWrapper) Size() int {
	return 32
}
func (p *poseidongoldWrapper) BlockSize() int {
	return 64
}

func bytesToUints(buf []byte) ([8]uint64, error) {
	if len(buf) > 64 {
		return [8]uint64{}, fmt.Errorf("the underlying data is too large. Expected less than 64 bytes but got %d", len(buf))
	}

	var result [8]uint64
	for i := 0; i < len(buf); i += 8 {
		var val uint64
		for j := 0; j < 8 && i+j < len(buf); j++ {
			val |= uint64(buf[i+j]) << (8 * j)
		}
		result[i/8] = val
	}

	return result, nil
}

func Uint64ArrayToBytes(arr [4]uint64) []byte {
	buf := new(bytes.Buffer)
	for k := range arr {
		// this serialization seems a little weird, but it matches the values from erigion as far as I can tell.
		// echo -n "0000" | xxd -r -p | ./out/polycli hash poseidongold
		// This returns
		// c71603f33a1144ca7953db0ab48808f4c4055e3364a246c33c18a9786cb0b359
		// which matches
		// https://github.com/0xPolygonHermez/cdk-erigon/blob/92e7d22534bc103203a4a699b8de60a4a5df4c3e/smt/pkg/utils/utils.go#L24
		//
		// This doesn't quite match the order from the tests in hermez and plonky
		// https://github.com/0xPolygonHermez/goldilocks/blob/a0a516ac38145049a527e6c822e11c8df108d6d1/tests/tests.cpp#L1515-L1518
		// https://github.com/0xPolygonZero/plonky2/blob/349beae1431ecffc1bf8c044d6c00e2bf194b74a/plonky2/src/hash/poseidon_goldilocks.rs#L468
		err := binary.Write(buf, binary.BigEndian, arr[4-k-1])
		if err != nil {
			// This error handling is mainly for demonstration, it should never occur for fixed size and types.
			fmt.Println("binary.Write failed:", err)
		}
	}
	return buf.Bytes()
}
