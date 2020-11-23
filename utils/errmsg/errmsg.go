/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package errmsg

const (
	Success = 200
	Error   = 500

	// common errors
	ErrorNumberOfParams = 501
	ErrorParamType      = 502

	// code = 1000... channel errors
	ErrorGetChannel  = 1001
	ErrorGetChannels = 1002

	// code = 2000... block errors
	ErrorGetBlockByDataHash  = 2001
	ErrorGetBlockByBlockNum  = 2002
	ErrorGetBlockByTimeRange = 2003

	// code = 3000... transaction errors
	ErrorGetCreateMsp      = 3001
	ErrorGetTxByTxHash     = 3002
	ErrorGetTxByBlockNum   = 3003
	ErrorGetTxsByTimeRange = 3004
	ErrorGetCountByOrg     = 3005
	ErrorGetCountByWeek    = 3006
	ErrorGetCountByDay     = 3007
	ErrorGetCountByMonth   = 3008

	// code = 4000... schema errors
	ErrorGetSchema = 4001

	// code = 5000... conflict result errors
	ErrorGetConflictFunction = 5001
	ErrorGetSourceTxs        = 5002

	// code = 6000... document history errors
	ErrorGetHistoryByTimeRange = 6001
	ErrorGetHistoryByKey       = 6002
)

var codeMsg = map[int]string{
	Success:                    "OK",
	Error:                      "FAIL",
	ErrorNumberOfParams:        "missing parameters",
	ErrorParamType:             "param type mismatch",
	ErrorGetChannel:            "failed to get channel",
	ErrorGetChannels:           "failed to get channel list",
	ErrorGetBlockByDataHash:    "failed to get block by data hash",
	ErrorGetBlockByBlockNum:    "failed to get block by block number",
	ErrorGetBlockByTimeRange:   "failed to get blocks by time range",
	ErrorGetCreateMsp:          "failed to get create msp from txs",
	ErrorGetTxByTxHash:         "failed to get tx by tx hash",
	ErrorGetTxByBlockNum:       "failed to get txs by block number",
	ErrorGetTxsByTimeRange:     "failed to get txs by time range",
	ErrorGetSchema:             "failed to get schemas",
	ErrorGetConflictFunction:   "failed to get conflict functions",
	ErrorGetSourceTxs:          "failed to get source txs",
	ErrorGetHistoryByTimeRange: "failed to get history by time range",
	ErrorGetHistoryByKey:       "failed to get history by key",
	ErrorGetCountByOrg:         "failed to get count of txs group by orgs",
	ErrorGetCountByWeek:        "failed to get count of txs group by week",
	ErrorGetCountByDay:         "failed to get count of txs group by day",
	ErrorGetCountByMonth:       "failed to get count of txs group by month",
}

func GetErrMsg(code int) string {
	return codeMsg[code]
}
