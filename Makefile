ifeq ($(OS),Windows_NT)
    RM=-del
    APP_BINARY_PATH=.\bin\cli-player.exe
else
    RM=-rm
    APP_BINARY_PATH=./bin/cli-player
endif

all: build test
 
build:
	go build -o ./bin/ ./...

test:
	go test -v ./...
 
run:
	go build -o ${APP_BINARY_PATH} main.go
	./${APP_BINARY_PATH}
 
 .PHONY: clean
clean:
	go clean
	$(RM) ${APP_BINARY_PATH}