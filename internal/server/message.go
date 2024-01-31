package server

type Message struct {
	text []byte
}

func NewMessage(b []byte) *Message {
	return &Message{
		text: b,
	}
}
