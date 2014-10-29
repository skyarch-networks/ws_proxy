package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

func notificationHandler(w http.ResponseWriter, r *http.Request) {
	// params
	vars := mux.Vars(r)
	id := vars["id"]
	kind := vars["kind"]

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	} else {
		log.Println("upgraded")
	}

	quit := make(chan bool)
	go writer(ws, kind, id, quit)
	reader(ws, quit)
}

func writer(ws *websocket.Conn, kind string, id string, quit <-chan bool) {
	defer ws.Close()

	pingTicker := time.NewTicker(pingPeriod)
	defer pingTicker.Stop()

	defer log.Println("disconnect websocket")

	quitSub := make(chan bool)
	defer func() { quitSub <- true }()

	chs := sub(kind, id, (<-chan bool)(quitSub))
	for {
		select {
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
			log.Println("ping")

		case v := <-chs:
			err := ws.WriteMessage(websocket.TextMessage, v)
			if err != nil {
				log.Println(err)
				return
			}
			log.Printf("send data: %v", string(v))

		case <-quit:
			return
		}
	}
}

func reader(ws *websocket.Conn, quit chan<- bool) {
	defer ws.Close()

	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		log.Println("pong")
		return ws.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Printf("reader error: %v", err)
			quit <- true
			break
		}
	}
}
