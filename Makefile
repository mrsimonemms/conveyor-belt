CONFIG_FILE ?= ./examples/basic/config.yaml
LOG_LEVEL ?= debug

run:
	@go run . run \
		--config ${CONFIG_FILE} \
		--log-level ${LOG_LEVEL}
.PHONY: run
