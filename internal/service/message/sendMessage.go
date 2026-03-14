package message

import (
	"context"

	"github.com/mrlexus21/chat-server/internal/model"
)

func (s *serv) SendMessage(ctx context.Context, message *model.Message) (int64, error) {
	chatID, err := s.chatUserRepository.FindUserChat(message.User.ID)
	if err != nil {
		return 0, err
	}

	if err = s.messageRepository.Create(ctx, chatID, message.User.ID, message.Text, message.User.Name); err != nil {
		return 0, err
	}

	return chatID, nil
}
