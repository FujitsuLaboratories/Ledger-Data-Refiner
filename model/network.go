/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package model

import (
	"gorm.io/gorm"
)

type Node struct {
	ID          int    `json:"id" gorm:"primarykey"`
	Name        string `json:"name" gorm:"type:varchar(256);default null;"`
	Url         string `json:"url" gorm:"type:varchar(256);default null;"`
	MSP         string `json:"msp" gorm:"type:varchar(256);default null;"`
	ChannelName string `json:"channel_name" gorm:"type:varchar(256);default null;"`
}

func (n *Node) TableName() string {
	return "network"
}

func InsertNodes(session *gorm.DB, nodes []Node) error {
	if session == nil {
		session = db
	}

	for _, node := range nodes {
		err := session.Exec("INSERT INTO network(name, url, msp, channel_name) select ?,?,?,? WHERE NOT EXISTS (SELECT * FROM network WHERE name=? and url = ? and msp = ?)",
			node.Name, node.Url, node.MSP, node.ChannelName, node.Name, node.Url, node.MSP).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func CountNodes(channelId int) (int64, error) {
	var count int64
	err := db.Table("network").Select("count(id)").Where("channel_name = (select channel_name from channel where id = ?)", channelId).
		Find(&count).Error
	if err != nil {
		logger.WithField("error", err).Error("failed to get count of network nodes")
		return 0, err
	}

	return count, nil
}

func GetNodes(channelId int) ([]Node, error) {
	var nodes []Node
	err := db.Table("network").Select("name, url, msp").Where("channel_name = (select channel_name from channel where id = ?)", channelId).
		Find(&nodes).Error
	if err != nil {
		logger.WithField("error", err).Error("failed to get network nodes")
		return nil, err
	}

	return nodes, nil
}
