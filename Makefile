BINARY_NAME=rpcap

all: build test

build:
	go build -o $(BINARY_NAME)

test:
	go test -cover ./...

cyclo:
	gocyclo . capture output rplog

clean:
	@rm rpcap || true

dep:
	go get golang.org/x/crypto/ssh
	go get github.com/abiosoft/ishell
