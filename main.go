// Copyright 2017 Joel Scoble
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy
// of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
//

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	app = filepath.Base(os.Args[0]) // name of application
	cfg Config
)

type Config struct {
	Length  int      // line length
	Ext     Strings  // extension to filter on
	Exclude bool     // exclude the specified extension
	Include bool     // include the specified extension
	LogFile string   // output destination for logs; stderr is default
	f       *os.File // logfile handle for close; this will be nil if output is stderr
}

type Strings []string // a type to add support for string arrays to a flag arg

// implement the flag.Value iinterface
func (s *Strings) String() string {
	return fmt.Sprintf("%s")
}

func (s *Strings) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// Exists returns if the string is in the slice
func (s *Strings) Exists(v string) bool {
	for _, val := range *s {
		if v == val {
			return true
		}
	}
	return false
}

func init() {
	flag.IntVar(&cfg.Length, "length", 80, "line length")
	flag.Var(&cfg.Ext, "ext", "extension to filter on")
	flag.BoolVar(&cfg.Exclude, "exclude", false, "exclude the extensions")
	flag.BoolVar(&cfg.Include, "include", false, "include the extensions")
	flag.StringVar(&cfg.LogFile, "logfile", "stderr", "output destination for logs")
	log.SetPrefix(app + ": ")
}

func main() {
	flag.Usage = usage

	// Process flags
	FlagParse()
	paths := flag.Args()
	os.Exit(wrapperMain(paths))
}
