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
	"strings"
	"testing"

	"github.com/rochaporto/ezgliding/context"
)

// FIXME: should go away when we start passing the context explicitly
// to the runAirfield* functions.
func setupContext(ctx context.Context) {
	context.Ctx = ctx
}

func TestHelpFlags(t *testing.T) {
	e := "param=defaultvalue"

	var tf = flag.CommandLine
	var _ = tf.String("param", "defaultvalue", "message")
	r := helpFlags(tf)
	if !strings.Contains(r, e) {
		t.Errorf("Missing '%v'\nin help message\n'%v'", e, r)
	}
}
