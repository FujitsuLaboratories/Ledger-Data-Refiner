/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package model

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type Channel struct {
	ID               int       `json:"id" gorm:"primarykey"`
	NetworkName      string    `json:"network_name" gorm:"type:varchar(256)"`
	ClientName       string    `json:"client_name" gorm:"type:varchar(256)"`
	ChannelName      string    `json:"channel_name" gorm:"type:varchar(256)"`
	BlockCount       int64     `json:"block_count" gorm:"type:integer"`
	TxCount          int64     `json:"tx_count" gorm:"type:integer"`
	ChannelVersion   uint64    `json:"channel_version" gorm:"type:integer"`
	GenesisBlockHash string    `json:"genesis_block_hash" gorm:"type:varchar(256)"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// TableName implements Tabler
func (c *Channel) TableName() string {
	return "channel"
}

func InsertChannel(channel *Channel) (int, error) {
	logger.WithField("channel_name", channel.ChannelName).Debug("insert one channel")

	result := new(Channel)
	// check if the channel exists
	err := db.Select("id").Where(&Channel{ChannelName: channel.ChannelName, GenesisBlockHash: channel.GenesisBlockHash}).
		First(result).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, errors.Wrap(err, "error searching for specified channel")
	}

	if result.ID > 0 {
		logger.WithField("channel_name", channel.ChannelName).Info("channel exists")
		return result.ID, nil
	}

	// insert channel
	err = db.Create(channel).Error
	return channel.ID, err
}

// UpdateChannel updates channel with the specified channel id
func UpdateChannel(session *gorm.DB, block *Block, channelID int) error {
	if session == nil {
		session = db
	}

	channel := new(Channel)

	err := session.Where("id = ?", channelID).First(channel).Error
	if err != nil {
		return errors.Wrap(err, "failed to find the specified channel")
	}
	// update the count of block and tx
	channel.BlockCount += 1
	channel.TxCount += int64(block.TxCount)

	return session.Model(&Channel{}).Where("id = ?", channelID).
		Updates(map[string]interface{}{"block_count": channel.BlockCount, "tx_count": channel.TxCount}).Error
}

func GetChannelIDByNetworkAndClient(networkName, clientName string) (int, error) {
	channel := new(Channel)
	err := db.Select("id").Where("network_name = ? and client_name = ?", networkName, clientName).
		First(channel).Error
	if err != nil {
		return 0, err
	}

	return int(channel.ID), nil
}

func GetChannels() ([]Channel, error) {
	var channels []Channel
	if err := db.Find(&channels).Error; err != nil {
		return nil, err
	}
	return channels, nil
}

func GetChannelById(channelId int) (*Channel, error) {
	channel := new(Channel)
	err := db.Where("id = ?", channelId).First(channel).Error
	if err != nil {
		logger.WithField("error", err).Error("failed to find the specified channel by channel id")
		return nil, errors.Wrap(err, "failed to find the specified channel by channel id")
	}

	return channel, nil
}
