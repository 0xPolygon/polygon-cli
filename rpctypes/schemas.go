package rpctypes

import (
	_ "embed"
)

//go:embed schemas/rpcschemaethsyncing.json
var RPCSchemaEthSyncing string

//go:embed schemas/rpcschemaethblock.json
var RPCSchemaEthBlock string
