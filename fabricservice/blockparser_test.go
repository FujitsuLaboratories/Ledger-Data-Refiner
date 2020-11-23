/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package fabricservice

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBlockParser(t *testing.T) {

	fabClient := NewFabClient("Org1", "Admin", "User1", "mychannel")

	err := fabClient.Initialize()
	require.Nil(t, err)

	block, err := fabClient.ledger.QueryBlock(6)
	require.Nil(t, err)
	require.NotNil(t, block)

	number := GetBlockNumber(block)
	require.Equal(t, number, uint64(6))

	datahash := GetDataHash(block)
	require.NotNil(t, datahash)

	createTime, err := GetBlockCreateTime(block)
	require.Nil(t, err)
	t.Logf("block create time: %s\n", createTime)

	config, err := fabClient.ledger.QueryConfig()
	require.Nil(t, err)
	version := GetChannelVersion(config)
	t.Logf("channel version: %d\n", version)

	blockHash, err := GenerateBlockHash(block)
	require.Nil(t, err)
	block7, _ := fabClient.ledger.QueryBlock(7)
	require.Equal(t, blockHash, hex.EncodeToString(block7.Header.PreviousHash))

	invalidTxCount, err := GetInvalidTxCount(fabClient.ledger, block)
	require.Nil(t, err)
	require.NotEqual(t, 0, invalidTxCount)

	txCount := GetTxCount(block)
	t.Logf("the number of block txs: %d\n", txCount)

	filter := GetTxFilter(block, 0)
	t.Logf("tx filter: %d\n", filter)

	tx := GetTx(block, 0)
	t.Log(tx)

	preHash := GetPreviousHash(block)
	t.Logf("pre hash: %s\n", preHash)

	fabClient.Teardown()
}
