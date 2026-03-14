package message

import (
	"github.com/mrlexus21/chat-server/internal/repository"
	"github.com/mrlexus21/chat-server/internal/service"
)

type serv struct {
	messageRepository  repository.MessageRepository
	chatUserRepository repository.ChatUserRepository
}

func NewService(messageRepository repository.MessageRepository, chatUserRepository repository.ChatUserRepository) service.MessageService {
	return &serv{
		messageRepository:  messageRepository,
		chatUserRepository: chatUserRepository,
	}
}
