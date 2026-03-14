package chat

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/mrlexus21/chat-server/internal/model"
	"github.com/mrlexus21/chat-server/internal/repository"
)

const (
	chatsTableName       = "chats"
	chatsIdColumn        = "id"
	chatsCreatedAtColumn = "created_at"
)

type repo struct {
	db                 *pgxpool.Pool
	chatUserRepository repository.ChatUserRepository
}

func NewRepository(db *pgxpool.Pool, chatUserRepository repository.ChatUserRepository) repository.ChatRepository {
	return &repo{
		db:                 db,
		chatUserRepository: chatUserRepository,
	}
}

func (r *repo) Create(ctx context.Context, chat *model.Chat) (int64, error) {
	userIDs := make([]int64, len(chat.Users))
	for i, user := range chat.Users {
		userIDs[i] = user.ID
	}

	chatID, err := r.chatUserRepository.ChatByUsers(userIDs)
	if err != nil {
		return 0, err
	}

	if chatID == 0 {
		query, args, err := sq.Insert(chatsTableName).
			PlaceholderFormat(sq.Dollar).
			Columns(chatsCreatedAtColumn).
			Values(time.Now()).
			Suffix("RETURNING " + chatsIdColumn).
			ToSql()
		if err != nil {
			return 0, fmt.Errorf("failed to build chat insert: %v", err)
		}

		if err = r.db.QueryRow(ctx, query, args...).Scan(&chatID); err != nil {
			return 0, fmt.Errorf("failed to insert chat: %v", err)
		}

		for _, userID := range userIDs {
			if err = r.chatUserRepository.Create(ctx, chatID, userID); err != nil {
				return 0, err
			}
		}
	}

	return chatID, nil
}

func (r *repo) Delete(ctx context.Context, chatID int64) error {
	query, args, err := sq.Delete(chatsTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{chatsIdColumn: chatID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %v", err)
	}

	if _, err = r.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to delete chat: %v", err)
	}

	return nil
}
