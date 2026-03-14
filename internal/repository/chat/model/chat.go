package model

import "github.com/mrlexus21/chat-server/internal/repository/chat_user/model"

type Chat struct {
	ID    int64        `db:"id"`
	Users []model.User `db:""`
}
