package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("content-type", "application/json")
		fmt.Fprintln(res, `{"status": "ok", "message": "Minikube is running"}`)
		slog.Info("Ping request received")
	})

	slog.Info("Server running on port :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}