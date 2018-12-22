GO_BIN?=go

HMD?=./build/hmd
HMCLI?=./build/hmcli
HMD_HOME?=${HOME}/.hmd
HMCLI_HOME?=${HOME}/.hmcli

MNEMONIC?="hine volcano run comic tiger traffic attitude hockey depend cash clever curtain"
HDW_PATH?=m/44'/60'/0'/0

COMMIT_HASH:=$(shell git rev-parse --short HEAD)
BUILD_FLAGS?=-ldflags "-X github.com/bluele/hypermint/pkg/consts.GitCommit=${COMMIT_HASH}"

GO_BUILD_CMD=$(GO_BIN) build $(BUILD_FLAGS)
GO_TEST_FLAGS?=-v
GO_TEST_CMD=$(GO_BIN) test $(GO_TEST_FLAGS)

.PHONY: build

build: server cli

server:
	$(GO_BUILD_CMD) -o $(HMD) ./cmd/hmd

cli:
	$(GO_BUILD_CMD) -o $(HMCLI) ./cmd/hmcli

start:
	$(HMD) start --log_level="main:error" --home=$(HMD_HOME)

clean:
	@rm -rf $(HMD_HOME) $(HMCLI_HOME)

init: clean init-validator
	$(eval ADDR1 := $(shell $(HMCLI) new --password=password --silent --home=$(HMCLI_HOME) --mnemonic=$(MNEMONIC) --hdw_path="$(HDW_PATH)/1" ))
	$(eval ADDR2 := $(shell $(HMCLI) new --password=password --silent --home=$(HMCLI_HOME) --mnemonic=$(MNEMONIC) --hdw_path="$(HDW_PATH)/2" ))
	@$(HMD) init --address=$(ADDR1) --home=$(HMD_HOME)
	@echo export ADDR1='$(ADDR1)'
	@echo export ADDR2='$(ADDR2)'

init-validator:
	@$(HMD) tendermint init-validator --mnemonic=$(MNEMONIC) --hdw_path="$(HDW_PATH)/0"

test:
	$(GO_TEST_CMD) ./pkg/...
