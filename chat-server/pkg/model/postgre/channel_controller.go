package postgre

import (
	"github.com/harisheoran/my-chat-system/pkg/model"
	"gorm.io/gorm"
)

type ChannelController struct {
	DbConnection *gorm.DB
}

func (cc *ChannelController) InsertChannel(channel *model.Channel) error {
	result := cc.DbConnection.Create(&channel)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
