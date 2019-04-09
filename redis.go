package main

import (
        "os"
	"log"

	"github.com/garyburd/redigo/redis" // redigoが使いづらい気がする
)

func sub(kind, id string, quit <-chan bool) <-chan []byte {
	ch := make(chan []byte)

	go func() {
		chs, c := _sub(kind, id)
		defer c.Close()

		for {
			select {
			case <-quit:
				return
			case v := <-chs:
				ch <- v
			}
		}
	}()

	return (<-chan []byte)(ch)
}

func _sub(kind, id string) (<-chan []byte, redis.Conn) {
	ch := make(chan []byte)
	c, err := redis.Dial("tcp", os.Getenv("REDIS_HOST") + ":6379")
	if err != nil {
		panic(err)
	}

	go func() {
		defer c.Close()

		psc := redis.PubSubConn{c}
		psc.Subscribe(kind + "." + id)

		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				ch <- v.Data
				// case redis.Subscription:
			case redis.Error:
				log.Println(v)
				return
			case error:
				// has many case, connection closed
				log.Println(v)
				return
			default:
				log.Println(v)
			}
		}
	}()

	return ch, c
}
