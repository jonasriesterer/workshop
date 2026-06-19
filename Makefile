.PHONY: run build tidy sqlc fmt lint test test-v

run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

tidy:
	go mod tidy

sqlc:
	~/go/bin/sqlc generate

fmt:
	goimports -w .

lint:
	golangci-lint run ./...

test:
	go test ./...

test-v:
	go test -v ./...
