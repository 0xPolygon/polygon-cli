package publish

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"iter"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

const scannerBufferSize = 1024 * 1024

func getInputData(inputFileName *string, args []string) (iter.Seq[string], inputDataSource, error) {
	// firstly check and see if we have an input file
	if inputFileName != nil && *inputFileName != "" {
		// If we get here, we're going to assume the user
		// wants to load transactions from a file and we're not going to look
		// for other input sources
		return dataFromFile(*inputFileName)
	}

	// secondly check and see if we have any args
	if len(args) > 0 {
		// checks if any of the args start with 0x, if so, return only the args
		// that start with 0x
		txArgs := make([]string, 0)
		for _, arg := range args {
			if arg[:2] == "0x" {
				txArgs = append(txArgs, arg)
			}
		}

		// If we get here, we're going to assume the user wants to load transactions
		// from the command line and we're not going to look for other input sources
		if len(txArgs) > 0 {
			return dataFromArgs(txArgs)
		}
	}

	// if we get here, we're going to assume the user wants to load transactions
	// from stdin or from the command line
	return dataFromStdin()
}

func dataFromArgs(args []string) (iter.Seq[string], inputDataSource, error) {
	log.Info().Msg("Reading data from args")

	return func(yield func(string) bool) {
		for _, arg := range args {
			if !yield(arg) {
				return
			}
		}
	}, InputDataSourceArgs, nil
}

// dataFromFile returns an iterator that reads lines from a file
func dataFromFile(filename string) (iter.Seq[string], inputDataSource, error) {
	f, err := os.Open(filename)
	if err != nil {
		return func(yield func(string) bool) {}, InputDataSourceFile, err
	}
	log.Info().
		Str("filename", filename).
		Msg("Reading data from file")

	return func(yield func(string) bool) {
		// Ensure the file is closed after the function exits
		defer f.Close()
		s := bufio.NewScanner(f)
		sBuf := make([]byte, 0)
		s.Buffer(sBuf, scannerBufferSize)
		for s.Scan() {
			if !yield(s.Text()) {
				return
			}
		}
		if err := s.Err(); err != nil {
			log.Error().
				Err(err).
				Str("filename", filename).
				Msg("error scanning file")
		}
	}, InputDataSourceFile, nil
}

// dataFromStdin returns an iterator that reads lines from stdin
func dataFromStdin() (iter.Seq[string], inputDataSource, error) {
	fmt.Println("Reading data from Stdin, type the transactions you want to publish:")

	return func(yield func(string) bool) {
		s := bufio.NewScanner(os.Stdin)
		sBuf := make([]byte, 0)
		s.Buffer(sBuf, scannerBufferSize)
		for s.Scan() {
			if !yield(s.Text()) {
				return
			}
		}
		if err := s.Err(); err != nil {
			log.Error().
				Err(err).
				Msg("error scanning stdin")
		}
	}, InputDataSourceStdin, nil
}

// inputDataItemToTx converts an input data item that represents a transaction
// rlp hex encoded into a transaction
func inputDataItemToTx(inputDataItem string) (*types.Transaction, error) {
	tx := new(types.Transaction)

	inputDataItem = strings.TrimPrefix(inputDataItem, "0x")

	// Check if the string has an odd length
	if len(inputDataItem)%2 != 0 {
		// Prepend a '0' to make it even-length
		inputDataItem = "0" + inputDataItem
	}

	b, err := hex.DecodeString(inputDataItem)
	if err != nil {
		return nil, err
	}

	if err := tx.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return tx, nil
}
