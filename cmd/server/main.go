package main

import (
	"context"
	"flag"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mrlexus21/chat-server/internal/config"
	chat_v1 "github.com/mrlexus21/chat-server/pkg/chat/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"time"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	chat_v1.UnimplementedChatV1Server
	pool *pgxpool.Pool
}

func main() {
	flag.Parse()
	ctx := context.Background()

	// Считываем переменные окружения
	if err := config.Load(configPath); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	chat_v1.RegisterChatV1Server(s, &server{pool: pool})

	log.Printf("Server listening at %s", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *server) Create(ctx context.Context, _ *chat_v1.CreateRequest) (*chat_v1.CreateResponse, error) {
	chatID, err := s.insertChat(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("Created chat with id: %d", chatID)

	for userID := 1; userID <= 3; userID++ {
		if err := s.insertChatUser(ctx, chatID, int64(userID)); err != nil {
			return nil, err
		}
		log.Printf("Inserted chat user id: %d", userID)
	}

	return &chat_v1.CreateResponse{Id: chatID}, nil
}

func (s *server) insertChat(ctx context.Context) (int64, error) {
	query, args, err := sq.Insert("chats").
		PlaceholderFormat(sq.Dollar).
		Columns("created_at").
		Values(time.Now()).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build query for chat insert: %v", err)
	}

	var chatID int64
	if err := s.pool.QueryRow(ctx, query, args...).Scan(&chatID); err != nil {
		return 0, fmt.Errorf("failed to insert chat: %v", err)
	}

	return chatID, nil
}

func (s *server) insertChatUser(ctx context.Context, chatID, userID int64) error {
	query, args, err := sq.Insert("chat_users").
		PlaceholderFormat(sq.Dollar).
		Columns("chat_id", "user_id").
		Values(chatID, userID).
		Suffix("RETURNING user_id").
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query for chat users insert: %v", err)
	}

	var chatUserID int64
	if err := s.pool.QueryRow(ctx, query, args...).Scan(&chatUserID); err != nil {
		return fmt.Errorf("failed to insert chat users: %v", err)
	}

	return nil
}

func (s *server) Delete(ctx context.Context, req *chat_v1.DeleteRequest) (*emptypb.Empty, error) {
	tables := []string{"messages", "chat_users", "chats"}
	conditions := []map[string]interface{}{
		{"chat_id": req.GetId()},
		{"chat_id": req.GetId()},
		{"id": req.GetId()},
	}

	for i, table := range tables {
		query, args, err := sq.Delete(table).
			PlaceholderFormat(sq.Dollar).
			Where(sq.Eq(conditions[i])).
			ToSql()
		if err != nil {
			return nil, fmt.Errorf("failed to build query for %s: %v", table, err)
		}

		if _, err := s.pool.Exec(ctx, query, args...); err != nil {
			return nil, fmt.Errorf("failed to delete from %s: %v", table, err)
		}
	}

	log.Printf("Deleted chat: %d", req.GetId())
	return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *chat_v1.SendMessageRequest) (*emptypb.Empty, error) {
	const userID = 2

	query, args, err := sq.Select("chat_id").
		From("chat_users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"user_id": userID}).
		GroupBy("chat_id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query for chat users: %v", err)
	}

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to select chat users: %v", err)
	}
	defer rows.Close()

	var userChats []int64
	for rows.Next() {
		var chatID int64
		if err := rows.Scan(&chatID); err != nil {
			return nil, fmt.Errorf("failed to scan chat users: %v", err)
		}
		userChats = append(userChats, chatID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over chat users: %v", err)
	}

	for _, chatID := range userChats {
		if err := s.insertMessage(ctx, chatID, userID, req.GetText(), req.GetFrom()); err != nil {
			return nil, err
		}
	}

	return &emptypb.Empty{}, nil
}

func (s *server) insertMessage(ctx context.Context, chatID, userID int64, msgText, fromUser string) error {
	query, args, err := sq.Insert("messages").
		PlaceholderFormat(sq.Dollar).
		Columns("chat_id", "from_user_id", "msg_txt").
		Values(chatID, userID, msgText).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query for message insert: %v", err)
	}

	var messageID int64
	if err := s.pool.QueryRow(ctx, query, args...).Scan(&messageID); err != nil {
		return fmt.Errorf("failed to insert message: %v", err)
	}

	log.Printf("Send message from: %s to chat id: %d", fromUser, chatID)
	return nil
}
