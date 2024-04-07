build:
	@go build -o bin/vt cmd/main.go

run: build
	@./bin/vt