.PHONY: build
build:
	go build -v -o ./bin/ ./cmd/app/main.go

run:
	go run -v ./cmd/app/main.go

test:
	go test ./...

bench:
	go test -bench=. ./... -benchmem -run=^#

DEFAULT_GOAL: build