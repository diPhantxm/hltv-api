.PHONY: build
build:
	go build -v -o ./bin/ ./cmd/app/main.go

run:
	go run -v ./cmd/app/main.go

test:
	go test ./...

bench:
	go test -bench=. ./... -benchmem -run=^#

covertest:
	go test -coverprofile=cover.out internal/parsers
	go tool cover -html cover.out

DEFAULT_GOAL: build