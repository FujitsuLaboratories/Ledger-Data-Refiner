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

// @Description get chaincode by channelId
// @Tags Schema
// @Accept  json
// @Produce  json
// @Param channel_id path int true "channel id"
// @Success 200 {object} ChaincodeNamesSuccessResponse {"code":200, "msg:"OK", "data":{"name1","name2"}}
// @Failure 500 {object} FailedResponse {"code":1001, "msg:"failed to get chaincodes"}
// @Router /schema/getchaincodename/{channel_id} [get]
func GetChaincodeName(ctx *gin.Context) {
	channelId, _ := strconv.Atoi(ctx.Param("channel_id"))
	names, err := model.GetChaincodeNameByChannelId(channelId)
	if err != nil {
		// TODO
		return
	}

	returnMsg(ctx, 200, names)
}

// @Description get schemas by channelId and chaincode
// @Tags Schema
// @Accept  json
// @Produce  json
// @Param channel_id path int true "channel id"
// @Param chaincode_name path int true "chaincode name"
// @Success 200 {object} SchemasSuccessResponse {"code":200, "msg:"OK", "data":{schema1, schema2}}
// @Failure 500 {object} FailedResponse {"code":1001, "msg:"failed to get channel"}
// @Router /schema/getschemas/{channel_id}/{chaincode_name} [get]
func GetSchemas(ctx *gin.Context) {
	channelId, _ := strconv.Atoi(ctx.Param("channel_id"))
	chaincodeName := ctx.Param("chaincode_name")
	if chaincodeName == "" {
		returnMsg(ctx, errmsg.ErrorNumberOfParams, nil)
		return
	}

	schemas, err := model.GetSchemasByChaincodeAndChannelId(channelId, chaincodeName)
	if err != nil {
		returnMsg(ctx, errmsg.ErrorGetSchema, nil)
		return
	}

	returnMsg(ctx, errmsg.Success, schemas)
}

//------------------------------------------
// response examples (for swagger api doc)
type ChaincodeNamesSuccessResponse struct {
	Code string   `json:"code"`
	Msg  string   `json:"msg"`
	Data []string `json:"data"`
}

type SchemasSuccessResponse struct {
	Code string         `json:"code"`
	Msg  string         `json:"msg"`
	Data []model.Schema `json:"data"`
}
