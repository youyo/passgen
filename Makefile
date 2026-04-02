.PHONY: build test lint clean

build:
	go build -o passgen .

test:
	go test -v -race ./...

lint:
	go vet ./...

clean:
	rm -f passgen
