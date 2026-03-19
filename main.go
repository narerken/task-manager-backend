package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Working on server 8080"))
	})

	log.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
