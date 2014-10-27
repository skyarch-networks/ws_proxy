package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ws/notifications/{id}", notificationHandler)

	http.Handle("/", r)

	log.Println("start ws server")
	err := http.ListenAndServe(":3210", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
