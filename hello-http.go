package main

import (
	"fmt"
	"http"
)

func handler(c http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(c, "Hello, %s.", r.URL.Path[1:])
}

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(handler))
}
