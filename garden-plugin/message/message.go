package message

// Message serialises network plugin args
type Message struct {
	Command []byte
	Handle  []byte
	Data    []byte
}

func New(command, handle string, data []byte) Message {
	return Message{
		Command: []byte(command),
		Handle:  []byte(handle),
		Data:    data,
	}
}
