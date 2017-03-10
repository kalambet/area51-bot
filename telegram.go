package area51bot

import (
	"log"
)

type Message struct {
	Type     string
	UserName string
	Name     string
	Text     string
}

func PostNotification(m *Message) {
	log.Printf("HAHA %#v", *m)
}
