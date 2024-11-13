package main

import (
	"context"
	"github.com/fatih/color"
	chat_v1 "github.com/mrlexus21/chat-server/pkg/chat/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

const (
	address = "localhost:50052"
	userID  = 1
)

func main() {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Failed to close connection: %v", err)
		}
	}(conn)

	c := chat_v1.NewUserV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Create(ctx, &chat_v1.CreateRequest{
		Usernames: []string{"Alex", "Tom"},
	})
	if err != nil {
		log.Fatalf("Failed create chat: %v", err)
	}

	log.Printf(color.RedString("Create chat info:\n"), color.GreenString("%+v", r.String()))
}
