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
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

func loadConfig() (*Overview, error) {
	b, err := ioutil.ReadFile(*cfgFile)
	if err != nil {
		return nil, err
	}
	var o Overview
	if err := yaml.Unmarshal(b, &o); err != nil {
		return nil, err
	}
	return &o, nil
}

type Overview struct {
	BackgroundOrder     []StateType       `yaml:"backgroundOrder"`
	BackgroundStates    []StateType       `yaml:"backgroundStates"`
	ColumnOrder         []string          `yaml:"columnOrder"`
	FlagOrder           []StateType       `yaml:"flagOrder"`
	FlagStates          []StateType       `yaml:"flagStates"`
	OverviewColumns     []string          `yaml:"overviewColumns"`
	Presets             []*Preset         `yaml:"presets"`
	ShipLabelOrder      []NullableString  `yaml:"shipLabelOrder"`
	ShipLabels          []*ShipLabel      `yaml:"shipLabels"`
	StateBlinks         []*StateBlink     `yaml:"stateBlinks"`
	StateColorsNameList []*StateColorName `yaml:"stateColorsNameList"`
	TabSetup            []*TabSetup       `yaml:"tabSetup"`
	UserSettings        []*UserSetting    `yaml:"userSettings"`
}

type StateType int

func (st StateType) name() string {
	s, ok := stateTypes[st]
	if !ok {
		return "Unknown StateType"
	}
	return strings.TrimSpace(s)
}

func (st StateType) String() string {
	return fmt.Sprintf("%s (%d)", st.name(), int(st))
}

func (st StateType) MarshalYAML() (interface{}, error) {
	return fmt.Sprintf("%d %s %s", int(st), commentMarker, st.name()), nil
}

type ShipLabelState int

func (ss ShipLabelState) name() string {
	switch ss {
	case 0:
		return "Disabled"
	case 1:
		return "Enabled"
	default:
		return "Unknown ShipLabelState"
	}
}

func (ss ShipLabelState) String() string {
	return fmt.Sprintf("%s (%d)", ss.name(), int(ss))
}

func (ss ShipLabelState) MarshalYAML() (interface{}, error) {
	return fmt.Sprintf("%d %s %s", int(ss), commentMarker, ss.name()), nil
}

type NullableString string

func (ns NullableString) MarshalYAML() (interface{}, error) {
	if len(ns) == 0 {
		return nil, nil
	}
	return string(ns), nil
}
