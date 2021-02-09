/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package utils

import (
	"encoding/json"
	"github.com/pkg/errors"
	"time"
)

// StringToTime parse string to timestamp
func StringToTime(timeStr string) (time.Time, error) {
	timeTemplate := "2006-01-02 15:04:05"
	stamp, err := time.ParseInLocation(timeTemplate, timeStr, time.Local)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "failed to parse string to time")
	}

	return stamp, nil
}

func ToJson(arg interface{}) string {
	bytes, _ := json.Marshal(arg)
	return string(bytes)
}

// IsJson checks if the content is json
func IsJson(content string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(content), &js) == nil
}

func IsContain(str string, strs []string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}
	return false
}

// RemoveInvalidCharacters removes some characters, like `\u0001`,`\u001A`,...
func RemoveInvalidCharacters(str string) string {
	if str == "" {
		return ""
	}
	srcRunes := []rune(str)
	dstRunes := make([]rune, 0, len(srcRunes))
	// remove useless characters
	for _, c := range srcRunes {
		if c >= 0 && c <= 31 {
			continue
		}
		if c == 127 {
			continue
		}
		dstRunes = append(dstRunes, c)
	}

	return string(dstRunes)
}
