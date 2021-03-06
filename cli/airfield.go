// Copyright 2014 The ezgliding Authors.
//
// This file is part of ezgliding.
//
// ezgliding is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// ezgliding is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with ezgliding.  If not, see <http://www.gnu.org/licenses/>.
//
// Author: Ricardo Rocha <rocha.porto@gmail.com>

package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	commander "code.google.com/p/go-commander"
	"github.com/golang/glog"
	"github.com/rochaporto/ezgliding/airfield"
	"github.com/rochaporto/ezgliding/config"
	"github.com/rochaporto/ezgliding/plugin"
	"github.com/rochaporto/ezgliding/util"
)

// CmdAirfieldGet command gets airfield information and outputs the result.
var CmdAirfieldGet = &commander.Command{
	UsageLine: "airfield-get [options]",
	Short:     "gets airfield information",
	Long: `
Gets available airfield information according to the given parameters
` + "\n" + helpFlags(flag.CommandLine),
	Run:  runAirfieldGet,
	Flag: *flag.CommandLine,
}

// runAirfieldGet invokes the configured plugin and outputs airfield data.
func runAirfieldGet(cmd *commander.Command, args []string) {
	var err error
	cfg, _ := config.Get()
	afield, err := plugin.GetAirfielder("", cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get airfield plugin :: %v\n", err)
		return
	}
	tafter := time.Time{}
	if *after != "" {
		tafter, err = time.Parse("2006-01-02", *after)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get airfield :: %v\n", err)
			return
		}
	}
	airfields, err := afield.(airfield.Airfielder).GetAirfield(strings.Split(*region, ","), tafter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get airfield :: %v\n", err)
		return
		// FIXME: must return -1, but no way now to check this in test
	}
	glog.V(5).Infof("airfield get with args '%v' got %d results", args, len(airfields))
	glog.V(20).Infof("%+v", airfields)
	fmt.Printf("%v", util.Struct2CSV(airfields))
}

// CmdAirfieldPut command puts airfield information from a source to a destination.
var CmdAirfieldPut = &commander.Command{
	UsageLine: "airfield-put [options] destination",
	Short:     "puts airfield information",
	Long: `
Puts airfield information according to the given parameters
` + "\n" + helpFlags(flag.CommandLine),
	Run:  runAirfieldPut,
	Flag: *flag.CommandLine,
}

// runAirfieldPut invokes the configured plugins to put airfield data from source to dest.
func runAirfieldPut(cmd *commander.Command, args []string) {
	var err error
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "failed to put airfield data :: no destination given\n")
		return
	}
	cfg, _ := config.Get()
	pluginID := args[0]
	destPlugin, err := plugin.GetAirfielder(pluginID, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get plugin '%v' :: %v\n", pluginID, err)
		return
	}
	afield, _ := plugin.GetAirfielder("", cfg)
	airfields, err := afield.(airfield.Airfielder).GetAirfield(strings.Split(*region, ","), time.Time{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get airfield :: %v\n", err)
		return
	}
	glog.V(5).Infof("putting %v airfields", len(airfields))
	glog.V(20).Infof("%v", airfields)
	if len(airfields) > 0 {
		err = destPlugin.(airfield.Airfielder).PutAirfield(airfields)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to put airfields :: %v\n", err)
			return
		}
	}
	fmt.Printf("pushed %v airfields into %v\n", len(airfields), pluginID)
}
