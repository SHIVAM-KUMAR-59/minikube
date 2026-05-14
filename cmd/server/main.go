package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("content-type", "application/json")
		fmt.Fprintln(res, `{"status": "ok", "message": "minikube server is running"}`)
	})

	log.Println("Server running on port :8000")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}