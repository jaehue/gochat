package gochat

import "time"

type Message struct {
	Sender    string
	Text      string
	CreatedAt time.Time
}
