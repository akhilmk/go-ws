package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	bc "github.com/akhilmk/go-ws/broadcast"
	"github.com/akhilmk/go-ws/util"
	"github.com/gorilla/websocket"
)

var (
	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     util.CORSCheck,
	}

	bCaster *bc.Broadcast
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bCaster = bc.NewBroadcaster(ctx)
	setupHandlers()

	fmt.Println("server started " + util.APP_PORT)
	log.Fatal(http.ListenAndServe(util.APP_PORT, nil))

	// TODOs
	// handle ping-pong
	// dockerize
}

func setupHandlers() {
	http.Handle("/", homePage())
	http.HandleFunc("/ws", wsHomePage)
}

func homePage() http.Handler {
	return http.FileServer(http.Dir("./frontend"))
}

func wsHomePage(w http.ResponseWriter, r *http.Request) {
	wsConn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error ws upgrade")
		return
	}

	userName := r.URL.Query().Get("user")

	valid := bCaster.ValidateUser(wsConn, userName)
	if !valid {
		log.Printf("ValidateUser user not valid")
		return
	}

	ctx := util.GetCtxWithUserName(context.Background(), userName)
	bCaster.AddClient(ctx, wsConn)

	log.Println("ws new connection success.. user=" + userName)
}
