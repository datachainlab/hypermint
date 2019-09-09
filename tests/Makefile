GO_BIN?=go
GO_TEST_FLAGS?=-v -count=1
GO_TEST_CMD=$(GO_BIN) test $(GO_TEST_FLAGS)
VERSION ?= v0.2.0

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
