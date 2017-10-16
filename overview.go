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
	"strconv"
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
	BackgroundOrder  []StateType `yaml:"backgroundOrder"`
	BackgroundStates []StateType `yaml:"backgroundStates"`
	ColumnOrder      []string    `yaml:"columnOrder"`
	FlagOrder        []StateType `yaml:"flagOrder"`
	FlagStates       []StateType `yaml:"flagStates"`
	OverviewColumns  []string    `yaml:"overviewColumns"`
	Presets          []*Preset   `yaml:"presets"`
}

type Preset struct {
	Name              string
	AlwaysShownStates []StateType
	FilteredStates    []StateType
	Groups            []InvGroup
}

func (p *Preset) MarshalYAML() (interface{}, error) {
	return []interface{}{
		p.Name,
		[]interface{}{
			[]interface{}{"alwaysShownStates", p.AlwaysShownStates},
			[]interface{}{"filteredStates", p.FilteredStates},
			[]interface{}{"groups", p.Groups},
		},
	}, nil
}

func (p *Preset) UnmarshalYAML(f func(interface{}) error) error {
	v := []interface{}{
		"",
		[]interface{}{
			[]interface{}{"", make([]string, 0)},
			[]interface{}{"", make([]string, 0)},
			[]interface{}{"", make([]string, 0)},
		},
	}
	var err error
	if err = f(&v); err != nil {
		return err
	}
	p.Name = v[0].(string)
	for _, attr := range v[1].([]interface{}) {
		name, ns, err := parseAttribute(attr.([]interface{}))
		if err != nil {
			return fmt.Errorf("Preset %+q has bad %+q attribute: %s", p.Name, name, err)
		}
		switch name {
		case "alwaysShownStates":
			p.AlwaysShownStates = make([]StateType, len(ns))
			for i, n := range ns {
				p.AlwaysShownStates[i] = StateType(n)
			}
		case "filteredStates":
			p.FilteredStates = make([]StateType, len(ns))
			for i, n := range ns {
				p.FilteredStates[i] = StateType(n)
			}
		case "groups":
			p.Groups = make([]InvGroup, len(ns))
			for i, n := range ns {
				p.Groups[i] = InvGroup(n)
			}
		default:
			return fmt.Errorf("Preset %+q has unknown attribute: %+q", p.Name, name)
		}
	}
	return nil
}

func (p *Preset) String() string {
	return p.Name
}

func parseAttribute(entry []interface{}) (string, []int, error) {
	if len(entry) != 2 {
		return "", nil, fmt.Errorf("attribute has wrong length (got %d, expected 2): %+q",
			len(entry), entry)
	}
	name := entry[0].(string)
	vals := entry[1].([]interface{})
	var ns []int
	for _, val := range vals {
		switch val := val.(type) {
		case int:
			ns = append(ns, val)
		case string:
			ss := strings.SplitN(val, " ", 2)
			if len(ss) != 2 {
				return name, nil, fmt.Errorf("attribute value entry has format: %+q", val)
			}
			n, err := strconv.Atoi(ss[0])
			if err != nil {
				return name, nil, fmt.Errorf("attribute value not parsable as int: %+q", ss[0])
			}
			ns = append(ns, n)
		default:
			return name, nil, fmt.Errorf("attribute value entry has unexpected type: %T", val)
		}
	}
	return name, ns, nil
}
