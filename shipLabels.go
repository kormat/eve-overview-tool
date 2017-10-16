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

type ShipLabel struct {
	Name  NullableString
	Post  string
	Pre   string
	State ShipLabelState
	Type  NullableString
}

func (sl *ShipLabel) MarshalYAML() (interface{}, error) {
	attrs := [][]interface{}{
		{"post", sl.Post},
		{"pre", sl.Pre},
		{"state", sl.State},
		{"type", sl.Type},
	}
	return []interface{}{sl.Name, attrs}, nil
}

func (sl *ShipLabel) UnmarshalYAML(f func(interface{}) error) error {
	var v []interface{}
	var err error
	if err = f(&v); err != nil {
		return err
	}
	if len(v) != 2 {
		return fmt.Errorf(
			"ShipLabel has wrong number of entries (Expected: 2 Got: %d): %+q", len(v), v)
	}
	name, err := intfToString(v[0])
	if err != nil {
		return fmt.Errorf("ShipLabel name: %v", err)
	}
	sl.Name = NullableString(name)
	attrList, err := intfTointfSlice(v[1], 0)
	if err != nil {
		return fmt.Errorf("ShipLabel (%+q) attribute list: %v", sl.Name, err)
	}
	for _, attr := range attrList {
		name, attrVal, err := parseAttr(attr)
		if err != nil {
			return fmt.Errorf("ShipLabel (%+q) has bad attribute: %s", sl.Name, err)
		}
		if name == "state" {
			n, err := intfToInt(attrVal)
			if err != nil {
				return fmt.Errorf("ShipLabel (%+q) state attribute: %s", sl.Name, err)
			}
			sl.State = ShipLabelState(n)
			continue
		}
		// Otherwise it's a string.
		s, err := intfToString(attrVal)
		if err != nil {
			return fmt.Errorf("ShipLabel (%+q) %+q attribute: %s", sl.Name, name, err)
		}
		switch name {
		case "post":
			sl.Post = s
		case "pre":
			sl.Pre = s
		case "type":
			sl.Type = NullableString(s)
		default:
			return fmt.Errorf("ShipLabel (%+q) has unknown attribute: %+q", sl.Name, name)
		}
	}
	return nil
}
