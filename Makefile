GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

WIN_BINARY_NAME=balancer.exe
BUILD_PATH=bin

build_win:
	$(GOBUILD) -o $(BUILD_PATH)/$(WIN_BINARY_NAME) -v

deps:
	$(GOGET) -u github.com/gorilla/mux

.PHONY: build_win
