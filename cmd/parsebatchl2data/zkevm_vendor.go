package parsebatchl2data

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/rs/zerolog/log"
)

// L2BlockRaw is the raw representation of a L2 block.
type L2BlockRaw struct {
	DeltaTimestamp  uint32
	IndexL1InfoTree uint32
	Transactions    []L2TxRaw
}

// BatchRawV2 is the  representation of a batch of transactions.
type BatchRawV2 struct {
	Blocks []L2BlockRaw
}

type ForcedBatchRawV2 struct {
	Transactions []L2TxRaw
}

type L2TxRaw struct {
	Tx                   *types.Transaction
	EfficiencyPercentage uint8
}

const (
	FORKID_DRAGONFRUIT = 5

	changeL2Block        = uint8(0x0b)
	c0            uint64 = 192 // 192 is c0. This value is defined by the rlp protocol
	double               = 2
	ether155V            = 27
	etherPre155V         = 35
	// Decoding constants
	headerByteLength uint64 = 1
	sLength          uint64 = 32
	rLength          uint64 = 32
	vLength          uint64 = 1
	ff               uint64 = 255 // max value of rlp header
	shortRlp         uint64 = 55  // length of the short rlp codification
	f7               uint64 = 247 // 192 + 55 = c0 + shortRlp

	// EfficiencyPercentageByteLength is the length of the effective percentage in bytes
	EfficiencyPercentageByteLength uint64 = 1
)

var (
	// ErrBatchV2DontStartWithChangeL2Block is returned when the batch start directly with a trsansaction (without a changeL2Block)
	ErrBatchV2DontStartWithChangeL2Block = errors.New("batch v2 must start with changeL2Block before Tx (suspect a V1 Batch or a ForcedBatch?))")
	// ErrInvalidBatchV2 is returned when the batch is invalid.
	ErrInvalidBatchV2 = errors.New("invalid batch v2")
	// ErrInvalidRLP is returned when the rlp is invalid.
	ErrInvalidRLP = errors.New("invalid rlp codification")
	// ErrInvalidData is the error when the raw txs is unexpected
	ErrInvalidData = errors.New("invalid data")
)

func deserializeUint32(txsData []byte, pos int) (int, uint32, error) {
	if len(txsData)-pos < 4 { // nolint:gomnd
		return 0, 0, fmt.Errorf("can't get u32 because not enough data: %w", ErrInvalidBatchV2)
	}
	return pos + 4, uint32(txsData[pos])<<24 | uint32(txsData[pos+1])<<16 | uint32(txsData[pos+2])<<8 | uint32(txsData[pos+3]), nil // nolint:gomnd
}

// decodeBlockHeader decodes a block header from a byte slice.
//
//	Extract: 4 bytes for deltaTimestamp + 4 bytes for indexL1InfoTree
func decodeBlockHeader(txsData []byte, pos int) (int, *L2BlockRaw, error) {
	var err error
	currentBlock := &L2BlockRaw{}
	pos, currentBlock.DeltaTimestamp, err = deserializeUint32(txsData, pos)
	if err != nil {
		return 0, nil, fmt.Errorf("can't get deltaTimestamp: %w", err)
	}
	pos, currentBlock.IndexL1InfoTree, err = deserializeUint32(txsData, pos)
	if err != nil {
		return 0, nil, fmt.Errorf("can't get leafIndex: %w", err)
	}

	return pos, currentBlock, nil
}

func decodeRLPListLengthFromOffset(txsData []byte, offset int) (uint64, error) {
	txDataLength := uint64(len(txsData))
	num := uint64(txsData[offset])
	if num < c0 { // c0 -> is a empty data
		return 0, fmt.Errorf("first byte of tx (%x) is < 0xc0: %w", num, ErrInvalidRLP)
	}
	length := num - c0
	if length > shortRlp { // If rlp is bigger than length 55
		// n is the length of the rlp data without the header (1 byte) for example "0xf7"
		pos64 := uint64(offset)
		lengthInByteOfSize := num - f7
		if (pos64 + headerByteLength + lengthInByteOfSize) > txDataLength {
			return 0, fmt.Errorf("not enough data to get length: %w", ErrInvalidRLP)
		}

		n, err := strconv.ParseUint(hex.EncodeToString(txsData[pos64+1:pos64+1+lengthInByteOfSize]), 16, 64) // +1 is the header. For example 0xf7
		if err != nil {
			return 0, fmt.Errorf("error parsing length value: %w", err)
		}
		// TODO: RLP specifications says length = n ??? that is wrong??
		length = n + num - f7 // num - f7 is the header. For example 0xf7
	}
	return length + headerByteLength, nil
}
func RlpFieldsToLegacyTx(fields [][]byte, v, r, s []byte) (tx *types.LegacyTx, err error) {
	const (
		fieldsSizeWithoutChainID = 6
		fieldsSizeWithChainID    = 7
	)

	if len(fields) < fieldsSizeWithoutChainID {
		return nil, types.ErrTxTypeNotSupported
	}

	nonce := big.NewInt(0).SetBytes(fields[0]).Uint64()
	gasPrice := big.NewInt(0).SetBytes(fields[1])
	gas := big.NewInt(0).SetBytes(fields[2]).Uint64()
	var to *common.Address

	if len(fields[3]) != 0 {
		tmp := common.BytesToAddress(fields[3])
		to = &tmp
	}
	value := big.NewInt(0).SetBytes(fields[4])
	data := fields[5]

	txV := big.NewInt(0).SetBytes(v)
	if len(fields) >= fieldsSizeWithChainID {
		chainID := big.NewInt(0).SetBytes(fields[6])

		// a = chainId * 2
		// b = v - 27
		// c = a + 35
		// v = b + c
		//
		// same as:
		// v = v-27+chainId*2+35
		a := new(big.Int).Mul(chainID, big.NewInt(double))
		b := new(big.Int).Sub(new(big.Int).SetBytes(v), big.NewInt(ether155V))
		c := new(big.Int).Add(a, big.NewInt(etherPre155V))
		txV = new(big.Int).Add(b, c)
	}

	txR := big.NewInt(0).SetBytes(r)
	txS := big.NewInt(0).SetBytes(s)

	return &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gas,
		To:       to,
		Value:    value,
		Data:     data,
		V:        txV,
		R:        txR,
		S:        txS,
	}, nil
}
func decodeTxRLP(txsData []byte, offset int) (int, *L2TxRaw, error) {
	var err error
	length, err := decodeRLPListLengthFromOffset(txsData, offset)
	if err != nil {
		return 0, nil, fmt.Errorf("can't get RLP length (offset=%d): %w", offset, err)
	}
	endPos := uint64(offset) + length + rLength + sLength + vLength + EfficiencyPercentageByteLength
	if endPos > uint64(len(txsData)) {
		return 0, nil, fmt.Errorf("can't get tx because not enough data (endPos=%d lenData=%d): %w",
			endPos, len(txsData), ErrInvalidBatchV2)
	}
	fullDataTx := txsData[offset:endPos]
	dataStart := uint64(offset) + length
	txInfo := txsData[offset:dataStart]
	rData := txsData[dataStart : dataStart+rLength]
	sData := txsData[dataStart+rLength : dataStart+rLength+sLength]
	vData := txsData[dataStart+rLength+sLength : dataStart+rLength+sLength+vLength]
	efficiencyPercentage := txsData[dataStart+rLength+sLength+vLength]
	var rlpFields [][]byte
	err = rlp.DecodeBytes(txInfo, &rlpFields)
	if err != nil {
		log.Error().Bytes("fullDataTx", fullDataTx).Bytes("tx", txInfo).Bytes("txReceived", txsData).Err(err)
		return 0, nil, err
	}
	legacyTx, err := RlpFieldsToLegacyTx(rlpFields, vData, rData, sData)
	if err != nil {
		return 0, nil, err
	}

	l2Tx := &L2TxRaw{
		Tx:                   types.NewTx(legacyTx),
		EfficiencyPercentage: efficiencyPercentage,
	}

	return int(endPos), l2Tx, err
}

// DecodeBatchV2 decodes a batch of transactions from a byte slice.
func DecodeBatchV2(txsData []byte) (*BatchRawV2, error) {
	// The transactions is not RLP encoded. Is the raw bytes in this form: 1 byte for the transaction type (always 0b for changeL2Block) + 4 bytes for deltaTimestamp + for bytes for indexL1InfoTree
	var err error
	var blocks []L2BlockRaw
	var currentBlock *L2BlockRaw
	pos := int(0)
	for pos < len(txsData) {
		switch txsData[pos] {
		case changeL2Block:
			if currentBlock != nil {
				blocks = append(blocks, *currentBlock)
			}
			pos, currentBlock, err = decodeBlockHeader(txsData, pos+1)
			if err != nil {
				return nil, fmt.Errorf("pos: %d can't decode new BlockHeader: %w", pos, err)
			}
		// by RLP definition a tx never starts with a 0x0b. So, if is not a changeL2Block
		// is a tx
		default:
			if currentBlock == nil {
				_, _, err = decodeTxRLP(txsData, pos)
				if err == nil {
					// There is no changeL2Block but have a valid RLP transaction
					return nil, ErrBatchV2DontStartWithChangeL2Block
				} else {
					// No changeL2Block and no valid RLP transaction
					return nil, fmt.Errorf("no ChangeL2Block neither valid Tx, batch malformed : %w", ErrInvalidBatchV2)
				}
			}
			var tx *L2TxRaw
			pos, tx, err = decodeTxRLP(txsData, pos)
			if err != nil {
				return nil, fmt.Errorf("can't decode transactions: %w", err)
			}

			currentBlock.Transactions = append(currentBlock.Transactions, *tx)
		}
	}
	if currentBlock != nil {
		blocks = append(blocks, *currentBlock)
	}
	return &BatchRawV2{blocks}, nil
}

// DecodeTxs extracts Transactions for its encoded form
func DecodeTxs(txsData []byte, forkID uint64) ([]*types.Transaction, []byte, []uint8, error) {
	// Process coded txs
	var pos uint64
	var txs []*types.Transaction
	var efficiencyPercentages []uint8
	txDataLength := uint64(len(txsData))
	if txDataLength == 0 {
		return txs, txsData, nil, nil
	}
	for pos < txDataLength {
		num, err := strconv.ParseUint(hex.EncodeToString(txsData[pos:pos+1]), 16, 64)
		if err != nil {
			return []*types.Transaction{}, txsData, []uint8{}, err
		}
		// First byte is the length and must be ignored
		if num < c0 {
			return []*types.Transaction{}, txsData, []uint8{}, ErrInvalidData
		}
		length := num - c0
		if length > shortRlp { // If rlp is bigger than length 55
			// n is the length of the rlp data without the header (1 byte) for example "0xf7"
			if (pos + 1 + num - f7) > txDataLength {
				return []*types.Transaction{}, txsData, []uint8{}, err
			}
			var n uint64
			n, err = strconv.ParseUint(hex.EncodeToString(txsData[pos+1:pos+1+num-f7]), 16, 64) // +1 is the header. For example 0xf7
			if err != nil {
				return []*types.Transaction{}, txsData, []uint8{}, err
			}
			if n+num < f7 {
				return []*types.Transaction{}, txsData, []uint8{}, ErrInvalidData
			}
			length = n + num - f7 // num - f7 is the header. For example 0xf7
		}

		endPos := pos + length + rLength + sLength + vLength + headerByteLength

		if forkID >= FORKID_DRAGONFRUIT {
			endPos += EfficiencyPercentageByteLength
		}

		if endPos > txDataLength {
			err = fmt.Errorf("endPos %d is bigger than txDataLength %d", endPos, txDataLength)
			return []*types.Transaction{}, txsData, []uint8{}, err
		}

		if endPos < pos {
			err = fmt.Errorf("endPos %d is smaller than pos %d", endPos, pos)
			return []*types.Transaction{}, txsData, []uint8{}, err
		}

		if endPos < pos {
			err = fmt.Errorf("endPos %d is smaller than pos %d", endPos, pos)
			return []*types.Transaction{}, txsData, []uint8{}, err
		}

		fullDataTx := txsData[pos:endPos]
		dataStart := pos + length + headerByteLength
		txInfo := txsData[pos:dataStart]
		rData := txsData[dataStart : dataStart+rLength]
		sData := txsData[dataStart+rLength : dataStart+rLength+sLength]
		vData := txsData[dataStart+rLength+sLength : dataStart+rLength+sLength+vLength]

		if forkID >= FORKID_DRAGONFRUIT {
			efficiencyPercentage := txsData[dataStart+rLength+sLength+vLength : endPos]
			efficiencyPercentages = append(efficiencyPercentages, efficiencyPercentage[0])
		}

		pos = endPos

		// Decode rlpFields
		var rlpFields [][]byte
		err = rlp.DecodeBytes(txInfo, &rlpFields)
		if err != nil {
			log.Error().Bytes("fullDataTx", fullDataTx).Bytes("tx", txInfo).Bytes("txReceived", txsData).Err(err)
			return []*types.Transaction{}, txsData, []uint8{}, ErrInvalidData
		}

		legacyTx, err := RlpFieldsToLegacyTx(rlpFields, vData, rData, sData)
		if err != nil {
			return []*types.Transaction{}, txsData, []uint8{}, err
		}

		tx := types.NewTx(legacyTx)
		txs = append(txs, tx)
	}
	return txs, txsData, efficiencyPercentages, nil
}

func DecodeForcedBatchV2(txsData []byte) (*ForcedBatchRawV2, error) {
	txs, _, efficiencyPercentages, err := DecodeTxs(txsData, 7)
	if err != nil {
		return nil, err
	}
	// Sanity check, this should never happen
	if len(efficiencyPercentages) != len(txs) {
		return nil, fmt.Errorf("error decoding len(efficiencyPercentages) != len(txs). len(efficiencyPercentages)=%d, len(txs)=%d : %w", len(efficiencyPercentages), len(txs), ErrInvalidRLP)
	}
	forcedBatch := ForcedBatchRawV2{}
	for i := range txs {
		forcedBatch.Transactions = append(forcedBatch.Transactions, L2TxRaw{
			Tx:                   txs[i],
			EfficiencyPercentage: efficiencyPercentages[i],
		})
	}
	return &forcedBatch, nil
}
