BIN_DIR := bin

install:
	go mod download
	go mod tidy

run-server:
	@trap 'true' INT; go run ./cmd/server/; exit 0

run-gui:
	go run ./cmd/gui

run-cli:
	@trap 'true' INT; go run ./cmd/cli/; exit 0

server: run-server

cli: run-cli

gui: run-gui

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/server ./cmd/server
	go build -o $(BIN_DIR)/cli    ./cmd/cli
	go build -o $(BIN_DIR)/gui    ./cmd/gui

clean:
	rm -rf server.log
	rm -rf $(BIN_DIR)
	go clean ./...

lint:
	go vet ./...
	@test -z "$$(gofmt -l .)" || { echo "Files not formatted (run: gofmt -w .):"; gofmt -l .; exit 1; }

kill:
	-fuser -k 8080/tcp

help:
	@echo "install         download dependencies"
	@echo "run-server      start the server (port 8080)"
	@echo "run-client      start the CLI client"
	@echo "run-client-gui  start the GUI client"
	@echo "build           build the 3 binaries in $(BIN_DIR)/"
	@echo "lint            go vet + gofmt check"
	@echo "lint-strict     lint + golangci-lint"
	@echo "clean           remove binaries and logs"
	@echo "kill            free port 8080"

.PHONY: install run-server run-client run-client-gui server cli gui build \
        lint clean kill help
