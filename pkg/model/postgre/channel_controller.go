package postgre

import (
	"github.com/harisheoran/my-chat-system/pkg/model"
	"gorm.io/gorm"
)

type ChannelController struct {
	DbConnection *gorm.DB
}

// insert channel
func (cc *ChannelController) InsertChannel(channel *model.Channel) error {
	result := cc.DbConnection.Create(&channel)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// get all channels
func (cc *ChannelController) GetChannels() ([]model.Channel, error) {
	channels := []model.Channel{}

	result := cc.DbConnection.Find(&channels)

	if result.Error != nil {
		return channels, result.Error
	}

	return channels, nil
}
