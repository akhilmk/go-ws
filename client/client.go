package client

import (
	"context"
	"encoding/json"
	"log"

	"github.com/akhilmk/go-ws/event"
	"github.com/akhilmk/go-ws/util"
	"github.com/gorilla/websocket"
)

type Broadcaster interface {
	BroadCast(event.Event)
	RemoveClient(msg string)
}

type ClientWS interface {
	SendMessage(msg event.Event)
}

type client struct {
	clientId    string
	wsConn      *websocket.Conn
	sendChan    chan event.Event
	ctx         context.Context
	cancel      context.CancelFunc
	broadCaster Broadcaster
}

func NewClient(ctxReq context.Context, conn *websocket.Conn, bc Broadcaster) (string, ClientWS) {
	cId := util.GetUserNameFromContext(ctxReq)
	ctx, cancel := context.WithCancel(ctxReq)

	c := &client{
		clientId:    cId,
		wsConn:      conn,
		ctx:         ctx,
		cancel:      cancel,
		sendChan:    make(chan event.Event),
		broadCaster: bc,
	}

	// start socket operation parallel.
	go c.socketReader()
	go c.socketWriter()

	return cId, c
}

func (c client) SendMessage(msg event.Event) {
	// un-buffered channel helps to write one message to socket at a time.
	c.sendChan <- msg
}

func (c *client) socketReader() {

	defer func() {
		log.Printf("socketReader ended")
		c.broadCaster.RemoveClient(c.clientId)
		c.cancel()
	}()

	for {
		_, p, err := c.wsConn.ReadMessage()
		if err != nil {
			log.Printf("socketReader, error:%v", err)

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// todo handle unexpected close
				log.Printf("socketReader UnexpectedClose")
			}
			break
		}

		var eventMsg event.Event
		if err := json.Unmarshal(p, &eventMsg); err != nil {
			log.Printf("socketReader error marshalling event")
		}

		// if no errors, send message to broadcaster.
		c.broadCaster.BroadCast(eventMsg)
	}

}

func (c *client) socketWriter() {

	defer func() {
		log.Println("socketWriter ended..")
	}()

	for {
		select {
		case <-c.ctx.Done():
			log.Println("socketWriter ctx.Done()..")
			return
		case msg, ok := <-c.sendChan:

			// usually not happen. server side closing, let client knows closing the connection
			if !ok {
				err := c.wsConn.WriteMessage(websocket.CloseMessage, nil)
				if err != nil {
					log.Println("socketWriter closing connection")
				}
				return
			}

			eventMsg := event.Event{Type: event.EventSendMessage, Payload: msg.Payload}
			b, _ := json.Marshal(eventMsg)
			err := c.wsConn.WriteMessage(websocket.TextMessage, b)
			if err != nil {
				log.Printf("socketWriter error:%v", err)
			}
		}
	}

}
