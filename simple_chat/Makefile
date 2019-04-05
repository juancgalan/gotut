GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

WIN_BINARY_NAME=readln.exe
BINARY_NAME=readln
BUILD_PATH=bin

build_win:
	$(GOBUILD) -o $(BUILD_PATH)/$(WIN_BINARY_NAME) -v

build:
	$(GOBUILD) -o $(BUILD_PATH)/$(BINARY_NAME) -v

deps:
	$(GOGET) -u github.com/chzyer/readline

run: build
	$(BUILD_PATH)/$(BINARY_NAME)

.PHONY: build
