package main

import (
	"golang.org/x/net/websocket"
	"net/http"
	"log"
	"github.com/go-redis/redis"
)

var rDb = redis.NewClient(&redis.Options{
	Addr:     "127.0.0.1:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
});

func pushMsg(c *websocket.Conn, msg chan string) {
	k := <-msg
	lenInt64 := rDb.LLen(k)
	lenE := lenInt64.Val()
	listRs := rDb.LRange(k, 0, lenE)
	listRsE := listRs.Val()
	for i := 0; i < int(lenE); i++ {
		websocket.JSON.Send(c, listRsE[i])
	}
}

func readMsg(c *websocket.Conn, msg chan string) {
	m := make([]byte, 1024)
	n, err := c.Read(m)
	if err != nil {
		log.Fatal("Read Error" + err.Error())
		c.Close()
	}
	msg <- string(m[:n])
}

func handleConn(c *websocket.Conn) {
	msg := make(chan string, 1)
	for {
		readMsg(c, msg)
		pushMsg(c, msg)
	}
}

func main() {
	go http.Handle("/echo", websocket.Handler(handleConn))
	http.Handle("/", http.FileServer(http.Dir(".")))

	err := http.ListenAndServe(":888", nil)
	if err != nil {
		panic("Listen:" + err.Error())
	}
}
