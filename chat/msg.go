package chat

// DirectMsg represents a direct message from a peer to another
type DirectMsg struct {
	From string // peer ID of sender
	To   string // peer ID of receiver
	Data []byte // message bytes
}
