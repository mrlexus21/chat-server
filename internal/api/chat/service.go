package chat

import (
	"github.com/mrlexus21/chat-server/internal/service"
	desc "github.com/mrlexus21/chat-server/pkg/chat/v1"
)

type Implementation struct {
	desc.UnimplementedChatV1Server
	chatService    service.ChatService
	messageService service.MessageService
}

func NewImplementation(
	chatService service.ChatService,
	messageService service.MessageService) *Implementation {
	return &Implementation{
		chatService:    chatService,
		messageService: messageService,
	}
}
