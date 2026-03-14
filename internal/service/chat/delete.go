package chat

import (
	"context"
)

func (s *serv) Delete(ctx context.Context, chatID int64) error {
	err := s.chatRepository.Delete(ctx, chatID)

	return err
}
