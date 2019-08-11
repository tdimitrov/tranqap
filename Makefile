BINARY_NAME=tranqap

all: build test

build:
	go build -o $(BINARY_NAME) ./cmd/tranqap

install:
	go install ./cmd/tranqap

test:
	go test -cover ./...

cyclo:
	gocyclo . capture output tqlog

clean:
	@rm tranqap || true

dep:
	go get golang.org/x/crypto/ssh
	go get github.com/abiosoft/ishell
