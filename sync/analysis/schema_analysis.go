/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package analysis

import (
	"encoding/json"
	"github.com/FujitsuLaboratories/ledgerdata-refiner/utils"
	"github.com/pkg/errors"
	"reflect"
	"sort"
	"strings"
)

const (
	Equal = iota
	Contain
	BeContained
	Difference
)

type SchemaArray []string

func (s SchemaArray) Len() int {
	return len(s)
}

func (s SchemaArray) Less(i, j int) bool {
	return strings.Compare(s[i], s[j]) < 0
}

func (s SchemaArray) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// iterGetSchemaInArray extracts all the keys of a json object through a depth-first traversal,
// and sorted them in allKeys
func iterGetSchemaInArray(content map[string]interface{}, keys *SchemaArray, prefix string) {

	// iterate the map to get all the keys
	for key, value := range content {
		var (
			subKey map[string]interface{}
			ok     bool
		)
		// if the value is still a map, do it again
		if subKey, ok = value.(map[string]interface{}); ok {
			if prefix != "" {
				iterGetSchemaInArray(subKey, keys, prefix+key+".")
			} else {
				iterGetSchemaInArray(subKey, keys, key+".")
			}
		}
		// TODO use struct instead of string slice, the struct should implement interface{}
		if prefix != "" {
			*keys = append(*keys, prefix+key)
		} else {
			*keys = append(*keys, key)
		}
	}
}

// GetSchemaInArray extracts all the keys of a json object
func GetSchemaInArray(content string) (SchemaArray, error) {
	if content == "" {
		return nil, errors.New("The object is empty or undefined")
	}
	if !utils.IsJson(content) {
		return nil, errors.New("the object is not a JSON schema")
	}
	keys := &SchemaArray{}
	args := make(map[string]interface{})
	err := json.Unmarshal([]byte(content), &args)
	if err != nil {
		return nil, err
	}
	iterGetSchemaInArray(args, keys, "")
	sort.Sort(keys)
	return *keys, nil
}

// GetSchemaInJson extracts json schema
func GetSchemaInJson(content string) (string, error) {
	if content == "" {
		return "", errors.New("The object is empty or undefined")
	}
	if !utils.IsJson(content) {
		return "", errors.New("the object is not a JSON schema")
	}
	keys := make(map[string]interface{})
	args := make(map[string]interface{})
	err := json.Unmarshal([]byte(content), &args)
	if err != nil {
		return "", err
	}
	iterGetSchemaInJson(keys, args)
	jsonBytes, err := json.Marshal(keys)
	return string(jsonBytes), err
}

func iterGetSchemaInJson(keys map[string]interface{}, args map[string]interface{}) {
	for k, v := range args {
		if v == nil {
			continue
		}
		if subKey, ok := v.(map[string]interface{}); ok {
			keys[k] = make(map[string]interface{})
			iterGetSchemaInJson((keys[k]).(map[string]interface{}), subKey)
		} else if reflect.ValueOf(v).Kind() == reflect.Slice ||
			reflect.ValueOf(v).Kind() == reflect.Array {
			keys[k] = "[" + reflect.TypeOf(v.([]interface{})[0]).Name() + "]"
		} else {
			keys[k] = reflect.TypeOf(v).Name()
		}
	}
}

// SchemaCompare compares array a and array b
// "Equal":if a == b
// "Contain": if a contains b
// "Be contained" : if b contains a
// "Difference": if a doesn't contain b and isn't contained by b
func SchemaCompare(a, b SchemaArray) int {
	var beCompared SchemaArray
	var compared SchemaArray
	//
	if a.Len() >= b.Len() {
		beCompared, compared = b, a
	} else {
		beCompared, compared = a, b
	}

	for i := 0; i < beCompared.Len(); i++ {
		if strings.Compare(beCompared[i], compared[i]) != 0 {
			return Difference
		}
	}

	if a.Len() == b.Len() {
		return Equal
	}

	if a.Len() == beCompared.Len() {
		return BeContained
	} else {
		return Contain
	}
}
