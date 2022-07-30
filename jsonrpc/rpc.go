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
)

var (
	UnitWei       = big.NewInt(1)                                    // | 1 | 1 | wei | Wei
	UnitBabbage   = new(big.Int).Mul(UnitWei, big.NewInt(1000))      // | 1,000 | 10^3^ | Babbage | Kilowei or femtoether
	UnitLovelace  = new(big.Int).Mul(UnitBabbage, big.NewInt(1000))  // | 1,000,000 | 10^6^ | Lovelace | Megawei or picoether
	UnitShannon   = new(big.Int).Mul(UnitLovelace, big.NewInt(1000)) // | 1,000,000,000 | 10^9^ | Shannon | Gigawei or nanoether
	UnitSzabo     = new(big.Int).Mul(UnitShannon, big.NewInt(1000))  // | 1,000,000,000,000 | 10^12^ | Szabo | Microether or micro
	UnitFinney    = new(big.Int).Mul(UnitSzabo, big.NewInt(1000))    // | 1,000,000,000,000,000 | 10^15^ | Finney | Milliether or milli
	UnitEther     = new(big.Int).Mul(UnitFinney, big.NewInt(1000))   // | 1,000,000,000,000,000,000 | 10^18^ | Ether | Ether
	UnitGrand     = new(big.Int).Mul(UnitEther, big.NewInt(1000))    // | 1,000,000,000,000,000,000,000 | 10^21^ | Grand | Kiloether
	UnitMegaether = new(big.Int).Mul(UnitGrand, big.NewInt(1000))    // | 1,000,000,000,000,000,000,000,000 | 10^24^ | | Megaether

)

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
	body := RPCReq{}

	body.Method = method
	body.Params = params
	body.JSONRPC = "2.0"
	body.ID = c.Inc()

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
func HexToBigInt(raw any) (bi *big.Int, err error) {
	bi = big.NewInt(0)
	hexString, ok := raw.(string)
	if !ok {
		err = fmt.Errorf("Could not assert value %v as a string", raw)
		return
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

// HexToUint64 assumes that its input is a hex encoded string and it
// will attempt to convert this into a uint64
func HexToUint64(raw any) (uint64, error) {
	hexString, ok := raw.(string)
	if !ok {
		return 0, fmt.Errorf("Could not assert %v as a string", hexString)
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
