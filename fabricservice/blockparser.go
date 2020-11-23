/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package fabricservice

import (
	"encoding/hex"
	"github.com/FujitsuLaboratories/ledgerdata-refiner/fabricservice/utils"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/pkg/errors"
	"time"
)

func GetBlockNumber(block *common.Block) uint64 {
	return block.Header.Number
}

func GetDataHash(block *common.Block) string {
	return hex.EncodeToString(block.Header.DataHash)
}

func GetPreviousHash(block *common.Block) string {
	return hex.EncodeToString(block.Header.PreviousHash)
}

func GetTxCount(block *common.Block) int {
	return len(block.Data.Data)
}

func GetChannelVersion(config fab.ChannelCfg) uint64 {
	return config.Versions().Channel.Version
}

func GetBlockCreateTime(block *common.Block) (time.Time, error) {
	// The timestamp of the first transaction can be
	// considered as the creation time of the block
	firstTx := block.Data.Data[0]
	channelHeader, err := utils.GetChannelHeader(firstTx)
	if err != nil {
		return time.Time{}, err
	}

	timestamp := channelHeader.Timestamp
	return time.Unix(timestamp.Seconds, 0), nil
}

func GenerateBlockHash(block *common.Block) (string, error) {
	headerBytes, err := utils.HeaderBytes(block.Header)
	if err != nil {
		return "", err
	}

	sha256Bytes := utils.ComputeSHA256(headerBytes)
	return hex.EncodeToString(sha256Bytes), nil
}

func GetInvalidTxCount(lc *ledger.Client, block *common.Block) (uint64, error) {
	var count uint64

	for _, txBytes := range block.Data.Data {
		channelHeader, err := utils.GetChannelHeader(txBytes)
		if err != nil {
			return 0, err
		}

		txId := fab.TransactionID(channelHeader.TxId)
		tx, err := lc.QueryTransaction(txId)
		if err != nil {
			return 0, errors.Wrap(err, "error getting transaction")
		}

		if tx.ValidationCode != 0 {
			count++
		}
	}
	return count, nil
}

func GetTx(block *common.Block, i uint64) []byte {
	return block.Data.Data[i]
}

func GetTxFilter(block *common.Block, i uint64) int {
	return int(block.Metadata.Metadata[2][i])
}
