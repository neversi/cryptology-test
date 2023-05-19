-include .env

export GO111MODULE=on
export GOPROXY=https://proxy.golang.org
BUILD_ENVPARMS:=CGO_ENABLED=0

MOCKGEN_VERSION := 1.6.0

# install project dependencies
.PHONY: deps
deps:
	@echo 'install dependencies'
	go mod tidy -v

.PHONY: test
test: gen
	@echo 'running tests'
	go test -v -cover -race ./... -count=1

bench: gen
	@echo 'running benchmarks'
	 go test -bench=. -count 5 -benchmem ./...

.PHONY: lint
lint:
	@echo 'run golangci lint'
	@if ! bin/golangci-lint --version | grep -q $(LINTER_VERSION); \
		then curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v$(LINTER_VERSION); fi;
	bin/golangci-lint run --out-format=tab
	@echo

.PHONY: build-binary
build-binary:
	@echo 'build app $(APP_NAME)'
	$(shell $(BUILD_ENVPARMS) go build .)

.PHONY: build
build: deps build-binary

.PHONY:
build-image:
	@echo 'build docker image'
	docker build --no-cache --rm --progress=plain -t $(REGISTRY_HOST)/$(PROJECT_NAME)/$(APP_NAME):$(APP_VERSION) \
		-f $(DOCKERFILE_NAME) --build-arg GITHUB_CREDS=$(GITHUB_CREDS) --build-arg OUTPUT_BINARY=$(OUTPUT_BINARY) \
		--build-arg APP_VERSION=$(APP_VERSION) --build-arg BUILD_DIR=$(BUILD_DIR) .

.PHONY: install-mockgen gen
install-mockgen: deps
	go install github.com/golang/mock/mockgen@v$(MOCKGEN_VERSION)

gen: deps install-mockgen
	go generate ./...