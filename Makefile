.PHONY: build test lint install

build:
	go build -o rta .

test:
	go test ./...

lint:
	go vet ./...

install:
	go install .
