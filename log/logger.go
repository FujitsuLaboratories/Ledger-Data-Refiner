/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package log

import (
	"github.com/FujitsuLaboratories/ledgerdata-refiner/utils"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	// set log level
	setLevel(Logger)
	// set output format
	setFormat(Logger)
	switch utils.LogOutput {
	case "file":
		// set log cutting
		writer, err := rotatelogs.New(
			utils.LogStore+"%Y%m%d%H%M",
			rotatelogs.WithLinkName(utils.LogStore),
			rotatelogs.WithMaxAge(time.Duration(utils.LogMaxAge)),
			rotatelogs.WithRotationSize(utils.LogMaxSize),
		)
		if err != nil {
			panic(err)
		}
		Logger.SetOutput(writer)
	default:
		Logger.SetOutput(os.Stdout)
	}
}

func setLevel(logger *logrus.Logger) {
	switch utils.LogLevel {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	default:
		logger.SetLevel(logrus.DebugLevel)
	}
}

func setFormat(logger *logrus.Logger) {
	switch utils.LogFormat {
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{})
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{})
	}
}
