// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

// Package clienv does very simple fallback from flag to env var and then to flag default
//
//	Only arrays of strings are supported as output
package clienv

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Ix int // index into a Flagvars map

// Flagvar tracks which flags are set and which env vars they match with
type Flagvar struct {
	Name    string
	Default string
	Use     string
	Env     string
	RawVar  *string
	FinVar  *[]string
	present bool
}

type Flagvars map[Ix]Flagvar

func (fls *Flagvars) Register() {
	for _, fv := range *fls {
		flag.StringVar(fv.RawVar, fv.Name, fv.Default,
			fmt.Sprintf("%s, if not set reads %s envvar", fv.Use, fv.Env))
	}
}

func (fls *Flagvars) Fallback() {

	fname2ix := map[string]Ix{}
	for ix, v := range *fls {
		fname2ix[v.Name] = ix
	}

	var vf = func(f *flag.Flag) {
		ix, ok := fname2ix[f.Name]
		if !ok { // not one we care about
			return
		}
		v, _ := (*fls)[ix]
		*v.FinVar = strings.Fields(*v.RawVar)
		v.present = true
		(*fls)[ix] = v
		return
	}

	flag.Visit(vf)

	for _, v := range *fls {
		if !v.present {
			env := os.Getenv(v.Env)
			if len(env) > 0 {
				if strings.ToLower(env) == "none" {
					*v.FinVar = []string{}
					continue
				}
				*v.FinVar = strings.Fields(env) // env var value
			} else {
				*v.FinVar = strings.Fields(*v.RawVar) // flag default
			}
		}
	}

}

func (fls *Flagvars) String() string {
	s := ""
	for _, v := range *fls {
		s += fmt.Sprintf("# flag: %s, final: '%s' (cli: %t, env: '%s', cliraw: '%s')\n",
			v.Name, *v.FinVar, v.present, os.Getenv(v.Env), *v.RawVar)
	}
	return s
}
