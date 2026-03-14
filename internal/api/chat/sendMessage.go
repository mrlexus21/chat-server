package chat

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/mrlexus21/chat-server/internal/converter"
	desc "github.com/mrlexus21/chat-server/pkg/chat/v1"
)

func (i *Implementation) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	if req.Message.From == "" || req.Message.Text == "" || req.Message.Timestamp == nil {
		return &emptypb.Empty{}, fmt.Errorf("message fields 'from', 'text' and 'timestamp' are required")
	}

	message := converter.ToMessageFromDesc(req.Message)

	chatID, err := i.messageService.SendMessage(ctx, message)
	if err != nil {
		log.Printf("Error send message: %v", err)

		return &emptypb.Empty{}, err
	}

	log.Printf("Sent message from: %s to chat id: %d", message.User.Name, chatID)

	return &emptypb.Empty{}, nil
}
