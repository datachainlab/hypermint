BUILD_DIR?=../../build
CONFIG_DIR?=./config
VERSION?=$(shell cat ../../version)
RESOURCES?=$(wildcard $(CONFIG_DIR)/*.yaml)
RESOURCES_OPT?=$(addprefix -f,$(RESOURCES))
VALS_NUM?=4
# mnemonic "token dash time stand brisk fatal health honey frozen brown flight kitchen"
# m/44'/60'/0'/0
GENESIS_ADDR?=0x1221a0726d56aEdeA9dBe2522DdAE3Dd8ED0f36c

create:
	kubectl apply $(RESOURCES_OPT)

destory:
	kubectl delete $(RESOURCES_OPT)

gen-config:
	@rm -f ./config/*.yaml
	@bash ./gen_testnet.sh $(VALS_NUM)

destory-pvc:
	@for n in $$(kubectl get pvc -o go-template --template '{{range .items}}{{.metadata.labels.app}}{{"\n"}}{{end}}' | grep 'hm-validator'); do \
		kubectl delete pvc -l app=$$n; \
	done

build-image:
	sed -e "s/<VERSION>/${VERSION}/g" -e "s/<VALS_NUM>/$(VALS_NUM)/g" -e "s/<GENESIS_ADDR>/$(GENESIS_ADDR)/g" Dockerfile-test.tpl | docker build --no-cache=true -t bluele/hypermint-testnet-sentry:${VERSION} -

balance0:
	$(eval port := $(shell kubectl get services/hm-validator-0 -o go-template='{{- range .spec.ports }}{{- if (eq "rpc" .name) }}{{.nodePort}}{{- end}}{{- end}}'))
	@$(BUILD_DIR)/hmcli balance --address=$(GENESIS_ADDR) --node=tcp://localhost:$(port)
