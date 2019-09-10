lint-tools:
	rustup component add clippy
	rustup component add rustfmt

lint:
	cargo fmt -- --check
	cargo clippy -- -D warnings

test:
	cargo test

.PHONY: lint-tools lint test
