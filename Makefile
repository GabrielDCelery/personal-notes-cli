install:
	@go mod tidy

build: install
	@go build -o ./bin/pnotes
