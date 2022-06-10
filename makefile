.PHONY: build
build:
	go build -v -o ./bin/ ./cmd/main.go

run:
	go run -v ./cmd/main.go

test:
	go test ./...

bench:
	go test -bench=. ./... -benchmem -run=^#

DEFAULT_GOAL: build