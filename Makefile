include .env.example
export

run:
	go mod tidy && go mod download && \
	go run  ./cmd/app
.PHONY: run
