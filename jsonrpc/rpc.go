package jsonrpc

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/rs/zerolog/log"
)

type (
	Client struct {
		httpClient *http.Client
		username   string
		password   string
		hasAuth    bool
		counter    uint64
	}
	ChainClient struct {
		RPCClient  *Client
		URL        string
		PrivateKey *ecdsa.PrivateKey
		ChainID    *big.Int
	}

	RPCBase struct {
		JSONRPC string `json:"jsonrpc"`
		ID      uint64 `json:"id"`
	}
	RPCReq struct {
		Method string `json:"method"`
		Params []any  `json:"params"`
		RPCBase
	}
	RPCError struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
		Data    any    `json:"data"`
	}
	RPCResp struct {
		Result any      `json:"result"`
		Error  RPCError `json:"error"`
		RPCBase
	}
	RPCBlockResp struct {
		Result RawBlockResponse `json:"result"`
		Error  RPCError         `json:"error"`
		RPCBase
	}
	RPCReceiptResp struct {
		Result RawTxReceipt `json:"result"`
		Error  RPCError
		RPCBase
	}

	RawQuantityResponse string
	RawDataResponse     string
	RawData8Response    string
	RawData20Response   string
	RawData32Response   string
	RawData256Response  string

	RawTransactionResponse struct {
		// blockHash: DATA, 32 Bytes - hash of the block where this transaction was in. null when its pending.
		BlockHash RawData32Response `json:"blockHash"`

		// blockNumber: QUANTITY - block number where this transaction was in. null when its pending.
		BlockNumber RawQuantityResponse `json:"blockNumber"`

		// from: DATA, 20 Bytes - address of the sender.
		From RawData20Response `json:"from"`

		// gas: QUANTITY - gas provided by the sender.
		Gas RawQuantityResponse `json:"gas"`

		// gasPrice: QUANTITY - gas price provided by the sender in Wei.
		GasPrice RawQuantityResponse `json:"gasPrice"`

		// hash: DATA, 32 Bytes - hash of the transaction.
		Hash RawData32Response `json:"hash"`

		// input: DATA - the data send along with the transaction.
		Input RawDataResponse `json:"input"`

		// nonce: QUANTITY - the number of transactions made by the sender prior to this one.
		Nonce RawQuantityResponse `json:"nonce"`

		// to: DATA, 20 Bytes - address of the receiver. null when its a contract creation transaction.
		To RawData20Response `json:"to"`

		// transactionIndex: QUANTITY - integer of the transactions index position in the block. null when its pending.
		TransactionIndex RawQuantityResponse `json:"transactionIndex"`

		// value: QUANTITY - value transferred in Wei.
		Value RawQuantityResponse `json:"value"`

		// v: QUANTITY - ECDSA recovery id
		V RawQuantityResponse `json:"v"`

		// r: QUANTITY - ECDSA signature r
		R RawQuantityResponse `json:"r"`

		// s: QUANTITY - ECDSA signature s
		S RawQuantityResponse `json:"s"`
	}

	RawBlockResponse struct {
		// number: QUANTITY - the block number. null when its pending block.
		Number RawQuantityResponse `json:"number"`

		// hash: DATA, 32 Bytes - hash of the block. null when its pending block.
		Hash RawData32Response `json:"hash"`

		// parentHash: DATA, 32 Bytes - hash of the parent block.
		ParentHash RawData32Response `json:"parentHash"`

		// nonce: DATA, 8 Bytes - hash of the generated proof-of-work. null when its pending block.
		Nonce RawData8Response `json:"nonce"`

		// sha3Uncles: DATA, 32 Bytes - SHA3 of the uncles data in the block.
		SHA3Uncles RawData32Response `json:"sha3Uncles"`

		// logsBloom: DATA, 256 Bytes - the bloom filter for the logs of the block. null when its pending block.
		LogsBloom RawData256Response `json:"logsBloom"`

		// transactionsRoot: DATA, 32 Bytes - the root of the transaction trie of the block.
		TransactionsRoot RawData32Response `json:"transactionsRoot"`

		// stateRoot: DATA, 32 Bytes - the root of the final state trie of the block.
		StateRoot RawData32Response `json:"stateRoot"`

		// receiptsRoot: DATA, 32 Bytes - the root of the receipts trie of the block.
		ReceiptsRoot RawData32Response `json:"receiptsRoot"`

		// miner: DATA, 20 Bytes - the address of the beneficiary to whom the mining rewards were given.
		Miner RawData20Response `json:"miner"`

		// difficulty: QUANTITY - integer of the difficulty for this block.
		Difficulty RawQuantityResponse `json:"difficulty"`

		// totalDifficulty: QUANTITY - integer of the total difficulty of the chain until this block.
		TotalDifficulty RawQuantityResponse `json:"totalDifficulty"`

		// extraData: DATA - the "extra data" field of this block.
		ExtraData RawDataResponse `json:"extraData"`

		// size: QUANTITY - integer the size of this block in bytes.
		Size RawQuantityResponse `json:"size"`

		// gasLimit: QUANTITY - the maximum gas allowed in this block.
		GasLimit RawQuantityResponse `json:"gasLimit"`

		// gasUsed: QUANTITY - the total used gas by all transactions in this block.
		GasUsed RawQuantityResponse `json:"gasUsed"`

		// timestamp: QUANTITY - the unix timestamp for when the block was collated.
		Timestamp RawQuantityResponse `json:"timestamp"`

		// transactions: Array - Array of transaction objects, or 32 Bytes transaction hashes depending on the last given parameter.
		Transactions []RawTransactionResponse `json:"transactions"`

		// uncles: Array - Array of uncle hashes.
		Uncles []RawQuantityResponse `json:"uncles"`

		// baseFeePerGass: QUANTITY - fixed per block fee
		BaseFeePerGas RawQuantityResponse `json:"baseFeePerGas"`
	}
	RawTxReceipt struct {
		// transactionHash: DATA, 32 Bytes - hash of the transaction.
		TransactionHash RawData32Response `json:"transactionHash"`

		// transactionIndex: QUANTITY - integer of the transactions index position in the block.
		TransactionIndex RawQuantityResponse `json:"transactionIndex"`

		// blockHash: DATA, 32 Bytes - hash of the block where this transaction was in.
		BlockHash RawData32Response `json:"blockHash"`

		// blockNumber: QUANTITY - block number where this transaction was in.
		BlockNumber RawQuantityResponse `json:"blockNumber"`

		// from: DATA, 20 Bytes - address of the sender.
		From RawData20Response `json:"from"`

		// to: DATA, 20 Bytes - address of the receiver. null when its a contract creation transaction.
		To RawDataResponse `json:"to"`

		// cumulativeGasUsed : QUANTITY - The total amount of gas used when this transaction was executed in the block.
		CumulativeGasUsed RawQuantityResponse `json:"cumulativeGasUsed"`

		// gasUsed : QUANTITY - The amount of gas used by this specific transaction alone.
		GasUsed RawQuantityResponse `json:"gasUsed"`

		// contractAddress : DATA, 20 Bytes - The contract address created, if the transaction was a contract creation, otherwise null.
		ContractAddress RawData20Response `json:"contractAddress"`

		// logs: Array - Array of log objects, which this transaction generated.
		Logs []RawDataResponse `json:"logs"`

		// logsBloom: DATA, 256 Bytes - Bloom filter for light clients to quickly retrieve related logs. It also returns either :
		LogsBloom RawData256Response `json:"logsBloom"`

		// root : DATA 32 bytes of post-transaction stateroot (pre Byzantium)
		Root RawData32Response `json:"root"`

		// status: QUANTITY either 1 (success) or 0 (failure)
		Status RawQuantityResponse `json:"status"`
	}
)

var ()

func NewClient() *Client {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	httpClient := &http.Client{Transport: tr}

	c := new(Client)
	c.hasAuth = false
	c.httpClient = httpClient
	c.counter = 0
	return c
}

func NewChainClient(c *Client, url string, privateKey *ecdsa.PrivateKey, chainID *big.Int) *ChainClient {
	cc := new(ChainClient)
	cc.RPCClient = c
	cc.URL = url
	cc.PrivateKey = privateKey
	cc.ChainID = chainID

	return cc
}
func (cc *ChainClient) SendTx(txdata ethtypes.TxData) (*RPCResp, error) {
	return cc.RPCClient.SendTx(cc.URL, txdata, cc.PrivateKey, cc.ChainID)
}
func (cc *ChainClient) GetTxReceipt(hash string) (*RawTxReceipt, error) {
	body := cc.RPCClient.BuildRequest("eth_getTransactionReceipt", []any{hash})
	respBody, err := cc.RPCClient.doRequest(cc.URL, body)
	if err != nil {
		return nil, err
	}

	var r *RPCReceiptResp = new(RPCReceiptResp)
	err = json.Unmarshal(respBody, r)
	if err != nil {
		return nil, err
	}

	// we might get a null response if the transaction has been mined yet
	if string(r.Result.TransactionHash) == "" {
		return nil, nil
	}

	return &r.Result, nil
}

func (c *Client) SetTimeout(duration time.Duration) {
	c.httpClient.Timeout = duration
}

func (c *Client) SetAuth(auth string) {
	pieces := strings.SplitN(auth, ":", 2)
	c.username = pieces[0]
	c.password = pieces[1]
	c.hasAuth = true
}
func (c *Client) SetProxy(proxy, proxyAuth string) {
	t := c.httpClient.Transport.(*http.Transport)
	if proxyAuth != "" {
		pieces := strings.SplitN(proxyAuth, ":", 2)
		t.Proxy = http.ProxyURL(&url.URL{
			Scheme: "http",
			User:   url.UserPassword(pieces[0], pieces[1]),
			Host:   proxy,
		})
	} else {
		t.Proxy = http.ProxyURL(&url.URL{
			Scheme: "http",
			Host:   proxy,
		})
	}
}

func (c *Client) SetKeepAlive(shouldKeepAlive bool) {
	t := c.httpClient.Transport.(*http.Transport)
	if shouldKeepAlive {
		t.DisableKeepAlives = false
	} else {
		t.DisableKeepAlives = true
	}
}

func (c *Client) MakeRequest(url, method string, params []any) (*RPCResp, error) {
	body := c.BuildRequest(method, params)
	respBody, err := c.doRequest(url, body)
	if err != nil {
		return nil, err
	}
	var r *RPCResp
	err = json.Unmarshal(respBody, &r)
	if err != nil {
		log.Trace().Interface("resp bod", respBody).Msg("There was an error unmarshalling the json response")
		return nil, err
	}
	if r.Error.Code != 0 {
		return r, fmt.Errorf("RPC Error: %s", r.Error.Message)
	}
	return r, nil
}

func (c *Client) doRequest(url string, body any) ([]byte, error) {
	s, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(s)

	resp, err := c.httpClient.Post(url, "application/json", buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

func (c *Client) BuildRequest(method string, params []any) *RPCReq {
	body := new(RPCReq)

	body.Method = method
	body.Params = params
	body.JSONRPC = "2.0"
	body.ID = c.Inc()

	return body
}

func (c *Client) MakeRequestBatchGenric(url string, methods []string, params [][]any, r any) error {
	if len(methods) != len(params) {
		return fmt.Errorf("The number of methods doesn't match the number of parameter sets being passed")
	}

	if len(methods) < 1 {
		return fmt.Errorf("Need at least one method call in a batch")
	}
	requestBatch := make([]*RPCReq, len(methods))
	for i := range methods {
		requestBatch[i] = c.BuildRequest(methods[i], params[i])
	}

	respBody, err := c.doRequest(url, requestBatch)
	if err != nil {
		return err
	}
	err = json.Unmarshal(respBody, r)
	if err != nil {
		log.Trace().Interface("respbod", respBody).Msg("There was an error unmarshalling the json response")
		return err
	}
	return nil
}

// MakeRequestBatch will perform an RPC patch call. The response array
// will match the order of the input methods
func (c *Client) MakeRequestBatch(url string, methods []string, params [][]any) ([]RPCResp, error) {
	if len(methods) != len(params) {
		return nil, fmt.Errorf("The number of methods doesn't match the number of parameter sets being passed")
	}

	if len(methods) < 1 {
		return nil, fmt.Errorf("Need at least one method call in a batch")
	}
	requestBatch := make([]*RPCReq, len(methods))
	for i := range methods {
		requestBatch[i] = c.BuildRequest(methods[i], params[i])
	}

	respBody, err := c.doRequest(url, requestBatch)
	if err != nil {
		return nil, err
	}
	var r []RPCResp
	err = json.Unmarshal(respBody, &r)
	if err != nil {
		log.Trace().Interface("respbod", respBody).Msg("There was an error unmarshalling the json response")
		return nil, err
	}

	if len(r) != len(methods) {
		return nil, fmt.Errorf("Mismatch between response aray and methods")
	}
	mappedResponses := make(map[uint64]RPCResp)
	for _, v := range r {
		mappedResponses[v.ID] = v
	}

	matchedResponses := make([]RPCResp, len(r))
	for i, v := range requestBatch {
		matchedResponses[i] = mappedResponses[v.ID]
	}

	return matchedResponses, nil
}

func (c *Client) Inc() uint64 {
	c.counter = c.counter + 1
	return c.counter
}

func (c *Client) SendTxSimple(pk string, value *big.Int, gasPrice *big.Int, toAddress string, nonce uint64, chainID *big.Int, url string) error {
	privateKey, err := ethcrypto.HexToECDSA(pk)
	if err != nil {
		return err
	}

	gasLimit := uint64(21000) // in units

	sendTo := ethcommon.HexToAddress(toAddress)
	var data []byte
	// tx := ethtypes.NewTransaction(nonce, sendTo, value, gasLimit, gasPrice, data)
	lt := ethtypes.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &sendTo,
		Value:    value,
		Data:     data,
	}
	log.Trace().Interface("tx", lt).Msg("Generated transaction")

	tx := ethtypes.NewTx(&lt)

	signedTx, err := ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return err
	}

	// ts := ethtypes.Transactions{signedTx}
	// rawTx := hex.EncodeToString(ts.GetRlp(0))

	// fmt.Printf(rawTx) // f86...772
	var buf bytes.Buffer
	signedTx.EncodeRLP(&buf)

	resp, err := c.MakeRequest(url, "eth_sendRawTransaction", []any{"0x" + hex.EncodeToString(buf.Bytes())})
	if err != nil {
		log.Error().Err(err).Msg("There was an issue sending the raw transaction")
		return err
	}
	log.Trace().Interface("txresp", resp.Result).Msg("Sent transaction")

	resp, err = c.MakeRequest(url, "eth_getTransactionReceipt", []any{resp.Result})
	if err != nil {
		return err
	}
	log.Trace().Interface("txreceipt", resp.Result).Msg("Got transaction receipt")

	return nil
}

func (c *Client) SendTx(url string, txdata ethtypes.TxData, privateKey *ecdsa.PrivateKey, chainID *big.Int) (*RPCResp, error) {
	tx := ethtypes.NewTx(txdata)

	signedTx, err := ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Error().Err(err).Msg("could not sign transaction")
		return nil, err
	}

	var buf bytes.Buffer
	signedTx.EncodeRLP(&buf)

	resp, err := c.MakeRequest(url, "eth_sendRawTransaction", []any{"0x" + hex.EncodeToString(buf.Bytes())})
	if err != nil {
		log.Error().Err(err).Msg("There was an issue sending the raw transaction")
		return nil, err
	}
	return resp, nil
}

// HexToBigInt assumes that it's input is a hex encoded string and
// will try to convert it to a big int
func ConvHexToBigInt(raw any) (bi *big.Int, err error) {
	bi = big.NewInt(0)
	hexString, err := rawRespToString(raw)
	if err != nil {
		return nil, err
	}
	hexString = strings.Replace(hexString, "0x", "", -1)
	if len(hexString)%2 != 0 {
		hexString = "0" + hexString
	}

	rawGas, err := hex.DecodeString(hexString)
	if err != nil {
		log.Error().Err(err).Str("hex", hexString).Msg("Unable to decode hex string")
		return
	}
	bi.SetBytes(rawGas)
	return
}
func MustConvHexToBigInt(raw any) *big.Int {
	bi, err := ConvHexToBigInt(raw)
	if err != nil {
		panic(fmt.Sprintf("Failed to convert Hex to Big int: %v", err))
	}
	return bi
}

func rawRespToString(raw any) (string, error) {
	var hexString string
	switch v := raw.(type) {
	case RawQuantityResponse:
		hexString = string(v)
	case RawDataResponse:
		hexString = string(v)
	case RawData8Response:
		hexString = string(v)
	case RawData20Response:
		hexString = string(v)
	case RawData32Response:
		hexString = string(v)
	case RawData256Response:
		hexString = string(v)
	case string:
		hexString = v
	default:
		return "", fmt.Errorf("Could not assert %v as a string", raw)
	}
	return hexString, nil
}

// HexToUint64 assumes that its input is a hex encoded string and it
// will attempt to convert this into a uint64
func ConvHexToUint64(raw any) (uint64, error) {
	hexString, err := rawRespToString(raw)
	if err != nil {
		return 0, err
	}

	hexString = strings.Replace(hexString, "0x", "", -1)
	if len(hexString)%2 != 0 {
		hexString = "0" + hexString
	}

	result, err := strconv.ParseUint(hexString, 16, 64)
	if err != nil {
		return 0, err
	}
	return uint64(result), nil
}

func MustConvHexToUint64(raw any) uint64 {
	num, err := ConvHexToUint64(raw)
	if err != nil {
		panic(fmt.Sprintf("Failed to covert Hex to uint64: %v", err))
	}
	return num
}

func NewRawBlockResponseFromAny(raw any) (*RawBlockResponse, error) {
	topMap, ok := raw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("Unable to map raw response")
	}
	_ = topMap
	return nil, nil

}

func (r RawQuantityResponse) ToUint64() uint64 {
	hexString := strings.Replace(string(r), "0x", "", -1)
	if len(hexString)%2 != 0 {
		hexString = "0" + hexString
	}

	result, err := strconv.ParseUint(hexString, 16, 64)
	if err != nil {
		return 0
	}
	return uint64(result)

}
func (r RawQuantityResponse) ToInt64() int64 {
	hexString := strings.Replace(string(r), "0x", "", -1)
	if len(hexString)%2 != 0 {
		hexString = "0" + hexString
	}

	result, err := strconv.ParseUint(hexString, 16, 64)
	if err != nil {
		return 0
	}
	return int64(result)

}
