.PHONY: build
build:
	go build -v -o ./bin/ ./cmd/main.go

run:
	go run -v ./cmd/main.go

DEFAULT_GOAL: build