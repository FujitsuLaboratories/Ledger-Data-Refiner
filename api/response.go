/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package api

import (
	"github.com/FujitsuLaboratories/ledgerdata-refiner/log"
	"github.com/FujitsuLaboratories/ledgerdata-refiner/utils/errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
)

var logger = log.Logger

// Response defines a common response for web request
type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func returnMsg(ctx *gin.Context, code int, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Code: code,
		Data: data,
		Msg:  errmsg.GetErrMsg(code),
	})
}

//-----------------------------------------
// FAILED EXAMPLE (for swagger)
type FailedResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}
