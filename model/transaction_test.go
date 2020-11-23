/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package model

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func initDB(t *testing.T) {
	err := InitDB()
	require.Nil(t, err)
}

func TestInsertTx(t *testing.T) {
	initDB(t)
	tx := &Transaction{
		ChannelId:                   1,
		BlockNum:                    1,
		TxNum:                       6,
		TxHash:                      "sdfsad989vxcvx612xzcv",
		ChaincodeName:               "basic",
		TxType:                      "good",
		TxFilter:                    0,
		TxResponse:                  0,
		ReadSet:                     `{"a":100}`,
		WriteSet:                    `{"b":50}`,
		ReadKeyList:                 []string{"a", "b", "c"},
		WriteKeyList:                []string{"b", "c", "d"},
		ChaincodeFunction:           "a",
		ChaincodeFunctionParameters: "a,b,c",
		CreateMspId:                 "org1msp",
		EndorserMspId:               "org1msp",
		EndorserSignature:           "saf234kln523lkntmnvsoidug902q8052",
	}

	err := InsertTx(nil, tx)
	require.Nil(t, err)
}

func TestCountTxsByOrg(t *testing.T) {
	initDB(t)
	org, err := CountTxsByOrg()
	require.Nil(t, err)

	t.Log(org)
}

func TestCountTxsWeek(t *testing.T) {
	initDB(t)

	search := TransactionSearch{
		ChannelId: 1,
		StartTime: "2020-09-20",
		EndTime:   "2020-09-27",
	}

	countTxsWeek, err := CountTxsByWeek(search)
	require.Nil(t, err)

	t.Log(countTxsWeek)
}

func TestCountTxsDay(t *testing.T) {
	initDB(t)

	search := TransactionSearch{
		ChannelId: 1,
		StartTime: "2020-09-21",
	}

	countTxsDay, err := CountTxsByDay(search)
	require.Nil(t, err)

	t.Log(countTxsDay)
}

func TestCountTxsMonth(t *testing.T) {
	initDB(t)

	search := TransactionSearch{
		ChannelId: 1,
		StartTime: "2020-09-01",
		EndTime:   "2020-09-30",
	}

	countTxsDay, err := CountTxsByMonth(search)
	require.Nil(t, err)

	t.Log(countTxsDay)
}
