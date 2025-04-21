package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page!")
}

func main() {
	fmt.Println("Hello, World!")

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler).Methods(http.MethodGet)

	fmt.Println("Starting server on port: 5500")
	if err := http.ListenAndServe(":5500", r); err != nil {
		log.Fatal(err)
	}
}