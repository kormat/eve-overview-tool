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
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	invGroupPath    = "data/invGroups.csv"
	invGroupPathBz2 = invGroupPath + ".bz2"
	invGroupPathUrl = "https://www.fuzzwork.co.uk/dump/latest/invGroups.csv.bz2"
	filterStatePath = "data/filterStates.csv"
)

var groupsFile = flag.String("groups", "", fmt.Sprintf("Use external inventory groups CSV file."))
var stateFile = flag.String("states", "", "Use external filter states CSV file")

func loadGroups() (map[InvGroup]string, error) {
	var err error
	var bz2 bool
	var reader io.Reader
	if *groupsFile != "" {
		var f *os.File
		if f, err = os.Open(*groupsFile); err != nil {
			return nil, err
		}
		defer f.Close()
		reader = f
		bz2 = strings.HasSuffix(strings.ToLower(*groupsFile), "bz2")
	} else {
		data, err := Asset(invGroupPath)
		if err != nil {
			data, err = Asset(invGroupPathBz2)
			bz2 = true
		}
		if err != nil {
			log.Printf("Download the inventory groups file from %s", invGroupPathUrl)
			return nil, err
		}
		reader = bytes.NewReader(data)
	}
	if bz2 {
		reader = bzip2.NewReader(reader)
	}
	entries, err := loadCsvEntries(reader, 2)
	if err != nil {
		return nil, err
	}
	m := make(map[InvGroup]string, len(entries))
	for i := range entries {
		entry := &entries[i]
		m[InvGroup(entry.i)] = entry.s
	}
	return m, nil
}

func loadStates() (map[StateType]string, error) {
	var reader io.Reader
	if *stateFile != "" {
		f, err := os.Open(*stateFile)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		reader = f
	} else {
		data, err := Asset(filterStatePath)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(data)
	}
	entries, err := loadCsvEntries(reader, 1)
	if err != nil {
		return nil, err
	}
	m := make(map[StateType]string, len(entries))
	for i := range entries {
		entry := &entries[i]
		m[StateType(entry.i)] = entry.s
	}
	return m, nil
}

func loadCsvEntries(r io.Reader, nameIdx int) ([]csvEntry, error) {
	var entries []csvEntry
	csvr := csv.NewReader(r)
	csvr.FieldsPerRecord = nameIdx + 1
	csvr.LazyQuotes = true
	for {
		record, err := csvr.Read()
		if err == io.EOF {
			break
		}
		n, err := strconv.Atoi(record[0])
		if err != nil {
			if len(stateTypes) == 0 {
				// Skip the header line, if present.
				continue
			}
			return nil, err
		}
		entries = append(entries, csvEntry{n, record[nameIdx]})
	}
	return entries, nil
}

type csvEntry struct {
	i int
	s string
}
