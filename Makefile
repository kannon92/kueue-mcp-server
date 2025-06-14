APP_NAME := kueue-mcp-server
SRC := main.go

.PHONY: all build clean

all: build

build:
	go build -o $(APP_NAME) $(SRC)

clean:
	rm -f $(APP_NAME)