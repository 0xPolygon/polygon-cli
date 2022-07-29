package jsonrpc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
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

func (c *Client) MakeRawRequest() {

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
	return r, nil
}

func (c *Client) Inc() uint64 {
	c.counter = c.counter + 1
	return c.counter
}

// func init() {

// 	privateKey, err := ethcrypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
// 	if err != nil {
// 	}

// 	var nonce uint64 = 1

// 	value := big.NewInt(1000000000000000000) // in wei (1 eth)
// 	gasLimit := uint64(21000)                // in units
// 	gasPrice := big.NewInt(1)                // , err := client.SuggestGasPrice(context.Background())

// 	toAddress := ethcommon.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
// 	var data []byte
// 	tx := ethtypes.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

// 	chainID := big.NewInt(123) // client.NetworkID(context.Background())

// 	signedTx, err := ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(chainID), privateKey)
// 	if err != nil {
// 		// log.Fatal(err)
// 	}

// 	// ts := ethtypes.Transactions{signedTx}
// 	// rawTx := hex.EncodeToString(ts.GetRlp(0))

// 	// fmt.Printf(rawTx) // f86...772
// 	var buf bytes.Buffer
// 	signedTx.EncodeRLP(&buf)
// 	fmt.Println(hex.EncodeToString(buf.Bytes()))

// }

func (c *Client) SendTx(pk string, value *big.Int, gasPrice *big.Int, toAddress string, nonce uint64, chainID *big.Int, url string) error {
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
	fmt.Println(hex.EncodeToString(buf.Bytes()))

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
