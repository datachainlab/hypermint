GO_BIN?=go
GO_TEST_FLAGS?=-v -count=1
GO_TEST_CMD=$(GO_BIN) test $(GO_TEST_FLAGS)
VERSION ?= 0.1.0

build-image:
	docker build . -t bluele/go-rust:$(VERSION)

upload-image:
	docker push bluele/go-rust:$(VERSION)

build/contract_test.wasm:
	cd ./contract-test && cargo build --target=wasm32-unknown-unknown
	mkdir -p ./build/
	mv ./contract-test/target/wasm32-unknown-unknown/debug/contract_test.wasm ./build/contract_test.wasm

integration-test: build/contract_test.wasm
	@$(GO_TEST_CMD) ./integration/...

e2e-test: build/contract_test.wasm
	@$(GO_TEST_CMD) ./e2e/...
