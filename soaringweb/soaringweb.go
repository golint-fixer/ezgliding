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

// soaringweb provides functionality to fetch and parse airspace
// information, taking the international soaringweb db as input.
//
// Check the soaringweb website for more information on the data:
// http://soaringweb.org/Airspace/HomePage.html
//
package soaringweb

import (
	"fmt"
	"github.com/rochaporto/ezgliding/common"
	"io/ioutil"
	"net/http"
)

// List returns all latest airspace info available
func List(location string) ([]common.Airspace, error) {
	var content []byte
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
	// TODO: actual listing
	fmt.Println("%v", content)

	return nil, nil
}

// Fetch gets and returns the Airspace at the given location
func Fetch(location string) (*common.Airspace, error) {
	var content []byte

	resp, err := http.Get(location)
	// case http
	if err == nil {
		defer resp.Body.Close()
		content, err = ioutil.ReadAll(resp.Body)
	} else { // case file
		content, err = ioutil.ReadFile(location)
	}
	if err != nil {
		return nil, err
	}
	fmt.Println("%v", content)
	// TODO: actual fetching
	return nil, nil
}
