BUILD_DIR?=../build

init:
	kubectl apply -f ./deployment.yaml -f ./svc.yaml --record
	$(eval port := $(shell kubectl get services/hypermint -o go-template='{{(index .spec.ports 0).nodePort}}'))
	$(BUILD_DIR)/hmcli balance --address=${ADDR1} --node=tcp://localhost:$(port)

clean:
	kubectl delete -f ./deployment.yaml -f ./svc.yaml
