BUILD_DIR?=../build
CONFIG_DIR?=./config
VERSION?=$(shell cat ../../version)
RESOURCES?=$(wildcard $(CONFIG_DIR)/*.yaml)
RESOURCES_OPT?=$(addprefix -f,$(RESOURCES))
VALS_NUM?=4

create:
	kubectl apply $(RESOURCES_OPT)

destory:
	kubectl delete $(RESOURCES_OPT)

gen-config: build-image
	@rm -f ./config/*.yaml
	@bash ./gen_testnet.sh $(VALS_NUM)

build-image:
	sed -e "s/<VERSION>/${VERSION}/g" -e "s/<VALS_NUM>/$(VALS_NUM)/g" Dockerfile-test.tpl | docker build -t bluele/hypermint-testnet:${VERSION} -
