
build:
	cargo build --target=wasm32-unknown-unknown
	wasm-gc ./target/wasm32-unknown-unknown/debug/contract.wasm -o ./contract.min.wasm
