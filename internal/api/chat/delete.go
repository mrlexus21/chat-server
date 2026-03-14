package chat

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/mrlexus21/chat-server/pkg/chat/v1"
)

func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	if req.Id == 0 {
		return &emptypb.Empty{}, fmt.Errorf("chat id are required")
	}

	err := i.chatService.Delete(ctx, req.Id)
	if err != nil {
		log.Printf("Error deleted chat: %v", err)

		return &emptypb.Empty{}, err
	}

	log.Printf("Deleted chat with id: %d", req.Id)

	return &emptypb.Empty{}, nil
}
