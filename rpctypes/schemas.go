package rpctypes

import (
	_ "embed"
)

//go:embed jsonschemas/eth_syncing.json
var RPCSchemaEthSyncing string

//go:embed jsonschemas/eth_getBlockByNumber.json
var RPCSchemaEthBlock string

//go:embed jsonschemas/eth_accounts.json
var RPCSchemaAccountList string

//go:embed jsonschemas/rpcschemasigntxresponse.json
var RPCSchemaSignTxResponse string

//go:embed jsonschemas/eth_getTransactionByHash.json
var RPCSchemaEthTransaction string

//go:embed jsonschemas/eth_getTransactionReceipt.json
var RPCSchemaEthReceipt string

//go:embed jsonschemas/eth_getFilterChanges.json
var RPCSchemaEthFilter string

//go:embed jsonschemas/eth_feeHistory.json
var RPCSchemaEthFeeHistory string

//go:embed jsonschemas/eth_createAccessList.json
var RPCSchemaEthAccessList string

//go:embed jsonschemas/eth_getProof.json
var RPCSchemaEthProof string

//go:embed jsonschemas/rpcschemadebugtrace.json
var RPCSchemaDebugTrace string

//go:embed jsonschemas/rpcschemahexarray.json
var RPCSchemaHexArray string

//go:embed jsonschemas/debug_getBadBlocks.json
var RPCSchemaBadBlocks string

//go:embed jsonschemas/rpcschemadebugblock.json
var RPCSchemaDebugTraceBlock string
