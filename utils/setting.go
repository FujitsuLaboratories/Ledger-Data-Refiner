/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/viper"
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
	viper.SetConfigName("config")
	viper.SetConfigType("ini")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")
	viper.SetEnvPrefix("refiner")
	viper.AutomaticEnv()

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	loadDefault()
	loadServer()
	loadLog()
	loadDatabase()
	loadFabric()
}
func loadDefault() {
	viper.SetDefault("server.app_mode", "debug")
	viper.SetDefault("server.http_port", ":9999")

	viper.SetDefault("log.log_output", "stdout")
	viper.SetDefault("log.log_level", "debug")
	viper.SetDefault("log.log_format", "text")
	viper.SetDefault("log.log_store", "")
	viper.SetDefault("log.log_max_age", "432000")   // 5 * 24 * 60 * 60
	viper.SetDefault("log.log_max_size", "1048576") // 1024 * 1024

	viper.SetDefault("database.db_host", "localhost")
	viper.SetDefault("database.db_port", "5432")
	viper.SetDefault("database.db_user", "refiner")
	viper.SetDefault("database.db_password", "123456")
	viper.SetDefault("database.db_name", "ledgerdata_refiner")

	viper.SetDefault("fabric.org_name", "Org1")
	viper.SetDefault("fabric.org_admin", "Admin")
	viper.SetDefault("fabric.org_user", "User1")
	viper.SetDefault("fabric.channel_id", "mychannel")
	viper.SetDefault("fabric.network_name", "test-network")
	viper.SetDefault("fabric.client_name", "test-client")
}

func loadServer() {
	AppMode = viper.GetString("server.app_mode")
	HttpPort = viper.GetString("server.http_port")
}

func loadLog() {
	LogOutput = viper.GetString("log.log_output")
	LogLevel = viper.GetString("log.log_level")
	LogFormat = viper.GetString("log.log_format")
	LogStore = viper.GetString("log.log_store")
	LogStore = filepath.Join(os.Getenv("HOME"), LogStore)
	LogMaxAge, _ = strconv.ParseInt(viper.GetString("log.log_max_age"), 10, 64)
	LogMaxSize, _ = strconv.ParseInt(viper.GetString("log.log_max_size"), 10, 64)
}

func loadDatabase() {
	DBHost = viper.GetString("database.db_host")
	DBPort = viper.GetString("database.db_port")
	DBUser = viper.GetString("database.db_user")
	DBPassword = viper.GetString("database.db_password")
	DBName = viper.GetString("database.db_name")
}

func loadFabric() {
	OrgName = viper.GetString("fabric.org_name")
	OrgAdmin = viper.GetString("fabric.org_admin")
	OrgUser = viper.GetString("fabric.org_user")
	ChannelId = viper.GetString("fabric.channel_id")
	NetworkName = viper.GetString("fabric.network_name")
	ClientName = viper.GetString("fabric.client_name")
}
