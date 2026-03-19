.PHONY: build test lint fmt vet clean

build:
	go build ./...

test:
	go test -race -cover ./...

lint:
	golangci-lint run ./...

fmt:
	gofmt -w .

vet:
	go vet ./...

clean:
	go clean ./...

check: fmt vet test lint
