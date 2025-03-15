package cdk

import (
	"reflect"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
)

// customFilter is an auxiliary structs to allow filtering logs
// when querying the blockchain, but also filtering them after
// with a custom logic that is not allowed by the regular filter
// like filtering by the data.
//
// example:
//
//	rollupManagerFilter := customFilter{
//	    contractInstance: rollupManager.instance,
//	    contractABI:      rollupManagerABI,
//	    blockchainFilter: ethereum.FilterQuery{
//	        Addresses: []common.Address{
//	            rollupManagerArgs.rollupManagerAddress,
//	        },
//	        Topics: [][]common.Hash{
//	            nil, // no filter to topic 0,
//	            {common.BigToHash(big.NewInt(0).SetUint64(uint64(data.RollupID)))}, // filter topic 1 by RollupID
//	        },
//	    },
//	    postFilterFunc: func(logs []types.Log, contractInstance reflect.Value, contractABI *abi.ABI) []types.Log {
//	        filteredLogs := make([]types.Log, 0, len(logs))
//	        // filter logs from rollup-manager by RollupID
//	        for _, l := range logs {
//	            e, _ := contractABI.EventByID(l.Topics[0])
//	            // if the event is not found, there is no way to check if the log is related to a
//	            // rollup, this can happen in the case the contract was updated, but the ABI was not
//	            //
//	            // in order to avoid missing logs, we will not filter this log out
//	            if e == nil {
//	                filteredLogs = append(filteredLogs, l)
//	                continue
//	            }
//
//	            parseLogMethod, methodFound := contractInstance.Type().MethodByName(fmt.Sprintf("Parse%s", e.Name))
//	            if !methodFound {
//	                // if the method to parse the log is not found, there is no way to check if the
//	                // log is related to a rollup, this can happen in the case the contract was
//	                // updated, but the ABI was not
//	                //
//	                // in order to avoid missing logs, we will not filter this log out
//	                filteredLogs = append(filteredLogs, l)
//	            } else {
//	                parsedLogValues := parseLogMethod.Func.Call([]reflect.Value{contractInstance, reflect.ValueOf(l)})
//	                errValue := parsedLogValues[1].Interface()
//	                if errValue != nil {
//	                    log.Warn().Err(errValue.(error)).Msgf("Error parsing log %v", l)
//	                } else {
//	                    rollupIDField := parsedLogValues[0].FieldByName("RollupID")
//	                    // in case the RollupID field is not found, we filter this log out
//	                    if !rollupIDField.IsValid() {
//	                        continue
//	                    }
//
//	                    // if the event is parsable, has the RollupID field and the RollupID matches
//	                    // we keep the log
//	                    rollupID := uint32(rollupIDField.Uint())
//	                    if rollupID == data.RollupID {
//	                        filteredLogs = append(filteredLogs, l)
//	                    }
//	                }
//	            }
//	        }
//	        return filteredLogs
//	    },
//	}
type customFilter struct {
	// contract instance used to parse logs returned by the filter
	contractInstance reflect.Value
	// contract ABI used to parse logs returned by the filter
	contractABI *abi.ABI
	// filter used to query the blockchain
	blockchainFilter ethereum.FilterQuery
	// function used to filter logs after querying the blockchain
	postFilterFunc func(logs []types.Log, contractInstance reflect.Value, contractABI *abi.ABI) []types.Log
}
