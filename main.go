package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/notifications/{id}", notificationHandler)

	http.Handle("/", r)

	err := http.ListenAndServe(":3210", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
