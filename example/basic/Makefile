BUILD_DIR?=../../build
HMD?=$(BUILD_DIR)/hmd
HMCLI?=$(BUILD_DIR)/hmcli
HMD_HOME?=${HOME}/.hmd
HMCLI_HOME?=${HOME}/.hmcli

NODE_ADDR?=localhost:26657
CONTRACT_CODE?=./contract/contract.min.wasm

build:
	$(MAKE) -C ./contract build

deploy: build
	$(eval CONTRACT_ADDRESS := $(shell $(HMCLI) contract deploy --path=$(CONTRACT_CODE) --address=${ADDR1} --gas=1 --password=password --node=tcp://$(NODE_ADDR)))
	@echo "contract address is ${CONTRACT_ADDRESS}"

call:
	$(HMCLI) contract call --address=${ADDR1} --contract=$(CONTRACT_ADDRESS) --gas=1  --func="app_main" --password=password --node=tcp://$(NODE_ADDR)
