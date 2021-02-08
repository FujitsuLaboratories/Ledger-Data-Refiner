/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package fabricservice

import (
	"encoding/hex"
	"github.com/FujitsuLaboratories/ledgerdata-refiner/fabricservice/utils"
	refinerutil "github.com/FujitsuLaboratories/ledgerdata-refiner/utils"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/ledger/rwset"
	"github.com/hyperledger/fabric-protos-go/ledger/rwset/kvrwset"
	"github.com/hyperledger/fabric-protos-go/msp"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

type IKVRead struct {
	Key     string `json:"key"`
	Version *Version
}

type Version struct {
	BlockNum uint64 `json:"block_num"`
	TxNum    uint64 `json:"tx_num"`
}

type IKVWrite struct {
	Key      string `json:"key"`
	IsDelete bool   `json:"is_delete"`
	Value    string `json:"value"`
}

func transformFbricKVRead(read *kvrwset.KVRead) *IKVRead {
	iKVRead := new(IKVRead)
	iKVRead.Key = refinerutil.RemoveInvalidCharacters(read.Key)
	if read.Version != nil {
		iKVRead.Version = &Version{
			BlockNum: read.Version.BlockNum,
			TxNum:    read.Version.TxNum,
		}
	}

	return iKVRead
}

func transformFabricKVWrite(write *kvrwset.KVWrite) *IKVWrite {
	iKVWrite := new(IKVWrite)
	if write.Value != nil && len(write.Value) > 0 {
		iKVWrite.Value = refinerutil.RemoveInvalidCharacters(string(write.Value))
	}

	iKVWrite.Key = refinerutil.RemoveInvalidCharacters(write.Key)
	iKVWrite.IsDelete = write.IsDelete
	return iKVWrite
}

func GetTxHash(tx []byte) (string, error) {
	channelHeader, err := utils.GetChannelHeader(tx)
	if err != nil {
		return "", errors.Wrap(err, "err getting channel header")
	}

	return channelHeader.TxId, nil
}

func GetChannelId(tx []byte) (string, error) {
	channelHeader, err := utils.GetChannelHeader(tx)
	if err != nil {
		return "", errors.Wrap(err, "error getting channel header")
	}
	return channelHeader.ChannelId, nil
}

func GetTxType(tx []byte) (string, error) {
	channelHeader, err := utils.GetChannelHeader(tx)
	if err != nil {
		return "", errors.Wrap(err, "error getting channel header")
	}
	return common.HeaderType_name[channelHeader.Type], nil
}

func GetTxCreateTime(tx []byte) (time.Time, error) {
	channelHeader, err := utils.GetChannelHeader(tx)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "error getting channel header")
	}
	return time.Unix(channelHeader.Timestamp.Seconds, 0), nil
}

func GetChaincodeName(tx []byte) (string, error) {
	chaincodeAction, err := getChaincodeAction(tx)
	if err != nil {
		return "", errors.Wrap(err, "error getting chaincode action")
	}
	chaincodeName := refinerutil.RemoveInvalidCharacters(chaincodeAction.ChaincodeId.Name)
	return chaincodeName, nil
}

func GetResponseStatus(tx []byte) (int32, error) {
	chaincodeAction, err := getChaincodeAction(tx)
	if err != nil {
		return -1, errors.Wrap(err, "error getting chaincode action")
	}
	return chaincodeAction.Response.Status, nil
}

func GetCreatorMSPId(tx []byte) (string, error) {
	env, err := utils.GetEnvelopeFromBlock(tx)
	if err != nil {
		return "", err
	}
	if env == nil {
		return "", errors.New("nil envelope")
	}
	payload, err := utils.GetPayload(env)
	if err != nil {
		return "", errors.Wrap(err, "error extracting ChannelHeader from payload")
	}
	signatureHeader := &common.SignatureHeader{}
	err = proto.Unmarshal(payload.Header.SignatureHeader, signatureHeader)
	if err != nil {
		return "", errors.Wrap(err, "error extracting signature header")
	}
	mspContent := &msp.SerializedIdentity{}
	err = proto.Unmarshal(signatureHeader.Creator, mspContent)
	if err != nil {
		return "", err
	}
	mspId := refinerutil.RemoveInvalidCharacters(mspContent.Mspid)
	return mspId, nil
}

func GetEndorserMSPId(tx []byte) ([]string, error) {
	chaincodeActionPayload, err := getChaincodeActionPayload(tx)
	if err != nil {
		return nil, err
	}
	var endorserMSPIds []string
	for _, endorsement := range chaincodeActionPayload.Action.Endorsements {
		mspContent := &msp.SerializedIdentity{}
		err = proto.Unmarshal(endorsement.Endorser, mspContent)
		endorserMSPIds = append(endorserMSPIds, refinerutil.RemoveInvalidCharacters(mspContent.Mspid))
	}
	return endorserMSPIds, nil
}

func GetReadSet(tx []byte) ([]map[string]interface{}, error) {
	chaincodeAction, err := getChaincodeAction(tx)
	if err != nil {
		return nil, err
	}
	result := &rwset.TxReadWriteSet{}
	err = proto.Unmarshal(chaincodeAction.Results, result)
	if err != nil {
		return nil, err
	}
	var readSets []map[string]interface{}
	for _, set := range result.NsRwset {
		if set.Namespace != "lscc" && set.Namespace != "_lifecycle" {
			readSet := make(map[string]interface{})
			readSet["namespace"] = set.Namespace
			rwset := &kvrwset.KVRWSet{}
			err := proto.Unmarshal(set.Rwset, rwset)
			if err != nil {
				return nil, err
			}
			var reads []*IKVRead
			for _, read := range rwset.Reads {

				reads = append(reads, transformFbricKVRead(read))
			}
			readSet["set"] = reads
			readSets = append(readSets, readSet)
		}
	}
	return readSets, nil
}

func GetReadKeyList(tx []byte) ([]string, error) {
	chaincodeAction, err := getChaincodeAction(tx)
	if err != nil {
		return nil, err
	}
	result := &rwset.TxReadWriteSet{}
	err = proto.Unmarshal(chaincodeAction.Results, result)
	if err != nil {
		return nil, err
	}
	var readSet []string
	for _, set := range result.NsRwset {
		if set.Namespace != "lscc" && set.Namespace != "_lifecycle" {
			kvrwset := &kvrwset.KVRWSet{}
			err := proto.Unmarshal(set.Rwset, kvrwset)
			if err != nil {
				return nil, err
			}
			for _, read := range kvrwset.Reads {
				if read.Version == nil {
					readSet = append(readSet, refinerutil.RemoveInvalidCharacters(read.Key+"!#null"))
				} else {
					readSet = append(readSet, refinerutil.RemoveInvalidCharacters(read.Key+"!#"+strconv.Itoa(int(read.Version.BlockNum))+"_"+
						strconv.Itoa(int(read.Version.TxNum))))
				}
			}
		}
	}
	return readSet, nil
}

func GetWriteSet(tx []byte) ([]map[string]interface{}, error) {
	chaincodeAction, err := getChaincodeAction(tx)
	if err != nil {
		return nil, err
	}
	result := &rwset.TxReadWriteSet{}
	err = proto.Unmarshal(chaincodeAction.Results, result)
	if err != nil {
		return nil, err
	}
	var writeSets []map[string]interface{}
	for _, set := range result.NsRwset {
		if set.Namespace != "lscc" && set.Namespace != "_lifecycle" {
			writeSet := make(map[string]interface{})
			writeSet["namespace"] = set.Namespace
			rwset := &kvrwset.KVRWSet{}
			err := proto.Unmarshal(set.Rwset, rwset)
			if err != nil {
				return nil, err
			}
			var writes []*IKVWrite
			for _, write := range rwset.Writes {
				iKVWrite := transformFabricKVWrite(write)
				writes = append(writes, iKVWrite)
			}
			writeSet["set"] = writes
			writeSets = append(writeSets, writeSet)
		}
	}
	return writeSets, nil
}

func GetWriteKeyList(tx []byte) ([]string, error) {
	chaincodeAction, err := getChaincodeAction(tx)
	if err != nil {
		return nil, err
	}
	result := &rwset.TxReadWriteSet{}
	err = proto.Unmarshal(chaincodeAction.Results, result)
	if err != nil {
		return nil, err
	}
	var writeSet []string
	for _, set := range result.NsRwset {
		if set.Namespace != "lscc" && set.Namespace != "_lifecycle" {

			kvrwset := &kvrwset.KVRWSet{}
			err := proto.Unmarshal(set.Rwset, kvrwset)
			if err != nil {
				return nil, err
			}
			for _, write := range kvrwset.Writes {
				writeSet = append(writeSet, refinerutil.RemoveInvalidCharacters(write.Key))
			}
		}
	}
	return writeSet, nil
}

func GetEndorserSignature(tx []byte) ([]map[string]string, error) {
	chaincodeActionPayload, err := getChaincodeActionPayload(tx)
	if err != nil {
		return nil, err
	}
	var endorserSignatures []map[string]string
	for _, endorsement := range chaincodeActionPayload.Action.Endorsements {
		endorserSignature := make(map[string]string)
		endorserSignature["signature"] = hex.EncodeToString(endorsement.Signature)
		mspContent := &msp.SerializedIdentity{}
		err = proto.Unmarshal(endorsement.Endorser, mspContent)
		endorserSignature["msp_id"] = refinerutil.RemoveInvalidCharacters(mspContent.Mspid)
		endorserSignature["cerficate"] = refinerutil.RemoveInvalidCharacters(string(mspContent.IdBytes))
		endorserSignatures = append(endorserSignatures, endorserSignature)
	}
	return endorserSignatures, nil
}

func GetChaincodeFunction(tx []byte) (string, error) {
	invokeSpec, err := getTxAllArgs(tx)
	if err != nil {
		return "", err
	}
	if invokeSpec.ChaincodeSpec == nil {
		return "", nil
	}
	chaincodeFunction := refinerutil.RemoveInvalidCharacters(string(invokeSpec.ChaincodeSpec.Input.Args[0]))
	return chaincodeFunction, nil
}

func GetFunctionParameters(tx []byte) ([]string, error) {
	invokeSpec, err := getTxAllArgs(tx)
	if err != nil {
		return nil, err
	}
	if invokeSpec.ChaincodeSpec == nil {
		return nil, nil
	}
	var args []string
	for i := 1; i < len(invokeSpec.ChaincodeSpec.Input.Args); i++ {
		args = append(args, refinerutil.RemoveInvalidCharacters(string(invokeSpec.ChaincodeSpec.Input.Args[i])))
	}
	return args, nil
}

func getTxAllArgs(tx []byte) (*peer.ChaincodeInvocationSpec, error) {
	propPayload := &peer.ChaincodeProposalPayload{}
	chaincodeActionPayload, err := getChaincodeActionPayload(tx)
	if err != nil {
		return nil, err
	}
	if err := proto.Unmarshal(chaincodeActionPayload.ChaincodeProposalPayload, propPayload); err != nil {
		return nil, errors.Wrap(err, "error extracting ChannelHeader from payload")
	}
	invokeSpec := &peer.ChaincodeInvocationSpec{}
	err = proto.Unmarshal(propPayload.Input, invokeSpec)
	if err != nil {
		return nil, errors.Wrap(err, "error extracting ChannelHeader from payload")
	}
	return invokeSpec, nil
}

func getChaincodeActionPayload(tx []byte) (*peer.ChaincodeActionPayload, error) {
	env, err := utils.GetEnvelopeFromBlock(tx)
	if err != nil {
		return nil, err
	}
	if env == nil {
		return nil, errors.New("nil envelope")
	}
	payload, err := utils.GetPayload(env)
	if err != nil {
		return nil, errors.Wrap(err, "error extracting ChannelHeader from payload")
	}
	transaction, err := utils.GetTransaction(payload.Data)
	if err != nil {
		return nil, errors.Wrap(err, "error getting transaction")
	}
	chaincodeActionPayload, err := utils.GetChaincodeActionPayload(transaction.Actions[0].Payload)
	if err != nil {
		return nil, errors.Wrap(err, "error getting chaincodeActionPayload")
	}
	return chaincodeActionPayload, err
}

func getChaincodeAction(tx []byte) (*peer.ChaincodeAction, error) {
	chaincodeActionPayload, err := getChaincodeActionPayload(tx)
	if err != nil {
		return nil, err
	}
	proposalResponsePayload := &peer.ProposalResponsePayload{}
	err = proto.Unmarshal(chaincodeActionPayload.Action.ProposalResponsePayload, proposalResponsePayload)
	if err != nil {
		return nil, err
	}
	chaincodeAction := &peer.ChaincodeAction{}
	err = proto.Unmarshal(proposalResponsePayload.Extension, chaincodeAction)
	if err != nil {
		return nil, err
	}
	return chaincodeAction, nil
}
