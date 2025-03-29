BINARY_NAME=checker
SOURCE_DIR=./src
PKGS = 	$(shell find $(SOURCE_DIR) -name '*.go' )        \
		./main.go
BIN_DIR=./bin

# OS detection
UNAME_S := $(shell uname -s)
CURRENT_OS := unknown
CURRENT_ARCH := $(shell uname -m)

ifeq ($(OS),Windows_NT)
    CURRENT_OS := windows
    BINARY_EXTENSION := .exe
else ifeq ($(UNAME_S),Darwin)
    CURRENT_OS := darwin
else ifeq ($(UNAME_S),Linux)
    CURRENT_OS := linux
endif

# Build configuration
GOARCH := amd64
ifeq ($(CURRENT_ARCH),arm64)
    GOARCH := arm64
endif

# Common tasks
dependencies:
	go get -d ./...

fmt:
	gofmt -s -l -w $(PKGS)

vet:
	go vet ./...

lint:
	GO_GOLANGCI_LINT_CLI_LINT_MODE=project
	GO_GOLANGCI_LINT_ARGUMENTS=["./.."]
	golangci-lint run -c .golangci.yml

deploy:
	GOOS=$(OS) GOARCH=$(ARCH) go build -o $(BIN_NAME) ./main.go

# Generic build function
build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BIN_DIR)/$(BINARY_NAME)$(BINARY_EXTENSION) ./main.go

# Platform-specific builds
build-linux: GOOS=linux
build-linux: BINARY_EXTENSION=
build-linux: build

build-macos: GOOS=darwin
build-macos: BINARY_EXTENSION=
build-macos: build

build-windows: GOOS=windows
build-windows: BINARY_EXTENSION=.exe
build-windows: build

dev:
	go run ./main.go

clean:
	go clean
	rm -f $(BIN_DIR)/*

.PHONY: dependencies fmt vet lint staticcheck build build-linux build-macos build-windows dev clean
