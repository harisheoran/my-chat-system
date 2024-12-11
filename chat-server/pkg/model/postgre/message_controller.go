package postgre

import (
	"github.com/harisheoran/my-chat-system/pkg/model"
	"gorm.io/gorm"
)

type MessageController struct {
	DbConnection *gorm.DB
}

func (mc *MessageController) InsertMessage(message *model.Message) error {
	result := mc.DbConnection.Create(&message)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
