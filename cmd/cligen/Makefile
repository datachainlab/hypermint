
.PHONY: test
test: gen-example
	go build
	./cligen --outdir example --name token --package main --abi example/token.json
	(cd example && go build)

.PHONY: gen-example
gen-example:
	mkdir -p example
	(cd ../../hmcdk/genabi && cargo run ../../example/token/src/lib.rs) > example/token.json
