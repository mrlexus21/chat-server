package chat

import (
	"context"
	"fmt"
	"log"

	"github.com/mrlexus21/chat-server/internal/converter"
	"github.com/mrlexus21/chat-server/internal/model"
	desc "github.com/mrlexus21/chat-server/pkg/chat/v1"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	if req.UserNames == nil {
		return &desc.CreateResponse{Id: 0}, fmt.Errorf("chat user names are required")
	}

	chat := &model.Chat{
		ID:    0,
		Users: converter.ToUserNamesFromDesc(req.UserNames),
	}

	chatID, err := i.chatService.Create(ctx, chat)
	if err != nil {
		log.Printf("Error created chat: %v", err)

		return &desc.CreateResponse{Id: 0}, err
	}

	log.Printf("Created chat with id: %d", chatID)

	return &desc.CreateResponse{Id: chatID}, nil
}
