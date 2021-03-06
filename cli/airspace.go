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
	"strings"
	"time"

	commander "code.google.com/p/go-commander"
	"github.com/golang/glog"
	"github.com/rochaporto/ezgliding/config"
	"github.com/rochaporto/ezgliding/plugin"
)

// CmdGetAirspace command gets airspace information.
var CmdAirspaceGet = &commander.Command{
	UsageLine: "airspace-get [options]",
	Short:     "gets airspace information",
	Long: `
Gets latest airspace data corresponding to the given parameters.

Example:
  ezgliding airspace-get --region=FR
` + "\n" + helpFlags(flag.CommandLine),
	Run:  runAirspaceGet,
	Flag: *flag.CommandLine,
}

// runAirspaceGet invokes the configured plugin and outputs airspace data.
func runAirspaceGet(cmd *commander.Command, args []string) {
	var err error
	cfg, _ := config.Get()
	aspace, err := plugin.GetAirspacer("", cfg)
	if err != nil {
		glog.Errorf("failed to get airspacer plugin :: %v\n", err)
		return
	}
	airspaces, err := aspace.GetAirspace(strings.Split(*region, ","), time.Time{})
	if err != nil {
		glog.Errorf("failed to get airspace :: %v", err)
		// FIXME: must return -1, but no way now to check this in test
	}
	glog.V(5).Infof("airspace get with args '%v' got %d results", args, len(airspaces))
	glog.V(20).Infof("%+v", airspaces)
	for i := range airspaces {
		fmt.Printf("%+v\n", airspaces[i])
	}
}
