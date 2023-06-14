package rpctypes

import (
	_ "embed"
)

//go:embed schemas/rpcschemaethsyncing.json
var RPCSchemaEthSyncing string

//go:embed schemas/rpcschemaethblock.json
var RPCSchemaEthBlock string

//go:embed schemas/rpcschemaaccountlist.json
var RPCSchemaAccountList string

//go:embed schemas/rpcschemasigntxresponse.json
var RPCSchemaSignTxResponse string

//go:embed schemas/rpcschemaethtransaction.json
var RPCSchemaEthTransaction string

//go:embed schemas/rpcschemaethreceipt.json
var RPCSchemaEthReceipt string

//go:embed schemas/rpcschemafilterchanges.json
var RPCSchemaEthFilter string

//go:embed schemas/rpcschemaethfeehistory.json
var RPCSchemaEthFeeHistory string

//go:embed schemas/rpcschemaethaccesslist.json
var RPCSchemaEthAccessList string

//go:embed schemas/rpcschemaethproof.json
var RPCSchemaEthProof string
