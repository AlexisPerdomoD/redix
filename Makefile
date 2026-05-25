BINARY=redix
BUILD_DIR=build

.PHONY: all build clean run test lint fmt

all: build

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY) ./cmd/redix

run: build
	./$(BUILD_DIR)/$(BINARY)

clean:
	rm -rf $(BUILD_DIR)

test:
	go test ./...

lint:
	go vet ./...

fmt:
	go fmt ./...
