# Variables
BINARY_NAME = minik
SERVER_BINARY = minik-server
WORKER_BINARY = minik-worker

# Build for current OS
build:
	go build -o $(BINARY_NAME) cmd/minik/main.go
	go build -o $(SERVER_BINARY) cmd/server/main.go
	go build -o $(WORKER_BINARY) cmd/worker/main.go

# Build for all platforms
build-all:
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-darwin-arm64 cmd/minik/main.go
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin-amd64 cmd/minik/main.go
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 cmd/minik/main.go
	GOOS=darwin GOARCH=arm64 go build -o $(SERVER_BINARY)-darwin-arm64 cmd/server/main.go
	GOOS=darwin GOARCH=amd64 go build -o $(SERVER_BINARY)-darwin-amd64 cmd/server/main.go
	GOOS=linux GOARCH=amd64 go build -o $(SERVER_BINARY)-linux-amd64 cmd/server/main.go
	GOOS=darwin GOARCH=arm64 go build -o $(WORKER_BINARY)-darwin-arm64 cmd/worker/main.go
	GOOS=darwin GOARCH=amd64 go build -o $(WORKER_BINARY)-darwin-amd64 cmd/worker/main.go
	GOOS=linux GOARCH=amd64 go build -o $(WORKER_BINARY)-linux-amd64 cmd/worker/main.go

# Clean all binaries
clean:
	rm -f $(BINARY_NAME) $(SERVER_BINARY) $(WORKER_BINARY)
	rm -f $(BINARY_NAME)-* $(SERVER_BINARY)-* $(WORKER_BINARY)-*

# Install to /usr/local/bin
install: build
	sudo cp $(BINARY_NAME) /usr/local/bin/
	sudo cp $(SERVER_BINARY) /usr/local/bin/
	sudo cp $(WORKER_BINARY) /usr/local/bin/
	sudo rm -rf /usr/local/bin/dashboard
	sudo cp -r dashboard /usr/local/bin/dashboard
	sudo rm -rf /usr/local/bin/dashboard/node_modules
	cd /usr/local/bin/dashboard && sudo npm install
	sudo chown -R $(USER) /usr/local/bin/dashboard