TEST?=./...
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=verisignmdns

default: build

build: fmtcheck
	go install

test: fmtcheck
	go test $(TEST) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v -parallel 20 $(TESTARGS) -timeout 60s

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./main.go
	gofmt -s -w ./$(PKG_NAME)

# Currently required by tf-deploy compile
fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

lint:
	@echo "==> Checking source code against linters..."
	@golangci-lint run ./$(PKG_NAME)

tools:
	GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

.PHONY: build test testacc fmt fmtcheck lint tools test-compile
