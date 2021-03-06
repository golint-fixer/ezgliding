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

// Package welt2000 provides functionality to fetch and parse airfield and
// waypoint information, taking the welt release as input.
//
// Check the welt2000 website for more information on the data:
// 	http://www.segelflug.de/vereine/welt2000/
//
package welt2000

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/rochaporto/ezgliding/airfield"
	"github.com/rochaporto/ezgliding/spatial"
	"github.com/rochaporto/ezgliding/waypoint"
	"github.com/rochaporto/rss"
)

// ID for this plugin implementation.
const (
	ID string = "welt2000"
)

// Release contains info about a specific release
type Release struct {
	Date      time.Time
	Source    string
	Airfields []airfield.Airfield
	Waypoints []waypoint.Waypoint
}

const (
	// RSSURL is the default location of the welt2000 RSS feed
	RSSURL = "http://www.segelflug.de/vereine/welt2000/content/en/news/updates.xml"
	// ReleaseURL is the default location for the welt2000 latest release
	ReleaseURL = "http://www.segelflug.de/vereine/welt2000/download/WELT2000.TXT"
)

// Config holds all config information for the welt2000 plugin.
type Config struct {
	RSSURL     string
	ReleaseURL string
}

// Welt2000 is the plugin implementation to collect welt2000 data,
// for airfields and waypoints.
type Welt2000 struct {
	Config
}

// New returns a new instance of Welt2000.
func New(cfg Config) (*Welt2000, error) {
	if cfg.RSSURL == "" {
		cfg.RSSURL = RSSURL
	}
	if cfg.ReleaseURL == "" {
		cfg.ReleaseURL = ReleaseURL
	}
	wt := Welt2000{Config: cfg}
	glog.V(20).Infof("Plugin welt2000 initialized :: %+v", wt)
	return &wt, nil
}

// GetAirfield follows airfield.GetAirfield().
// FIXME: use region
func (wt *Welt2000) GetAirfield(regions []string, updatedSince time.Time) ([]airfield.Airfield, error) {
	glog.V(10).Infof("GetAirfield with regions %v and updatedSince %v", regions, updatedSince)
	releases, err := List(wt.RSSURL)
	if err != nil {
		return nil, err
	}
	release := releases[0]
	if !release.Date.After(updatedSince) {
		return release.Airfields, nil
	}

	// We 'update' the release source as the rss feed points to an update
	// summary page, not the actually release source.
	release.Source = wt.ReleaseURL
	err = release.Fetch()
	// Filter out entries not in the regions.
	// This could be done more efficiently, but for now we go with post-filter.
	m := make(map[string]bool) // Use a map as there's no slice.contains(str)
	for _, v := range regions {
		m[v] = true
	}
	var filtered []airfield.Airfield
	for _, a := range release.Airfields {
		if _, ok := m[a.Region]; ok {
			filtered = append(filtered, a)
		}
	}
	release.Airfields = filtered
	glog.V(10).Infof("GetAirfield for regions %v and updatedsince %v retrieved %d results",
		regions, updatedSince, len(release.Airfields))
	glog.V(20).Infof("%v", release.Airfields)
	return release.Airfields, err
}

// PutAirfield follows airfield.PutAirfield().
func (wt *Welt2000) PutAirfield(airfields []airfield.Airfield) error {
	return errors.New("not available for welt2000 plugin")
}

// GetWaypoint follows airfield.GetWaypoint().
// FIXME: use region
func (wt *Welt2000) GetWaypoint(regions []string, updatedSince time.Time) ([]waypoint.Waypoint, error) {
	glog.V(10).Infof("GetWaypoint with regions %v and updatedSince %v", regions, updatedSince)
	releases, err := List(wt.RSSURL)
	if err != nil {
		return nil, err
	}
	release := releases[0]
	if !release.Date.After(updatedSince) {
		return release.Waypoints, nil
	}

	// We 'update' the release source as the rss feed points to an update
	// summary page, not the actually release source.
	release.Source = wt.ReleaseURL
	err = release.Fetch()
	// Filter out entries not in the regions.
	// This could be done more efficiently, but for now we go with post-filter.
	m := make(map[string]bool) // Use a map as there's no slice.contains(str)
	for _, v := range regions {
		m[v] = true
	}
	var filtered []waypoint.Waypoint
	for _, a := range release.Waypoints {
		if _, ok := m[a.Region]; ok {
			filtered = append(filtered, a)
		}
	}
	release.Waypoints = filtered
	glog.V(10).Infof("GetWaypoint for regions %v and updatedsince %v retrieved %d results",
		regions, updatedSince, len(release.Waypoints))
	glog.V(20).Infof("%v", release.Waypoints)
	return release.Waypoints, err
}

// PutWaypoint follows waypoint.PutWaypoint().
func (wt *Welt2000) PutWaypoint(waypoints []waypoint.Waypoint) error {
	return errors.New("not available for welt2000 plugin")
}

// List checks the welt2000 rss feed and lists the releases found
func List(location string) ([]Release, error) {
	var content []byte

	glog.V(10).Infof("List for location %v", location)
	// case http
	resp, err := http.Get(location)
	if err == nil {
		defer resp.Body.Close()
		content, err = ioutil.ReadAll(resp.Body)
	} else { // case file
		resp, err := ioutil.ReadFile(location)
		if err != nil {
			return nil, err
		}
		content = resp
	}
	rss.Init()
	feed, err := rss.Parse(content)
	if err != nil {
		return nil, err
	}

	res := make([]Release, 10)
	for i, item := range feed.Items {
		res[i].Date = item.Date
		res[i].Source = item.Link
	}
	glog.V(10).Infof("List got %v releases", len(res))
	glog.V(20).Infof("%v", res)
	return res, nil
}

// Fetch gets and returns the Release at the given location
func Fetch(location string) (*Release, error) {
	glog.V(10).Infof("Fetch with location %v", location)
	r := Release{Source: location}
	err := r.Fetch()
	return &r, err
}

// Fetch fills up the Release object with data after parsing the content at Release.Source
func (r *Release) Fetch() error {
	glog.V(10).Infof("Release fetch :: %+v", r)
	var content []byte

	resp, err := http.Get(r.Source)
	// case http
	if err == nil {
		defer resp.Body.Close()
		content, err = ioutil.ReadAll(resp.Body)
	} else { // case file
		content, err = ioutil.ReadFile(r.Source)
	}
	if err != nil {
		return err
	}
	return r.Parse(content)
}

// Parse fills in the Release object by parsing r.Data
func (r *Release) Parse(content []byte) error {
	if content == nil || len(content) == 0 {
		return errors.New("No data available to parse")
	}

	lines := strings.Split(string(content), "\n")
	for i := range lines {
		switch {
		case len(lines[i]) == 0: // empty line
			continue
		case lines[i][0] == '$': // comment
			continue
		case lines[i][5] == '1' || lines[i][5] == '2': // airfield
			r.parseAirfield(lines[i])
		default: // waypoint
			r.parseWaypoint(lines[i])
		}
	}
	return nil
}

func (r *Release) parseAirfield(line string) error {
	afield := airfield.Airfield{Update: r.Date}
	if line[4] == '2' { // unclear airstrip
		afield.Flags |= airfield.UnclearAirstrip
		afield.ShortName = strings.Trim(line[0:4], " ")
	} else { // regular airstrip
		afield.ShortName = strings.Trim(line[0:5], " ")
	}
	afield.Name = strings.Trim(line[7:20], " ")
	if line[23] == '#' && line[24] != ' ' && string(line[24:28]) != "GLD!" { // ICAO available
		afield.ICAO = line[24:28]
		afield.ID = afield.ICAO
	} else {
		afield.ID = afield.ShortName
	}
	if line[23:27] == "*ULM" { // ultralight site
		afield.Flags |= airfield.ULMSite
	} else if line[5] == '2' { // outlanding
		afield.Flags |= airfield.Outlanding
		afield.Catalog, _ = strconv.Atoi(line[26:28])
	} else if line[20:24] == "GLD#" || line[23:28] == "#GLD" { // glider site
		afield.Flags |= airfield.GliderSite
	}
	afield.Flags |= r.runwayType2Bit(line[28])
	afield.Length, _ = strconv.Atoi(line[29:32])
	afield.Length *= 10
	afield.Runway = line[32:36]
	decimal, _ := strconv.ParseFloat(line[39:41], 64)
	afield.Frequency, _ = strconv.ParseFloat(line[36:39], 64)
	afield.Frequency += decimal * 0.01
	elevation := strings.Trim(line[41:45], " ")
	afield.Elevation, _ = strconv.Atoi(elevation)
	afield.Latitude = spatial.DMS2Decimal(line[45:52])
	afield.Longitude = spatial.DMS2Decimal(line[52:60])
	afield.Region = line[60:62]
	r.Airfields = append(r.Airfields, afield)
	return nil
}

func (r *Release) runwayType2Bit(t uint8) int {
	switch t {
	case 'A':
		return airfield.Asphalt
	case 'C':
		return airfield.Concrete
	case 'L':
		return airfield.Loam
	case 'S':
		return airfield.Sand
	case 'Y':
		return airfield.Clay
	case 'G':
		return airfield.Grass
	case 'V':
		return airfield.Gravel
	case 'D':
		return airfield.Dirt
	}
	return 0
}

func (r *Release) parseWaypoint(line string) error {
	waypoint := waypoint.Waypoint{
		Name: strings.Trim(line[0:6], " "), ID: strings.Trim(line[0:6], " "),
		Description: strings.Trim(line[7:41], " "),
		Latitude:    spatial.DMS2Decimal(line[45:52]), Longitude: spatial.DMS2Decimal(line[52:60]),
		Region: line[60:62], Update: r.Date,
	}
	elevation := strings.Trim(line[41:45], " ")
	waypoint.Elevation, _ = strconv.Atoi(elevation)
	r.Waypoints = append(r.Waypoints, waypoint)
	return nil
}
