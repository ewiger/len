.PHONY: build test clean

build:
	mkdir -p bin
	go build -o bin/len ./cmd/len-cli

test:
	go test ./...

clean:
	rm -rf bin