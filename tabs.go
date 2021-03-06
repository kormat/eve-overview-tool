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

type TabSetup struct {
	Id           int
	Bracket      NullableString
	Name         string
	Overview     string
	ShowAll      *bool
	ShowNone     *bool
	ShowSpecials *bool
}

func (ts *TabSetup) MarshalYAML() (interface{}, error) {
	attrs := [][]interface{}{
		{"bracket", ts.Bracket},
		{"name", ts.Name},
		{"overview", ts.Overview},
	}
	if ts.ShowAll != nil {
		attrs = append(attrs, []interface{}{"showAll", ts.ShowAll})
	}
	if ts.ShowNone != nil {
		attrs = append(attrs, []interface{}{"showNone", ts.ShowNone})
	}
	if ts.ShowSpecials != nil {
		attrs = append(attrs, []interface{}{"showSpecials", ts.ShowSpecials})
	}
	return []interface{}{ts.Id, attrs}, nil
}

func (ts *TabSetup) UnmarshalYAML(f func(interface{}) error) error {
	var v []interface{}
	var err error
	if err = f(&v); err != nil {
		return err
	}
	if len(v) != 2 {
		return fmt.Errorf(
			"TabSetup has wrong number of entries (Expected: 2 Got: %d): %+q", len(v), v)
	}
	ts.Id, err = intfToInt(v[0])
	if err != nil {
		return fmt.Errorf("TabSetup id: %v", err)
	}
	attrList, err := intfTointfSlice(v[1], 0)
	if err != nil {
		return fmt.Errorf("TabSetup (%d) attribute list: %v", ts.Id, err)
	}
	for _, attr := range attrList {
		name, attrVal, err := parseAttr(attr)
		if err != nil {
			return fmt.Errorf("TabSetup (%d) has bad attribute: %s", ts.Id, err)
		}
		var b bool
		switch name {
		case "bracket":
			s, err := intfToString(attrVal)
			if err != nil {
				return fmt.Errorf("TabSetup (%d) attribute %+q: %s", ts.Id, name, err)
			}
			ts.Bracket = NullableString(s)
		case "name":
			if ts.Name, err = intfToString(attrVal); err != nil {
				return fmt.Errorf("TabSetup (%d) attribute %+q: %s", ts.Id, name, err)
			}
		case "overview":
			if ts.Overview, err = intfToString(attrVal); err != nil {
				return fmt.Errorf("TabSetup (%d) attribute %+q: %s", ts.Id, name, err)
			}
		case "showAll":
			if b, err = intfToBool(attrVal); err != nil {
				return fmt.Errorf("TabSetup (%d) attribute %+q: %s", ts.Id, name, err)
			}
			ts.ShowAll = &b
		case "showNone":
			if b, err = intfToBool(attrVal); err != nil {
				return fmt.Errorf("TabSetup (%d) attribute %+q: %s", ts.Id, name, err)
			}
			ts.ShowNone = &b
		case "showSpecials":
			if b, err = intfToBool(attrVal); err != nil {
				return fmt.Errorf("TabSetup (%d) attribute %+q: %s", ts.Id, name, err)
			}
			ts.ShowSpecials = &b
		default:
			return fmt.Errorf("TabSetup (%d) has unknown attribute: %+q", ts.Id, name)
		}
	}
	return nil
}
