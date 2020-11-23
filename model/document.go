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

type Document struct {
	ID            int    `json:"id" gorm:"primarykey"`
	ChannelId     int    `json:"channel_id" gorm:"type:integer;default null"`
	ChaincodeName string `json:"chaincode_name" gorm:"type:varchar(256);default null"`
	SchemaId      int    `json:"schema_id" gorm:"type:integer;default null"`
	Key           string `json:"key" gorm:"type:varchar(256);default null"`
	Content       string `json:"content" gorm:"type:json;default null"`
	IsDelete      bool   `json:"is_delete" gorm:"type:boolean;"`
	//Version       *Version `gorm:"embedded"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Version struct {
	BlockNum uint64 `json:"block_num" gorm:"type:integer;default null"`
	TxNum    uint64 `json:"tx_num" gorm:"integer;default null"`
}

func (d *Document) TableName() string {
	return "document"
}

type DocumentSearch struct {
	ChannelId  int      `json:"channel_id"`
	SchemaId   int      `json:"schema_id"`
	Selects    []string `json:"selects"`
	Conditions []string `json:"conditions"`
	PageNum    int      `json:"page_num"`
	PageCount  int      `json:"page_count"`
}

func InsertDocument(doc *Document) error {
	logger.Debug("insert one document")
	return db.Create(doc).Error
}

func UpdateDocSchema(session *gorm.DB, schemaId, docId int) error {
	if session == nil {
		session = db
	}
	return session.Model(&Document{}).Where("id = ?", docId).Updates(Document{SchemaId: schemaId}).Error
}

func UpdateDocStatus(session *gorm.DB, doc *Document) error {
	if session == nil {
		session = db
	}
	var oldDoc Document
	err := session.Where("key = ? and channel_id = ? and chaincode_name = ?", doc.Key, doc.ChannelId, doc.ChaincodeName).First(&oldDoc).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.Wrap(err, "failed to fetch the specified document")
	} else if err == gorm.ErrRecordNotFound {
		err = session.Create(doc).Error
		if err != nil {
			return errors.Wrap(err, "failed to insert document into the database")
		}
	} else {
		err = db.Model(&Document{}).Where("id = ?", oldDoc.ID).Updates(Document{Content: doc.Content,
			IsDelete: doc.IsDelete}).Error
		if err != nil {
			return errors.Wrap(err, "failed to update document")
		}
	}

	return nil
}

func GetSchemaUndeterminedDocs() ([]Document, error) {
	document := make([]Document, 0, 100)
	err := db.Where("schema_id = 0").Limit(100).Find(&document).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.Wrap(err, "failed to get undetermined documents")
	}

	return document, nil
}

func AdvancedQuery(search DocumentSearch) ([]map[string]interface{}, error) {
	tx := db.Table("document").
		Where("channel_id = ? and schema_id = ?", search.ChannelId, search.SchemaId)
	selectCondition := ""
	if len(search.Selects) > 0 {
		firstSelect := true
		for _, s := range search.Selects {
			if firstSelect {
				selectCondition += "content::json#>>" + s + ","
				firstSelect = false
			} else {
				selectCondition += " content::json#>>" + s + ","
			}
		}
	}

	tx = tx.Select(selectCondition + " key as Doc_Key, content as Doc_Value, channel_id, chaincode_name, is_delete")

	if len(search.Conditions) > 0 {
		for _, c := range search.Conditions {
			tx = tx.Where("content::json#>>" + c)
		}
	}

	if search.PageCount <= 0 {
		search.PageCount = 15
	}

	if search.PageNum <= 0 {
		search.PageNum = 1
	}

	var results = make([]map[string]interface{}, 0)
	err := tx.Find(&results).Error
	if err != nil {
		logger.WithField("error", err).Error("failed to doing advanced query")
		return nil, err
	}

	return results, nil
}

func CountOfAdvancedQuery(channelId, schemaId int, conditions []string) (int64, error) {
	tx := db.Table("document").Select("count(id)").
		Where("channel_id = ? and schema_id = ?", channelId, schemaId)

	if len(conditions) > 0 {
		for _, c := range conditions {
			tx = tx.Where("content::json#>>" + c)
		}
	}

	var count int64
	err := tx.Find(&count).Error
	if err != nil {
		logger.WithField("error", err).Error("failed to get count of advanced query")
		return 0, err
	}

	return count, nil
}
