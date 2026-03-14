package repository

import (
	"context"

	"github.com/mrlexus21/chat-server/internal/model"
)

type ChatRepository interface {
	Create(ctx context.Context, chat *model.Chat) (int64, error)
	Delete(ctx context.Context, chatID int64) error
}

type ChatUserRepository interface {
	ChatByUsers(userIDs []int64) (int64, error)
	FindUserChat(userID int64) (int64, error)
	Create(ctx context.Context, chatID, userID int64) error
}

type MessageRepository interface {
	Create(ctx context.Context, chatID, userID int64, msgText, userName string) error
}
