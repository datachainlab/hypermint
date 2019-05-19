VERSION ?= 0.1.0

build-image:
	docker build . -t bluele/go-rust:$(VERSION)

upload-image:
	docker push bluele/go-rust:$(VERSION)

fmt:
	cd contract/contract-test && cargo fmt
