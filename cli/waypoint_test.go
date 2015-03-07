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
	"errors"
	"flag"
	"testing"
	"time"

	"github.com/rochaporto/ezgliding/config"
	"github.com/rochaporto/ezgliding/mock"
	"github.com/rochaporto/ezgliding/plugin"
	"github.com/rochaporto/ezgliding/waypoint"
)

// ExampleWaypointGet uses the mock waypoint implementation to query data and
// verify waypoint-get works. First, no region is passed. Second, a region but
// no updatedAfter is passed. Finally, both region and updatedAfter are given.
func ExampleWaypointGet() {
	plugin.Register("mockwaypointget", &mock.Mock{
		GetWaypointF: func(regions []string, updatedSince time.Time) ([]waypoint.Waypoint, error) {
			waypoints := []waypoint.Waypoint{
				waypoint.Waypoint{
					ID: "MockID1", Name: "MockName",
					Description: "MockDescription", Region: "FR", Flags: 0,
					Elevation: 2000, Latitude: 32.533, Longitude: 100.376,
					Update: time.Date(2014, 02, 01, 0, 0, 0, 0, time.UTC),
				},
				waypoint.Waypoint{
					ID: "MockID2", Name: "MockName",
					Description: "MockDescription", Region: "CH", Flags: 0,
					Elevation: 2000, Latitude: 32.533, Longitude: 100.376,
					Update: time.Date(2014, 02, 02, 0, 0, 0, 0, time.UTC),
				},
				waypoint.Waypoint{
					ID: "MockID3", Name: "MockName",
					Description: "MockDescription", Region: "CH", Flags: 0,
					Elevation: 2000, Latitude: 32.533, Longitude: 100.376,
					Update: time.Date(2014, 02, 03, 0, 0, 0, 0, time.UTC),
				},
			}
			result := []waypoint.Waypoint{}
			for _, waypoint := range waypoints {
				b := false
				for _, r := range regions {
					if waypoint.Region == r {
						b = true
					}
				}
				if waypoint.Update.After(updatedSince) && b {
					result = append(result, waypoint)
				}
			}
			return result, nil
		},
	},
	)
	config.Set(config.Config{Global: config.Global{Waypointer: "mockwaypointget"}})
	flag.Set("region", "CH")
	flag.Set("after", "2014-02-02")
	runWaypointGet(CmdWaypointGet, []string{})
	// Output:
	// ID,Name,Description,Region,Flags,Elevation,Latitude,Longitude,Update
	// MockID3,MockName,MockDescription,CH,0,2000,32.533,100.376,2014-02-03 00:00:00 +0000 UTC
}

func TestWaypointGetFailed(t *testing.T) {
	plugin.Register("mockwaypointgetfailed", &mock.Mock{
		GetWaypointF: func(regions []string, updatedSince time.Time) ([]waypoint.Waypoint, error) {
			return nil, errors.New("mock testing get waypoint failed")
		},
	},
	)
	config.Set(config.Config{Global: config.Global{Waypointer: "mockwaypointgetfailed"}})
	flag.Set("region", "")
	flag.Set("after", "")
	runWaypointGet(CmdWaypointGet, []string{})
}

func TestWaypointGetBadAfter(t *testing.T) {
	plugin.Register("mockwaypointbadafter", &mock.Mock{
		GetWaypointF: func(regions []string, updatedSince time.Time) ([]waypoint.Waypoint, error) {
			return nil, nil
		},
	},
	)
	config.Set(config.Config{Global: config.Global{Waypointer: "mockwaypointbadafter"}})
	flag.Set("after", "22-00-11")
	flag.Set("region", "")
	runWaypointGet(CmdWaypointGet, []string{})
}

// ExampleWaypointPut uses the mock waypoint implementation to push data and
// verify waypoint-put works.
func ExampleWaypointPut() {
	plugin.Register("mockwaypoint", &mock.Mock{
		GetWaypointF: func(regions []string, updatedSince time.Time) ([]waypoint.Waypoint, error) {
			return []waypoint.Waypoint{
				waypoint.Waypoint{ID: "MockID", Name: "MockName", Description: "MockDescription",
					Region: "FR", Flags: 0, Elevation: 2000, Latitude: 32.533, Longitude: 100.576},
			}, nil
		},
	},
	)
	plugin.Register("mockwaypointput", &mock.Mock{
		PutWaypointF: func(waypoints []waypoint.Waypoint) error {
			return nil
		},
	},
	)
	config.Set(config.Config{Global: config.Global{Waypointer: "mockwaypoint"}})
	runWaypointPut(CmdWaypointPut, []string{"mockwaypointput"})
	// Output:
	// pushed 1 waypoints into mockwaypointput
}

func TestWaypointPutFailed(t *testing.T) {
	plugin.Register("mockwaypointbadput", mock.Mock{
		PutWaypointF: func(waypoints []waypoint.Waypoint) error {
			return errors.New("mock testing put waypoint failed")
		},
	})
	plugin.Register("mockwaypoint", mock.Mock{
		GetWaypointF: func(regions []string, updatedSince time.Time) ([]waypoint.Waypoint, error) {
			return []waypoint.Waypoint{waypoint.Waypoint{}}, nil
		},
	},
	)
	config.Set(config.Config{Global: config.Global{Waypointer: "mockwaypoint"}})
	runWaypointPut(CmdWaypointPut, []string{"mockwaypointbadput"})
}

func TestWaypointPutBadGet(t *testing.T) {
	plugin.Register("mockwaypointbadget", mock.Mock{
		GetWaypointF: func(regions []string, updatedSince time.Time) ([]waypoint.Waypoint, error) {
			return nil, errors.New("mock testing get waypoint failed")
		},
	},
	)
	config.Set(config.Config{Global: config.Global{Waypointer: "mockwaypointbadget"}})
	runWaypointPut(CmdWaypointPut, []string{"mockwaypointbadget"})
}

func TestWaypointPutBadPluginID(t *testing.T) {
	runWaypointPut(CmdWaypointPut, []string{"wpnonexisting"})
}

func TestWaypointPutBadArgNumber(t *testing.T) {
	runWaypointPut(CmdWaypointPut, []string{})
}
