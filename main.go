package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	bc "github.com/akhilmk/go-ws/broadcast"
	"github.com/gorilla/websocket"
)

var (
	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	bCaster *bc.Broadcast
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bCaster = bc.NewBroadcaster(ctx)
	setupHandlers()

	fmt.Println("server started :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func setupHandlers() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsHomePage)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "message broadcaster service started..")
}

func wsHomePage(w http.ResponseWriter, r *http.Request) {
	wsConn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error ws upgrade")
		return
	}
	log.Println("ws new connection success..")
	bCaster.AddClient(context.Background(), wsConn)
}
