package util

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/0xPolygon/polygon-cli/rpctypes"
	"github.com/cenkalti/backoff"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/clique"
	"github.com/ethereum/go-ethereum/crypto"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"

	"github.com/ethereum/go-ethereum/core/types"
)

type (
	simpleRPCTransaction struct {
		Hash string `json:"hash"`
	}
	simpleRPCBlock struct {
		Number       string                 `json:"number"`
		Transactions []simpleRPCTransaction `json:"transactions"`
	}
	txpoolStatus struct {
		Pending any `json:"pending"`
		Queued  any `json:"queued"`
	}
)

func Ecrecover(block *types.Block) ([]byte, error) {
	header := block.Header()
	sigStart := len(header.Extra) - ethcrypto.SignatureLength
	if sigStart < 0 || sigStart > len(header.Extra) {
		return nil, fmt.Errorf("unable to recover signature")
	}
	signature := header.Extra[sigStart:]
	pubkey, err := ethcrypto.Ecrecover(clique.SealHash(header).Bytes(), signature)
	if err != nil {
		return nil, err
	}
	signer := ethcrypto.Keccak256(pubkey[1:])[12:]

	return signer, nil
}

func EcrecoverTx(tx *types.Transaction) ([]byte, error) {
	chainID := tx.ChainId()
	signer := types.LatestSignerForChainID(chainID)
	from, err := types.Sender(signer, tx)
	if err != nil {
		return nil, err
	}
	return from.Bytes(), nil
}

func GetBlockRange(ctx context.Context, from, to uint64, c *ethrpc.Client, onlyTxHashes bool) ([]*json.RawMessage, error) {
	blms := make([]ethrpc.BatchElem, 0)
	for i := from; i <= to; i = i + 1 {
		r := new(json.RawMessage)
		var err error
		blms = append(blms, ethrpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{"0x" + strconv.FormatUint(i, 16), !onlyTxHashes},
			Result: r,
			Error:  err,
		})
	}
	log.Trace().Uint64("start", from).Uint64("end", to).Msg("Fetching block range")

	err := c.BatchCallContext(ctx, blms)
	if err != nil {
		log.Error().Err(err).Msg("RPC issue fetching blocks")
		return nil, err
	}
	blocks := make([]*json.RawMessage, 0)

	for _, b := range blms {
		if b.Error != nil {
			return nil, b.Error
		}
		blocks = append(blocks, b.Result.(*json.RawMessage))

	}

	return blocks, nil
}

func GetBlockRangeInPages(ctx context.Context, from, to, pageSize uint64, c *ethrpc.Client, onlyTxHashes bool) ([]*json.RawMessage, error) {
	var allBlocks []*json.RawMessage

	for i := from; i <= to; i += pageSize {
		end := min(i+pageSize-1, to)

		blocks, err := GetBlockRange(ctx, i, end, c, onlyTxHashes)
		if err != nil {
			return nil, err
		}

		allBlocks = append(allBlocks, blocks...)
	}

	return allBlocks, nil
}

var getReceiptsByBlockIsSupported *bool

func GetReceipts(ctx context.Context, rawBlocks []*json.RawMessage, c *ethrpc.Client, batchSize uint64) ([]*json.RawMessage, error) {
	if getReceiptsByBlockIsSupported == nil {
		err := c.CallContext(ctx, nil, "eth_getBlockReceipts", "0x0")
		supported := err == nil
		getReceiptsByBlockIsSupported = &supported
	}

	if getReceiptsByBlockIsSupported != nil && *getReceiptsByBlockIsSupported {
		return getReceiptsByBlock(ctx, rawBlocks, c, batchSize)
	}

	return getReceiptsByTx(ctx, rawBlocks, c, batchSize)
}

func getReceiptsByBlock(ctx context.Context, rawBlocks []*json.RawMessage, c *ethrpc.Client, batchSize uint64) ([]*json.RawMessage, error) {
	var startBlock *string
	batchElements := make([]ethrpc.BatchElem, 0, len(rawBlocks))
	for _, rawBlock := range rawBlocks {
		var block simpleRPCBlock
		err := json.Unmarshal(*rawBlock, &block)
		if err != nil {
			return nil, err
		}
		batchElements = append(batchElements, ethrpc.BatchElem{
			Method: "eth_getBlockReceipts",
			Args:   []interface{}{block.Number},
			Result: new([]*json.RawMessage),
		})
		if startBlock == nil {
			startBlock = &block.Number
		}
	}
	if len(batchElements) == 0 {
		log.Debug().Int("Length of BatchElem", len(batchElements)).Msg("BatchElem is empty")
		return nil, nil
	}

	var start uint64 = 0
	for {
		last := false
		end := start + batchSize
		if int(end) >= len(batchElements) {
			last = true
			end = uint64(len(batchElements))
		}

		log.Trace().Str("startblock", *startBlock).Uint64("start", start).Uint64("end", end).Msg("Fetching receipt range")
		err := c.BatchCallContext(ctx, batchElements[start:end])
		if err != nil {
			log.Error().Err(err).Uint64("start", start).Uint64("end", end).Msg("RPC issue fetching receipts, have you checked the batch size limit of the RPC endpoint and adjusted the --batch-size flag?")
			break
		}
		start = end
		if last {
			break
		}
	}

	receipts := make([]*json.RawMessage, 0)
	for _, b := range batchElements {
		if b.Error != nil {
			log.Error().Err(b.Error).
				Interface("blockNumber", b.Args[0]).
				Msg("Block response err")
			return nil, b.Error
		}
		if b.Result == nil || reflect.ValueOf(b.Result).IsNil() {
			continue
		}
		rs := *(b.Result.(*[]*json.RawMessage))
		receipts = append(receipts, rs...)
	}
	if len(receipts) == 0 {
		log.Info().Msg("No receipts have been fetched")
		return nil, nil
	}
	log.Info().Int("blocks", len(rawBlocks)).Int("receipts", len(receipts)).Msg("Fetched tx receipts")
	return receipts, nil
}

func getReceiptsByTx(ctx context.Context, rawBlocks []*json.RawMessage, c *ethrpc.Client, batchSize uint64) ([]*json.RawMessage, error) {
	txHashes := make([]string, 0)
	txHashMap := make(map[string]string, 0)
	for _, rb := range rawBlocks {
		var block simpleRPCBlock
		err := json.Unmarshal(*rb, &block)
		if err != nil {
			return nil, err
		}
		for _, tx := range block.Transactions {
			txHashes = append(txHashes, tx.Hash)
			txHashMap[tx.Hash] = block.Number
		}
	}
	if len(txHashes) == 0 {
		return nil, nil
	}

	blms := make([]ethrpc.BatchElem, 0)
	blmsBlockMap := make(map[int]string, 0)
	for i, tx := range txHashes {
		r := new(json.RawMessage)
		var err error
		blms = append(blms, ethrpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []interface{}{tx},
			Result: r,
			Error:  err,
		})
		blmsBlockMap[i] = txHashMap[tx]
	}

	if len(blms) == 0 {
		log.Debug().Int("Length of BatchElem", len(blms)).Msg("BatchElem is empty")
		return nil, nil
	}

	var start uint64 = 0
	for {
		last := false
		end := start + batchSize
		if int(end) > len(blms) {
			last = true
			end = uint64(len(blms))
		}

		log.Trace().Str("startblock", blmsBlockMap[int(start)]).Uint64("start", start).Uint64("end", end).Msg("Fetching tx receipt range")
		// json: cannot unmarshal object into Go value of type []rpc.jsonrpcMessage
		// The error occurs when we call batchcallcontext with a single transaction for some reason.
		// polycli dumpblocks -c 1 http://127.0.0.1:9209/ 34457958 34458108
		// To handle this I'm making an exception when start and end are equal to make a single call.
		if start == end {
			log.Trace().Int("length", len(blmsBlockMap)).Msg("Test Jesse")
			if len(blmsBlockMap) == int(start) {
				start = start - 1
			}
			err := c.CallContext(ctx, &blms[start].Result, "eth_getTransactionReceipt", blms[start].Args[0])
			if err != nil {
				log.Error().Err(err).Uint64("start", start).Uint64("end", end).Msg("RPC issue fetching single receipt")
				return nil, err
			}
			break
		}

		err := c.BatchCallContext(ctx, blms[start:end])
		if err != nil {
			log.Error().Err(err).Str("randtx", txHashes[0]).Uint64("start", start).Uint64("end", end).Msg("RPC issue fetching receipts, have you checked the batch size limit of the RPC endpoint and adjusted the --batch-size flag?")
			break
		}
		start = end
		if last {
			break
		}
	}

	receipts := make([]*json.RawMessage, 0)

	for _, b := range blms {
		if b.Error != nil {
			log.Error().Err(b.Error).Msg("Block response err")
			return nil, b.Error
		}
		receipts = append(receipts, b.Result.(*json.RawMessage))
	}
	if len(receipts) == 0 {
		log.Info().Msg("No receipts have been fetched")
		return nil, nil
	}
	log.Info().Int("hashes", len(txHashes)).Int("receipts", len(receipts)).Msg("Fetched tx receipts")
	return receipts, nil
}

func GetTxPoolStatus(rpc *ethrpc.Client) (uint64, uint64, error) {
	var status = new(txpoolStatus)
	err := rpc.Call(status, "txpool_status")
	if err != nil {
		return 0, 0, err
	}
	pendingCount, err := tryCastToUint64(status.Pending)
	if err != nil {
		return 0, 0, err
	}
	queuedCount, err := tryCastToUint64(status.Queued)
	if err != nil {
		return pendingCount, 0, err
	}

	return pendingCount, queuedCount, nil
}

func GetZkEVMBatches(rpc *ethrpc.Client) (uint64, uint64, uint64, error) {
	trustedBatches, err := getZkEVMBatch(rpc, trusted)
	if err != nil {
		return 0, 0, 0, err
	}

	virtualBatches, err := getZkEVMBatch(rpc, virtual)
	if err != nil {
		return trustedBatches, 0, 0, err
	}

	verifiedBatches, err := getZkEVMBatch(rpc, verified)
	if err != nil {
		return trustedBatches, virtualBatches, 0, err
	}

	return trustedBatches, virtualBatches, verifiedBatches, nil
}

func GetForkID(rpc *ethrpc.Client) (uint64, error) {
	var raw interface{}
	if err := rpc.Call(&raw, "zkevm_getForkId"); err != nil {
		return 0, err
	}
	forkID, err := hexutil.DecodeUint64(fmt.Sprintf("%v", raw))
	if err != nil {
		return 0, err
	}
	return forkID, nil
}

func GetRollupAddress(rpc *ethrpc.Client) (string, error) {
	var raw interface{}
	if err := rpc.Call(&raw, "zkevm_getRollupAddress"); err != nil {
		return "", err
	}
	rollupAddress := fmt.Sprintf("%v", raw)

	return rollupAddress, nil
}

func GetRollupManagerAddress(rpc *ethrpc.Client) (string, error) {
	var raw interface{}
	if err := rpc.Call(&raw, "zkevm_getRollupManagerAddress"); err != nil {
		return "", err
	}
	rollupManagerAddress := fmt.Sprintf("%v", raw)

	return rollupManagerAddress, nil
}

type batch string

const (
	trusted  batch = "zkevm_batchNumber"
	virtual  batch = "zkevm_virtualBatchNumber"
	verified batch = "zkevm_verifiedBatchNumber"
)

func getZkEVMBatch(rpc *ethrpc.Client, batchType batch) (uint64, error) {
	var raw interface{}
	if err := rpc.Call(&raw, string(batchType)); err != nil {
		return 0, err
	}
	batch, err := hexutil.DecodeUint64(fmt.Sprintf("%v", raw))
	if err != nil {
		return 0, err
	}
	return batch, nil
}

func tryCastToUint64(val any) (uint64, error) {
	switch t := val.(type) {
	case float64:
		return uint64(t), nil
	case string:
		return convHexToUint64(t)
	default:
		return 0, fmt.Errorf("the value %v couldn't be marshalled to uint64", t)

	}
}
func convHexToUint64(hexString string) (uint64, error) {
	hexString = strings.TrimPrefix(hexString, "0x")
	if len(hexString)%2 != 0 {
		hexString = "0" + hexString
	}

	result, err := strconv.ParseUint(hexString, 16, 64)
	if err != nil {
		return 0, err
	}
	return uint64(result), nil
}

// BlockUntilSuccessfulFn is designed to wait until a specified number of Ethereum blocks have been
// mined, periodically checking for the completion of a given function within each block interval.
type BlockUntilSuccessfulFn func(ctx context.Context, c *ethclient.Client, f func() error) error

func BlockUntilSuccessful(ctx context.Context, c *ethclient.Client, retryable func() error) error {
	// this function use to be very complicated (and not work). I'm dumbing this down to a basic time based retryable which should work 99% of the time
	b := backoff.WithContext(backoff.WithMaxRetries(backoff.NewConstantBackOff(5*time.Second), 24), ctx)
	return backoff.Retry(retryable, b)
}

func WrapDeployedCode(deployedBytecode string, storageBytecode string) string {
	deployedBytecode = strings.ToLower(strings.TrimPrefix(deployedBytecode, "0x"))
	storageBytecode = strings.ToLower(strings.TrimPrefix(storageBytecode, "0x"))

	codeCopySize := len(deployedBytecode) / 2
	codeCopyOffset := (len(storageBytecode) / 2) + 13 + 8 // 13 for CODECOPY + 8 for RETURN

	return fmt.Sprintf(
		"0x%s"+ // storage initialization code
			"63%08x"+ // PUSH4 to indicate the size of the data that should be copied into memory
			"63%08x"+ // PUSH4 to indicate the offset in the call data to start the copy
			"6000"+ // PUSH1 00 to indicate the destination offset in memory
			"39"+ // CODECOPY
			"63%08x"+ // PUSH4 to indicate the size of the data to be returned from memory
			"6000"+ // PUSH1 00 to indicate that it starts from offset 0
			"f3"+ // RETURN
			"%s", // CODE starts here.
		storageBytecode, codeCopySize, codeCopyOffset, codeCopySize, deployedBytecode)
}

func GetHexString(data any) string {
	var result string
	if reflect.TypeOf(data).Kind() == reflect.Float64 {
		result = fmt.Sprintf("%x", int64(data.(float64)))
	} else if reflect.TypeOf(data).Kind() == reflect.String {
		if strings.HasPrefix(data.(string), "0x") {
			result = strings.TrimPrefix(data.(string), "0x")
		} else {
			result = data.(string)
		}
	} else {
		log.Fatal().Any("data", data).Msg("unknown storage data type")
	}
	if len(result)%2 != 0 {
		result = "0" + result
	}
	return strings.ToLower(result)
}

func GetChainID(ctx context.Context, ec *ethrpc.Client) (*big.Int, error) {
	var chainIDHex string
	err := ec.CallContext(ctx, &chainIDHex, "eth_chainId")
	if err != nil {
		return nil, err
	}
	chainID, err := hexutil.DecodeBig(chainIDHex)
	if err != nil {
		return nil, err
	}
	return chainID, nil
}

// HeaderByBlockNumber retrieves a block header using rpc.BlockNumber.
// If blockNum is nil, it defaults to latest block.
// Examples:
//   - HeaderByBlockNumber(ctx, client, nil) // latest
//   - HeaderByBlockNumber(ctx, client, &rpc.LatestBlockNumber)
//   - HeaderByBlockNumber(ctx, client, &rpc.FinalizedBlockNumber)
//   - num := rpc.BlockNumber(12345); HeaderByBlockNumber(ctx, client, &num)
func HeaderByBlockNumber(ctx context.Context, ec *ethrpc.Client, blockNum *ethrpc.BlockNumber) (*types.Header, error) {
	var blockParam string
	if blockNum == nil {
		blockParam = ethrpc.LatestBlockNumber.String()
	} else {
		blockParam = blockNum.String()
	}

	var raw json.RawMessage
	err := ec.CallContext(ctx, &raw, "eth_getBlockByNumber", blockParam, false)
	if err != nil {
		return nil, err
	}

	var block types.Header
	err = json.Unmarshal(raw, &block)
	if err != nil {
		return nil, err
	}
	return &block, nil
}

// GetSenderFromTx recovers the sender address from a transaction without using types.Transaction
func GetSenderFromTx(ctx context.Context, tx rpctypes.PolyTransaction) (common.Address, error) {
	// Get transaction type
	txType := tx.Type()

	// For non-standard transaction types, we assume the sender is already set
	if txType > 2 {
		return tx.From(), nil
	}

	// Get transaction fields
	chainID := tx.ChainID()
	nonce := tx.Nonce()
	value := tx.Value()
	gas := tx.Gas()
	to := tx.To()
	data := tx.Data()
	v := tx.V()
	r := tx.R()
	s := tx.S()

	// Calculate the signing hash based on transaction type
	var sigHash []byte
	var err error

	switch txType {
	case 0: // Legacy transaction
		sigHash, err = calculateLegacySigningHash(chainID, nonce, tx.GasPrice(), gas, to, value, data)
	case 1: // EIP-2930 (Access List)
		// For now, we can try with empty access list
		// If you need full support, you'll need to add AccessList to PolyTransaction interface
		sigHash, err = calculateEIP2930SigningHash(chainID, nonce, tx.GasPrice(), gas, to, value, data, []interface{}{})
	case 2: // EIP-1559
		maxPriorityFee := big.NewInt(int64(tx.MaxPriorityFeePerGas()))
		maxFee := big.NewInt(int64(tx.MaxFeePerGas()))
		sigHash, err = calculateEIP1559SigningHash(chainID, nonce, maxPriorityFee, maxFee, gas, to, value, data)
	default:
		return common.Address{}, fmt.Errorf("unsupported transaction type: %d (0x%x)", txType, txType)
	}

	if err != nil {
		return common.Address{}, fmt.Errorf("failed to calculate signing hash: %w", err)
	}

	// Normalize v value for recovery
	var recoveryID byte
	if txType == 0 {
		// Legacy transaction with EIP-155
		if chainID > 0 {
			// EIP-155: v = chainId * 2 + 35 + {0,1}
			// Extract recovery id: recoveryID = v - (chainId * 2 + 35)
			vBig := new(big.Int).Set(v)
			vBig.Sub(vBig, big.NewInt(35))
			vBig.Sub(vBig, new(big.Int).Mul(big.NewInt(int64(chainID)), big.NewInt(2)))
			recoveryID = byte(vBig.Uint64())
		} else {
			// Pre-EIP-155: v is 27 or 28
			recoveryID = byte(v.Uint64() - 27)
		}
	} else {
		// EIP-2930 and EIP-1559: v is 0 or 1 (or 27/28)
		vVal := v.Uint64()
		if vVal >= 27 {
			recoveryID = byte(vVal - 27)
		} else {
			recoveryID = byte(vVal)
		}
	}

	// Validate recoveryID
	if recoveryID > 1 {
		return common.Address{}, fmt.Errorf("invalid recovery id: %d (v=%s, chainID=%d, type=%d)", recoveryID, v.String(), chainID, txType)
	}

	// Build signature in the [R || S || V] format (65 bytes)
	sig := make([]byte, 65)
	// Use FillBytes to ensure proper padding with leading zeros
	r.FillBytes(sig[0:32])
	s.FillBytes(sig[32:64])
	sig[64] = recoveryID

	// Recover public key from signature using go-ethereum's crypto package
	pubKey, err := crypto.Ecrecover(sigHash, sig)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to recover public key: %w", err)
	}

	// Derive address from public key
	// The public key returned by Ecrecover is 65 bytes: [0x04 || X || Y]
	// We hash the X and Y coordinates (skip first byte) and take last 20 bytes
	hash := crypto.Keccak256(pubKey[1:])
	address := common.BytesToAddress(hash[12:])

	return address, nil
}

// calculateLegacySigningHash calculates the signing hash for legacy (type 0) transactions
func calculateLegacySigningHash(chainID uint64, nonce uint64, gasPrice *big.Int, gas uint64, to common.Address, value *big.Int, data []byte) ([]byte, error) {
	var items []interface{}

	// Handle contract creation (to = zero address)
	var toPtr *common.Address
	if to != (common.Address{}) {
		toPtr = &to
	}

	if chainID > 0 {
		// EIP-155: RLP([nonce, gasPrice, gas, to, value, data, chainId, 0, 0])
		items = []interface{}{
			nonce,
			gasPrice,
			gas,
			toPtr,
			value,
			data,
			chainID,
			uint(0),
			uint(0),
		}
	} else {
		// Pre-EIP-155: RLP([nonce, gasPrice, gas, to, value, data])
		items = []interface{}{
			nonce,
			gasPrice,
			gas,
			toPtr,
			value,
			data,
		}
	}

	encoded, err := rlp.EncodeToBytes(items)
	if err != nil {
		return nil, fmt.Errorf("failed to RLP encode legacy transaction: %w", err)
	}
	return crypto.Keccak256(encoded), nil
}

// calculateEIP2930SigningHash calculates the signing hash for EIP-2930 (type 1) transactions
func calculateEIP2930SigningHash(chainID uint64, nonce uint64, gasPrice *big.Int, gas uint64, to common.Address, value *big.Int, data []byte, accessList []interface{}) ([]byte, error) {
	var toPtr *common.Address
	if to != (common.Address{}) {
		toPtr = &to
	}

	// EIP-2930: keccak256(0x01 || rlp([chainId, nonce, gasPrice, gas, to, value, data, accessList]))
	items := []interface{}{
		chainID,
		nonce,
		gasPrice,
		gas,
		toPtr,
		value,
		data,
		accessList,
	}

	encoded, err := rlp.EncodeToBytes(items)
	if err != nil {
		return nil, fmt.Errorf("failed to RLP encode EIP-2930 transaction: %w", err)
	}

	// Prepend transaction type byte (0x01)
	typedData := append([]byte{0x01}, encoded...)
	return crypto.Keccak256(typedData), nil
}

// calculateEIP1559SigningHash calculates the signing hash for EIP-1559 (type 2) transactions
func calculateEIP1559SigningHash(chainID uint64, nonce uint64, maxPriorityFee, maxFee *big.Int, gas uint64, to common.Address, value *big.Int, data []byte) ([]byte, error) {
	var toPtr *common.Address
	if to != (common.Address{}) {
		toPtr = &to
	}

	// EIP-1559: keccak256(0x02 || rlp([chainId, nonce, maxPriorityFeePerGas, maxFeePerGas, gas, to, value, data, accessList]))
	items := []interface{}{
		chainID,
		nonce,
		maxPriorityFee,
		maxFee,
		gas,
		toPtr,
		value,
		data,
		[]interface{}{}, // empty access list
	}

	encoded, err := rlp.EncodeToBytes(items)
	if err != nil {
		return nil, fmt.Errorf("failed to RLP encode EIP-1559 transaction: %w", err)
	}

	// Prepend transaction type byte (0x02)
	typedData := append([]byte{0x02}, encoded...)
	return crypto.Keccak256(typedData), nil
}
