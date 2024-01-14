# TODO: Update Binary BINARY_NAME
BINARY_NAME=certwatch
DOC_PATH=docs/

# This can be passed at the command line when running make
# like so: `BIN_DIR=/your/path/here make`
# If this value is not specified it will default to your GOBIN
# setting, and if that is empty it will default to a local `bin/`
# directory.
BIN_DIR ?= $(shell go env GOBIN)

ifeq ($(BIN_DIR),) 
	BIN_DIR = bin/
endif

# This is the default target when `make` is run without arguments
# will update CLI docs on every build. You can remove the docs if you
# want.
install: docs
	go install

# Build will create the BIN_DIR first if it does not exist
# Then create a binary inside of the destination.
# IMPORTANT: Could potentially require root or admin privileges
# depending on the path. It's best to install somewhere
# with user permissions.
build:
	mkdir -p $(BIN_DIR) && go build -o $(BIN_DIR)/$(BINARY_NAME)

run: 
	go run .

.PHONY: docs
docs:
	$(BINARY_NAME) docs

test: 
	go test -v ./...
