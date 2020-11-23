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

func AdvanceQuery(ctx *gin.Context) {
	var documentSearch model.DocumentSearch
	if err := ctx.ShouldBind(&documentSearch); err != nil {
		logger.Error(err)
		returnMsg(ctx, errmsg.ErrorParamType, nil)
		return
	}

	count, err := model.CountOfAdvancedQuery(documentSearch.ChannelId, documentSearch.SchemaId, documentSearch.Conditions)
	if err != nil {
		returnMsg(ctx, errmsg.Error, nil)
		return
	}

	var result []map[string]interface{}

	if count > 0 {
		result, err = model.AdvancedQuery(documentSearch)
		if err != nil {
			returnMsg(ctx, errmsg.Error, nil)
			return
		}
	}

	data := map[string]interface{}{
		"count":  count,
		"result": result,
	}

	returnMsg(ctx, errmsg.Success, data)
}
