.PHONY: run build clean

# Run the server in development
run:
	go run cmd/server/main.go

# Build the binary
build:
	go build -o bin/server cmd/server/main.go

# Clean build artifacts
clean:
	rm -rf bin/
