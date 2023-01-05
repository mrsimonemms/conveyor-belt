CALLBACK_URL ?= https://eosv8e8x84ccn8d.m.pipedream.net?stage=completed
CONFIG_FILE ?= ./examples/basic/config.yaml
LOG_LEVEL ?= debug
SERVER_ADDRESS ?= http://localhost:3000
WEBHOOK ?= /webhook/basic

run:
	@go run . run \
		--config ${CONFIG_FILE} \
		--log-level ${LOG_LEVEL}
.PHONY: run

trigger:
	curl -iL -X POST -H "X-Callback-URL: ${CALLBACK_URL}" ${SERVER_ADDRESS}${WEBHOOK}
.PHONY: trigger
