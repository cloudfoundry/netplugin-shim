package message

// Message serialises network plugin args
type Message struct {
	Command string
	Data    interface{}
}
