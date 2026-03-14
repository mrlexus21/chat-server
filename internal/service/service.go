package service

import (
	"context"

	"github.com/mrlexus21/chat-server/internal/model"
)

type ChatService interface {
	Create(ctx context.Context, chat *model.Chat) (int64, error)
	Delete(ctx context.Context, chatID int64) error
}

type MessageService interface {
	SendMessage(ctx context.Context, message *model.Message) (int64, error)
}
