package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func foo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("foo"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/notifications/{id}", notificationHandler)
	r.HandleFunc("/", foo)

	http.Handle("/", r)

	err := http.ListenAndServe(":3210", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
