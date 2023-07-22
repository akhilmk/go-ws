package broadcast

import (
	"context"
	"log"
	"sync"

	"github.com/akhilmk/go-ws/client"
	"github.com/gorilla/websocket"
)

// keep all clients here, so Broadcast can inform all clients.
type Broadcast struct {
	clientsWs     map[string]client.ClientWS
	lock          *sync.RWMutex
	ctx           context.Context
	broadcastChan chan string
}

func NewBroadcaster(ctx context.Context) *Broadcast {
	bc := &Broadcast{
		clientsWs:     make(map[string]client.ClientWS),
		lock:          &sync.RWMutex{},
		ctx:           ctx,
		broadcastChan: make(chan string),
	}
	go bc.listenClients()
	return bc
}

func (bc *Broadcast) AddClient(ctxReq context.Context, conn *websocket.Conn) {

	cId, c := client.NewClient(ctxReq, conn, bc)

	// keep all clients in Broadcast.
	bc.lock.Lock()
	bc.clientsWs[cId] = c
	bc.lock.Unlock()
}

func (bc Broadcast) BroadCast(msg string) {
	bc.broadcastChan <- msg
}

func (bc Broadcast) RemoveClient(clientId string) {
	bc.lock.Lock()
	delete(bc.clientsWs, clientId)
	bc.lock.Unlock()
}

func (bc *Broadcast) listenClients() {
	log.Printf("listenClients start")

	for {
		select {
		case <-bc.ctx.Done():
			log.Printf("listenClients returned")
			return
		case msg, ok := <-bc.broadcastChan:
			if ok {

				// todo configure number of routine, in a way to handle max number of clients per second.
				for _, wsC := range bc.clientsWs {
					go func(wc client.ClientWS) {
						wc.SendMessage(string(msg))
					}(wsC)
				}
			}
		}
	}
}
