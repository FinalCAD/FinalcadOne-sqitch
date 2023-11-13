CMD_FOLDER=./cmd/sqitch-config
INTERNAL_FOLDER=./internal
PKG_FOLDER=./pkg
OUTPUT_FOLDER=./bin/

all: deps fmt build

# Install dependencies
deps:
	@go mod tidy

# build
build: fmt
	@go build -o $(OUTPUT_FOLDER) $(CMD_FOLDER)

# Run tests
test:
	@go test -v $(INTERNAL_FOLDER)/... ./cmd/... $(PKG_FOLDER)/...

# Lint
fmt:
	@go fmt ./...

# Clean
clean:
	@rm -rf $(OUTPUT_FOLDER)

# Run
run: build
	@go run ./...
