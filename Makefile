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
	curl -iL -X POST ${SERVER_ADDRESS}${WEBHOOK}
.PHONY: trigger
