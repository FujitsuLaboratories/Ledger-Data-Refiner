/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package analysis

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetSchemaInArray(t *testing.T) {
	testStrs := []string{
		"",
		"not a json",
		"{\"name\":\"jack\", \"age\":\"12\"}",
		"{\"name\":{\"name\":\"jack\", \"age\":\"12\"}, \"age\":{\"month\":12, \"day\":12}}",
	}

	key, err := GetSchemaInArray(testStrs[0])
	require.NotNil(t, err)

	key, err = GetSchemaInArray(testStrs[1])
	require.NotNil(t, err)

	key, err = GetSchemaInArray(testStrs[2])
	require.Nil(t, err)
	t.Log(key)

	key, err = GetSchemaInArray(testStrs[3])
	require.Nil(t, err)
	t.Log(key)
}

func TestGetSchemaInJson(t *testing.T) {
	testStrs := []string{
		"",
		"not a json",
		"{\"name\":\"jack\", \"age\":\"12\"}",
		"{\"name\":{\"name\":\"jack\", \"age\":\"12\"}, \"age\":{\"month\":12, \"day\":12}}",
		"{\"name\":[12,45,12], \"age\":{\"month\":12, \"day\":12}}",
	}

	key, err := GetSchemaInJson(testStrs[0])
	require.NotNil(t, err)

	key, err = GetSchemaInJson(testStrs[1])
	require.NotNil(t, err)

	key, err = GetSchemaInJson(testStrs[2])
	require.Nil(t, err)
	t.Log(key)

	key, err = GetSchemaInJson(testStrs[3])
	require.Nil(t, err)
	bytes, _ := json.Marshal(&key)
	t.Log(string(bytes))

	key, err = GetSchemaInJson(testStrs[4])
	require.Nil(t, err)
	bytes, _ = json.Marshal(&key)
	t.Log(string(bytes))
}

func TestSchemaCompare(t *testing.T) {
	testStrs := []string{
		"{\"name\":{\"name\":\"jack\", \"age\":\"12\"}, \"age\":{\"month\":12, \"day\":12}}",
		"{\"name\":[12,45,12], \"age\":{\"month\":12, \"day\":12}}",
	}

	key1, err := GetSchemaInArray(testStrs[0])
	require.Nil(t, err)

	key2, err := GetSchemaInArray(testStrs[1])
	require.Nil(t, err)

	compare := SchemaCompare(key1, key2)
	t.Log(compare)
}
