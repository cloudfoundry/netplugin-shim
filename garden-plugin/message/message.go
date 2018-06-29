package message

import "fmt"

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

func (m *Message) String() string {
	return fmt.Sprintf(`{"Command":"%s","Handle":"%s", "Data":"%s"}`, string(m.Command), string(m.Handle), string(m.Data))
}
