package client

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Broadcaster interface {
	BroadCast(msg string)
	RemoveClient(msg string)
}

type ClientWS interface {
	SendMessage(msg string)
}

type client struct {
	clientId    string
	wsConn      *websocket.Conn
	sendChan    chan string // todo make message Event struct based or []byte on requirement
	ctx         context.Context
	cancel      context.CancelFunc
	broadCaster Broadcaster
}

func NewClient(ctxReq context.Context, conn *websocket.Conn, bc Broadcaster) (string, ClientWS) {
	cId := uuid.NewString()
	ctx, cancel := context.WithCancel(ctxReq)

	c := &client{
		clientId:    cId,
		wsConn:      conn,
		ctx:         ctx,
		cancel:      cancel,
		sendChan:    make(chan string),
		broadCaster: bc,
	}

	// start socket operation parallel.
	go c.socketReader()
	go c.socketWriter()

	return cId, c
}

func (c client) SendMessage(msg string) {
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

		// if no errors, send message to broadcaster.
		c.broadCaster.BroadCast(string(p))
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

			err := c.wsConn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Printf("socketWriter error:%v", err)
			}
		}
	}

}
