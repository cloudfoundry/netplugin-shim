package message

// Message serialises network plugin args
type Message struct {
	Command []byte
	Handle  []byte
	Data    []byte
}
