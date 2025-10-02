package parsebatchl2data

import (
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
)

var (
	//go:embed usage.md
	usage         string
	inputFileName *string
)

var ParseBatchL2Data = &cobra.Command{
	Use:     "parse-batch-l2-data [flags]",
	Aliases: []string{"parsebatchl2data"},
	Short:   "Convert batch l2 data into an ndjson stream",
	Long:    usage,
	RunE: func(cmd *cobra.Command, args []string) error {

		rawData, err := getInputData(cmd, args)
		if err != nil {
			return err
		}
		batchL2Data := strings.TrimSpace(strings.TrimPrefix(string(rawData), "0x"))
		rawBatchL2Data, err := hex.DecodeString(batchL2Data)
		if err != nil {
			log.Err(err).Msg("Unable to hex decode batch l2 data")
			return err
		}

		rawBatch, err := DecodeBatchV2(rawBatchL2Data)
		if err != nil {
			log.Error().Err(err).Msg("unable to decode l2 batch data")
			tryRawBatch(rawBatchL2Data)
			return err
		}

		blocks := rawBatch.Blocks

		for _, l2RawBlock := range blocks {
			blockData := struct {
				IndexL1InfoTree uint32
				DeltaTimestamp  uint32
			}{l2RawBlock.IndexL1InfoTree, l2RawBlock.DeltaTimestamp}
			blockDataBytes, err := json.Marshal(blockData)
			if err != nil {
				log.Err(err).Msg("unable to marshal block data")
				return err
			}
			fmt.Println(string(blockDataBytes))

			for i := range l2RawBlock.Transactions {
				if err := printTxData(&l2RawBlock.Transactions[i]); err != nil {
					log.Error().Err(err).Int("tx_index", i).Msg("Failed to print transaction data")
				}
			}
		}
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	flagSet := ParseBatchL2Data.Flags()
	inputFileName = flagSet.String("file", "", "Provide a file with the key information ")
}

func getInputData(cmd *cobra.Command, args []string) ([]byte, error) {
	if inputFileName != nil && *inputFileName != "" {
		return os.ReadFile(*inputFileName)
	}

	if len(args) > 1 {
		concat := strings.Join(args[1:], " ")
		return []byte(concat), nil
	}

	return io.ReadAll(os.Stdin)
}

func printTxData(rawL2Tx *L2TxRaw) error {
	signer := types.NewEIP155Signer(rawL2Tx.Tx.ChainId())
	sender, err := signer.Sender(rawL2Tx.Tx)
	if err != nil {
		log.Error().Err(err).Msg("unable to reccover sender")
		return err
	}
	jsonTx, err := rawL2Tx.Tx.MarshalJSON()
	if err != nil {
		log.Error().Err(err).Msg("unable to json marshal tx")
		return err
	}
	txMap := make(map[string]string, 0)
	err = json.Unmarshal(jsonTx, &txMap)
	if err != nil {
		log.Error().Err(err).Msg("unable to remarshal json tx")
		return err
	}
	txMap["from"] = sender.String()
	jsonTx, err = json.Marshal(txMap)
	if err != nil {
		log.Error().Err(err).Msg("unable to marhshal tx with from")
		return err
	}

	fmt.Println(string(jsonTx))
	return nil
}

func tryRawBatch(rawBatchL2Data []byte) {
	rawBatch, err := DecodeForcedBatchV2(rawBatchL2Data)
	if err != nil {
		log.Error().Err(err).Msg("unable to decode raw l2 batch data")
		return
	}
	for i, t := range rawBatch.Transactions {
		if err := printTxData(&t); err != nil {
			log.Error().Err(err).Int("tx_index", i).Msg("Failed to print transaction data in forced batch")
		}
	}
}
