.PHONY: run build tidy

run:
	go run ./cmd/surose-os

build:
	GOOS=linux GOARCH=amd64 go build -o surose-os ./cmd/surose-os

tidy:
	go mod tidy
