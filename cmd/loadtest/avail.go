package loadtest

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	_ "embed"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	gssignature "github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	gstypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

func availLoop(ctx context.Context, c *gsrpc.SubstrateAPI) error {
	var err error

	ltp := inputLoadTestParams
	log.Trace().Interface("Input Params", ltp).Msg("Params")

	routines := *ltp.Concurrency
	requests := *ltp.Requests
	currentNonce := uint64(0) // *ltp.CurrentNonce
	chainID := new(big.Int).SetUint64(*ltp.ChainID)
	privateKey := ltp.ECDSAPrivateKey
	mode := *ltp.Mode

	_ = chainID
	_ = privateKey

	meta, err := c.RPC.State.GetMetadataLatest()
	if err != nil {
		return err
	}

	genesisHash, err := c.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return err
	}

	key, err := gstypes.CreateStorageKey(meta, "System", "Account", ltp.FromAvailAddress.PublicKey, nil)
	if err != nil {
		log.Error().Err(err).Msg("Could not create storage key")
		return err
	}

	var accountInfo gstypes.AccountInfo
	ok, err := c.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		log.Error().Err(err).Msg("Could not load storage")
		return err
	}
	if !ok {
		err = fmt.Errorf("loaded storage is not okay")
		log.Error().Err(err).Msg("Loaded storage is not okay")
		return err
	}

	currentNonce = uint64(accountInfo.Nonce)

	rl := rate.NewLimiter(rate.Limit(*ltp.RateLimit), 1)
	if *ltp.RateLimit <= 0.0 {
		rl = nil
	}

	var currentNonceMutex sync.Mutex

	var i int64

	var wg sync.WaitGroup
	for i = 0; i < routines; i = i + 1 {
		log.Trace().Int64("routine", i).Msg("Starting Thread")
		wg.Add(1)
		go func(i int64) {
			var j int64
			var startReq time.Time
			var endReq time.Time

			for j = 0; j < requests; j = j + 1 {

				if rl != nil {
					err = rl.Wait(ctx)
					if err != nil {
						log.Error().Err(err).Msg("Encountered a rate limiting error")
					}
				}
				currentNonceMutex.Lock()
				myNonceValue := currentNonce
				currentNonce = currentNonce + 1
				currentNonceMutex.Unlock()

				localMode := mode
				// if there are multiple modes, iterate through them, 'r' mode is supported here
				if len(mode) > 1 {
					localMode = string(mode[int(i+j)%(len(mode))])
				}
				// if we're doing random, we'll just pick one based on the current index
				if localMode == loadTestModeRandom {
					localMode = validLoadTestModes[int(i+j)%(len(validLoadTestModes)-1)]
				}
				// this function should probably be abstracted
				switch localMode {
				case loadTestModeTransaction:
					startReq, endReq, err = loadtestAvailTransfer(ctx, c, myNonceValue, meta, genesisHash)
				case loadTestModeDeploy:
					startReq, endReq, err = loadtestNotImplemented(ctx, c, myNonceValue)
				case loadTestModeCall:
					startReq, endReq, err = loadtestNotImplemented(ctx, c, myNonceValue)
				case loadTestModeFunction:
					startReq, endReq, err = loadtestNotImplemented(ctx, c, myNonceValue)
				case loadTestModeInc:
					startReq, endReq, err = loadtestNotImplemented(ctx, c, myNonceValue)
				case loadTestModeStore:
					startReq, endReq, err = loadtestAvailStore(ctx, c, myNonceValue, meta, genesisHash)
				default:
					log.Error().Str("mode", mode).Msg("We've arrived at a load test mode that we don't recognize")
				}
				recordSample(i, j, err, startReq, endReq, myNonceValue)
				if err != nil {
					log.Trace().Err(err).Msg("Recorded an error while sending transactions")
				}

				log.Trace().Int64("routine", i).Str("mode", localMode).Int64("request", j).Msg("Request")
			}
			wg.Done()
		}(i)

	}
	log.Trace().Msg("Finished starting go routines. Waiting..")
	wg.Wait()
	return nil

}

func initAvailTestParams(ctx context.Context, c *gsrpc.SubstrateAPI) error {
	toAddr, err := gstypes.NewMultiAddressFromHexAccountID(*inputLoadTestParams.ToAddress)
	if err != nil {
		log.Error().Err(err).Msg("Unable to create new multi address")
		return err
	}

	if *inputLoadTestParams.PrivateKey == codeQualityPrivateKey {
		// Avail keys can use the same seed but the way the key is derived is different
		*inputLoadTestParams.PrivateKey = codeQualitySeed
	}

	kp, err := gssignature.KeyringPairFromSecret(*inputLoadTestParams.PrivateKey, uint8(*inputLoadTestParams.ChainID))
	if err != nil {
		log.Error().Err(err).Msg("Could not create key pair")
		return err
	}

	amt, err := hexToBigInt(*inputLoadTestParams.HexSendAmount)
	if err != nil {
		log.Error().Err(err).Msg("Couldn't parse send amount")
		return err
	}

	rv, err := c.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		log.Error().Err(err).Msg("Couldn't get runtime version")
		return err
	}

	inputLoadTestParams.AvailRuntime = rv
	inputLoadTestParams.SendAmount = amt
	inputLoadTestParams.FromAvailAddress = &kp
	inputLoadTestParams.ToAvailAddress = &toAddr
	return nil
}

func loadtestAvailTransfer(ctx context.Context, c *gsrpc.SubstrateAPI, nonce uint64, meta *gstypes.Metadata, genesisHash gstypes.Hash) (t1 time.Time, t2 time.Time, err error) {
	ltp := inputLoadTestParams

	toAddr := *ltp.ToAvailAddress
	if *ltp.ToRandom {
		pk := make([]byte, 32)
		_, err = randSrc.Read(pk)
		if err != nil {
			// For some reason weren't able to read the random data
			log.Error().Msg("Sending to random is not implemented for substrate yet")
		} else {
			toAddr = gstypes.NewMultiAddressFromAccountID(pk)
		}

	}

	gsCall, err := gstypes.NewCall(meta, "Balances.transfer", toAddr, gstypes.NewUCompact(ltp.SendAmount))
	if err != nil {
		return
	}

	ext := gstypes.NewExtrinsic(gsCall)
	rv := ltp.AvailRuntime
	kp := *inputLoadTestParams.FromAvailAddress

	o := gstypes.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                gstypes.ExtrinsicEra{IsMortalEra: false, IsImmortalEra: true},
		GenesisHash:        genesisHash,
		Nonce:              gstypes.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                gstypes.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	err = ext.Sign(kp, o)
	if err != nil {
		return
	}

	t1 = time.Now()
	defer func() { t2 = time.Now() }()
	_, err = c.RPC.Author.SubmitExtrinsic(ext)
	return
}

func loadtestAvailStore(ctx context.Context, c *gsrpc.SubstrateAPI, nonce uint64, meta *gstypes.Metadata, genesisHash gstypes.Hash) (t1 time.Time, t2 time.Time, err error) {
	ltp := inputLoadTestParams

	inputData := make([]byte, *ltp.ByteCount)
	_, _ = hexwordRead(inputData)

	gsCall, err := gstypes.NewCall(meta, "DataAvailability.submit_data", gstypes.NewBytes([]byte(inputData)))
	if err != nil {
		return
	}

	// Create the extrinsic
	ext := gstypes.NewExtrinsic(gsCall)

	rv := ltp.AvailRuntime

	kp := *inputLoadTestParams.FromAvailAddress

	o := gstypes.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                gstypes.ExtrinsicEra{IsMortalEra: false, IsImmortalEra: true},
		GenesisHash:        genesisHash,
		Nonce:              gstypes.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                gstypes.NewUCompactFromUInt(100),
		TransactionVersion: rv.TransactionVersion,
	}
	// Sign the transaction using Alice's default account
	err = ext.Sign(kp, o)
	if err != nil {
		return
	}

	// Send the extrinsic
	t1 = time.Now()
	defer func() { t2 = time.Now() }()
	_, err = c.RPC.Author.SubmitExtrinsic(ext)
	return
}
