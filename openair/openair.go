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

// Package openair provides functionality for parsing airspace information defined
// in the OpenAir format.
//
// The format specification is available at:
// 	http://www.winpilot.com/UsersGuide/UserAirspace.asp
//
package openair

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/rochaporto/ezgliding/airspace"
)

// Local temporary storage for airspace pen/brush types
// Key is the airspace class, value in the Pen object
var airspacePens = map[byte]airspace.Pen{}

// Fetch gets and returns the airspace definitions at the given location
// Both http URIs and local (relative or absolute) paths are supported.
func Fetch(location string) ([]airspace.Airspace, error) {
	var content []byte

	resp, err := http.Get(location)
	// case http
	if err == nil {
		defer resp.Body.Close()
		content, err = ioutil.ReadAll(resp.Body)
	} else { // case file
		content, err = ioutil.ReadFile(location)
		if err != nil {
			return nil, err
		}
	}
	return Parse(content)
}

// Parse parses the content given, retrieving the corresponding array
// of Airspace objects.
func Parse(content []byte) ([]airspace.Airspace, error) {
	items := strings.Split(string(content), "*\n")
	result := []airspace.Airspace{}
	for i := range items {
		airspace, found, err := parseSingle([]byte(items[i]))
		if err != nil {
			return nil, err
		}
		if found {
			result = append(result, airspace)
		}
	}
	return result, nil
}

// styleToAirspace converts the given value to the corresponding enum
// value in airspace.PenStyle.
// Supported values are 0 (Solid), 1 (Dash), 5 (None). Any other value
// will also return None.
func styleToAirspace(value int) airspace.PenStyle {
	switch value {
	case 0:
		return airspace.Solid
	case 1:
		return airspace.Dash
	default:
		return airspace.None
	}
}

// parseSingle expects a single airspace definition in the content object,
// returning the corresponding Airspace object.
// It is usually called by Parse.
func parseSingle(content []byte) (airspace.Airspace, bool, error) {
	var aspace airspace.Airspace
	var x string
	var clockwise bool
	found := false

	lines := strings.Split(string(content), "\n")
	for i := range lines {
		line := strings.Trim(lines[i], " ")
		if len(line) == 0 || line[0] == '*' { // comment or empty
			continue
		}
		elems := strings.SplitN(line, " ", 2)
		key, value := elems[0], elems[1]
		switch key {
		case "AC":
			aspace.Class = value[0]
			aspace.Pen = airspacePens[aspace.Class]
		case "AN":
			found = true
			aspace.Name = strings.Trim(value, " ")
		case "AL":
			aspace.Floor = value
		case "AH":
			aspace.Ceiling = value
		case "DA":
			values := strings.Split(value, ",")
			angleStart, _ := strconv.ParseFloat(values[1], 64)
			angleEnd, _ := strconv.ParseFloat(values[2], 64)
			radius, _ := strconv.ParseFloat(values[0], 64)
			aspace.Segments = append(aspace.Segments,
				airspace.Segment{Type: airspace.Arc, X: x, Clockwise: clockwise,
					Radius: radius, AngleStart: angleStart, AngleEnd: angleEnd})
		case "DB":
			values := strings.Split(value, ",")
			aspace.Segments = append(aspace.Segments,
				airspace.Segment{Type: airspace.Arc, X: x, Clockwise: clockwise,
					Coordinate1: values[0], Coordinate2: values[1]})
		case "DC":
			radius, _ := strconv.ParseFloat(value, 64)
			aspace.Segments = append(aspace.Segments,
				airspace.Segment{Type: airspace.Circle, X: x, Clockwise: clockwise,
					Radius: radius})
		case "DP":
			aspace.Segments = append(aspace.Segments,
				airspace.Segment{Type: airspace.Polygon, X: x, Clockwise: clockwise, Coordinate1: value})
		case "V":
			splitequals := strings.Split(value, "=")
			varkey := splitequals[0]
			varvalue := strings.Trim(splitequals[1], " ")
			switch varkey {
			case "X":
				x = varvalue
			case "D":
				clockwise = (varvalue == "+")
			}
		case "SP": // pen to draw (including color)
			split := strings.Split(value, ",")
			width, _ := strconv.Atoi(split[1])
			r, _ := strconv.Atoi(split[2])
			g, _ := strconv.Atoi(split[3])
			b, _ := strconv.Atoi(split[4])
			airspacePens[aspace.Class] = airspace.Pen{
				Style: airspace.Solid, Width: width,
				Color:       color.RGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: 1.0},
				InsideColor: color.RGBA64{},
			}
		case "SB": // brush to draw (color) TODO:
			//split := strings.Split(value, ",")
			//r, _ := strconv.Atoi(split[0])
			//g, _ := strconv.Atoi(split[1])
			//b, _ := strconv.Atoi(split[2])
			//airspacePens[airspace.Class].InsideColor = color.RGBA64{}
		default:
			return airspace.Airspace{}, false, fmt.Errorf("Unrecognized key '%v' in '%v'", key, line)
		}
	}

	return aspace, found, nil
}
