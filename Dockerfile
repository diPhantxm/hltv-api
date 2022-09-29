FROM golang:latest

WORKDIR /hltv-api

COPY ./ ./

RUN go mod download
RUN go build -v -o ./bin/ ./cmd/app/main.go

CMD ["./bin/main"]