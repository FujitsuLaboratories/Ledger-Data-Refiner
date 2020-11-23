/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package api

import (
	"github.com/FujitsuLaboratories/ledgerdata-refiner/model"
	"github.com/FujitsuLaboratories/ledgerdata-refiner/utils/errmsg"
	"github.com/gin-gonic/gin"
)

// @Description get documents history by doc key
// @Tags Document history
// @Accept  json
// @Produce  json
// @Param channel_id body int true "channel id"
// @Param sort_by body string false "sort by"
// @Param doc_key body string true "document key"
// @Param create_msp body string false "create msp"
// @Param chaincode_name body string false "chaincode name"
// @Param page_num body int false "page number"
// @Param page_count body int false "page count"
// @Success 200 {object} HistoriesSuccessResponse {"code":200, "msg:"OK", "data":[dh1,dh2]}
// @Failure 500 {object} FailedResponse {"code":1002, "msg:"failed to get document history"}
// @Router /history/gethistorybykey [post]
func GetHistoryByKey(ctx *gin.Context) {
	var historySearch model.HistorySearch
	if err := ctx.ShouldBind(&historySearch); err != nil {
		returnMsg(ctx, errmsg.ErrorParamType, nil)
		return
	}

	count, err := model.GetHistoryCountByKey(historySearch)
	if err != nil {
		returnMsg(ctx, errmsg.ErrorGetHistoryByKey, nil)
		return
	}

	var historyByKey []model.HistoryWithTransaction
	if count != 0 {
		historyByKey, err = model.GetHistoryByKey(historySearch)
		if err != nil {
			returnMsg(ctx, errmsg.ErrorGetHistoryByKey, nil)
			return
		}
	}

	data := map[string]interface{}{
		"count":   count,
		"history": historyByKey,
	}

	returnMsg(ctx, errmsg.Success, data)
}

// @Description get documents history by time range
// @Tags Document history
// @Accept  json
// @Produce  json
// @Param channel_id body int true "channel id"
// @Param sort_by body string false "sort by"
// @Param start_time body string false "start time"
// @Param end_time body string false "end time"
// @Param page_num body int false "page number"
// @Param page_count body int false "page count"
// @Success 200 {object} HistoriesSuccessResponse {"code":200, "msg:"OK", "data":[dh1,dh2]}
// @Failure 500 {object} FailedResponse {"code":1002, "msg:"failed to get document history"}
// @Router /history/gethistorybytimerange [post]
func GetHistoryByTimerRange(ctx *gin.Context) {
	var historySearch model.HistorySearch
	if err := ctx.ShouldBind(&historySearch); err != nil {
		returnMsg(ctx, errmsg.ErrorParamType, nil)
		return
	}

	count, err := model.GetHistoryCountByTimeRange(historySearch)
	if err != nil {
		returnMsg(ctx, errmsg.ErrorGetHistoryByTimeRange, nil)
		return
	}

	var historyWithTransactions []model.HistoryWithTransaction

	if count != 0 {
		historyWithTransactions, err = model.GetHistoryByTimeRange(historySearch)
		if err != nil {
			returnMsg(ctx, errmsg.ErrorGetHistoryByTimeRange, nil)
			return
		}
	}

	data := map[string]interface{}{
		"count":   count,
		"history": historyWithTransactions,
	}

	returnMsg(ctx, errmsg.Success, data)
}

//------------------------------------------
// response examples (for swagger api doc)
type HistoriesSuccessResponse struct {
	Code string                         `json:"code"`
	Msg  string                         `json:"msg"`
	Data []model.HistoryWithTransaction `json:"data"`
}
