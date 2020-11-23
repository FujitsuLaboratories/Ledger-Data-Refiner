/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package api

import (
	"github.com/FujitsuLaboratories/ledgerdata-refiner/model"
	"github.com/FujitsuLaboratories/ledgerdata-refiner/utils/errmsg"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetNodes(ctx *gin.Context) {
	param := ctx.Param("channel_id")
	channelId, err := strconv.Atoi(param)
	if err != nil {
		returnMsg(ctx, errmsg.ErrorParamType, nil)
		return
	}

	nodes, err := model.GetNodes(channelId)
	if err != nil {
		returnMsg(ctx, errmsg.Error, nil)
		return
	}

	returnMsg(ctx, errmsg.Success, nodes)
}

func GetParams(ctx *gin.Context) {
	param := ctx.Param("channel_id")
	channelId, err := strconv.Atoi(param)
	if err != nil {
		returnMsg(ctx, errmsg.ErrorParamType, nil)
		return
	}
	//if error not nil, the default value of count is 0
	txSearch := model.TransactionSearch{ChannelId: channelId}
	txCount, _ := model.GetTxCountByTimeRange(txSearch)

	blockSearch := model.BlockSearch{ChannelId: channelId}
	blockCount, _ := model.GetBlockCountByTimeRange(blockSearch)

	nodeCount, _ := model.CountNodes(channelId)
	chaincodeCount, _ := model.CountChaincodeByChannelId(channelId)

	returnMsg(ctx, errmsg.Success, gin.H{
		"tx_num":        txCount,
		"block_num":     blockCount,
		"node_num":      nodeCount,
		"chaincode_num": chaincodeCount,
	})
}
