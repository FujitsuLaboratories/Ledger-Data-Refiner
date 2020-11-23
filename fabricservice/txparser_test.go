/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package fabricservice

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestTxParser(t *testing.T) {
	fabClient := NewFabClient("Org1", "Admin", "User1", "mychannel")

	err := fabClient.Initialize()
	require.Nil(t, err)

	block, err := fabClient.ledger.QueryBlock(7)
	require.Nil(t, err)
	require.NotNil(t, block)

	txBytes := block.Data.Data[0]

	txHash, err := GetTxHash(txBytes)
	require.Nil(t, err)
	t.Logf("TX HASH: %s\n", txHash)

	channelId, err := GetChannelId(txBytes)
	require.Nil(t, err)
	t.Logf("CHANNEL ID: %s\n", channelId)

	txType, err := GetTxType(txBytes)
	require.Nil(t, err)
	t.Logf("TX TYPE: %s\n", txType)

	creatTime, err := GetTxCreateTime(txBytes)
	require.Nil(t, err)
	t.Logf("CREATE TIME: %v\n", creatTime)

	chaincodeName, err := GetChaincodeName(txBytes)
	require.Nil(t, err)
	t.Logf("CHAINCODE NAME: %s\n", chaincodeName)

	status, err := GetResponseStatus(txBytes)
	require.Nil(t, err)
	t.Logf("TX RESPONSE: %d\n", status)

	mspId, err := GetCreatorMSPId(txBytes)
	require.Nil(t, err)
	t.Logf("MSP ID: %s\n", mspId)

	endorserMSPId, err := GetEndorserMSPId(txBytes)
	require.Nil(t, err)
	t.Logf("ENDORSER MSP ID: %s\n", endorserMSPId)

	readSet, err := GetReadSet(txBytes)
	require.Nil(t, err)
	t.Logf("READ SET: %v\n", readSet)

	keyList, err := GetReadKeyList(txBytes)
	require.Nil(t, err)
	t.Logf("READ KEY LIST: %v\n", keyList)

	writeSet, err := GetWriteSet(txBytes)
	require.Nil(t, err)
	t.Logf("WRITE SET: %v\n", writeSet)

	writeKeyList, err := GetWriteKeyList(txBytes)
	require.Nil(t, err)
	t.Logf("WRITE KEY LIST: %v\n", writeKeyList)

	signature, err := GetEndorserSignature(txBytes)
	require.Nil(t, err)
	t.Logf("ENDORSER SIGNATURE: %v\n", signature)

	chaincodeFunction, err := GetChaincodeFunction(txBytes)
	require.Nil(t, err)
	t.Logf("CHAINCODE FUNCTION: %v\n", chaincodeFunction)

	functionParameters, err := GetFunctionParameters(txBytes)
	require.Nil(t, err)
	t.Logf("FUNCTION PARAMETERS: %v\n", functionParameters)

	fabClient.Teardown()
}

func TestGetFunctionParameters(t *testing.T) {
	fabClient := NewFabClient("Org1", "Admin", "User1", "mychannel")

	err := fabClient.Initialize()
	require.Nil(t, err)

	block, err := fabClient.ledger.QueryBlock(3)
	require.Nil(t, err)
	require.NotNil(t, block)

	txBytes := block.Data.Data[0]
	parameters, err := GetFunctionParameters(txBytes)
	require.Nil(t, err)
	str := strings.Join(parameters, ",")
	srcRunes := []rune(str)
	dstRunes := make([]rune, 0, len(srcRunes))
	// remove useless characters
	for _, c := range srcRunes {
		if c >= 0 && c <= 31 {
			continue
		}
		if c == 127 {
			continue
		}
		dstRunes = append(dstRunes, c)
	}

	result := string(dstRunes)
	t.Log(result)
	fabClient.Teardown()
}
