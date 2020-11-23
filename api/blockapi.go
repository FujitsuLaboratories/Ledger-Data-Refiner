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

// @Description get block by data hash
// @Tags Block
// @Accept  json
// @Produce  json
// @Param channel_id path int true "channel id"
// @Param data_hash path string true "data hash"
// @Success 200 {object} BlockSuccessResponse {"code":200, "msg:"OK", "data":{block}}
// @Failure 500 {object} FailedResponse {"code":2001, "msg:"failed to get block"}
// @Router /block/getblockbydatahash/{channel_id}/{data_hash} [get]
func GetBlockByDataHash(ctx *gin.Context) {
	channelIdStr := ctx.Param("channel_id")
	dataHash := ctx.Param("data_hash")
	if channelIdStr == "" || dataHash == "" {
		returnMsg(ctx, errmsg.ErrorNumberOfParams, nil)
		return
	}
	channelId, _ := strconv.Atoi(channelIdStr)

	block, err := model.GetBlockByDataHash(channelId, dataHash)
	if err != nil {
		returnMsg(ctx, errmsg.ErrorGetBlockByDataHash, nil)
		return
	}

	returnMsg(ctx, 200, block)
}

// @Description get block by block number
// @Tags Block
// @Accept  json
// @Produce  json
// @Param channel_id path int true "channel id"
// @Param block_num path int true "block number"
// @Success 200 {object} BlockSuccessResponse {"code":200, "msg:"OK", "data":{block}}
// @Failure 500 {object} FailedResponse {"code":2002, "msg:"failed to get block"}
// @Router /block/getblockbyblocknum/{channel_id}/{block_num} [get]
func GetBlockByBlockNum(ctx *gin.Context) {
	channelIdStr := ctx.Param("channel_id")
	blockNumStr := ctx.Param("block_num")
	if channelIdStr == "" || blockNumStr == "" {
		returnMsg(ctx, errmsg.ErrorNumberOfParams, nil)
		return
	}
	channelId, _ := strconv.Atoi(channelIdStr)
	blockNum, err := strconv.ParseInt(blockNumStr, 10, 64)
	if err != nil {
		returnMsg(ctx, errmsg.ErrorParamType, nil)
		return
	}

	block, err := model.GetBlockByBlockNumber(channelId, blockNum)
	if err != nil {
		returnMsg(ctx, errmsg.ErrorGetBlockByBlockNum, nil)
		return
	}

	returnMsg(ctx, 200, block)
}

// @Description get blocks by time range
// @Tags Block
// @Accept  json
// @Produce  json
// @Param channel_id body int true "channel id"
// @Param sort_by body string false "sort by"
// @Param start_time body string false "start time"
// @Param end_time body string false "end time"
// @Param page_num body int false "page number"
// @Param page_count body int false "page count"
// @Success 200 {object} BlocksSuccessResponse {"code":200, "msg:"OK", "data":{block}}
// @Failure 500 {object} FailedResponse {"code":2001, "msg:"failed to get block"}
// @Router /block/getblocksbytimerange [post]
func GetBlockByTimeRange(ctx *gin.Context) {
	var search model.BlockSearch
	if err := ctx.ShouldBind(&search); err != nil {
		logger.WithField("error", err).Error("failed to unmarshal params form body")
		returnMsg(ctx, errmsg.ErrorParamType, nil)
		return
	}

	blockCount, err := model.GetBlockCountByTimeRange(search)
	if err != nil {
		returnMsg(ctx, errmsg.ErrorParamType, nil)
		return
	}

	var blocks []model.Block

	if blockCount > 0 {
		blocks, err = model.GetBlocksByTimeRange(search)
		if err != nil {
			returnMsg(ctx, errmsg.ErrorParamType, nil)
			return
		}
	}

	data := map[string]interface{}{
		"count":  blockCount,
		"blocks": blocks,
	}
	returnMsg(ctx, errmsg.Success, data)
}

// --------------------------------------------------------
// response example (for swagger api doc)
type BlocksSuccessResponse struct {
	Code string        `json:"code"`
	Msg  string        `json:"msg"`
	Data []model.Block `json:"data"`
}

type BlockSuccessResponse struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data model.Block `json:"data"`
}
