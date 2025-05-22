package publish

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"iter"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
)

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
			return dataFromArgs(args)
		}
	}

	// if we get here, we're going to assume the user wants to load transactions
	// from stdin or from the command line
	return dataFromStdin()
}

func dataFromArgs(args []string) (iter.Seq[string], inputDataSource, error) {
	fmt.Println("Reading data from args")

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
	fmt.Println("Reading data from file: ", filename)

	return func(yield func(string) bool) {
		s := bufio.NewScanner(f)
		for s.Scan() {
			if !yield(s.Text()) {
				return
			}
		}
		// Ensure the file is closed after the function exits
		defer f.Close()
		if err := s.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "error scanning file:", err)
		}
	}, InputDataSourceFile, nil
}

// dataFromStdin returns an iterator that reads lines from stdin
func dataFromStdin() (iter.Seq[string], inputDataSource, error) {
	fmt.Println("Reading data from Stdin, type the transactions you want to publish:")

	return func(yield func(string) bool) {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			if !yield(s.Text()) {
				return
			}
		}
		if err := s.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "error scanning stdin:", err)
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
