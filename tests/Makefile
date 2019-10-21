GO_BIN?=go
GO_TEST_FLAGS?=-v -count=1
GO_TEST_CMD=$(GO_BIN) test $(GO_TEST_FLAGS)
VERSION ?= v0.2.0
TOP_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))/..

build-image:
	docker build . -t hypermint/go-rust:$(VERSION)

upload-image:
	docker push hypermint/go-rust:$(VERSION)

%.wasm:
	cd ./contracts/$* && cargo build --target=wasm32-unknown-unknown
	mkdir -p ./build/
	mv ./contracts/$*/target/wasm32-unknown-unknown/debug/$@ ./build/$@

integration-test: contract_test.wasm
	@$(GO_TEST_CMD) ./integration/...

e2e-test: contract_test.wasm external_contract_test.wasm
	@$(GO_TEST_CMD) ./e2e/...

cligen-test: contract_test.wasm
	make -C .. abigen cligen
	(cd $(TOP_DIR)/hmcdk/genabi && cargo run $(TOP_DIR)/tests/contracts/contract_test/src/lib.rs) > cligen/contract_test.json
	$(TOP_DIR)/build/cligen --outdir cligen --name contract_test --package cligen --abi cligen/contract_test.json
	@$(GO_TEST_CMD) ./cligen/...
