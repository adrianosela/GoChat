package chat

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

const (
	writeWait      = 10 * time.Second    // Time allowed to write a message to the peer
	pongWait       = 60 * time.Second    // Time allowed to read the next pong message from the peer
	pingPeriod     = (pongWait * 9) / 10 // Send pings to peer with this period. Must be less than pongWait
	maxMessageSize = 512                 // Maximum message size allowed from peer
)

// Peer is a middleman between the websocket connection and the contoller.
type Peer struct {
	ID      string
	Ctrl    *Controller
	WSConn  *websocket.Conn
	MsgChan chan []byte
}

// NewPeer is the constructor for the Peer abstraction of a websockets client
func NewPeer(ctrl *Controller, conn *websocket.Conn, outChan chan []byte) *Peer {
	return &Peer{
		ID:      uuid.Must(uuid.NewV4()).String(),
		Ctrl:    ctrl,
		WSConn:  conn,
		MsgChan: outChan,
	}
}

func (p *Peer) enroll() {
	p.Ctrl.RegisterChan <- p
	go p.writer()
	go p.reader()
}

func (p *Peer) leave() {
	p.Ctrl.DeregisterChan <- p
	p.WSConn.Close()
}

func (p *Peer) reader() {
	defer p.leave()
	p.WSConn.SetReadLimit(maxMessageSize)
	p.WSConn.SetReadDeadline(time.Now().Add(pongWait))
	p.WSConn.SetPongHandler(func(string) error { p.WSConn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, jsonMsg, err := p.WSConn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WS connection was closed unexpectedly: %s", err)
			}
			break
		}
		var msg Msg
		if err = json.Unmarshal(jsonMsg, &msg); err != nil {
			log.Printf("WS connection was closed unexpectedly: %s", err)
			break
		}
		msg.From = p.ID
		// FOR NOW... if form "to" field is empty, broadcast, else send to direct msg chan
		if msg.To != "" {
			p.Ctrl.DirectMsgChan <- &msg
			return
		}
		p.Ctrl.BroadcastChan <- []byte(msg.Data)
	}
}

func (p *Peer) writer() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		p.WSConn.Close()
	}()
	for {
		select {
		case message, ok := <-p.MsgChan:
			p.WSConn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				p.WSConn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := p.WSConn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			// Add queued chat messages to the current websocket message.
			n := len(p.MsgChan)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write(<-p.MsgChan)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			p.WSConn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := p.WSConn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
