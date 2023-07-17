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
	go bc.socketReader(cId, bc.wsConns[cId])
}

func (bc broadcaster) socketReader(cId string, wsConn *websocket.Conn) {
	log.Println("ws reader sarted..")
	go func() {
		for {
			_, p, err := wsConn.ReadMessage()
			if err != nil {
				log.Printf("ws msg read error, client removed, error:%v", err)
				delete(bc.wsConns, cId)
				break
			}
			log.Printf("ws reader msg:%v", string(p))
			bc.socketWriter(cId, p)
		}
	}()
}

func (bc broadcaster) socketWriter(cId string, msg []byte) {
	for _, wsConn := range bc.wsConns {
		go func(wsConn *websocket.Conn) {
			err := wsConn.WriteMessage(1, msg)
			if err != nil {
				log.Printf("ws msg write error,client removed, error:%v", err)
				delete(bc.wsConns, cId)
			}
		}(wsConn)
	}
}
