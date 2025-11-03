.PHONY: all validator example test clean

all: validator example

validator:
	go build -o build/cobo-auth-data-validator ./cmd

example:
	go build -o build/cobo-auth-data-validator-example ./example

test:
	go test -v ./...

clean:
	go clean
	@rm -r build
