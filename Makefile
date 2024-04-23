build:
	@go build -o bin/blog .

run:
	@make build && ./bin/blog