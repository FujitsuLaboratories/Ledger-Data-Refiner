/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"github.com/FujitsuLaboratories/ledgerdata-refiner/model"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewSyncService(t *testing.T) {
	syncService, err := NewSyncService()
	require.Nil(t, err)
	err = model.InitDB()
	require.Nil(t, nil)
	syncService.syncChannel()

}
