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

type Schema struct {
	ID            int       `json:"id" gorm:"primarykey"`
	ChannelId     int       `json:"channel_id" gorm:"type:integer"`
	ChaincodeName string    `json:"chaincode_name" gorm:"type:varchar(256)"`
	SchemaArray   string    `json:"schema_array" gorm:"type:json"`
	SchemaJSON    string    `json:"schema_json" gorm:"type:json"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (s *Schema) TableName() string {
	return "schema"
}

func GetSchemas() (schemas []Schema, err error) {
	err = db.Find(&schemas).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.Wrap(err, "failed to get schemas from database")
	}
	return
}

func InsertSchema(session *gorm.DB, schema *Schema) error {
	if session == nil {
		session = db
	}
	logger.Debug("insert one schema")
	return session.Create(schema).Error
}

func UpdateSchema(session *gorm.DB, schemaId int, schemaInArray, schemaJSON string) error {
	if session == nil {
		session = db
	}
	return session.Model(&Schema{}).Where("id = ?", schemaId).
		Updates(Schema{SchemaArray: schemaInArray, SchemaJSON: schemaJSON}).Error
}

func GetChaincodeNameByChannelId(channelId int) ([]string, error) {
	var chaincodeNames []string
	err := db.Table("schema").Select("chaincode_name").Where("channel_id = ?", channelId).Group("chaincode_name").Find(&chaincodeNames).Error
	if err != nil {
		logger.WithField("error", err).Error("failed to get chaincode name by channelId")
		return nil, err
	}

	return chaincodeNames, nil
}

func CountChaincodeByChannelId(channelId int) (int64, error) {
	var count int64
	err := db.Table("schema").Select("count(distinct(chaincode_name))").
		Where("channel_id = ?", channelId).
		Find(&count).Error
	if err != nil {
		logger.WithField("error", err).Error("failed to get count of chaincode")
		return 0, err
	}

	return count, nil
}

func GetSchemasByChaincodeAndChannelId(channelId int, chaincodeName string) ([]Schema, error) {
	var schemas []Schema
	err := db.Where("chaincode_name = ? and channel_id = ?", chaincodeName, channelId).Find(&schemas).Error
	if err != nil {
		logger.WithField("error", err).Error("failed to get schemas by chaincode and channelId")
		return nil, err
	}

	return schemas, nil
}
