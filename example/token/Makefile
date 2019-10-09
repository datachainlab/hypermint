BUILD_DIR?=../../build
HMD?=$(BUILD_DIR)/hmd
HMCLI?=$(BUILD_DIR)/hmcli

NODE_ADDR?=localhost:26657
CONTRACT_CODE?=./token.min.wasm

deploy: build
	$(eval CONTRACT_ADDRESS := $(shell $(HMCLI) contract deploy --path=$(CONTRACT_CODE) --address=${ADDR1} --gas=1 --password=password --node=tcp://$(NODE_ADDR)))
	@echo "contract address is ${CONTRACT_ADDRESS}"

transfer:
	$(HMCLI) contract call --address=${ADDR1} --contract=$(CONTRACT_ADDRESS) --gas=1  --func="transfer" --args ${ADDR2} --args 10 --argtypes=address,int64 --password=password --node=tcp://$(NODE_ADDR)

build:
	cargo build --target=wasm32-unknown-unknown
	wasm-gc ./target/wasm32-unknown-unknown/debug/token.wasm -o $(CONTRACT_CODE)

run: build
	$(HMD) vm --path=$(CONTRACT_CODE)

clean:
	rm -rf ./target $(CONTRACT_CODE)
