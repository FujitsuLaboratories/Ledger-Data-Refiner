/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package model

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCountNodes(t *testing.T) {
	initDB(t)
	nodes, err := CountNodes(1)
	require.Nil(t, err)

	t.Log(nodes)
}

func TestGetNodes(t *testing.T) {
	initDB(t)
	nodes, err := GetNodes(1)
	require.Nil(t, err)

	t.Log(nodes)
}

func TestInsertNodes(t *testing.T) {
	initDB(t)
	node := Node{
		Name:        "aaa",
		Url:         "localhost:7051",
		MSP:         "Org1MSP",
		ChannelName: "testchannel",
	}

	err := InsertNodes(nil, []Node{node})
	require.Nil(t, err)
}
