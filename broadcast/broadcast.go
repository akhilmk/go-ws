package broadcast

import (
	"context"
	"log"
	"sync"

	"github.com/akhilmk/go-ws/client"
	"github.com/akhilmk/go-ws/event"
	"github.com/gorilla/websocket"
)

// keep all clients here, so Broadcast can inform all clients.
type Broadcast struct {
	clientsWs     map[string]client.ClientWS
	lock          *sync.RWMutex
	ctx           context.Context
	broadcastChan chan event.Event
}

func NewBroadcaster(ctx context.Context) *Broadcast {
	bc := &Broadcast{
		clientsWs:     make(map[string]client.ClientWS),
		lock:          &sync.RWMutex{},
		ctx:           ctx,
		broadcastChan: make(chan event.Event),
	}
	go bc.listenAndBroadcast()
	return bc
}

func (bc *Broadcast) ValidateUser(conn *websocket.Conn, user string) bool {
	bc.lock.Lock()
	_, exist := bc.clientsWs[user]
	bc.lock.Unlock()

	valid := true
	if exist {
		// close connection gracefully.
		cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "user-name-exist")
		if err := conn.WriteMessage(websocket.CloseMessage, cm); err != nil {
			log.Printf("ValidateUser WriteMessage error")
		}
		conn.Close()
		valid = false
	}
	return valid
}
func (bc *Broadcast) AddClient(ctxReq context.Context, conn *websocket.Conn) {
	// get new client
	cId, c := client.NewClient(ctxReq, conn, bc)

	// keep all clients in Broadcast.
	bc.lock.Lock()
	bc.clientsWs[cId] = c
	bc.lock.Unlock()
}

func (bc Broadcast) BroadCast(msg event.Event) {
	bc.broadcastChan <- msg
}

func (bc Broadcast) RemoveClient(clientId string) {
	bc.lock.Lock()
	delete(bc.clientsWs, clientId)
	bc.lock.Unlock()
}

func (bc *Broadcast) listenAndBroadcast() {
	log.Printf("listenClients start")

	for {
		select {
		case <-bc.ctx.Done():
			log.Printf("listenClients returned")
			return
		case msg, ok := <-bc.broadcastChan:
			if ok {

				// todo configure number of routine, in a way to handle max number of clients per second.
				for _, wsCl := range bc.clientsWs {
					go func(wsClient client.ClientWS) {

						wsClient.SendMessage(msg)
					}(wsCl)
				}
			}
		}
	}
}
