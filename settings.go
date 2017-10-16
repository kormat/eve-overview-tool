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

type UserSetting struct {
	Name string
	Val  bool
}

func (us *UserSetting) MarshalYAML() (interface{}, error) {
	return []interface{}{us.Name, us.Val}, nil
}

func (us *UserSetting) UnmarshalYAML(f func(interface{}) error) error {
	var v []interface{}
	var err error
	if err = f(&v); err != nil {
		return err
	}
	if len(v) != 2 {
		return fmt.Errorf(
			"UserSetting has wrong number of entries (Expected: 2 Got: %d): %+q", len(v), v)
	}
	us.Name, err = intfToString(v[0])
	if err != nil {
		return fmt.Errorf("UserSetting name: %v", err)
	}
	us.Val, err = intfToBool(v[1])
	if err != nil {
		return fmt.Errorf("UserSetting (%+q) value: %v", us.Name, err)
	}
	return nil
}
