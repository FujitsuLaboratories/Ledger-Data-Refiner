/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package model

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type Block struct {
	ID             int       `json:"id" gorm:"primarykey"`
	ChannelId      int       `json:"channel_id" gorm:"type:integer"`
	BlockNum       uint64    `json:"block_num" gorm:"type:integer"`
	DataHash       string    `json:"data_hash" gorm:"type:varchar(256)"`
	BlockHash      string    `json:"block_hash" gorm:"type:varchar(256)"`
	PreviousHash   string    `json:"previous_hash" gorm:"type:varchar(256)"`
	TxCount        int       `json:"tx_count" gorm:"type:integer"`
	InvalidTxCount uint64    `json:"invalid_tx_count" gorm:"type:integer"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (b *Block) TableName() string {
	return "block"
}

type BlockSearch struct {
	ChannelId int    `json:"channel_id"`
	PageNum   int    `json:"page_num"`
	PageCount int    `json:"page_count"`
	SortBy    string `json:"sort_by"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// InsertBlock inserts a block into the database.
// If session is not nil, you should commit the session by yourself (session.Commit())
func InsertBlock(session *gorm.DB, block *Block) error {
	if session == nil {
		session = db
	}
	logger.WithField("block_id", block.BlockNum).Debug("insert one block")
	return session.Create(block).Error
}

func GetMaxBlockHeight(channelId int) uint64 {
	var block Block
	err := db.Select("block_num").Order("block_num desc").Where("channel_id = ?", channelId).Find(&block).Error
	if err != nil {
		logger.WithField("error", err).Debug("failed to get max block height")
		return 0
	}
	return block.BlockNum
}

func GetBlockByDataHash(channelId int, dataHash string) (*Block, error) {
	block := new(Block)

	err := db.Where("channel_id = ? and data_hash = ?", channelId, dataHash).First(block).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find the specified block by channel id and data hash")
	}

	return block, nil
}

func GetBlockByBlockNumber(channelId int, blockNumber int64) (*Block, error) {
	block := new(Block)

	err := db.Where("channel_id = ? and block_num = ?", channelId, blockNumber).First(block).Error
	if err != nil {
		logger.WithField("error", err).Error("failed to find the specified block by channel id and block number")
		return nil, errors.Wrap(err, "failed to find the specified block by channel id and block number")
	}

	return block, nil
}

func GetBlocksByTimeRange(search BlockSearch) ([]Block, error) {
	var blocks []Block
	condition := ""

	if search.StartTime != "" {
		condition += fmt.Sprintf(" and created_at >= '%s' ", search.StartTime)
	}
	if search.EndTime != "" {
		condition += fmt.Sprintf(" and created_at <= '%s' ", search.EndTime)
	}
	if search.PageNum == 0 {
		search.PageNum = 1
	}
	if search.PageCount == 0 {
		search.PageCount = 15
	}
	if search.SortBy == "" {
		search.SortBy = "desc"
	}

	err := db.Select("block_num, data_hash, block_hash, previous_hash, tx_count, invalid_tx_count, created_at").
		Where("channel_id = ? "+condition, search.ChannelId).Order("block_num " + search.SortBy).Limit(search.PageCount).
		Offset((search.PageNum - 1) * search.PageCount).Find(&blocks).Error
	if err != nil {
		logger.WithField("error", err).Error("faild to get blocks by time range")
		return nil, err
	}

	return blocks, nil
}

func GetBlockCountByTimeRange(search BlockSearch) (int64, error) {
	var count int64
	tx := db.Table("block").Select("count(id)").Where("channel_id = ?", search.ChannelId)

	if search.StartTime != "" {
		tx = tx.Where("created_at >= ?", search.StartTime)
	}
	if search.EndTime != "" {
		tx = tx.Where("created_at <= ?", search.EndTime)
	}
	err := tx.Find(&count).Error
	if err != nil {
		logger.WithField("error", err).Error("failed to get the number of blocks")
		return 0, err
	}
	return count, nil
}
