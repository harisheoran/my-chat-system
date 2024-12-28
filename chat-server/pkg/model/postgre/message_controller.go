package postgre

import (
	"github.com/harisheoran/my-chat-system/pkg/model"
	"gorm.io/gorm"
)

type MessageController struct {
	DbConnection *gorm.DB
}

func (mc *MessageController) BulkInsertMessage(messages *[]model.Message) error {
	result := mc.DbConnection.Create(&messages)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (mc *MessageController) GetLastMessages() ([]model.Message, error) {
	var messages = []model.Message{}
	result := mc.DbConnection.Order("created_at desc").Limit(20).Find(&messages)

	if result.Error != nil {
		return messages, result.Error
	}

	return messages, nil
}
