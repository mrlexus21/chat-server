package converter

import (
	"github.com/mrlexus21/chat-server/internal/model"
	desc "github.com/mrlexus21/chat-server/pkg/chat/v1"
)

func ToUserMessageFromDesc(message *desc.Message) model.User {
	return model.User{
		ID:   1,
		Name: message.From,
	}
}

func ToMessageFromDesc(message *desc.Message) *model.Message {
	return &model.Message{
		User:      ToUserMessageFromDesc(message),
		Text:      message.Text,
		CreatedAt: message.Timestamp.AsTime(),
	}
}
