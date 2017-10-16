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

type Preset struct {
	Name              string
	AlwaysShownStates *presetStates
	FilteredStates    *presetStates
	Groups            *presetGroups
}

func (p *Preset) MarshalYAML() (interface{}, error) {
	attrs := []interface{}{}
	if p.AlwaysShownStates != nil {
		attrs = append(attrs, p.AlwaysShownStates)
	}
	if p.FilteredStates != nil {
		attrs = append(attrs, p.FilteredStates)
	}
	if p.Groups != nil {
		attrs = append(attrs, p.Groups)
	}
	return []interface{}{p.Name, attrs}, nil
}

func (p *Preset) UnmarshalYAML(f func(interface{}) error) error {
	var v []interface{}
	var err error
	if err = f(&v); err != nil {
		return err
	}
	if len(v) != 2 {
		return fmt.Errorf(
			"Preset has wrong number of entries (Expected: 2 Got: %d): %+q", len(v), v)
	}
	p.Name, err = intfToString(v[0])
	if err != nil {
		return fmt.Errorf("Preset name: %v", err)
	}
	attrList, err := intfTointfSlice(v[1], 0)
	if err != nil {
		return fmt.Errorf("Preset (%+q) attribute list: %v", p.Name, err)
	}
	for _, attr := range attrList {
		name, ns, err := parseIntListAttr(attr)
		if err != nil {
			return fmt.Errorf("Preset (%+q) has bad attribute: %s", p.Name, err)
		}
		switch name {
		case "alwaysShownStates":
			p.AlwaysShownStates = newPresetState(name, ns)
		case "filteredStates":
			p.FilteredStates = newPresetState(name, ns)
		case "groups":
			p.Groups = newPresetGroup(ns)
		default:
			return fmt.Errorf("Preset %+q has unknown attribute: %+q", p.Name, name)
		}
	}
	return nil
}

type presetStates struct {
	Name   string
	States []StateType
}

func newPresetState(name string, ns []int) *presetStates {
	ps := &presetStates{Name: name}
	ps.States = make([]StateType, len(ns))
	for i, n := range ns {
		ps.States[i] = StateType(n)
	}
	return ps
}

func (ps *presetStates) MarshalYAML() (interface{}, error) {
	if ps == nil {
		return nil, nil
	}
	return []interface{}{ps.Name, ps.States}, nil
}

type presetGroups struct {
	Groups []InvGroup
}

func newPresetGroup(ns []int) *presetGroups {
	pg := &presetGroups{Groups: make([]InvGroup, len(ns))}
	for i, n := range ns {
		pg.Groups[i] = InvGroup(n)
	}
	return pg
}

func (pg *presetGroups) MarshalYAML() (interface{}, error) {
	if pg == nil {
		return nil, nil
	}
	return []interface{}{"groups", pg.Groups}, nil
}
