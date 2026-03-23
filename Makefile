.PHONY: build test docs docs-serve clean

build:
	mkdir -p bin
	go build -o bin/len ./cmd/len-cli

test:
	go test ./...

docs:
	mkdocs build --clean

docs-serve:
	mkdocs serve

clean:
	rm -rf bin doc/lang-html