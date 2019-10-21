
.PHONY: test
test:
	(cd ../../hmcdk/genabi && cargo run ../../example/token/src/lib.rs) > example/token.json
	go build
	./cligen --outdir example --name token --package main --abi example/token.json
	(cd example && go build)
