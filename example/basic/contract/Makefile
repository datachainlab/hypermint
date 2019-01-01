BUILD_DIR?=../../../build
HMD?=$(BUILD_DIR)/hmd
HMCLI?=$(BUILD_DIR)/hmcli

build:
	cargo build --target=wasm32-unknown-unknown
	wasm-gc ./target/wasm32-unknown-unknown/debug/contract.wasm -o ./contract.min.wasm

run: build
	$(HMD) vm --path=./contract.min.wasm
