/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package fabricservice

import (
	"github.com/FujitsuLaboratories/ledgerdata-refiner/log"
	"github.com/FujitsuLaboratories/ledgerdata-refiner/model"
	fabledger "github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
)

var logger = log.Logger

const defaultConfigPath = "config/connection-config.yaml"

type FabClient struct {
	configPath string
	orgName    string
	orgAdmin   string
	orgUser    string
	channelId  string
	sdk        *fabsdk.FabricSDK
	ledger     *fabledger.Client
	hasInit    bool
}

// NewFabClient creates a new fabric client
func NewFabClient(orgName, orgAdmin, orgUser, channelID string) *FabClient {
	return &FabClient{
		configPath: defaultConfigPath,
		orgName:    orgName,
		orgAdmin:   orgAdmin,
		orgUser:    orgUser,
		channelId:  channelID,
		hasInit:    false,
	}
}

// Initialize inits fabric sdk and ledger client.
func (fc *FabClient) Initialize() error {
	if fc.hasInit {
		return nil
	}

	// init fabric client sdk
	logger.Info("Init fabric client")
	sdk, err := fabsdk.New(config.FromFile(fc.configPath))
	if err != nil {
		logger.WithField("error", err).Error("failed to init fabric sdk")
		return errors.Wrap(err, "failed to init fabric sdk")
	}

	fc.sdk = sdk

	// get network peers and orderers
	backend, err := sdk.Config()
	if err != nil {
		return errors.Wrap(err, "failed to get network config")
	}
	endpointConfig, err := fab.ConfigFromBackend(backend)
	if err != nil {
		return errors.Wrap(err, "failed to get endpoint config")
	}

	var nodes []model.Node
	peers := endpointConfig.NetworkPeers()
	if len(peers) > 0 {
		for _, peer := range peers {
			nodes = append(nodes, model.Node{Name: peer.GRPCOptions["ssl-target-name-override"].(string), Url: peer.URL, MSP: peer.MSPID, ChannelName: fc.channelId})
		}
	}

	orderers := endpointConfig.ChannelOrderers(fc.channelId)
	if len(orderers) > 0 {
		for _, orderer := range orderers {
			nodes = append(nodes, model.Node{Name: orderer.GRPCOptions["ssl-target-name-override"].(string), Url: orderer.URL, MSP: "OrdererMSP", ChannelName: fc.channelId})
		}
	}

	session := model.BeginSession()
	err = model.InsertNodes(session, nodes)
	if err != nil {
		logger.WithField("error", err).Error("failed to store node info")
		session.Rollback()
		return errors.Wrap(err, "failed to store node info")
	}
	session.Commit()

	// init fabric ledger client
	ctx := sdk.ChannelContext(fc.channelId, fabsdk.WithUser(fc.orgUser))
	ledger, err := fabledger.New(ctx)
	if err != nil {
		logger.WithField("error", err).Error("failed to init fabric ledger client")
		return errors.Wrap(err, "failed to init fabric ledger client")
	}

	fc.ledger = ledger
	fc.hasInit = true
	return nil
}

func (fc *FabClient) SDK() *fabsdk.FabricSDK {
	return fc.sdk
}

func (fc *FabClient) Ledger() *fabledger.Client {
	return fc.ledger
}

func (fc *FabClient) GetChannelId() string {
	return fc.channelId
}

// Teardown disconnect from fabric
func (fc *FabClient) Teardown() {
	logger.Info("disconnect from fabric")
	fc.sdk.Close()
}
