/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package model

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateDocumentHistory(t *testing.T) {
	initDB(t)
	dh := &DocumentHistory{
		TxId:     1,
		Key:      "sdf",
		Content:  `{"a":12,"b":"bb"}`,
		IsDelete: false,
	}

	err := InsertDocumentHistory(nil, dh)
	require.Nil(t, err)
}

func TestGetHistoryByTimeRange(t *testing.T) {
	initDB(t)
	s := HistorySearch{ChannelId: 1}
	transactions, err := GetHistoryByTimeRange(s)
	require.Nil(t, err)

	t.Log(len(transactions))
}

func TestGetHistoryByKey(t *testing.T) {
	initDB(t)
	s := HistorySearch{ChannelId: 1, DocKey: "asset6"}
	count, err := GetHistoryByKey(s)
	require.Nil(t, err)

	t.Log(count)
}

func TestGetHistoryCountByKey(t *testing.T) {
	initDB(t)
	s := HistorySearch{ChannelId: 1, DocKey: "asset6"}
	count, err := GetHistoryCountByKey(s)
	require.Nil(t, err)

	t.Log(count)
}
