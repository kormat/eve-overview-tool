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

type StateBlink struct {
	Name string
	Val  bool
}

func (sb *StateBlink) MarshalYAML() (interface{}, error) {
	return []interface{}{sb.Name, sb.Val}, nil
}

func (sb *StateBlink) UnmarshalYAML(f func(interface{}) error) error {
	var v []interface{}
	var err error
	if err = f(&v); err != nil {
		return err
	}
	if len(v) != 2 {
		return fmt.Errorf(
			"StateBlink has wrong number of entries (Expected: 2 Got: %d): %+q", len(v), v)
	}
	sb.Name, err = intfToString(v[0])
	if err != nil {
		return fmt.Errorf("StateBlink name: %v", err)
	}
	sb.Val, err = intfToBool(v[1])
	if err != nil {
		return fmt.Errorf("StateBlink (%+q) value: %v", sb.Name, err)
	}
	return nil
}

type StateColorName struct {
	Name string
	Val  string
}

func (sc *StateColorName) MarshalYAML() (interface{}, error) {
	return []interface{}{sc.Name, sc.Val}, nil
}

func (sc *StateColorName) UnmarshalYAML(f func(interface{}) error) error {
	var v []interface{}
	var err error
	if err = f(&v); err != nil {
		return err
	}
	if len(v) != 2 {
		return fmt.Errorf(
			"StateColorName has wrong number of entries (Expected: 2 Got: %d): %+q", len(v), v)
	}
	sc.Name, err = intfToString(v[0])
	if err != nil {
		return fmt.Errorf("StateColorName name: %v", err)
	}
	sc.Val, err = intfToString(v[1])
	if err != nil {
		return fmt.Errorf("StateColorName (%+q) value: %v", sc.Name, err)
	}
	return nil
}
