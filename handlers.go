package main

import (
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func notificationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		log.Println(err)
		return
	}

	quit := make(chan int)
	go writer(ws, id, quit)
	reader(ws, quit)
}

func writer(ws *websocket.Conn, id string, quit chan int) {
	defer ws.Close()

	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		// panic?
		panic(err)
	}
	defer c.Close()

	psc := redis.PubSubConn{c}
	psc.Subscribe("notifications." + id)

	loop := true
	go func() {
		<-quit
		loop = false
	}()
	for loop {
		switch v := psc.Receive().(type) {
		case redis.Message:
			err = ws.WriteMessage(websocket.TextMessage, v.Data)
			if err != nil {
				log.Fatal(err)
			}
			// case redis.Subscription:
			// case error:
		}
	}
}

func reader(ws *websocket.Conn, quit chan int) {
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			quit <- 0
			break
		}
	}
}
