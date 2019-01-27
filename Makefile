PROJECT_NAME = githubby

CURRENT_DIR = $(shell pwd)
GO_CMD = go
GO_BUILD = $(GO_CMD) build

BINARY_EXT = .out
BINARY_NAME = $(PROJECT_NAME)$(BINARY_EXT)

install:
	go get github.com/teejays/clog

build:
	$(GO_BUILD) -o $(BINARY_NAME) main.go

run-dev:
	./$(BINARY_NAME) --dir $(CURRENT_DIR) --dev --wait-max 30 --max 3

run:
	./$(BINARY_NAME)
