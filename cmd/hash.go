/*
Copyright © 2022 Polygon <engineering@polygon.technology>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/md4"       //nolint:staticcheck
	"golang.org/x/crypto/ripemd160" //nolint:staticcheck
	"golang.org/x/crypto/sha3"
)

var (
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
	}
	inputFileName *string
)

// hashCmd represents the hash command
var hashCmd = &cobra.Command{
	Use:   fmt.Sprintf("hash [%s]", strings.Join(supportedHashFunctions, "|")),
	Short: "Simple command line tools for common crypto hashing functions",
	Long: `Hash Functions

This is a simple command line tool to run common hash functions. If
the --file argument is provided, we'll attempt to read the file and
hash it. If additionall data is provided on the command line, that
will be hashed. Otherwise, we'll attempted to read the standard input.
`,
	Run: func(cmd *cobra.Command, args []string) {
		data, err := getInputData(cmd, args)
		if err != nil {
			cmd.PrintErrf("There was a nerror reading input for hashing: %s", err.Error())
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
		for _, v := range supportedHashFunctions {
			if v == args[0] {
				return nil
			}
		}

		return fmt.Errorf("the name %s is not recognized. Please use one of the following: %s", args[0], strings.Join(supportedHashFunctions, ","))
	},
}

func init() {
	rootCmd.AddCommand(hashCmd)
	flagSet := hashCmd.PersistentFlags()
	inputFileName = flagSet.String("file", "", "Provide a filename to read and hash")
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
	}
	var h hash.Hash
	return h, fmt.Errorf("unable to create a hash function for %s", name)
}

func getInputData(cmd *cobra.Command, args []string) ([]byte, error) {
	// first check and see if we have an input file
	if inputFileName != nil && *inputFileName != "" {
		// If we get here, we're going to assume the user
		// wants to hash a file and we're not going to look
		// for other input sources
		return os.ReadFile(*inputFileName)
	}

	// This is a little tricky. If a user provdes multiple args that aren't quoted, it could be confusing
	if len(args) > 1 {
		concat := strings.Join(args[1:], " ")
		return []byte(concat), nil
	}

	return io.ReadAll(os.Stdin)
}
