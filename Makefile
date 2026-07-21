.PHONY: run test lint build

run:
	go run cmd/server/main.go

test:
	go test ./... -v -count=1

lint:
	go vet ./...

build:
	go build -o bin/subscriptions-service ./cmd/server
