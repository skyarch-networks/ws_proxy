package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	// kind: notifications, cook, etc...
	// TODO: kind and id validate
	r.HandleFunc("/ws/{kind}/{id}", notificationHandler)

	http.Handle("/", r)

	log.Println("start ws server")
	err := http.ListenAndServe(":3210", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
