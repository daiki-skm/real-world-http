package main

import (
	"io"
	"log"
	"net/http"
)

func base() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("receive request")
		io.WriteString(w, "Hello from origin server")
	})
	log.Println("Origin server start at :9001")
	log.Fatalln(http.ListenAndServe(":9001", nil))
}

func main() {
	base()
}
