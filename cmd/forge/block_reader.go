package forge

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/maticnetwork/polygon-cli/proto/gen/pb"
	"github.com/maticnetwork/polygon-cli/rpctypes"
	"google.golang.org/protobuf/proto"
)

type (
	BlockReader interface {
		ReadBlock() (rpctypes.PolyBlock, error)
	}
	JSONBlockReader struct {
		scanner *bufio.Scanner
	}
	ProtoBlockReader struct {
		file   *os.File
		offset int64
	}
)

// OpenBlockReader returns a block reader object which can be used to read the
// file. It will return a mode specific block reader.
func OpenBlockReader(file string, mode string) (BlockReader, error) {
	blockFile, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("unable to open %s blocks file: %w", file, err)
	}

	switch mode {
	case "json":
		maxCapacity := 5 * 1024 * 1024
		buf := make([]byte, maxCapacity)
		scanner := bufio.NewScanner(blockFile)
		scanner.Buffer(buf, maxCapacity)

		br := JSONBlockReader{
			scanner: scanner,
		}
		return &br, nil

	case "proto":
		br := ProtoBlockReader{
			file: blockFile,
		}
		return &br, nil

	default:
		return nil, fmt.Errorf("invalid mode: %s", mode)
	}
}

func (blockReader *JSONBlockReader) ReadBlock() (rpctypes.PolyBlock, error) {
	if !blockReader.scanner.Scan() {
		return nil, BlockReadEOF
	}

	rawBlockBytes := blockReader.scanner.Bytes()
	var raw rpctypes.RawBlockResponse
	err := json.Unmarshal(rawBlockBytes, &raw)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal file block: %w - %s", err, string(rawBlockBytes))
	}
	return rpctypes.NewPolyBlock(&raw), nil
}

func (blockReader *ProtoBlockReader) ReadBlock() (rpctypes.PolyBlock, error) {
	// reading the length of the encoded item before reading each item
	buf := make([]byte, 4)
	if _, err := blockReader.file.ReadAt(buf, blockReader.offset); err != nil {
		return nil, err
	}
	itemSize := binary.LittleEndian.Uint32(buf)
	blockReader.offset += 4

	// reading the actual encoded item
	item := make([]byte, itemSize)
	if _, err := blockReader.file.ReadAt(item, blockReader.offset); err != nil {
		return nil, err
	}

	block := &pb.Block{}
	if err := proto.Unmarshal(item, block); err != nil {
		return nil, err
	}

	blockReader.offset += int64(itemSize)

	txs := []rpctypes.RawTransactionResponse{}
	for _, tx := range block.Transactions {
		to := ""
		if tx.To != nil {
			to = *tx.To
		}

		txs = append(txs, rpctypes.RawTransactionResponse{
			BlockHash:        rpctypes.RawData32Response(tx.BlockHash),
			BlockNumber:      rpctypes.RawQuantityResponse(tx.BlockNumber),
			From:             rpctypes.RawData20Response(tx.From),
			Gas:              rpctypes.RawQuantityResponse(tx.Gas),
			GasPrice:         rpctypes.RawQuantityResponse(tx.GasPrice),
			Hash:             rpctypes.RawData32Response(tx.Hash),
			Input:            rpctypes.RawDataResponse(tx.Input),
			Nonce:            rpctypes.RawQuantityResponse(tx.Nonce),
			To:               rpctypes.RawData20Response(to),
			TransactionIndex: rpctypes.RawQuantityResponse(tx.TransactionIndex),
			Value:            rpctypes.RawQuantityResponse(tx.Value),
			V:                rpctypes.RawQuantityResponse(tx.V),
			R:                rpctypes.RawQuantityResponse(tx.R),
			S:                rpctypes.RawQuantityResponse(tx.S),
			Type:             rpctypes.RawQuantityResponse(tx.Type),
		})
	}

	uncles := []rpctypes.RawData32Response{}
	for _, uncle := range block.Uncles {
		uncles = append(uncles, rpctypes.RawData32Response(uncle))
	}

	raw := rpctypes.RawBlockResponse{
		Number:           rpctypes.RawQuantityResponse(block.Number),
		Hash:             rpctypes.RawData32Response(block.Hash),
		ParentHash:       rpctypes.RawData32Response(block.ParentHash),
		Nonce:            rpctypes.RawData8Response(block.Nonce),
		SHA3Uncles:       rpctypes.RawData32Response(block.Sha3Uncles),
		LogsBloom:        rpctypes.RawData256Response(block.LogsBloom),
		TransactionsRoot: rpctypes.RawData32Response(block.TransactionsRoot),
		StateRoot:        rpctypes.RawData32Response(block.StateRoot),
		ReceiptsRoot:     rpctypes.RawData32Response(block.ReceiptsRoot),
		Miner:            rpctypes.RawData20Response(block.Miner),
		Difficulty:       rpctypes.RawQuantityResponse(block.Difficulty),
		TotalDifficulty:  rpctypes.RawQuantityResponse(block.TotalDifficulty),
		ExtraData:        rpctypes.RawDataResponse(block.ExtraData),
		Size:             rpctypes.RawQuantityResponse(block.Size),
		GasLimit:         rpctypes.RawQuantityResponse(block.GasLimit),
		GasUsed:          rpctypes.RawQuantityResponse(block.GasUsed),
		Timestamp:        rpctypes.RawQuantityResponse(block.Timestamp),
		Transactions:     txs,
		Uncles:           uncles,
		BaseFeePerGas:    rpctypes.RawQuantityResponse(block.BaseFeePerGas),
	}

	return rpctypes.NewPolyBlock(&raw), nil
}

func ReadProtoFromFile(filepath string) ([][]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	var offset int64
	content := make([][]byte, 0)

	for {
		// reading the length of the encoded item before reading each item
		buf := make([]byte, 4)
		if _, err := file.ReadAt(buf, offset); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		itemSize := binary.LittleEndian.Uint32(buf)
		offset += 4

		// reading the actual encoded item
		item := make([]byte, itemSize)
		if _, err := file.ReadAt(item, offset); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		content = append(content, item)
		offset += int64(itemSize)
	}

	return content, nil
}
