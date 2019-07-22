GO_BIN?=go
GO_TEST_FLAGS?=-v
GO_TEST_CMD=$(GO_BIN) test $(GO_TEST_FLAGS)

test:
	$(GO_TEST_CMD) ./transaction/...
	$(MAKE) -C ./contract test
