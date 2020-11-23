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

// @Description get channel list
// @Tags Channel
// @Accept  json
// @Produce  json
// @Success 200 {object} ChannelsSuccessResponse {"code":200, "msg:"OK", "data":[channel1,channel2]}
// @Failure 500 {object} FailedResponse {"code":1002, "msg:"failed to get channel list"}
// @Router /channel/getchannels [get]
func GetChannels(ctx *gin.Context) {
	channels, err := model.GetChannels()
	if err != nil {
		returnMsg(ctx, errmsg.ErrorGetChannels, nil)
		return
	}
	returnMsg(ctx, errmsg.Success, channels)
}

// @Description get channel by channel ID
// @Tags Channel
// @Accept  json
// @Produce  json
// @Param channel_id path int true "channel id"
// @Success 200 {object} ChannelSuccessResponse {"code":200, "msg:"OK", "data":{channel}}
// @Failure 500 {object} FailedResponse {"code":1001, "msg:"failed to get channel"}
// @Router /channel/getchannelbyid/{channel_id} [get]
func GetChannelById(ctx *gin.Context) {
	param := ctx.Param("channel_id")
	channelId, _ := strconv.Atoi(param)
	logger.WithField("channel_id", channelId).Debug("get channel by channel id")
	channel, err := model.GetChannelById(channelId)
	if err != nil {
		returnMsg(ctx, errmsg.ErrorGetChannel, nil)
		return
	}

	returnMsg(ctx, errmsg.Success, channel)
}

//------------------------------------------
// response examples (for swagger api doc)
type ChannelsSuccessResponse struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data []model.Channel `json:"data"`
}

type ChannelSuccessResponse struct {
	Code string        `json:"code"`
	Msg  string        `json:"msg"`
	Data model.Channel `json:"data"`
}
