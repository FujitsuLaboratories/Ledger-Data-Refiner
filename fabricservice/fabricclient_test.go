/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package fabricservice

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFabClient_All(t *testing.T) {
	fabClient := NewFabClient("Org1", "Admin", "User1", "mychannel")

	err := fabClient.Initialize()
	require.Nil(t, err)

	// init again, it should be no error
	err = fabClient.Initialize()
	require.Nil(t, err)

	require.NotNil(t, fabClient.SDK())

	// close the connection
	fabClient.Teardown()
}

func TestFabClient_QueryBlocks(t *testing.T) {

	fabClient := NewFabClient("Org1", "Admin", "User1", "mychannel")

	err := fabClient.Initialize()
	require.Nil(t, err)

	// query block by block number
	blockByNumber, err := fabClient.ledger.QueryBlock(1)
	require.Nil(t, err)
	require.NotNil(t, blockByNumber)
	require.Equal(t, blockByNumber.Header.Number, uint64(2))

	// query block by block hash
	blockByHash, err := fabClient.ledger.QueryBlockByHash(blockByNumber.Header.PreviousHash)
	require.Nil(t, err)
	require.NotNil(t, blockByHash)

	fabClient.Teardown()
}
