/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"context"
	"github.com/FujitsuLaboratories/ledgerdata-refiner/fabricservice"
	"github.com/FujitsuLaboratories/ledgerdata-refiner/log"
	"github.com/FujitsuLaboratories/ledgerdata-refiner/model"
	"github.com/FujitsuLaboratories/ledgerdata-refiner/utils"
	"sync"
	"time"
)

var logger = log.Logger

const (
	syncInterval = 5
)

type SyncService struct {
	clientName        string
	networkName       string
	timer             *time.Timer
	fabClient         *fabricservice.FabClient
	cancel            context.CancelFunc
	storage           map[string][]SchemaStorage
	channelId         int
	firstSync         bool
	syncChannelStatus bool
	syncSchemaStatus  bool
	lostHeight        map[uint64]struct{} // Track those blocks that failed to sync
	mtx               sync.Mutex
}

type SchemaStorage struct {
	schemaId    int
	schemaArray string
}

func NewSyncService() (*SyncService, error) {
	timer := time.NewTimer(syncInterval * time.Second)
	// stop the timer, we don't need it right now
	timer.Stop()

	// connect to fabric and init it
	fabClient := fabricservice.NewFabClient(utils.OrgName, utils.OrgAdmin, utils.OrgUser,
		utils.ChannelId)
	err := fabClient.Initialize()
	if err != nil {
		return nil, err
	}

	syncService := &SyncService{
		timer:       timer,
		fabClient:   fabClient,
		networkName: utils.NetworkName,
		clientName:  utils.ClientName,
		lostHeight:  make(map[uint64]struct{}),
		firstSync:   true,
	}

	return syncService, nil
}

func (sync *SyncService) Initialize() {
	ctx, cancel := context.WithCancel(context.Background())
	sync.cancel = cancel

	go sync.syncWork(ctx)

}

// markLostHeight records lost heights
func (sync *SyncService) markLostHeight(height uint64) {
	sync.mtx.Lock()
	defer sync.mtx.Unlock()
	sync.lostHeight[height] = struct{}{}
}

func (sync *SyncService) removeLostHeight(height uint64) {
	sync.mtx.Lock()
	defer sync.mtx.Unlock()
	delete(sync.lostHeight, height)
}

func (sync *SyncService) getLostHeights() (heights []uint64) {
	sync.mtx.Lock()
	defer sync.mtx.Unlock()
	if len(sync.lostHeight) > 0 {
		for height := range sync.lostHeight {
			heights = append(heights, height)
		}
	} else {
		logger.Debug("there no lost blocks to sync")
	}
	return
}

func (sync *SyncService) syncLostBlocks() {
	heights := sync.getLostHeights()
	for _, height := range heights {
		session := model.BeginSession()
		err := sync.syncBlock(session, height)
		if err == nil {
			sync.removeLostHeight(height)
		}
		session.Commit()
	}
}

func (sync *SyncService) syncWork(ctx context.Context) {
	sync.timer.Reset(syncInterval * time.Second)
OUT_LOOP:
	for {
		select {
		case <-ctx.Done():
			break OUT_LOOP
		case <-sync.timer.C:
			sync.timer.Reset(syncInterval * time.Second)
			if !sync.syncChannelStatus {
				sync.syncChannel()
			}
			if !sync.syncSchemaStatus {
				sync.schemaAnalysis()
			}
			// handle lost blocks
			if len(sync.lostHeight) > 0 {
				go sync.syncLostBlocks()
			}
		}
	}
}

func (sync *SyncService) Shutdown() {
	sync.fabClient.Teardown()
	sync.cancel()
}
