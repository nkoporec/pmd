include .env.example
export

run:
	go mod tidy && go mod download && \
	go run main.go $(CMD)
.PHONY: run
