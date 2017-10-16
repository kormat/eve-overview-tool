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
	invGroupPath    = "invGroups.csv"
	invGroupPathBz2 = invGroupPath + ".bz2"
	invGroupPathUrl = "https://www.fuzzwork.co.uk/dump/latest/invGroups.csv.bz2"
)

var typesFile = flag.String("types", "", fmt.Sprintf("Inventory groups CSV file. Will try "+
	"`%s` and `%s.bz2` in order if not specified", invGroupPath, invGroupPathBz2))
var stateFile = flag.String("states", "filterStates.csv", "Filter states CSV file")

func loadGroups() (map[InvGroup]string, error) {
	var f *os.File
	var err error
	var bz2 bool
	if *typesFile != "" {
		if f, err = os.Open(*typesFile); err != nil {
			return nil, err
		}
		bz2 = strings.HasSuffix(strings.ToLower(*typesFile), "bz2")
	} else {
		f, err = os.Open(invGroupPath)
		if os.IsNotExist(err) {
			f, err = os.Open(invGroupPathBz2)
			bz2 = true
		}
		if os.IsNotExist(err) {
			log.Printf("Download the inventory groups file from %s", invGroupPathUrl)
			return nil, err
		}
		if err != nil {
			return nil, err
		}
	}
	defer f.Close()
	var reader io.Reader = f
	if bz2 {
		reader = bzip2.NewReader(f)
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
	f, err := os.Open(*stateFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	entries, err := loadCsvEntries(f, 1)
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
