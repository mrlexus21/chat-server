package converter

import (
	"github.com/mrlexus21/chat-server/internal/model"
	desc "github.com/mrlexus21/chat-server/pkg/chat/v1"
)

func ToUserNamesFromDesc(userNames *desc.UserNames) []model.User {
	id := int64(0)
	users := make([]model.User, 0, len(userNames.Names))
	for _, name := range userNames.Names {
		id++
		users = append(users, model.User{
			ID:   id,
			Name: name,
		})
	}

	return users
}
