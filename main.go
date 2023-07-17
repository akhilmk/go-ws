package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/akhilmk/go-ws/boradcaster"
	"github.com/gorilla/websocket"
)

var (
	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	bCaster boradcaster.Boradcaster
)

func main() {
	bCaster = boradcaster.NewBoradcaster()
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
	wsUpgrader.CheckOrigin = func(r *http.Request) bool { return true }
	wsConn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error ws upgrade")
		return
	}
	log.Println("ws new connection success..")
	bCaster.AddClient(wsConn)
}
