BINARY_NAME=rpcap

all: build test

build:
	go build -o $(BINARY_NAME)

test:
	go test -cover ./...

clean:
	rm rpcap

dep:
	go get github.com/abiosoft/ishell