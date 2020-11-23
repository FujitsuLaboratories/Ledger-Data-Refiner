/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

const apiPort =  "http://"+ window.location.host + "/api/";
//channel url
const channelIdUrl = apiPort + "channel/getchannels";

//chaincodeName  url
const chaincodeNameUrl = apiPort + "schema/getchaincodename/"

//dashboard url
const txStatisticUrl = apiPort + "getparams/"
const lineChartUrl = apiPort + "tx/gettxcountbyday"
const lineChartUrl2 = apiPort + "tx/gettxcountbyweek"
const lineChartUrl3 = apiPort + "tx/gettxcountbymonth"
const pieChartUrl = apiPort + "tx/gettxcountbyorg"
const networkUrl = apiPort + "node/getnodes/"

// creator url
const transactionCreateUrl = apiPort + "tx/getcreatemspfromtxs/"

// block url
const blockTimeUrl = apiPort + "block/getblocksbytimerange";
const blockNumberUrl = apiPort + "block/getblockbyblocknum/"
const blockDataHashUrl = apiPort + "block/getblockbydatahash/"

// transtion url
const transactionTimeUrl = apiPort + "tx/gettxsbytimerange"
const transactionBlockNumberUrl = apiPort + "tx/gettxbyblocknum"
const transactionHashUrl = apiPort + "tx/gettxbytxhash/"
const transactionHashDetailsUrl = apiPort + "tx/gettxbytxhash/"

// history url
const historyTimeUrl = apiPort + "history/gethistorybytimerange"
const historyKeyUrl = apiPort + "history/gethistorybykey"

// advanced url
const advanceSchemaUrl = apiPort + "schema/getschemas/"
const advanceTableUrl = apiPort + "document/advancedquery"

// conflict url
const conflictFunctionUrl = apiPort + "/getconflictfunction"
const conflictSourceUrl = apiPort + "/getsourcetxs"
const conflictTransactionUrl = apiPort + "/getconflicttxs"
