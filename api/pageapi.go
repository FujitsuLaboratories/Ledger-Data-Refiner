/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ToIndex(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "dashboard.html", nil)
}

func ToDashBoard(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "dashboard.html", nil)
}

func ToBlockPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "blocks.html", nil)
}

func ToTransactionPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "transaction.html", nil)
}

func ToHistoryPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "history.html", nil)
}

func ToAdvancePage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "advance.html", nil)
}

func ToError404Page(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "error4.html", nil)
}

func ToError503Page(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "error5.html", nil)
}
