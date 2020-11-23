/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package main

import (
	"github.com/FujitsuLaboratories/ledgerdata-refiner/log"
	"github.com/FujitsuLaboratories/ledgerdata-refiner/model"
	"github.com/FujitsuLaboratories/ledgerdata-refiner/router"
	"github.com/FujitsuLaboratories/ledgerdata-refiner/sync"
	"os"
)

var logger = log.Logger

// @title Ledgerdata refiner api doc
// @version 1.0

// @contact.name API Support
// @contact.url http://www.swagger.io/support

// @licence.name Apache License 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1:30052
// @BasePath /api
func main() {
	err := model.InitDB()
	if err != nil {
		logger.WithField("error", err).Error("failed to init database")
		os.Exit(1)
	}

	syncService, err := sync.NewSyncService()
	if err != nil {
		logger.Error("failed to init synchronization service")
		os.Exit(1)
	}
	syncService.Initialize()

	if err := router.InitRouter(); err != nil {
		logger.WithField("error", err).Error("failed to init gin")
		os.Exit(1)
	}
}
