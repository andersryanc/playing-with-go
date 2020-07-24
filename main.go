package main

import (
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
)

func main() {
	port := flag.String("port", "8080", "the port to listen on")
	flag.Parse()

	fmt.Println("Listening @ http://localhost:" + *port)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
