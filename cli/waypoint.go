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

	commander "code.google.com/p/go-commander"
)

// CmdWaypointGet command gets waypoint information and outputs the result.
var CmdWaypointGet = &commander.Command{
	UsageLine: "waypoint-get [options]",
	Short:     "gets waypoint information",
	Long: `
Gets waypoint information according to the given parameters.
` + "\n" + helpFlags(flag.CommandLine),
	Run: func(cmd *commander.Command, args []string) {
	},
	Flag: *flag.CommandLine,
}
