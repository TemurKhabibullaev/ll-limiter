.PHONY: test run build

test:
	go test ./... -race

run:
	RATE_PER_SEC=50 BURST=100 PORT=8080 go run ./cmd/limiter

build:
	go build -o bin/limiter ./cmd/limiter
