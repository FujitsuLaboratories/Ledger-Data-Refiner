/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package model

import (
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type Transaction struct {
	ID                          int            `json:"id" gorm:"primarykey"`
	ChannelId                   int            `json:"channel_id" gorm:"type:integer"`
	BlockNum                    uint64         `json:"block_num" gorm:"type:integer"`
	TxNum                       uint64         `json:"tx_num" gorm:"type:integer"`
	TxHash                      string         `json:"tx_hash" gorm:"type:varchar(256)"`
	ChaincodeName               string         `json:"chaincode_name" gorm:"type:varchar(256)"`
	TxType                      string         `json:"tx_type" gorm:"type:varchar(256)"`
	TxFilter                    int            `json:"tx_filter" gorm:"type:integer"`
	TxResponse                  int32          `json:"tx_response" gorm:"type:integer"`
	ReadSet                     string         `json:"read_set" gorm:"type:json"`
	WriteSet                    string         `json:"write_set" gorm:"type:json"`
	ReadKeyList                 pq.StringArray `json:"read_key_list" gorm:"type:text[]"`
	WriteKeyList                pq.StringArray `json:"write_key_list" gorm:"type:text[]"`
	ChaincodeFunction           string         `json:"chaincode_function" gorm:"type:varchar(256)"`
	ChaincodeFunctionParameters string         `json:"chaincode_function_parameters" gorm:"type:varchar(256)"`
	CreateMspId                 string         `json:"create_msp_id" gorm:"type:varchar(256)"`
	EndorserMspId               string         `json:"endorser_msp_id" gorm:"type:varchar(256)"`
	EndorserSignature           string         `json:"endorser_signature" gorm:"type:text"`
	CreatedAt                   time.Time      `json:"created_at"`
	UpdatedAt                   time.Time      `json:"updated_at"`
}

func (t *Transaction) TableName() string {
	return "transaction"
}

type TransactionSearch struct {
	ChannelId int    `json:"channel_id"`
	StartTime string `json:"start_time"`
	BlockNum  uint64 `json:"block_num"`
	EndTime   string `json:"end_time"`
	PageNum   int    `json:"page_num"`
	PageCount int    `json:"page_count"`
	SortBy    string `json:"sort_by"`
	CreateMSP string `json:"create_msp"`
}

func InsertTx(session *gorm.DB, tx *Transaction) error {
	logger.Debug("insert one transaction")
	if session == nil {
		session = db
	}
	return session.Create(tx).Error
}

func GetCreateMspFromTxs(channelId int) ([]string, error) {
	var createMSPs []string
	err := db.Model(&Transaction{}).Select("create_msp_id").Where("channel_id = ?", channelId).Group("create_msp_id").Find(&createMSPs).Error
	if err != nil {
		logger.WithField("error", err).Error("failed to get craeteMSP from txs")
		return nil, errors.Wrap(err, "failed to get craeteMSP from txs")
	}

	return createMSPs, nil
}

func GetTxByTxHash(channelId int, txHash string) (*Transaction, error) {
	tx := new(Transaction)

	err := db.Where("channel_id = ? and tx_hash = ?", channelId, txHash).First(tx).Error
	if err != nil {
		logger.WithField("error", err).Error("failed to get tx by tx hash")
		return nil, err
	}

	return tx, nil
}

func GetTxCountByBlockNum(channelId int, blockNum uint64) (int64, error) {
	var count int64
	err := db.Table("transaction").Select("count(id)").
		Where("channel_id = ? and block_num = ?", channelId, blockNum).
		Find(&count).Error
	if err != nil {
		logger.WithField("error", err).Error("failed to get the count of tx by block num")
		return 0, err
	}

	return count, nil
}

func GetTxsByBlockNum(search TransactionSearch) ([]Transaction, error) {
	if search.PageNum == 0 {
		search.PageNum = 1
	}
	if search.PageCount == 0 {
		search.PageCount = 15
	}
	if search.SortBy == "" {
		search.SortBy = "desc"
	}
	var txs []Transaction
	err := db.Where("channel_id = ? and block_num = ?", search.ChannelId, search.BlockNum).Offset((search.PageNum - 1) * search.PageCount).Limit(search.PageCount).
		Order("tx_num " + search.SortBy).Find(&txs).Error
	if err != nil {
		logger.WithField("error", err).Error("failed to get tx by block number")
		return nil, err
	}

	return txs, nil
}

func GetTxsByTimeRange(search TransactionSearch) ([]Transaction, error) {
	var txs []Transaction
	if search.PageNum == 0 {
		search.PageNum = 1
	}
	if search.PageCount == 0 {
		search.PageCount = 15
	}
	if search.SortBy == "" {
		search.SortBy = "desc"
	}

	tx := db.Where("channel_id = ? ", search.ChannelId).Order("block_num " + search.SortBy).Order("tx_num " + search.SortBy).
		Limit(search.PageCount).Offset((search.PageNum - 1) * search.PageCount)

	if search.StartTime != "" {
		tx = tx.Where("created_at >= ?", search.StartTime)
	}
	if search.EndTime != "" {
		tx = tx.Where("created_at <= ?", search.EndTime)
	}
	if search.CreateMSP != "" {
		tx = tx.Where("create_msp_id = ?", search.CreateMSP)
	}

	err := tx.Find(&txs).Error
	if err != nil {
		logger.WithField("error", err).Error("failed to get txs by time range")
		return nil, err
	}

	return txs, nil
}

func GetTxCountByTimeRange(search TransactionSearch) (int64, error) {
	var count int64

	tx := db.Table("transaction").Select("count(id)").Where("channel_id = ?", search.ChannelId)

	if search.StartTime != "" {
		tx = tx.Where("created_at >= ?", search.StartTime)
	}
	if search.EndTime != "" {
		tx = tx.Where("created_at <= ?", search.EndTime)
	}
	if search.CreateMSP != "" {
		tx = tx.Where("create_msp_id = ?", search.CreateMSP)
	}

	err := tx.Find(&count).Error
	if err != nil {
		logger.WithField("error", err).Error("failed to get tx count by time range")
		return 0, err
	}

	return count, nil
}

func CountTxsByOrg() ([]map[string]interface{}, error) {
	datas := make([]map[string]interface{}, 0)

	err := db.Table("transaction").Select("create_msp_id as name, count(id) as value").Group("create_msp_id").Find(&datas).Error
	if err != nil {
		logger.WithField("error", err).Error("failed to get count of txs by organizations")
		return nil, err
	}

	return datas, nil
}

func CountTxsByWeek(search TransactionSearch) ([]int64, error) {
	var week []int64

	err := db.Debug().Raw("SELECT COALESCE(b.c,0) from (select to_char(b, 'YYYY-MM-DD') as time from generate_series(to_timestamp(?, 'YYYY-MM-DD'), to_timestamp(?, 'YYYY-MM-DD'), '1 day') AS b GROUP BY time order by time) as a full outer join (select count(t2.id) as c, to_char(t2.created_at,'YYYY-MM-DD') as time from transaction t2 where t2.channel_id = ? group by time ) as b on a.time=b.time order by a.time",
		search.StartTime, search.EndTime, search.ChannelId, search.EndTime, search.StartTime).Find(&week).Error
	if err != nil {
		logger.WithField("error", err).Error("failed to get count of txs by week")
		return nil, err
	}
	return week, nil
}

func CountTxsByDay(search TransactionSearch) ([]int64, error) {
	var day []int64
	start := search.StartTime + " 00"
	end := search.StartTime + " 23"
	err := db.Debug().Raw("SELECT COALESCE(b.c,0) from (select to_char(b, 'YYYY-MM-DD hh24') as time from generate_series(to_timestamp(?, 'YYYY-MM-DD hh24'), to_timestamp(?, 'YYYY-MM-DD hh24'), '1 hour') AS b GROUP BY time order by time) as a full outer join (select count(t2.id) as c, to_char(t2.created_at,'YYYY-MM-DD hh24') as time from transaction t2 where t2.channel_id = ? group by time ) as b on a.time=b.time where a.time is not null order by a.time",
		start, end, search.ChannelId).Find(&day).Error
	if err != nil {
		logger.WithField("error", err).Error("failed to get count of txs by day")
		return nil, err
	}
	return day, nil
}

func CountTxsByMonth(search TransactionSearch) ([]int64, error) {
	var month []int64
	err := db.Raw("SELECT COALESCE(b.c,0) from (select to_char(b, 'YYYY-MM-DD') as time from generate_series(to_timestamp(?, 'YYYY-MM-DD'), to_timestamp(?, 'YYYY-MM-DD'), '1 day') AS b GROUP BY time order by time) as a full outer join (select count(t2.id) as c, to_char(t2.created_at,'YYYY-MM-DD') as time from transaction t2 where t2.channel_id = ? group by time) as b on a.time=b.time order by a.time\n",
		search.StartTime, search.EndTime, search.ChannelId).Find(&month).Error
	if err != nil {
		logger.WithField("error", err).Error("failed to get count of txs by day")
		return nil, err
	}
	return month, nil
}
