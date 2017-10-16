// Copyright 2017 Stephen Shirley
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
)

func parseIntListAttr(v interface{}) (string, []int, error) {
	vlist, err := intfTointfSlice(v, 2)
	if err != nil {
		return "", nil, fmt.Errorf("Int list attribute: %v", err)
	}
	name, err := intfToString(vlist[0])
	if err != nil {
		return "", nil, fmt.Errorf("Int list attribute name: %v", err)
	}
	vals, err := intfTointfSlice(vlist[1], 0)
	if err != nil {
		return name, nil, fmt.Errorf("Int list attribute (%+q) value list: %v", name, err)
	}
	var ns []int
	for _, val := range vals {
		n, err := intfToInt(val)
		if err != nil {
			return name, nil, fmt.Errorf(
				"Int list attribute (%+q) value entry: %v", name, err)
		}
		ns = append(ns, n)
	}
	return name, ns, nil
}

func parseAttr(v interface{}) (string, interface{}, error) {
	vlist, err := intfTointfSlice(v, 2)
	if err != nil {
		return "", "", fmt.Errorf("Attribute: %v", err)
	}
	name, err := intfToString(vlist[0])
	if err != nil {
		return "", "", fmt.Errorf("Attribute name: %v", err)
	}
	return name, vlist[1], nil
}

func intfTointfSlice(v interface{}, n int) ([]interface{}, error) {
	vlist, ok := v.([]interface{})
	if !ok {
		return nil, fmt.Errorf("type is not []interface{} (%T): %+q", v, v)
	}
	if n > 0 && len(vlist) != n {
		return nil, fmt.Errorf("wrong number of entries (Expected: %d Got: %d): %+q",
			n, len(vlist), vlist)
	}
	return vlist, nil
}

func intfToString(v interface{}) (string, error) {
	switch v := v.(type) {
	case nil:
		// Handle the case when the config contains `null` instead of empty string.
		return "", nil
	case string:
		return v, nil
	default:
		return "", fmt.Errorf("type is not string (%T): %+q", v, v)
	}
}

func intfToInt(v interface{}) (int, error) {
	n, ok := v.(int)
	if !ok {
		return -1, fmt.Errorf("type is not int (%T): %+q", v, v)
	}
	return n, nil
}

func intfToBool(v interface{}) (bool, error) {
	b, ok := v.(bool)
	if !ok {
		return b, fmt.Errorf("type is not bool (%T): %+q", v, v)
	}
	return b, nil
}
