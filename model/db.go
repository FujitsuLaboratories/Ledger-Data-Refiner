/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package model

import (
	"fmt"
	"github.com/FujitsuLaboratories/ledgerdata-refiner/log"
	"github.com/FujitsuLaboratories/ledgerdata-refiner/utils"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

var (
	db     *gorm.DB
	dbErr  error
	logger = log.Logger
)

func InitDB() error {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=disable TimeZone=Asia/Shanghai",
		utils.DBUser, utils.DBPassword, utils.DBName, utils.DBPort, utils.DBHost)
	db, dbErr = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if dbErr != nil {
		return errors.Wrap(dbErr, "failed to connect to database")
	}

	sqlDB, dbErr := db.DB()
	if dbErr != nil {
		return errors.Wrap(dbErr, "failed to get database object")
	}

	// set the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)
	// set the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(10)
	// set the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Second * 10)

	dbErr = db.AutoMigrate(&Channel{}, &Block{}, &Transaction{},
		&DocumentHistory{}, &Document{}, &Schema{}, &Node{})
	if dbErr != nil {
		return errors.Wrap(dbErr, "failed to automatically create tables")
	}

	return nil
}

// Manual begin transaction
func BeginSession() *gorm.DB {
	return db.Begin()
}
