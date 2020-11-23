/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package model

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGetMaxBlockHeight(t *testing.T) {
	initDB(t)
	height := GetMaxBlockHeight(1)
	t.Log(height)
}

func TestGetBlockCountByTimeRange(t *testing.T) {
	initDB(t)
	s := BlockSearch{ChannelId: 1, StartTime: "2020-09-21 15:34:06", EndTime: time.Now().String()}
	count, err := GetBlockCountByTimeRange(s)
	require.Nil(t, err)

	t.Log(count)
}

func TestGetBlocksByTimeRange(t *testing.T) {
	initDB(t)
	s := BlockSearch{ChannelId: 1}
	timeRange, err := GetBlocksByTimeRange(s)
	require.Nil(t, err)

	t.Log(len(timeRange))
}
