##@ Clients
HOST?=127.0.0.1
PORT?=8545
CHAIN_ID?=1337
LOADTEST_ACCOUNT?=0x85da99c8a7c2c95964c8efd687e95e632fc533d6
LOADTEST_FUNDING_AMOUNT_ETH?=1010000000

.PHONY: geth
geth: ## Start a local geth node.
	geth \
		--dev \
		--http \
		--http.addr ${HOST} \
		--http.port $(PORT) \
		--http.api admin,debug,web3,eth,txpool,personal,miner,net \
		--verbosity 5 \
		--rpc.gascap 50000000 \
		--rpc.txfeecap 0 \
		--miner.gaslimit 10 \
		--miner.gasprice 1 \
		--gpo.blocks 1 \
		--gpo.percentile 1 \
		--gpo.maxprice 10 \
		--gpo.ignoreprice 2 \
		--dev.gaslimit 100000000000

.PHONY: anvil
anvil: ## Start a local anvil node.
	anvil \
		--host ${HOST} \
		--port $(PORT) \
		--chain-id ${CHAIN_ID} \
		--balance 999999999999999

.PHONY: fund
fund: ## Fund the loadtest account with 100k ETH.
	eth_coinbase=$$(curl -s -H 'Content-Type: application/json' -d '{"jsonrpc": "2.0", "id": 2, "method": "eth_accounts", "params": []}' http://${HOST}:${PORT} | jq -r ".result[0]"); \
	hex_funding_amount=$$(echo "obase=16; ${LOADTEST_FUNDING_AMOUNT_ETH}*10^18" | bc); \
	echo $$eth_coinbase $$hex_funding_amount; \
	curl \
		-H "Content-Type: application/json" \
		-d '{"jsonrpc":"2.0", "method":"eth_sendTransaction", "params":[{"from": "'$$eth_coinbase'","to": "${LOADTEST_ACCOUNT}","value": "0x'$$hex_funding_amount'"}], "id":1}' \
		-s \
		http://${HOST}:${PORT} | jq

.PHONY: loadtest
loadtest: fund ## Run random loadtest against a local RPC.
	sleep 2
	go run -race main.go loadtest \
		--verbosity 600 \
		--rpc-url http://${HOST}:$(PORT) \
		--chain-id ${CHAIN_ID} \
		--mode random \
		--concurrency 1 \
		--requests 200 \
		--rate-limit 100
