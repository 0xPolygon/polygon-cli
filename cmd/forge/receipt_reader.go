package forge

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/maticnetwork/polygon-cli/rpctypes"
)

type (
	ReceiptReader interface {
		ReadReceipt() (*rpctypes.RawTxReceipt, error)
	}
	JSONReceiptReader struct {
		scanner *bufio.Scanner
	}
)

// OpenReceiptReader returns a receipt reader object which can be used to read the
// file. It will return a mode specific receipt reader.
func OpenReceiptReader(file string, mode string) (ReceiptReader, error) {
	receiptsFile, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("unable to open %s receipts file: %w", file, err)
	}

	switch mode {
	case "json":
		buf := make([]byte, maxCapacity)
		scanner := bufio.NewScanner(receiptsFile)
		scanner.Buffer(buf, maxCapacity)

		receiptsReader := JSONReceiptReader{
			scanner: scanner,
		}
		return &receiptsReader, nil

	default:
		return nil, fmt.Errorf("invalid mode: %s", mode)
	}
}

func (receiptsReader *JSONReceiptReader) ReadReceipt() (*rpctypes.RawTxReceipt, error) {
	if !receiptsReader.scanner.Scan() {
		return nil, BlockReadEOF
	}

	rawTxBytes := receiptsReader.scanner.Bytes()
	var raw rpctypes.RawTxReceipt
	err := json.Unmarshal(rawTxBytes, &raw)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal file receipt: %w - %s", err, string(rawTxBytes))
	}

	return &raw, nil
}
