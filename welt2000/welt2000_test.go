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

package welt2000

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/rochaporto/ezgliding/airfield"
	"github.com/rochaporto/ezgliding/waypoint"
)

var GMT, _ = time.LoadLocation("GMT")

type GetAirfieldTest struct {
	t   string
	r   string
	rss string
	rg  string
	d   time.Time
	rs  []airfield.Airfield
	err bool
}

var getAirfieldTests = []GetAirfieldTest{
	{"basic get airfield",
		"./t/test-release-basic.txt",
		"./t/test-releases-list.xml",
		"FR",
		time.Time{},
		[]airfield.Airfield{
			airfield.Airfield{
				ID: "HABER", Name: "HABERE POC69", ShortName: "HABER", Region: "FR",
				ICAO: "", Flags: airfield.GliderSite | airfield.Grass, Catalog: 0,
				Length: 0, Elevation: 1113, Runway: "0119", Frequency: 122.5,
				Latitude: 46.26972222222222, Longitude: 6.463333333333334,
				Update: time.Date(2014, time.February, 24, 12, 0, 0, 0, GMT),
			},
		},
		false,
	},
	{"get airfield with updated since",
		"./t/test-release-basic.txt",
		"./t/test-releases-list.xml",
		"FR",
		time.Date(2014, time.February, 25, 0, 0, 0, 0, time.UTC),
		[]airfield.Airfield{},
		false,
	},
	{"get airfield missing rss",
		"./t/test-release-basic.txt",
		"./t/missing-release-list.xml",
		"FR",
		time.Time{},
		[]airfield.Airfield{},
		true,
	},
	{"get airfield missing release",
		"./t/missing-release.txt",
		"./t/test-releases-list.xml",
		"FR",
		time.Time{},
		[]airfield.Airfield{},
		true,
	},
	{"get airfield with 0 values for region",
		"./t/test-release-basic.txt",
		"./t/test-releases-list.xml",
		"ZZ",
		time.Time{},
		[]airfield.Airfield{},
		false,
	},
}

func TestGetAirfield(t *testing.T) {
	for _, test := range getAirfieldTests {

		cfg := Config{}
		cfg.RSSURL = test.rss
		cfg.ReleaseURL = test.r
		plugin, err := New(cfg)
		if err != nil {
			t.Errorf("failed to initialize plugin :: %v", err)
			continue
		}

		var airfields []airfield.Airfield
		airfields, err = plugin.GetAirfield([]string{test.rg}, test.d)
		if err != nil && test.err {
			continue
		} else if err != nil {
			t.Errorf("failed to get airfield :: %v", err)
			continue
		}

		if len(airfields) != len(test.rs) {
			t.Errorf("got %v airfields but expected %v in test '%v'", len(airfields), len(test.rs), test.t)
			continue
		}

		for i := range airfields {
			var airfield = airfields[i]
			var expected = test.rs[i]
			if !reflect.DeepEqual(airfield, expected) {
				t.Errorf("expected %v got %v", expected, airfield)
				continue
			}
		}
	}
}

func TestPutAirfield(t *testing.T) {
	plugin, err := New(Config{})
	if err != nil {
		t.Errorf("failed to initialize plugin :: %v", err)
		return
	}

	err = plugin.PutAirfield(nil)
	if err == nil {
		t.Errorf("PutAirfield should return error for welt2000")
	}
}

type GetWaypointTest struct {
	t   string
	r   string
	rss string
	rg  string
	d   time.Time
	rs  []waypoint.Waypoint
	err bool
}

var getWaypointTests = []GetWaypointTest{
	{"basic get waypoint",
		"./t/test-release-basic.txt",
		"./t/test-releases-list.xml",
		"CH",
		time.Time{},
		[]waypoint.Waypoint{
			waypoint.Waypoint{
				ID: "FURKAP", Name: "FURKAP", Description: "FURKAPASS PASSHOEHE", Elevation: 2432,
				Latitude: 46.57277777777778, Longitude: 8.415277777777778, Region: "CH",
				Update: time.Date(2014, time.February, 24, 12, 0, 0, 0, GMT),
			},
		},
		false,
	},
	{"get waypoint with updated since",
		"./t/test-release-basic.txt",
		"./t/test-releases-list.xml",
		"CH",
		time.Date(2014, time.February, 25, 0, 0, 0, 0, time.UTC),
		[]waypoint.Waypoint{},
		false,
	},
	{"get waypoint missing rss",
		"./t/test-release-basic.txt",
		"./t/missing-release-list.xml",
		"CH",
		time.Time{},
		[]waypoint.Waypoint{},
		true,
	},
	{"get waypoint missing release",
		"./t/missing-release.txt",
		"./t/test-releases-list.xml",
		"CH",
		time.Time{},
		[]waypoint.Waypoint{},
		true,
	},
	{"get waypoint with 0 values for region",
		"./t/test-release-basic.txt",
		"./t/test-releases-list.xml",
		"ZZ",
		time.Time{},
		[]waypoint.Waypoint{},
		false,
	},
}

func TestGetWaypoint(t *testing.T) {
	for i := range getWaypointTests {
		test := getWaypointTests[i]

		cfg := Config{}
		cfg.RSSURL = test.rss
		cfg.ReleaseURL = test.r
		plugin, err := New(cfg)
		if err != nil {
			t.Errorf("failed to initialize plugin :: %v", err)
		}

		var waypoints []waypoint.Waypoint
		waypoints, err = plugin.GetWaypoint([]string{test.rg}, test.d)
		if err != nil && test.err {
			continue
		} else if err != nil {
			t.Errorf("failed to get waypoint :: %v", err)
			continue
		}

		if len(waypoints) != len(test.rs) {
			t.Errorf("got %v waypoints but expected %v in test '%v'", len(waypoints), len(test.rs), test.t)
			continue
		}

		for i := range waypoints {
			var waypoint = waypoints[i]
			var expected = test.rs[i]
			if !reflect.DeepEqual(waypoint, expected) {
				t.Errorf("got wrong waypoint. %v instead of %v", waypoint, expected)
				continue
			}
		}
	}
}

func TestPutWaypoint(t *testing.T) {
	cfg := Config{}
	plugin, err := New(cfg)
	if err != nil {
		t.Errorf("failed to initialize plugin :: %v", err)
	}

	err = plugin.PutWaypoint(nil)
	if err == nil {
		t.Errorf("put waypoint should returned error for welt2000")
	}
}

func TestListLocal(t *testing.T) {
	releases, err := List("./t/test-releases-list.xml")
	if err != nil {
		t.Errorf("failed to list releases :: %v", err)
	}
	if len(releases) < 1 {
		t.Errorf("got wrong number of releases :: %v", len(releases))
	}
}

func TestListHTTP(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, _ := ioutil.ReadFile("./t/test-releases-list.xml")
		io.WriteString(w, string(resp))
	}))
	defer ts.Close()

	releases, err := List(ts.URL)
	if err != nil {
		t.Errorf("failed to list releases from http endpoint :: %v", err)
	}
	if len(releases) < 1 {
		t.Errorf("got wrong number of releases :: %v", len(releases))
	}
}

func TestListEmpty(t *testing.T) {
	_, err := List("")
	if err == nil {
		t.Errorf("list empty string should give error")
	}
}

func TestListMissing(t *testing.T) {
	_, err := List("./nonexisting.file")
	if err == nil {
		t.Errorf("list non existing file should give error")
	}
}

func TestListBrokenFeed(t *testing.T) {
	_, err := List("./t/test-releases-broken.xml")
	if err == nil {
		t.Errorf("parsing a broken rss feed should have failed")
	}
}

func TestFetchLocal(t *testing.T) {
	release, err := Fetch("./t/test-release-basic.txt")
	if err != nil {
		t.Errorf("failed to fetch release from local :: %v", err)
	}
	if len(release.Airfields) < 1 || len(release.Waypoints) < 1 {
		t.Errorf("got wrong number of airfields (%v) or waypoints (%v)", len(release.Airfields), len(release.Waypoints))
	}
}

func TestFetchHTTP(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, _ := ioutil.ReadFile("./t/test-release-basic.txt")
		io.WriteString(w, string(resp))
	}))
	defer ts.Close()

	release, err := Fetch(ts.URL)
	if err != nil {
		t.Errorf("failed to fetch release from http endpoint :: %v", err)
	}
	if len(release.Airfields) < 1 || len(release.Waypoints) < 1 {
		t.Errorf("got wrong number of airfields or waypoints :: %v :: %v", len(release.Airfields), len(release.Waypoints))
	}
}

func TestFetchEmpty(t *testing.T) {
	_, err := Fetch("")
	if err == nil {
		t.Errorf("fetching an empty string should return error")
	}
}

func TestFetchMissing(t *testing.T) {
	_, err := Fetch("nonexisting.release")
	if err == nil {
		t.Errorf("fetching a non existing release should return error")
	}
}

func TestParseNil(t *testing.T) {
	r := Release{}
	err := r.Parse(nil)
	if err == nil {
		t.Errorf("parsing a nil value should return error")
	}
}

func TestParseEmpty(t *testing.T) {
	r := Release{}
	err := r.Parse([]byte{})
	if err == nil {
		t.Errorf("parsing an empty value should return an error")
	}
}

func TestParseComment(t *testing.T) {
	r := Release{}
	r.Parse([]byte("$ this is a comment line"))
	if len(r.Airfields) > 0 || len(r.Waypoints) > 0 {
		t.Errorf("parsing a comment line should fill up airfields or waypoints")
	}
}

func TestParseAirfield(t *testing.T) {
	r := Release{}
	r.Parse([]byte("ANNEM1 ANNEMASSE       #LFLIA129123012587 494N461131E0061606FRQ0"))
	afield := r.Airfields[0]

	expected := airfield.Airfield{ID: "LFLI", ShortName: "ANNEM",
		Name: "ANNEMASSE", ICAO: "LFLI", Flags: 0 | airfield.Asphalt,
		Length: 1290, Runway: "1230", Frequency: 125.87, Elevation: 494,
		Latitude: 46.19194444444444, Longitude: 6.2683333333333335, Region: "FR"}
	if afield != expected {
		t.Errorf("failed to parse airfield :: expected %v got %v", expected, afield)
	}
}

func TestParseUnclearAirstrip(t *testing.T) {
	r := Release{}
	r.Parse([]byte("AMBL21 AMBLETEUSE AERO #   ?G       1      32N504901E0013658FRQ0"))
	afield := r.Airfields[0]
	if afield.Flags&airfield.UnclearAirstrip != 1 {
		t.Errorf("parse failed for unclear airstrip")
	}
}

func TestParseGliderSite(t *testing.T) {
	r := Release{}
	// case GLD#
	r.Parse([]byte("CHALA1 CHALAIS      GLD#LFIHG 83072512350  88N451605E0000058FRQ0"))
	afield := r.Airfields[0]
	if afield.Flags&airfield.GliderSite == 0 {
		t.Errorf("parse failed for glider site")
	}
	// case GLD#GLD
	r.Parse([]byte("HABER1 HABERE POC69 GLD#GLD!G 980119122501113N461611E0062748FRP3"))
	afield = r.Airfields[0]
	if afield.Flags&airfield.GliderSite == 0 {
		t.Errorf("parse failed for glider site")
	}
}

func TestParseULMSite(t *testing.T) {
	r := Release{}
	r.Parse([]byte("CERVE2 CERVENS UL      *ULM!G 28052312350 619N461713E0062638FRQ0"))
	afield := r.Airfields[0]
	if afield.Flags&airfield.ULMSite == 0 {
		t.Errorf("parse failed for ulm site")
	}
}
func TestParseAsphalt(t *testing.T) {
	r := Release{}
	r.Parse([]byte("ANNEM1 ANNEMASSE       #LFLIA129123012587 494N461131E0061606FRQ0"))
	afield := r.Airfields[0]
	if afield.Flags&airfield.Asphalt == 0 {
		t.Errorf("parse failed for asphalt airstrip")
	}
}

func TestParseConcrete(t *testing.T) {
	r := Release{}
	r.Parse([]byte("ANNEM1 ANNEMASSE       #LFLIC129123012587 494N461131E0061606FRQ0"))
	afield := r.Airfields[0]
	if afield.Flags&airfield.Concrete == 0 {
		t.Errorf("parse failed for concrete airstrip")
	}
}

func TestParseLoam(t *testing.T) {
	r := Release{}
	r.Parse([]byte("ANNEM1 ANNEMASSE       #LFLIL129123012587 494N461131E0061606FRQ0"))
	afield := r.Airfields[0]
	if afield.Flags&airfield.Loam == 0 {
		t.Errorf("parse failed for loam airstrip")
	}
}

func TestParseSand(t *testing.T) {
	r := Release{}
	r.Parse([]byte("ANNEM1 ANNEMASSE       #LFLIS129123012587 494N461131E0061606FRQ0"))
	afield := r.Airfields[0]
	if afield.Flags&airfield.Sand == 0 {
		t.Errorf("parse failed for sand airstrip")
	}
}

func TestParseClay(t *testing.T) {
	r := Release{}
	r.Parse([]byte("ANNEM1 ANNEMASSE       #LFLIY129123012587 494N461131E0061606FRQ0"))
	afield := r.Airfields[0]
	if afield.Flags&airfield.Clay == 0 {
		t.Errorf("parse failed for asphalt airstrip")
	}
}

func TestParseGrass(t *testing.T) {
	r := Release{}
	r.Parse([]byte("ANNEM1 ANNEMASSE       #LFLIG129123012587 494N461131E0061606FRQ0"))
	afield := r.Airfields[0]
	if afield.Flags&airfield.Grass == 0 {
		t.Errorf("parse failed for grass airstrip")
	}
}

func TestParseGravel(t *testing.T) {
	r := Release{}
	r.Parse([]byte("ANNEM1 ANNEMASSE       #LFLIV129123012587 494N461131E0061606FRQ0"))
	afield := r.Airfields[0]
	if afield.Flags&airfield.Gravel == 0 {
		t.Errorf("parse failed for gravel airstrip")
	}
}

func TestParseDirt(t *testing.T) {
	r := Release{}
	r.Parse([]byte("ANNEM1 ANNEMASSE       #LFLID129123012587 494N461131E0061606FRQ0"))
	afield := r.Airfields[0]
	if afield.Flags&airfield.Dirt == 0 {
		t.Errorf("parse failed for dirt airstrip")
	}
}

func TestParseUnknownRunwayType(t *testing.T) {
	r := Release{}
	r.Parse([]byte("ANNEM1 ANNEMASSE       #LFLIO129123012587 494N461131E0061606FRQ0"))
	afield := r.Airfields[0]
	if afield.Flags != 0 {
		t.Errorf("parse failed for invalid runway type")
	}
}

func TestParseCatalogNumber(t *testing.T) {
	r := Release{}
	r.Parse([]byte("BONVI2 BONNEVILLE      *FL53S 400523      450N460441E0062310FRP0"))
	afield := r.Airfields[0]
	if afield.Catalog != 53 || afield.Flags&airfield.Outlanding == 0 {
		t.Errorf("parse failed for outlanding catalog number")
	}
}

func TestParseWaypoint(t *testing.T) {
	r := Release{}
	r.Parse([]byte("FURKAP FURKAPASS PASSHOEHE               2432N463422E0082455CHQ6"))
	wpoint := r.Waypoints[0]
	expected := waypoint.Waypoint{
		Name: "FURKAP", ID: "FURKAP", Description: "FURKAPASS PASSHOEHE",
		Latitude: 46.57277777777778, Longitude: 8.415277777777778, Elevation: 2432, Region: "CH",
	}
	if wpoint != expected {
		t.Errorf("parse failed for waypoint: got %v instead of %v", wpoint, expected)
	}
}

func BenchmarkFetch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Fetch("./t/test-release-bench.txt")
		if err != nil {
			b.Errorf("Failed to fetch release :: %v", err)
		}
	}
}
