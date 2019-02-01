package chat

import (
	"fmt"
	"log"
)

// Controller represents the central controller for a chat service
type Controller struct {
	OpenSessions   map[string]*Peer // Map of peer id to peer object
	BroadcastChan  chan []byte      // Channel to broadcast a message to all peers
	RegisterChan   chan *Peer       // Channel for registration requests
	DeregisterChan chan *Peer       // Channel for de-registration requests
	DirectMsgChan  chan *Msg	// Channel for messages from a peer to another peer
}

// NewController is the constructor for a chat controller
func NewController() *Controller {
	return &Controller{
		OpenSessions:   make(map[string]*Peer),
		BroadcastChan:  make(chan []byte),
		RegisterChan:   make(chan *Peer),
		DeregisterChan: make(chan *Peer),
		DirectMsgChan:  make(chan *Msg),
	}
}

// Start starts a chat service controller which handles messages from its channels
// TODO: add peer to peer communication
// TODO: add second layer of encyption for connection
func (c *Controller) Start() {
	for {
		select {
		// register new peer
		case newPeer := <-c.RegisterChan:
			c.registerPeer(newPeer)
			// deregister an existing peer
		case leavingPeer := <-c.DeregisterChan:
			c.deregisterPeer(leavingPeer)
			// broadcast message to all peers
		case msgBytes := <-c.BroadcastChan:
			c.broadcastMessage(msgBytes)
			// send direct message from peer to peer
		case msgObject := <-c.DirectMsgChan:
			c.sendMessage(msgObject)
		}
	}
}

func (c *Controller) registerPeer(p *Peer) {
	p.MsgChan <- []byte(welcomeMessage(p.ID))
	log.Printf("Peer %s joined server\n", p.ID)
	c.OpenSessions[p.ID] = p
}

func (c *Controller) deregisterPeer(p *Peer) {
	if _, ok := c.OpenSessions[p.ID]; ok {
		delete(c.OpenSessions, p.ID)
		close(p.MsgChan)
	}
	log.Printf("Peer %s left server\n", p.ID)
}

func (c *Controller) broadcastMessage(m []byte) {
	log.Printf("Broadcasting message to %d peers\n", len(c.OpenSessions))
	for peerID := range c.OpenSessions {
		select {
		case c.OpenSessions[peerID].MsgChan <- m:
		default:
			close(c.OpenSessions[peerID].MsgChan)
			delete(c.OpenSessions, peerID)
		}
	}
}

func (c *Controller) sendMessage(msg *Msg) {
	log.Printf("Directing message from peer %s to peer %s\n", msg.From, msg.To)
	if to, ok := c.OpenSessions[msg.To]; ok {
		select {
		case c.OpenSessions[to.ID].MsgChan <- []byte(fmt.Sprintf("%s [from peer %s]", msg.Data, msg.From)):
		default:
			close(c.OpenSessions[to.ID].MsgChan)
			delete(c.OpenSessions, to.ID)
		}
	}
}

func welcomeMessage(id string) string {
	return fmt.Sprintf("Welcome to the anonymous chat server! Your peer id is: \n%s", id)
}
