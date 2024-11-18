PROJECT_ROOT ?= .
GOPATH=$(shell go env GOPATH)
GO_BASE_CMD = GOTOOLCHAIN=go1.21.5 GO111MODULE=on GOSUMDB=sum.golang.org GOPROXY="https://proxy.golang.org,direct" go

.PHONY: test
test:
	rm -rf $(PROJECT_ROOT)/coverage.out
	@echo "  >  Running tests with coverage"
	$(GO_BASE_CMD) test -race -coverprofile=coverage.out -coverpkg=./... $(PROJECT_ROOT)/...
	$(GO_BASE_CMD) tool cover -func=coverage.out
