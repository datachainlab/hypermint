GO_BIN?=go
GO_TEST_FLAGS?=-v -count=1
GO_TEST_CMD=$(GO_BIN) test $(GO_TEST_FLAGS)
VERSION ?= 0.1.0

build-image:
	docker build . -t bluele/go-rust:$(VERSION)

upload-image:
	docker push bluele/go-rust:$(VERSION)

%.wasm:
	cd ./contracts/$* && cargo build --target=wasm32-unknown-unknown
	mkdir -p ./build/
	mv ./contracts/$*/target/wasm32-unknown-unknown/debug/$@ ./build/$@

integration-test: contract_test.wasm
	@$(GO_TEST_CMD) ./integration/...

e2e-test: contract_test.wasm external_contract_test.wasm
	@$(GO_TEST_CMD) ./e2e/...
