package message

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/mrlexus21/chat-server/internal/repository"
)

const (
	messagesTableName       = "messages"
	messagesChatsIdColumn   = "chat_id"
	messagesUserIdColumn    = "user_id"
	messagesTextColumn      = "msg_text"
	messagesCreatedAtColumn = "created_at"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.MessageRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, chatID, userID int64, msgText, userName string) error {
	query, args, err := sq.Insert(messagesTableName).
		PlaceholderFormat(sq.Dollar).
		Columns(messagesChatsIdColumn, messagesUserIdColumn, messagesTextColumn, messagesCreatedAtColumn).
		Values(chatID, userID, msgText, time.Now()).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build message insert: %v", err)
	}

	if _, err = r.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to insert message: %v", err)
	}

	return nil
}
