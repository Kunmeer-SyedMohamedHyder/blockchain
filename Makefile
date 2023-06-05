SHELL := /bin/bash

# Wallets
# Kennedy: 0xF01813E4B85e178A83e29B8E7bF26BD830a25f32
# Pavel: 0xdd6B972ffcc631a62CAE1BB9d80b7ff429c8ebA4
# Ceasar: 0xbEE6ACE826eC3DE1B6349888B9151B92522F7F76
# Baba: 0x6Fe6CF3c8fF57c58d24BfC869668F48BCbDb3BD9
# Ed: 0xa988b1866EaBF72B4c53b592c97aAD8e4b9bDCC0
# Miner1: 0xFef311483Cc040e1A89fb9bb469eeB8A70935EF8
# Miner2: 0xb8Ee4c7ac4ca3269fEc242780D7D960bd6272a61
#
# Run two miners
# make up
# make up2
#
# Bookeeping transactions
# curl -il -X GET http://localhost:8080/v1/genesis/list
# curl -il -X GET http://localhost:9080/v1/node/status
# curl -il -X GET http://localhost:8080/v1/accounts/list
# curl -il -X GET http://localhost:8080/v1/tx/uncommitted/list
# curl -il -X GET http://localhost:8080/v1/blocks/list
# curl -il -X GET http://localhost:9080/v1/node/block/list/1/latest
#
# Wallet Stuff
# go run app/wallet/cli/main.go generate
# go run app/wallet/cli/main.go account -a kennedy
# go run app/wallet/cli/main.go balance -a kennedy

# ==============================================================================
# Local support

scratch:
	go run app/scratch/main.go

up:
	go run app/services/node/main.go -race | go run app/tooling/logfmt/main.go

tidy:
	go mod tidy
	go mod vendor
