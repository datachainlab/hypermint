
.PHONY: test
test:
	go build
	./abigen --outdir example --name token --package main --abi example/token.json --mock
	(cd example && go test)

.PHONY: gen-example
gen-example:
	mkdir -p example
	(cd ../../hmcdk/genabi && cargo run ../../example/token/src/lib.rs) > example/token.json
