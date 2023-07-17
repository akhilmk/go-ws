package boradcaster

import (
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func NewBoradcaster() Boradcaster {
	bc := broadcaster{
		wsConns: make(map[string]*websocket.Conn),
	}
	return bc
}

type Boradcaster interface {
	AddClient(conn *websocket.Conn)
}

type broadcaster struct {
	wsConns map[string]*websocket.Conn
}

func (bc broadcaster) AddClient(conn *websocket.Conn) {
	cId := uuid.NewString()
	bc.wsConns[cId] = conn
	go bc.socketReader(bc.wsConns[cId])
}

func (bc broadcaster) socketReader(wsConn *websocket.Conn) {
	log.Println("ws reader sarted..")
	go func() {
		for {
			_, p, err := wsConn.ReadMessage()
			if err != nil {
				log.Printf("ws msg read error, error:%v", err)
				break
			}
			log.Printf("ws reader msg:%v", string(p))
			bc.socketWriter(p)
		}
	}()
}

func (bc broadcaster) socketWriter(msg []byte) {
	for _, wsConn := range bc.wsConns {
		go func(wsConn *websocket.Conn) {
			err := wsConn.WriteMessage(1, msg)
			if err != nil {
				log.Printf("ws msg write error, error:%v", err)
			}
		}(wsConn)
	}
}
