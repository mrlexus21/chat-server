package chat

import (
	"context"

	"github.com/mrlexus21/chat-server/internal/model"
)

func (s *serv) Create(ctx context.Context, chat *model.Chat) (int64, error) {
	chatID, err := s.chatRepository.Create(ctx, chat)
	if err != nil {
		return 0, err
	}

	return chatID, nil
}
