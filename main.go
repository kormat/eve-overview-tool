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
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

const commentMarker = "EOTCOMMENTMARKER"

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
	fmt.Printf("%s", string(unescapeComments(b)))
}

var quotesRx = regexp.MustCompile(`^(?P<start>^\s*(- )+)'?(?P<entry>\d+ ` + commentMarker + ` .+)$`)

func unescapeComments(b []byte) []byte {
	var out bytes.Buffer
	s := bufio.NewScanner(bytes.NewReader(b))
	for s.Scan() {
		line := s.Bytes()
		matches := quotesRx.FindStringSubmatch(string(line))
		if len(matches) > 0 {
			out.WriteString(matches[1])
			out.WriteString(strings.Replace(matches[3], commentMarker, "#", 1))
		} else {
			out.Write(line)
		}
		out.WriteString("\r\n")
	}
	return out.Bytes()
}
