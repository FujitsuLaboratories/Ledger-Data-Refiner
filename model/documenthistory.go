/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package model

import (
	"gorm.io/gorm"
	"time"
)

type DocumentHistory struct {
	ID        int       `json:"id" gorm:"primarykey"`
	TxId      int       `json:"tx_id" gorm:"type:integer"`
	Key       string    `json:"key" gorm:"type:varchar(256)"`
	Content   string    `json:"content" gorm:"type:json"`
	IsDelete  bool      `json:"is_delete" gorm:"type:boolean"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type HistorySearch struct {
	ChannelId     int    `json:"channel_id"`
	DocKey        string `json:"doc_key"`
	CreateMSP     string `json:"create_msp"`
	ChaincodeName string `json:"chaincode_name"`
	PageNum       int    `json:"page_num"`
	PageCount     int    `json:"page_count"`
	SortBy        string `json:"sort_by"`
	StartTime     string `json:"start_time"`
	EndTime       string `json:"end_time"`
}

type HistoryWithTransaction struct {
	DocumentHistory
	CreatedAt time.Time `json:"created_at"`
	Transaction
}

func (dh *DocumentHistory) TableName() string {
	return "document_history"
}

// InsertDocumentHistory inserts one document history into database
func InsertDocumentHistory(session *gorm.DB, dh *DocumentHistory) error {
	if session == nil {
		session = db
	}
	return session.Create(dh).Error
}

func GetHistoryByKey(search HistorySearch) ([]HistoryWithTransaction, error) {
	var hwt []HistoryWithTransaction

	if search.PageCount == 0 {
		search.PageCount = 15
	}
	if search.PageNum == 0 {
		search.PageNum = 1
	}

	tx := db.Table("transaction as tran").Select("tran.*, hist.tx_id, hist.key, hist.content, hist.is_delete").
		Joins("join document_history as hist on hist.tx_id = tran.id").
		Where("tran.channel_id = ? and hist.key = ?", search.ChannelId, search.DocKey).
		Limit(search.PageCount).Offset((search.PageNum - 1) * search.PageCount)

	if search.CreateMSP != "" {
		tx = tx.Where("tran.create_msp_id = ?", search.CreateMSP)
	}

	if search.ChaincodeName != "" {
		tx = tx.Where("tran.chaincode_name = ?", search.ChaincodeName)
	}

	if search.SortBy == "" {
		tx = tx.Order("tran.block_num DESC")
	} else {
		tx = tx.Order("tran.block_num " + search.SortBy)
	}

	err := tx.Find(&hwt).Error
	if err != nil {
		logger.WithField("error", err).Error("failed to get history by key")
		return nil, err
	}

	return hwt, nil
}

func GetHistoryCountByKey(search HistorySearch) (int64, error) {
	var count int64

	if search.PageCount == 0 {
		search.PageCount = 15
	}
	if search.PageNum == 0 {
		search.PageNum = 1
	}

	tx := db.Table("transaction as tran").Select("count(tran.id)").
		Joins("join document_history as hist on hist.tx_id = tran.id").
		Where("tran.channel_id = ? and hist.key = ?", search.ChannelId, search.DocKey)

	if search.CreateMSP != "" {
		tx = tx.Where("tran.create_msp_id = ?", search.CreateMSP)
	}

	if search.ChaincodeName != "" {
		tx = tx.Where("tran.chaincode_name = ?", search.ChaincodeName)
	}

	err := tx.Debug().Find(&count).Error
	if err != nil {
		logger.WithField("error", err).Error("failed to get the count of history by key")
		return 0, err
	}

	return count, nil
}

func GetHistoryByTimeRange(search HistorySearch) ([]HistoryWithTransaction, error) {

	var hwt []HistoryWithTransaction

	if search.PageCount == 0 {
		search.PageCount = 15
	}
	if search.PageNum == 0 {
		search.PageNum = 1
	}

	tx := db.Table("transaction as tran").Select("tran.*, hist.tx_id, hist.key, hist.content, hist.is_delete").
		Joins("join document_history as hist on hist.tx_id = tran.id").Where("tran.channel_id = ?", search.ChannelId).
		Limit(search.PageCount).Offset(search.PageCount * (search.PageNum - 1))

	if search.SortBy == "" {
		search.SortBy = " order by tran.block_num DESC, tran.tx_num DESC"
		tx.Order("tran.block_num DESC").Order("tran.tx_num DESC")
	} else {
		tx.Order("tran.block_num " + search.SortBy).Order("tran.tx_num " + search.SortBy)
	}

	if search.StartTime != "" {
		tx = tx.Where("tran.created_at >= ?", search.StartTime)
	}
	if search.EndTime != "" {
		tx = tx.Where("tran.created_at <= ?", search.EndTime)
	}

	err := tx.Scan(&hwt).Error
	if err != nil {
		logger.WithField("error", err).Error("failed to get history by time range")
		return nil, err
	}

	return hwt, nil
}

func GetHistoryCountByTimeRange(search HistorySearch) (int64, error) {
	var count int64

	tx := db.Select("count(tran.id)").Table("transaction as tran").Joins("join  document_history as hist on hist.tx_id = tran.id").
		Where("tran.channel_id = ?", search.ChannelId)

	if search.StartTime != "" {
		tx = tx.Where("tran.created_at >= ?", search.StartTime)
	}
	if search.EndTime != "" {
		tx = tx.Where("tran.created_at <= ?", search.EndTime)
	}

	err := tx.Find(&count).Error
	if err != nil {
		logger.WithField("error", err).Error("failed to get the count of history")
		return 0, err
	}

	return count, nil
}
