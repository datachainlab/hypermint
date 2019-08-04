GO_BIN?=go

BUILD_DIR?=./build
HMD?=$(BUILD_DIR)/hmd
HMCLI?=$(BUILD_DIR)/hmcli
HMD_HOME?=${HOME}/.hmd
HMCLI_HOME?=${HOME}/.hmcli

MNEMONIC?="token dash time stand brisk fatal health honey frozen brown flight kitchen"
HDW_PATH?=m/44'/60'/0'/0
HDW_VALIDATOR_IDX?=0

COMMIT_HASH:=$(shell git rev-parse --short HEAD)
VERSION:=$(shell cat version)
BUILD_FLAGS?=-ldflags "-X github.com/bluele/hypermint/pkg/consts.GitCommit=${COMMIT_HASH} -X github.com/bluele/hypermint/pkg/consts.Version=${VERSION}"

GO_BUILD_CMD=$(GO_BIN) build $(BUILD_FLAGS)
GO_TEST_FLAGS?=-v
GO_TEST_CMD=$(GO_BIN) test $(GO_TEST_FLAGS)

INCLUDE = -I=. -I=${GOPATH}/src -I=${GOPATH}/src/github.com/gogo/protobuf/protobuf

export GO111MODULE = on

.PHONY: build

build: server cli

server:
	$(GO_BUILD_CMD) -o $(HMD) ./cmd/hmd

cli:
	$(GO_BUILD_CMD) -o $(HMCLI) ./cmd/hmcli

release-build:
	$(GO_BUILD_CMD) -o $(HMD)_$(GOOS)_$(GOARCH)   ./cmd/hmd
	$(GO_BUILD_CMD) -o $(HMCLI)_$(GOOS)_$(GOARCH) ./cmd/hmcli

fmt:
	cd ./hmc && cargo fmt
	$(MAKE) -C ./tests fmt

start:
	$(HMD) start --log_level="main:error" --home=$(HMD_HOME)

clean:
	@rm -rf $(HMD_HOME) $(HMCLI_HOME)

init: clean init-validator
	$(eval ADDR1 := $(shell $(HMCLI) new --password=password --silent --home=$(HMCLI_HOME) --mnemonic=$(MNEMONIC) --hdw_path="$(HDW_PATH)/1" ))
	$(eval ADDR2 := $(shell $(HMCLI) new --password=password --silent --home=$(HMCLI_HOME) --mnemonic=$(MNEMONIC) --hdw_path="$(HDW_PATH)/2" ))
	@$(HMD) init --address=$(ADDR1) --home=$(HMD_HOME)
	@echo export ADDR1='$(ADDR1)'
	@echo export ADDR2='$(ADDR2)'

init-validator:
	@$(HMD) tendermint init-validator --home=$(HMD_HOME) --mnemonic=$(MNEMONIC) --hdw_path="$(HDW_PATH)/$(HDW_VALIDATOR_IDX)"

create-node:
	@$(HMD) create --home=$(HMD_HOME) --mnemonic=$(MNEMONIC) --hdw_path="$(HDW_PATH)/$(HDW_VALIDATOR_IDX)" --genesis=$(GENESIS)

test:
	$(GO_TEST_CMD) ./pkg/...

integration-test:
	$(MAKE) -C ./tests integration-test

e2e-test:
	$(MAKE) -C ./tests e2e-test

build-image:
	docker build . -t bluele/hypermint:${VERSION}

build-linux:
	docker run -v $(PWD):/go/src/github.com/bluele/hypermint -it golang:1.11.4-stretch make build

########################################
### Protobuf

protoc_all: protoc_proof

%.pb.go: %.proto
	## If you get the following error,
	## "error while loading shared libraries: libprotobuf.so.14: cannot open shared object file: No such file or directory"
	## See https://stackoverflow.com/a/25518702
	## Note the $< here is substituted for the %.proto
	## Note the $@ here is substituted for the %.pb.go
	protoc $(INCLUDE) $< --gogo_out=Mgoogle/protobuf/timestamp.proto=github.com/golang/protobuf/ptypes/timestamp,plugins=grpc:.

protoc_proof: pkg/proof/proof.pb.go

get_protoc:
	@# https://github.com/google/protobuf/releases
	curl -L https://github.com/google/protobuf/releases/download/v3.6.1/protobuf-cpp-3.6.1.tar.gz | tar xvz && \
		cd protobuf-3.6.1 && \
		DIST_LANG=cpp ./configure && \
		make && \
		make check && \
		sudo make install && \
		sudo ldconfig && \
		cd .. && \
		rm -rf protobuf-3.6.1
