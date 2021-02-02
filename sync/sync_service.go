/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"github.com/FujitsuLaboratories/ledgerdata-refiner/fabricservice"
	"github.com/FujitsuLaboratories/ledgerdata-refiner/model"
	"github.com/FujitsuLaboratories/ledgerdata-refiner/utils"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
)

const (
	endorserTransaction = "ENDORSER_TRANSACTION"
	lifeCycleChaincode  = "_lifecycle"
)

func (sync *SyncService) syncChannel() {
	logger.Debug("start syncing channel")
	sync.syncChannelStatus = true
	if sync.firstSync {
		// get the first block of fabric network
		firstBlock, err := sync.fabClient.Ledger().QueryBlock(0)
		if err != nil {
			// if syncing channel failed, panic
			logger.WithField("error", err).Error("failed to sync channel")
			panic(err)
		}

		channelCfg, err := sync.fabClient.Ledger().QueryConfig()
		if err != nil {
			logger.WithField("error", err).Error("failed to fetch channel config")
			panic(err)
		}

		generateBlockHash, _ := fabricservice.GenerateBlockHash(firstBlock)

		channel := &model.Channel{
			NetworkName:      sync.networkName,
			ClientName:       sync.clientName,
			ChannelName:      sync.fabClient.GetChannelId(),
			BlockCount:       0,
			TxCount:          0,
			ChannelVersion:   fabricservice.GetChannelVersion(channelCfg),
			GenesisBlockHash: generateBlockHash,
		}
		logger.WithField("channel", channel).Debug("start inserting channel into the database")
		sync.channelId, err = model.InsertChannel(channel)
		if err != nil {
			logger.WithField("error", err).Error("failed to insert channel into the database")
			panic(err)
		}

		sync.firstSync = false
	}

	logger.Info("start syncing blocks")
	sync.syncBlocks()
	sync.syncChannelStatus = false
}

// if a block failed to sync, all operations of it should be rollback
func (sync *SyncService) syncBlocks() {
	channelInfo, err := sync.fabClient.Ledger().QueryInfo()
	if err != nil {
		logger.Error("failed to fetch channel info when syncing blocks")
		return
	}
	// fetch the latest height of fabric network
	blockHeight := channelInfo.BCI.Height - 1
	logger.Debugf("now the block of fabric is %d", blockHeight)
	syncedMaxBlockHeight := model.GetMaxBlockHeight(sync.channelId)
	logger.Debugf("now the block height of db is %d", syncedMaxBlockHeight)
	if blockHeight != 0 && blockHeight > syncedMaxBlockHeight {
		if syncedMaxBlockHeight != 0 {
			syncedMaxBlockHeight++
		}
		for i := syncedMaxBlockHeight; i <= blockHeight; i++ {
			logger.WithField("block_id", i).Debug("start syncing block")
			// begin session
			session := model.BeginSession()
			// sync blocks
			_ = sync.syncBlock(session, i)
			session.Commit()

			logger.WithField("block_id", i).Debug("finish syncing block")
		}
		logger.Info("finish syncing blocks")
	} else {
		// nothing to sync
		logger.WithField("synced_block_height", syncedMaxBlockHeight).Info("no blocks need to be synchronized")
	}
}

func (sync *SyncService) syncBlock(session *gorm.DB, blockNum uint64) error {
	block, err := sync.fabClient.Ledger().QueryBlock(blockNum)
	if err != nil {
		sync.handleBlockError(session, blockNum, "failed to fetch the specified block", err)
		return err
	}

	blockHash, _ := fabricservice.GenerateBlockHash(block)
	invalidTxCount, _ := fabricservice.GetInvalidTxCount(sync.fabClient.Ledger(), block)
	createTime, _ := fabricservice.GetBlockCreateTime(block)

	newBlock := &model.Block{
		ChannelId:      sync.channelId,
		BlockNum:       blockNum,
		DataHash:       fabricservice.GetDataHash(block),
		BlockHash:      blockHash,
		PreviousHash:   fabricservice.GetPreviousHash(block),
		TxCount:        fabricservice.GetTxCount(block),
		InvalidTxCount: invalidTxCount,
		CreatedAt:      createTime,
	}

	logger.Debug("start syncing txs in this block")
	err = sync.syncTxs(session, block, blockNum)
	if err != nil {
		sync.handleBlockError(session, blockNum, "failed to insert block into the database", err)
		return err
	}
	logger.Debug("finish syncing txs in this block")
	err = model.InsertBlock(session, newBlock)
	if err != nil {
		sync.handleBlockError(session, blockNum, "failed to insert block into the database", err)
		return err
	}
	logger.WithField("block", newBlock).Debug("finish inserting block into the database")
	// update channel after syncing the block
	err = model.UpdateChannel(session, newBlock, sync.channelId)
	if err != nil {
		sync.handleBlockError(session, blockNum, "failed to update channel", err)
		return err
	}
	return nil
}

func (sync *SyncService) syncTxs(session *gorm.DB, block *common.Block, blockNum uint64) error {
	txCount := fabricservice.GetTxCount(block)
	logger.Debugf("there are %d txs in this block", txCount)

	if txCount > 0 {
		for i := 0; i < txCount; i++ {
			logger.Debugf("start sync tx %d in this block", i)
			err := sync.syncTx(session, block, blockNum, uint64(i))
			if err != nil {
				logger.Debug("failed to insert tx into database")
				return err
			}
			logger.Debugf("start sync tx %d in this block", i)
		}
	}
	return nil
}

func (sync *SyncService) syncTx(session *gorm.DB, block *common.Block, blockNum, txNum uint64) error {
	tx := fabricservice.GetTx(block, txNum)

	chaincodeName, _ := fabricservice.GetChaincodeName(tx)
	txHash, _ := fabricservice.GetTxHash(tx)
	txType, _ := fabricservice.GetTxType(tx)
	txFilter := fabricservice.GetTxFilter(block, txNum)
	responseStatus, _ := fabricservice.GetResponseStatus(tx)
	readSet, _ := fabricservice.GetReadSet(tx)
	writeSet, _ := fabricservice.GetWriteSet(tx)
	readKeyList, _ := fabricservice.GetReadKeyList(tx)
	writeKeyList, _ := fabricservice.GetWriteKeyList(tx)
	chaincodeFunction, _ := fabricservice.GetChaincodeFunction(tx)
	functionParameters, _ := fabricservice.GetFunctionParameters(tx)
	creatorMSPId, _ := fabricservice.GetCreatorMSPId(tx)
	endorserMSPId, _ := fabricservice.GetEndorserMSPId(tx)
	endorserSignature, _ := fabricservice.GetEndorserSignature(tx)
	creatTime, _ := fabricservice.GetTxCreateTime(tx)

	txRow := &model.Transaction{
		ChannelId:                   sync.channelId,
		BlockNum:                    blockNum,
		TxNum:                       txNum,
		TxHash:                      txHash,
		ChaincodeName:               chaincodeName,
		TxType:                      txType,
		TxFilter:                    txFilter,
		TxResponse:                  responseStatus,
		ReadSet:                     utils.ToJson(readSet),
		WriteSet:                    utils.ToJson(writeSet),
		ReadKeyList:                 readKeyList,
		WriteKeyList:                writeKeyList,
		ChaincodeFunction:           chaincodeFunction,
		ChaincodeFunctionParameters: strings.Join(functionParameters, ","),
		CreateMspId:                 creatorMSPId,
		EndorserSignature:           utils.ToJson(endorserSignature),
		CreatedAt:                   creatTime,
	}

	if endorserMSPId != nil {
		txRow.EndorserMspId = strings.Join(endorserMSPId, ",")
	}

	// insert tx into the database
	err := model.InsertTx(session, txRow)
	if err != nil {
		return errors.Wrap(err, "failed to insert tx into the database")
	}
	logger.Debugf("sync tx, the content is %v", txRow)
	logger.Debug("start syncing docs in this tx")
	// handle documents of this tx
	if txRow.TxType == endorserTransaction && txRow.ChaincodeName != lifeCycleChaincode {
		err = sync.syncDocs(session, txRow, writeSet)
		if err != nil {
			return err
		}
	}
	logger.Debug("finish syncing docs in this tx")

	return nil
}

func (sync *SyncService) syncDocs(session *gorm.DB, tx *model.Transaction, writeSet []map[string]interface{}) error {
	if len(writeSet) > 0 {
		for _, m := range writeSet {
			writes := m["set"].([]*fabricservice.IKVWrite)
			if len(writes) > 0 {
				for _, write := range writes {
					logger.WithField("block_id", tx.BlockNum).Debug("start syncing document")
					err := sync.syncDocHistory(session, tx.ID, write)
					if err != nil {
						return err
					}
					err = sync.syncDocStatus(session, tx.ChaincodeName, write)
					if err != nil {
						return err
					}
					logger.WithField("block_id", tx.BlockNum).Debug("finish syncing document")
				}
			}
		}
	}
	return nil
}

func (sync *SyncService) syncDocHistory(session *gorm.DB, txId int, write *fabricservice.IKVWrite) error {
	docRow := &model.DocumentHistory{
		TxId:     txId,
		Key:      write.Key,
		Content:  write.Value,
		IsDelete: write.IsDelete,
	}
	err := model.InsertDocumentHistory(session, docRow)
	if err != nil {
		return errors.Wrap(err, "failed to insert document history into the database")
	}
	logger.WithField("document_history", docRow).Debug("insert document_history table")
	return nil
}

func (sync *SyncService) syncDocStatus(session *gorm.DB, chaincodeName string, write *fabricservice.IKVWrite) error {
	docRow := &model.Document{
		ChannelId:     sync.channelId,
		ChaincodeName: chaincodeName,
		Key:           write.Key,
		Content:       write.Value,
		IsDelete:      write.IsDelete,
	}

	if err := model.UpdateDocStatus(session, docRow); err != nil {
		return err
	}

	return nil
}

func (sync *SyncService) handleBlockError(session *gorm.DB, blockNum uint64, reason string, err error) {
	logger.WithFields(logrus.Fields{
		"block_id": blockNum,
		"error":    err,
	}).Error(reason)
	logger.WithField("block_id", blockNum).Debug("mark failed block")
	sync.markLostHeight(blockNum)
	session.Rollback()
}
