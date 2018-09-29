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
	"path"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

const commentMarker = "EOTCOMMENT"
const allGroupPreset = "all"

var cfgFile = flag.String("f", "", "Overview file to operate on")
var updGroups = flag.Bool("update-groups", false, "Update groups/ using an 'All' preset.")

var invGroups map[InvGroupId]*InvGroup
var invCategories map[InvCategoryId]string
var stateTypes map[StateType]string

func main() {
	var err error
	flag.Parse()
	if *cfgFile == "" {
		log.Printf("ERROR: No overview file specified.")
		os.Exit(1)
	}
	if invCategories, err = loadCategories(); err != nil {
		log.Printf("ERROR: unable to load inventory categories CSV file: %s", err)
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
	if *updGroups {
		if err = updateGroups(o); err != nil {
			log.Printf("ERROR: unable to update groups/: %s", err)
			os.Exit(1)
		}
		return
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

func updateGroups(o *Overview) error {
	var p *Preset
	for i := range o.Presets {
		if strings.ToLower(o.Presets[i].Name) == allGroupPreset {
			p = o.Presets[i]
			break
		}
	}
	if p == nil {
		return fmt.Errorf("No 'All' preset found")
	}
	// Make a list of inventory group IDs per category.
	cats := make(map[InvCategoryId][]InvGroupId)
	for _, invG := range p.Groups.Groups {
		cat := invGroups[invG].Cat
		cats[cat] = append(cats[cat], invG)
	}
	for cat, invgs := range cats {
		if err := updateGroup(catToFilename(cat), invgs); err != nil {
			return err
		}
	}
	return nil
}

func updateGroup(catFile string, invgids []InvGroupId) error {
	f, err := os.OpenFile(path.Join("groups", catFile), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	for _, invgid := range invgids {
		invg := invGroups[invgid]
		// Do manual marshalling here, as it's just easier. Indent by 8 spaces
		// to allow easy use in creating an overview.
		l := fmt.Sprintf("        - %d # %s -- %s\n", invg.Id, invg.Cat, invg.Name)
		if _, err = f.WriteString(l); err != nil {
			return err
		}
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

func catToFilename(cat InvCategoryId) string {
	n := strings.ToLower(invCategories[cat])
	return strings.Replace(n, " ", "_", -1) + ".yaml"
}
