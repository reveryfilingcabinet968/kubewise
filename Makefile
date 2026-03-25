BINARY := kubectl-whatif
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GOFLAGS := -ldflags "-X main.version=$(VERSION) -X main.buildDate=$(BUILD_DATE)"

.PHONY: build install clean test-unit test-integration lint test-all fmt tidy vendor

build:
	@mkdir -p bin
	go build $(GOFLAGS) -o bin/$(BINARY) ./cmd/kubectl-whatif/

install: build
	cp bin/$(BINARY) $(shell go env GOPATH)/bin/

clean:
	rm -rf bin/
	rm -f coverage.out

test-unit:
	go test ./pkg/... ./internal/... ./scenarios/... -short -count=1 -race

test-integration:
	kind create cluster --name kubewise-test --config kind-config.yaml || true
	kubectl apply -f testdata/manifests/ --context kind-kubewise-test
	go test ./... -run Integration -count=1 -race -timeout 5m
	kind delete cluster --name kubewise-test

lint:
	golangci-lint run ./...

test-all: lint test-unit

fmt:
	gofmt -w -s .
	goimports -w .

tidy:
	go mod tidy

vendor: tidy
	go mod vendor

coverage:
	go test ./pkg/... ./internal/... -coverprofile=coverage.out -count=1
	go tool cover -func=coverage.out
