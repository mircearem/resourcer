build: 
	@go build -o bin/resources

run: build
	@./bin/resources