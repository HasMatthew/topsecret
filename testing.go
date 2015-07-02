package main

import (
	"io"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

func hello1(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world1!")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	mux.HandleFunc("/1", hello1)
	http.ListenAndServe(":8000", mux)

}
