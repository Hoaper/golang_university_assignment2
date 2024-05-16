package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/ws", HandleConnections)
	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
