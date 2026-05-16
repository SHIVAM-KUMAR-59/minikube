package main

import (
	"flag"
	"log/slog"

	"github.com/SHIVAM-KUMAR-59/minikube/internal/worker"
)

func main() {
	// Take input for node ID and server URL from command line arguments.
	nodeIDFlag := flag.String("node-id", "node-1", "Unique identifier for the worker node")
	serverUrl := flag.String("server-url", "http://localhost:8080", "URL of the API server")
	flag.Parse()
	
	// Create a worker with that node ID
	worker, err := worker.NewWorker(*serverUrl, *nodeIDFlag)
	if err != nil {
		slog.Error("Error creating worker", "error", err)
		return
	}

	// Start the worker to periodically check for scheduled pods and run them.
	worker.Start()

	// Block forever
	select {}
}