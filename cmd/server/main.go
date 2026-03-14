// Package main реализует gRPC-сервер для управления чатами, включая создание, удаление и отправку сообщений.
package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	chatApi "github.com/mrlexus21/chat-server/internal/api/chat"
	"github.com/mrlexus21/chat-server/internal/config"
	chatRepository "github.com/mrlexus21/chat-server/internal/repository/chat"
	chatUserRepository "github.com/mrlexus21/chat-server/internal/repository/chat_user"
	messageRepository "github.com/mrlexus21/chat-server/internal/repository/message"
	chatService "github.com/mrlexus21/chat-server/internal/service/chat"
	messageService "github.com/mrlexus21/chat-server/internal/service/message"
	desc "github.com/mrlexus21/chat-server/pkg/chat/v1"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
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

	messageRepo := messageRepository.NewRepository(pool)
	chatUserRepo := chatUserRepository.NewRepository(pool)
	chatRepo := chatRepository.NewRepository(pool, chatUserRepo)

	chatServ := chatService.NewService(chatRepo)
	messageServ := messageService.NewService(messageRepo, chatUserRepo)

	desc.RegisterChatV1Server(s, chatApi.NewImplementation(chatServ, messageServ))

	log.Printf("Server listening at %s", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
