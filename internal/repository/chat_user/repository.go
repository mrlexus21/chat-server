package chat_user

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/mrlexus21/chat-server/internal/repository"
)

const (
	chatsTableName = "chats"
	chatsIdColumn  = "id"

	chatsUsersTableName     = "chat_users"
	chatsUsersChatsIdColumn = "chat_id"
	chatsUsersUserIdColumn  = "user_id"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.ChatUserRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, chatID, userID int64) error {
	query, args, err := sq.Insert(chatsUsersTableName).
		PlaceholderFormat(sq.Dollar).
		Columns(chatsUsersChatsIdColumn, chatsUsersUserIdColumn).
		Values(chatID, userID).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build chat users insert: %v", err)
	}

	if _, err = r.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to insert chat users: %v", err)
	}

	return nil
}

func (r *repo) ChatByUsers(userIDs []int64) (int64, error) {
	subQuery, subArgs, err := sq.Select(chatsUsersChatsIdColumn).
		From(chatsUsersTableName).
		Where(sq.Eq{chatsUsersUserIdColumn: userIDs}).
		GroupBy(chatsUsersChatsIdColumn).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build subquery: %v", err)
	}

	query, args, err := sq.Select("c."+chatsIdColumn).
		From(chatsTableName+" c").
		Join(chatsUsersTableName+" cu ON c."+chatsIdColumn+" = cu."+chatsUsersChatsIdColumn).
		Where(fmt.Sprintf("c.%s IN (%s)", chatsIdColumn, subQuery), subArgs...).
		GroupBy("c." + chatsIdColumn).
		Having(fmt.Sprintf("COUNT(DISTINCT cu.%s) = %d", chatsUsersUserIdColumn, len(userIDs))).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %v", err)
	}

	var chatID int64
	err = r.db.QueryRow(context.Background(), query, args...).Scan(&chatID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return chatID, nil
		}

		return 0, fmt.Errorf("failed to get chat by users: %v", err)
	}

	return chatID, nil
}

func (r *repo) FindUserChat(userID int64) (int64, error) {
	query, args, err := sq.Select(chatsUsersChatsIdColumn).
		From(chatsUsersTableName).
		Where(sq.Eq{chatsUsersUserIdColumn: userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %v", err)
	}

	var chatID int64
	err = r.db.QueryRow(context.Background(), query, args...).Scan(&chatID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("not found chat by user %d: %v", userID, err)
		}

		return 0, fmt.Errorf("failed to get chat by user %d: %v", userID, err)
	}

	return chatID, nil
}
