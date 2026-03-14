package model

import (
	"time"

	"github.com/mrlexus21/chat-server/internal/repository/chat_user/model"
)

type Message struct {
	ID        int64      `db:"id"`
	ChatID    int64      `db:"chat_id"`
	User      model.User `db:""`
	Text      string     `db:"msg_text"`
	CreatedAt time.Time  `db:"created_at"`
}
