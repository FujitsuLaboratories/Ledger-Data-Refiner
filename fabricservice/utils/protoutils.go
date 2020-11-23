/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package utils

import (
	"crypto/sha256"
	"encoding/asn1"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/pkg/errors"
	"math"
)

const (
	SHA256 = "SHA256"
)

func GetEnvelopeFromBlock(data []byte) (*common.Envelope, error) {
	env := new(common.Envelope)

	if err := proto.Unmarshal(data, env); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling envelope")
	}
	return env, nil
}

func GetPayload(env *common.Envelope) (*common.Payload, error) {
	payload := new(common.Payload)
	if err := proto.Unmarshal(env.Payload, payload); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal payload")
	}
	return payload, nil
}

func ComputeSHA256(data []byte) []byte {
	bytes := sha256.Sum256(data)
	return bytes[:]
}

func GetTransaction(txBytes []byte) (*peer.Transaction, error) {
	tx := new(peer.Transaction)
	if err := proto.Unmarshal(txBytes, tx); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal transaction")
	}

	return tx, nil
}

func GetChannelHeader(txBytes []byte) (*common.ChannelHeader, error) {
	envelope, err := GetEnvelopeFromBlock(txBytes)
	if err != nil {
		return nil, err
	}
	if envelope == nil {
		return nil, errors.New("nil envelope")
	}

	// get payload from envelope
	payload, err := GetPayload(envelope)
	if err != nil {
		return nil, err
	}

	// extract channel header from payload
	channelHeaderBytes := payload.Header.ChannelHeader
	channelHeader := new(common.ChannelHeader)
	if err := proto.Unmarshal(channelHeaderBytes, channelHeader); err != nil {
		return nil, errors.Wrap(err, "error extracting ChannelHeader from payload")
	}

	return channelHeader, nil
}

// use to generate block hash
func HeaderBytes(header *common.BlockHeader) ([]byte, error) {
	asn1Header := struct {
		Number       int64
		PreviousHash []byte
		DataHash     []byte
	}{
		PreviousHash: header.PreviousHash,
		DataHash:     header.DataHash,
	}

	if header.Number > uint64(math.MaxInt64) {
		return nil, errors.New("golang does not currently support encoding uint64 to asn1")
	} else {
		asn1Header.Number = int64(header.Number)
	}

	return asn1.Marshal(asn1Header)
}

type SHA256Opts struct{}

func (s *SHA256Opts) Algorithm() string {
	return SHA256
}

func GetChaincodeProposalPayload(bytes []byte) (*peer.ChaincodeProposalPayload, error) {
	cpp := &peer.ChaincodeProposalPayload{}
	err := proto.Unmarshal(bytes, cpp)
	return cpp, errors.Wrap(err, "error unmarshaling ChaincodeProposalPayload")
}

func GetChaincodeActionPayload(capBytes []byte) (*peer.ChaincodeActionPayload, error) {
	cap := &peer.ChaincodeActionPayload{}
	err := proto.Unmarshal(capBytes, cap)
	return cap, errors.Wrap(err, "error unmarshaling ChaincodeActionPayload")
}
