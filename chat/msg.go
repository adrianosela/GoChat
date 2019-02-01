package chat

// Msg represents a direct message from a peer to another
type Msg struct {
	From string `json:"from,omitempty"` // peer ID of sender
	To   string `json:"to,omitempty"`   // peer ID of receiver
	Data string `json:"data"`           // message body
}
