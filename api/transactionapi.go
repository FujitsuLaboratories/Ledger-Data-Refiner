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

// @Description get create MSP list from txs
// @Tags Transaction
// @Accept  json
// @Produce  json
// @Param channel_id path int true "channel id"
// @Success 200 {object} MspSuccessResponse {"code":200, "msg:"OK", "data":[msp1,msp2,msp3]}
// @Failure 500 {object} FailedResponse {"code":3001, "msg:"failed to get create msp list"}
// @Router /tx/getcreatemspfromtxs/{channel_id} [get]
func GetCreateMSPFromTxs(ctx *gin.Context) {
	channelId, _ := strconv.Atoi(ctx.Param("channel_id"))

	createMspFromTxs, err := model.GetCreateMspFromTxs(channelId)
	if err != nil {
		returnMsg(ctx, errmsg.ErrorGetCreateMsp, nil)
		return
	}

	returnMsg(ctx, errmsg.Success, createMspFromTxs)
}

// @Description get txs by block number
// @Tags Transaction
// @Accept  json
// @Produce  json
// @Param channel_id path int true "channel id"
// @Param page_num path int false "page num"
// @Param page_count path int false "page count"
// @Param sort_by path string false "sort by(DESC ASC)"
// @Param block_num path int true "block number"
// @Success 200 {object} TxsSuccessResponse {"code":200, "msg:"OK", "data":[tx1,tx2,tx3]}
// @Failure 500 {object} FailedResponse {"code":3001, "msg:"failed to get txs"}
// @Router /tx/gettxbyblocknum [post]
func GetTxsByBlockNum(ctx *gin.Context) {
	var search model.TransactionSearch
	if err := ctx.ShouldBind(&search); err != nil {
		logger.WithField("error", err).Error("failed to unmarshal params form body")
		returnMsg(ctx, errmsg.ErrorParamType, nil)
		return
	}

	if search.ChannelId <= 0 || search.BlockNum < 0 {
		returnMsg(ctx, errmsg.ErrorNumberOfParams, nil)
		return
	}

	count, err := model.GetTxCountByBlockNum(search.ChannelId, search.BlockNum)
	if err != nil {
		returnMsg(ctx, errmsg.ErrorGetTxByBlockNum, nil)
		return
	}

	var transactions []model.Transaction

	if count > 0 {
		transactions, err = model.GetTxsByBlockNum(search)
		if err != nil {
			returnMsg(ctx, errmsg.ErrorGetTxByBlockNum, nil)
			return
		}
	}

	data := map[string]interface{}{
		"count": count,
		"txs":   transactions,
	}

	returnMsg(ctx, errmsg.Success, data)
}

// @Description get tx by tx hash
// @Tags Transaction
// @Accept  json
// @Produce  json
// @Param channel_id path int true "channel id"
// @Param tx_hash path string true "tx hash"
// @Success 200 {object} TxSuccessResponse {"code":200, "msg:"OK", "data":{tx}}
// @Failure 500 {object} FailedResponse {"code":3002, "msg:"failed to get tx"}
// @Router /tx/gettxbytxhash/{channel_id}/{tx_hash} [get]
func GetTxByTxHash(ctx *gin.Context) {
	channelId, err := strconv.Atoi(ctx.Param("channel_id"))
	if err != nil {
		returnMsg(ctx, errmsg.ErrorParamType, nil)
		return
	}

	txHash := ctx.Param("tx_hash")
	if txHash == "" {
		returnMsg(ctx, errmsg.ErrorNumberOfParams, nil)
		return
	}

	transaction, err := model.GetTxByTxHash(channelId, txHash)
	if err != nil {
		returnMsg(ctx, errmsg.ErrorGetTxByTxHash, nil)
		return
	}

	returnMsg(ctx, errmsg.Success, transaction)
}

// @Description get transactions by time range
// @Tags Transaction
// @Accept  json
// @Produce  json
// @Param channel_id body int true "channel id"
// @Param sort_by body string false "sort by"
// @Param start_time body string false "start time"
// @Param end_time body string false "end time"
// @Param page_num body int false "page number"
// @Param page_count body int false "page count"
// @Param create_msp body string false "create msp"
// @Success 200 {object} TxsSuccessResponse {"code":200, "msg:"OK", "data":{tx1, tx2}}
// @Failure 500 {object} FailedResponse {"code":3002, "msg:"failed to get tx"}
// @Router /tx/gettxsbytimerange [post]
func GetTxsByTimeRange(ctx *gin.Context) {
	var search model.TransactionSearch
	if err := ctx.ShouldBind(&search); err != nil {
		logger.WithField("error", err).Error("failed to unmarshal params form body")
		returnMsg(ctx, errmsg.ErrorParamType, nil)
		return
	}

	count, err := model.GetTxCountByTimeRange(search)
	if err != nil {
		returnMsg(ctx, errmsg.ErrorGetTxsByTimeRange, nil)
		return
	}

	var txs []model.Transaction

	if count > 0 {
		txs, err = model.GetTxsByTimeRange(search)
		if err != nil {
			returnMsg(ctx, errmsg.ErrorGetTxsByTimeRange, nil)
			return
		}
	}

	data := map[string]interface{}{
		"count": count,
		"txs":   txs,
	}

	returnMsg(ctx, errmsg.Success, data)
}

// @Description get count of txs of orgs
// @Tags Transaction
// @Accept  json
// @Produce  json
// @Success 200 {object} TxsSuccessResponse {"code":200, "msg:"OK", "data":{tx1, tx2}}
// @Failure 500 {object} FailedResponse {"code":3002, "msg:"failed to get tx"}
// @Router /tx/gettxcountbyorg [get]
func GetTxCountByOrg(ctx *gin.Context) {

	countTxsByOrg, err := model.CountTxsByOrg()
	if err != nil {
		returnMsg(ctx, errmsg.ErrorGetCountByOrg, nil)
		return
	}

	returnMsg(ctx, errmsg.Success, countTxsByOrg)
}

// @Description get count of transactions of week
// @Tags Transaction
// @Accept  json
// @Produce  json
// @Param channel_id body int true "channel id"
// @Param start_time body string true "start time"
// @Param end_time body string true "end time"
// @Success 200 {object} TxsWeekSuccessResponse {"code":200, "msg:"OK", "data":{tx1, tx2}}
// @Failure 500 {object} FailedResponse {"code":3002, "msg:"failed to get tx"}
// @Router /tx/gettxcountbyweek [post]
func GetTxCountByWeek(ctx *gin.Context) {
	var search model.TransactionSearch
	if err := ctx.ShouldBind(&search); err != nil {
		logger.WithField("error", err).Error("failed to unmarshal params form body")
		returnMsg(ctx, errmsg.ErrorParamType, nil)
		return
	}

	if search.StartTime == "" || search.EndTime == "" ||
		search.ChannelId == 0 {
		returnMsg(ctx, errmsg.ErrorNumberOfParams, nil)
		return
	}

	countTxsWeek, err := model.CountTxsByWeek(search)
	if err != nil {
		returnMsg(ctx, errmsg.ErrorGetCountByWeek, nil)
		return
	}

	returnMsg(ctx, errmsg.Success, countTxsWeek)
}

// @Description get count of transactions of day
// @Tags Transaction
// @Accept  json
// @Produce  json
// @Param channel_id body int true "channel id"
// @Param start_time body string true "start time"
// @Success 200 {object} TxsWeekSuccessResponse {"code":200, "msg:"OK", "data":{tx1, tx2}}
// @Failure 500 {object} FailedResponse {"code":3002, "msg:"failed to get tx"}
// @Router /tx/gettxcountbyday [post]
func GetTxCountByDay(ctx *gin.Context) {
	var search model.TransactionSearch
	if err := ctx.ShouldBind(&search); err != nil {
		logger.WithField("error", err).Error("failed to unmarshal params form body")
		returnMsg(ctx, errmsg.ErrorParamType, nil)
		return
	}

	if search.StartTime == "" || search.ChannelId == 0 {
		returnMsg(ctx, errmsg.ErrorNumberOfParams, nil)
		return
	}

	countTxsDay, err := model.CountTxsByDay(search)
	if err != nil {
		returnMsg(ctx, errmsg.ErrorGetCountByDay, nil)
		return
	}

	returnMsg(ctx, errmsg.Success, countTxsDay)
}

// @Description get count of transactions of month
// @Tags Transaction
// @Accept  json
// @Produce  json
// @Param channel_id body int true "channel id"
// @Param start_time body string true "start time"
// @Param end_time body string true "end time"
// @Success 200 {object} TxsWeekSuccessResponse {"code":200, "msg:"OK", "data":{tx1, tx2}}
// @Failure 500 {object} FailedResponse {"code":3002, "msg:"failed to get tx"}
// @Router /tx/gettxcountbymonth [post]
func GetTxCountByMonth(ctx *gin.Context) {
	var search model.TransactionSearch
	if err := ctx.ShouldBind(&search); err != nil {
		logger.WithField("error", err).Error("failed to unmarshal params form body")
		returnMsg(ctx, errmsg.ErrorParamType, nil)
		return
	}

	if search.StartTime == "" || search.EndTime == "" ||
		search.ChannelId == 0 {
		returnMsg(ctx, errmsg.ErrorNumberOfParams, nil)
		return
	}

	countTxsDay, err := model.CountTxsByMonth(search)
	if err != nil {
		returnMsg(ctx, errmsg.ErrorGetCountByMonth, nil)
		return
	}

	returnMsg(ctx, errmsg.Success, countTxsDay)
}

// --------------------------------------------------------
// response example (for swagger api doc)
type TxsSuccessResponse struct {
	Code string              `json:"code"`
	Msg  string              `json:"msg"`
	Data []model.Transaction `json:"data"`
}

type TxsWeekSuccessResponse struct {
	Code string  `json:"code"`
	Msg  string  `json:"msg"`
	Data []int64 `json:"data"`
}

type MspSuccessResponse struct {
	Code string   `json:"code"`
	Msg  string   `json:"msg"`
	Data []string `json:"data"`
}

type TxSuccessResponse struct {
	Code string            `json:"code"`
	Msg  string            `json:"msg"`
	Data model.Transaction `json:"data"`
}
