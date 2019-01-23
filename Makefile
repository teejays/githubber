PROJECT_NAME = githubby

GO_CMD = go
GO_BUILD = $(GO_CMD) build

BINARY_EXT = .out
BINARY_NAME = $(PROJECT_NAME)$(BINARY_EXT)

install:
	go get github.com/teejays/clog

build:
	$(GO_BUILD) -o $(BINARY_NAME) main.go

run:
	./$(BINARY_NAME)
