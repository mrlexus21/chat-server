package model

import "time"

type Message struct {
	ID        int64
	ChatID    int64
	User      User
	Text      string
	CreatedAt time.Time
}
