package chat

import "log"

// Controller represents the central controller for a chat service
type Controller struct {
	OpenSessions   map[string]*Peer // Map of peer to session data (now just a bool)
	BroadcastChan  chan []byte      // Channel to broadcast a message to all peers
	RegisterChan   chan *Peer       // Channel for registration requests
	DeregisterChan chan *Peer       // Channel for de-registration requests
}

// NewController is the constructor for a chat controller
func NewController() *Controller {
	return &Controller{
		OpenSessions:   make(map[string]*Peer),
		BroadcastChan:  make(chan []byte),
		RegisterChan:   make(chan *Peer),
		DeregisterChan: make(chan *Peer),
	}
}

// Start starts a chat service controller which handles messages from its channels
// TODO: add peer to peer communication
// TODO: add second layer of encyption for connection
func (c *Controller) Start() {
	for {
		select {
		// register new peer
		case peer := <-c.RegisterChan:
			log.Printf("Peer %s joined server\n", peer.ID)
			c.OpenSessions[peer.ID] = peer
		// deregister an existing peer
		case peer := <-c.DeregisterChan:
			if _, ok := c.OpenSessions[peer.ID]; ok {
				delete(c.OpenSessions, peer.ID)
				close(peer.OutboundMsgChan)
			}
			log.Printf("Peer %s left server\n", peer.ID)
		// broadcast message to all peers
		case msg := <-c.BroadcastChan:
			log.Printf("Broadcasting message to %d peers\n", len(c.OpenSessions))
			for peerID := range c.OpenSessions {
				select {
				case c.OpenSessions[peerID].OutboundMsgChan <- msg:
				default:
					close(c.OpenSessions[peerID].OutboundMsgChan)
					delete(c.OpenSessions, peerID)
				}
			}
		}
	}
}
