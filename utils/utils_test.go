/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStringToTime(t *testing.T) {
	str := "2020-09-15 11:34:15"
	time, err := StringToTime(str)
	require.Nil(t, err)

	t.Log(time.String())
}

func TestToJson(t *testing.T) {
	var a []map[string]string
	b := map[string]string{"b": "b"}
	c := map[string]string{"c": "c"}
	a = append(a, b)
	a = append(a, c)

	json := ToJson(a)
	t.Log(json)
}

func TestIsJson(t *testing.T) {
	jsons := []string{
		`{"a":"1"}`,
		`{aa":1}}`,
		"1",
		`{"aa":[1,2]]}`,
	}

	for i, json := range jsons {
		b := IsJson(json)
		if i == 0 {
			require.True(t, b)
			continue
		}
		require.False(t, b)
	}
}
