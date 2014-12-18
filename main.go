package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	var logDir string
	logDefault := "!!!log default string!!!"
	flag.StringVar(&logDir, "log", logDefault, "log directory")
	flag.Parse()

	if logDir != logDefault {
		f, err := os.OpenFile(logDir, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			panic(err)
		}
		log.SetOutput(f)
		defer f.Close()
	}

	r := mux.NewRouter()
	// kind: notifications, cook, etc...
	// TODO: kind and id validate
	r.HandleFunc("/ws/{kind}/{id}", wsHandler)

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
