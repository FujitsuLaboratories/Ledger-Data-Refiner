/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package router

import (
	"fmt"
	"github.com/FujitsuLaboratories/ledgerdata-refiner/api"
	_ "github.com/FujitsuLaboratories/ledgerdata-refiner/docs"
	"github.com/FujitsuLaboratories/ledgerdata-refiner/log"
	"github.com/FujitsuLaboratories/ledgerdata-refiner/utils"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var logger = log.Logger

func InitRouter() error {
	// default debug mode
	gin.SetMode(utils.AppMode)
	r := gin.New()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("web/*")

	// middleware
	//cors,logger
	r.Use(corsMiddleware(), ginLogger(), gin.Recovery())

	// page apis
	{
		r.GET("/", api.ToIndex)
		r.GET("/dashboard.html", api.ToDashBoard)
		r.GET("/blocks.html", api.ToBlockPage)
		r.GET("/transaction.html", api.ToTransactionPage)
		r.GET("/advance.html", api.ToAdvancePage)
		r.GET("/history.html", api.ToHistoryPage)
		r.GET("/error4.html.html", api.ToError404Page)
		r.GET("/error5.html.html", api.ToError503Page)
	}
	globalGroup := r.Group("/api")
	{
		globalGroup.GET("getparams/:channel_id", api.GetParams)
	}
	// channel apis
	channelApis := globalGroup.Group("/channel")
	{
		channelApis.GET("getchannels", api.GetChannels)
		channelApis.GET("getchannelbyid/:channel_id", api.GetChannelById)
	}

	// block apis
	blockApis := globalGroup.Group("/block")
	{
		blockApis.GET("getblockbydatahash/:channel_id/:data_hash", api.GetBlockByDataHash)
		blockApis.GET("getblockbyblocknum/:channel_id/:block_num", api.GetBlockByBlockNum)
		blockApis.POST("getblocksbytimerange", api.GetBlockByTimeRange)
	}

	// schema apis
	schemaApis := globalGroup.Group("/schema")
	{
		schemaApis.GET("getchaincodename/:channel_id", api.GetChaincodeName)
		schemaApis.GET("getschemas/:channel_id/:chaincode_name", api.GetSchemas)
	}

	// tx apis
	txApis := globalGroup.Group("/tx")
	{
		txApis.GET("getcreatemspfromtxs/:channel_id", api.GetCreateMSPFromTxs)
		txApis.GET("gettxbytxhash/:channel_id/:tx_hash", api.GetTxByTxHash)
		txApis.GET("gettxcountbyorg", api.GetTxCountByOrg)
		txApis.POST("gettxbyblocknum", api.GetTxsByBlockNum)
		txApis.POST("gettxsbytimerange", api.GetTxsByTimeRange)
		txApis.POST("gettxcountbyweek", api.GetTxCountByWeek)
		txApis.POST("gettxcountbyday", api.GetTxCountByDay)
		txApis.POST("gettxcountbymonth", api.GetTxCountByMonth)
	}

	// history apis
	historyApis := globalGroup.Group("/history")
	{
		historyApis.POST("gethistorybykey", api.GetHistoryByKey)
		historyApis.POST("gethistorybytimerange", api.GetHistoryByTimerRange)
	}

	// document apis
	documentApis := globalGroup.Group("/document")
	{
		documentApis.POST("advancedquery", api.AdvanceQuery)
	}

	nodeApis := globalGroup.Group("/node")
	{
		nodeApis.GET("getnodes/:channel_id", api.GetNodes)
	}

	// add swagger api docs
	url := ginSwagger.URL(fmt.Sprintf("http://127.0.0.1%s/swagger/doc.json", utils.HttpPort))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	if err := r.Run(utils.HttpPort); err != nil {
		return err
	}

	return nil
}
