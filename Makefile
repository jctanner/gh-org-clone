.PHONY: build clean install test

BINARY_NAME=gh-org-clone
BUILD_DIR=bin

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME)

clean:
	go clean
	rm -rf $(BUILD_DIR)

install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) $(HOME)/bin/$(BINARY_NAME)

test:
	go test ./...

all: clean build
