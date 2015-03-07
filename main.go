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

package main

import (
	"os"

	commander "code.google.com/p/go-commander"
	"github.com/golang/glog"
	"github.com/rochaporto/ezgliding/cli"
)

func exit(c int) {
	glog.Flush()
	os.Exit(c)
}

func main() {
	defer glog.Flush()

	c := commander.Commander{
		Name: "ezgliding",
		Commands: []*commander.Command{
			cli.CmdAirfieldGet,
			cli.CmdAirfieldPut,
			cli.CmdAirspaceGet,
			cli.CmdFlightGet,
			cli.CmdWaypointGet,
			cli.CmdWaypointPut,
			cli.CmdWeb,
		},
	}

	if len(os.Args) == 1 {
		os.Args = append(os.Args, "help")
	}
	if err := c.Run(os.Args[1:]); err != nil {
		glog.Errorf("Failed running command %q: %v", os.Args[1:], err)
		exit(-1)
	}
}
