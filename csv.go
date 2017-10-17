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
	"bytes"
	"compress/bzip2"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	invGroupPath    = "data/invGroups.csv"
	invGroupPathUrl = "https://www.fuzzwork.co.uk/dump/latest/invGroups.csv.bz2"
	invCatPath      = "data/invCategories.csv"
	invCatPathUrl   = "https://www.fuzzwork.co.uk/dump/latest/invCategories.csv.bz2"
	filterStatePath = "data/filterStates.csv"
)

var catFile = flag.String("categories", "",
	fmt.Sprintf("Use external inventory categories CSV file."))
var groupsFile = flag.String("groups", "", fmt.Sprintf("Use external inventory groups CSV file."))
var stateFile = flag.String("states", "", "Use external filter states CSV file")

type readerCloser struct {
	*bytes.Reader
}

func (cc readerCloser) Close() error {
	return nil
}

func loadFile(path, assetPath string) (io.Reader, error) {
	var err error
	var bz2 bool
	var data []byte
	if path != "" {
		if data, err = ioutil.ReadFile(path); err != nil {
			return nil, err
		}
		bz2 = strings.HasSuffix(strings.ToLower(path), "bz2")
	} else {
		if data, err = Asset(assetPath); err != nil {
			data, err = Asset(assetPath + ".bz2")
			if err != nil {
				return nil, err
			}
			bz2 = true
		}
	}
	reader := bytes.NewReader(data)
	if bz2 {
		return bzip2.NewReader(reader), nil
	}
	return reader, nil
}

type InvCategoryId int

func (ic InvCategoryId) name() string {
	s, ok := invCategories[ic]
	if !ok {
		return "Unknown InvCategory"
	}
	return strings.TrimSpace(s)
}

func (ic InvCategoryId) String() string {
	return fmt.Sprintf("%s (%d)", ic.name(), int(ic))
}

func loadCategories() (map[InvCategoryId]string, error) {
	reader, err := loadFile(*catFile, invCatPath)
	if err != nil {
		return nil, fmt.Errorf("Unable to load inventory categories CSV file: %v", err)
	}
	records, err := loadCsvEntries(reader, 4)
	if err != nil {
		return nil, err
	}
	m := make(map[InvCategoryId]string, len(records))
	for _, record := range records {
		id, err := strconv.Atoi(record[0])
		if err != nil {
			if len(m) == 0 {
				// Skip the header line, if present.
				continue
			}
			return nil, err
		}
		m[InvCategoryId(id)] = record[1]
	}
	return m, nil
}

type InvGroupId int

func (ig InvGroupId) name() string {
	g, ok := invGroups[ig]
	if !ok {
		return "Unknown InvGroup"
	}
	return fmt.Sprintf("%s -- %s", g.Cat, strings.TrimSpace(g.Name))
}

func (ig InvGroupId) String() string {
	return fmt.Sprintf("%s (%d)", ig.name(), int(ig))
}

func (ig InvGroupId) MarshalYAML() (interface{}, error) {
	return fmt.Sprintf("%d %s %s", int(ig), commentMarker, ig.name()), nil
}

type InvGroup struct {
	Id   InvGroupId
	Cat  InvCategoryId
	Name string
}

func loadGroups() (map[InvGroupId]*InvGroup, error) {
	reader, err := loadFile(*groupsFile, invGroupPath)
	if err != nil {
		return nil, fmt.Errorf("Unable to load inventory groups CSV file: %v", err)
	}
	records, err := loadCsvEntries(reader, 9)
	if err != nil {
		return nil, err
	}
	m := make(map[InvGroupId]*InvGroup, len(records))
	for _, record := range records {
		id, err := strconv.Atoi(record[0])
		if err != nil {
			if len(m) == 0 {
				// Skip the header line, if present.
				continue
			}
			return nil, err
		}
		catId, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, err
		}
		g := &InvGroup{Id: InvGroupId(id), Cat: InvCategoryId(catId), Name: record[2]}
		m[g.Id] = g
	}
	return m, nil
}

func loadStates() (map[StateType]string, error) {
	reader, err := loadFile(*stateFile, filterStatePath)
	if err != nil {
		return nil, fmt.Errorf("Unable to load filter states CSV file: %v", err)
	}
	records, err := loadCsvEntries(reader, 2)
	if err != nil {
		return nil, err
	}
	m := make(map[StateType]string, len(records))
	for _, record := range records {
		id, err := strconv.Atoi(record[0])
		if err != nil {
			if len(m) == 0 {
				// Skip the header line, if present.
				continue
			}
			return nil, err
		}
		m[StateType(id)] = record[1]
	}
	return m, nil
}

func loadCsvEntries(r io.Reader, nFields int) ([][]string, error) {
	csvr := csv.NewReader(r)
	csvr.FieldsPerRecord = nFields
	return csvr.ReadAll()
}
