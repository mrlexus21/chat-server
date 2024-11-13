package main

import (
	"context"
	"crypto/rand"
	"fmt"
	chat_v1 "github.com/mrlexus21/chat-server/pkg/chat/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"math/big"
	"net"
)

const grpcPort = 50052

type server struct {
	chat_v1.UnimplementedUserV1Server
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	chat_v1.RegisterUserV1Server(s, &server{})

	log.Printf("Server listening at %s", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *server) Create(_ context.Context, req *chat_v1.CreateRequest) (*chat_v1.CreateResponse, error) {
	log.Printf("Create chat: %v", req.String())

	id, err := rand.Int(rand.Reader, big.NewInt(1<<63-1))
	if err != nil {
		log.Fatalf("Failed user id generate: %v", err)
		return &chat_v1.CreateResponse{
			Id: 0,
		}, nil
	}
	return &chat_v1.CreateResponse{Id: id.Int64()}, nil
}

func (s *server) Delete(_ context.Context, req *chat_v1.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("Delete chat: %d", req.GetId())

	return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(_ context.Context, req *chat_v1.SendMessageRequest) (*emptypb.Empty, error) {
	log.Printf("Send message from: %s", req.GetFrom())

	return &emptypb.Empty{}, nil
}
