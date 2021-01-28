/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package utils

import (
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
)

// TODO use a struct ??
var (
	// server config
	AppMode  string
	HttpPort string
	// log config
	LogOutput  string
	LogLevel   string
	LogFormat  string
	LogStore   string
	LogMaxAge  int64
	LogMaxSize int64
	// database config
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// fabric config
	NetworkName string
	ClientName  string
	OrgName     string
	OrgAdmin    string
	OrgUser     string
	ChannelId   string
)

func init() {
	file, err := ini.Load("config/config.ini")
	if err != nil {
		file, err = ini.Load("../config/config.ini")
		if err != nil {
			file, err = ini.Load("../../config/config.ini")
			if err != nil {
				panic("failed to load ini config. err: " + err.Error())
			}
		}
	}

	loadServer(file)
	loadLog(file)
	loadDatabase(file)
	loadFabric(file)
}

func loadServer(file *ini.File) {
	AppMode = file.Section("server").Key("app_mode").MustString("debug")
	HttpPort = file.Section("server").Key("http_port").MustString(":9999")
}

func loadLog(file *ini.File) {
	LogOutput = file.Section("log").Key("log_output").MustString("stdout")
	LogLevel = file.Section("log").Key("log_level").MustString("debug")
	LogFormat = file.Section("log").Key("log_format").MustString("text")
	LogStore = file.Section("log").Key("log_store").MustString("")
	LogStore = filepath.Join(os.Getenv("HOME"), LogStore)
	LogMaxAge = file.Section("log").Key("log_max_age").MustInt64(5 * 24 * 60 * 60)
	LogMaxSize = file.Section("log").Key("log_max_size").MustInt64(1048576)
}

func loadDatabase(file *ini.File) {
	DBHost = file.Section("database").Key("db_host").MustString("localhost")
	DBPort = file.Section("database").Key("db_port").MustString("5432")
	DBUser = file.Section("database").Key("db_user").MustString("refiner")
	DBPassword = file.Section("database").Key("db_password").MustString("123456")
	DBName = file.Section("database").Key("db_name").MustString("ledgerdata_refiner")
}

func loadFabric(file *ini.File) {
	OrgName = file.Section("fabric").Key("org_name").MustString("Org1")
	OrgAdmin = file.Section("fabric").Key("org_admin").MustString("Admin")
	OrgUser = file.Section("fabric").Key("org_user").MustString("User1")
	ChannelId = file.Section("fabric").Key("channel_id").MustString("mychannel")
	NetworkName = file.Section("fabric").Key("network_name").MustString("test-network")
	ClientName = file.Section("fabric").Key("client_name").MustString("test-client")
}
