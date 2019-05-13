contract:
	cd ./contract-test && cargo build --target=wasm32-unknown-unknown

test: contract
	go test -v . -count=1
