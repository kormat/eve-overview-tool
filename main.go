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
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var cfgFile = flag.String("f", "", "Overview file to operate on")

var invGroups map[InvGroup]string
var stateTypes map[StateType]string

func main() {
	var err error
	flag.Parse()
	if *cfgFile == "" {
		log.Printf("ERROR: No overview file specified.")
		os.Exit(1)
	}
	if invGroups, err = loadGroups(); err != nil {
		log.Printf("ERROR: unable to load inventory types CSV file: %s", err)
		os.Exit(1)
	}
	if stateTypes, err = loadStates(); err != nil {
		log.Printf("ERROR: unable to load filter state types CSV file: %s", err)
		os.Exit(1)
	}
	var o *Overview
	if o, err = loadConfig(); err != nil {
		log.Printf("ERROR: unable to load overview file: %s", err)
		os.Exit(1)
	}
	b, err := yaml.Marshal(o)
	if err != nil {
		log.Printf("ERROR: unable to marshal back to yaml: %s", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", string(b))
}

type InvGroup int

func (ig InvGroup) name() string {
	s, ok := invGroups[ig]
	if !ok {
		s = "Unknown InvGroup"
	}
	return s
}

func (ig InvGroup) String() string {
	return fmt.Sprintf("%s (%d)", ig.name(), int(ig))
}

func (ig InvGroup) MarshalYAML() (interface{}, error) {
	return fmt.Sprintf("%d # %s", int(ig), ig.name()), nil
}

type StateType int

func (st StateType) name() string {
	s, ok := stateTypes[st]
	if !ok {
		s = "Unknown StateType"
	}
	return s
}

func (st StateType) String() string {
	return fmt.Sprintf("%s (%d)", st.name(), int(st))
}

func (st StateType) MarshalYAML() (interface{}, error) {
	return fmt.Sprintf("%d # %s", int(st), st.name()), nil
}
