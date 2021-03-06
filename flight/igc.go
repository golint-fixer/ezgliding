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

package flight

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rochaporto/ezgliding/spatial"
)

const (
	// TimeFormat is the golang time.Parse format for IGC time
	TimeFormat = "150405"
	// DateFormat is the golang time.Parse format for IGC time
	DateFormat = "020106"
)

// ParseIGC returns a Flight object corresponding to the given content.
// content should be a text string in the IGC format.
func ParseIGC(content string) (Flight, error) {
	f := NewFlight()
	var err error
	p := IGCParser{}
	lines := strings.Split(content, "\n")
	for i := range lines {
		line := strings.TrimSpace(lines[i])
		// ignore empty lines
		if len(strings.Trim(line, " ")) < 1 {
			continue
		}
		switch line[0] {
		case 'A':
			err = p.parseA(line, &f)
		case 'B':
			err = p.parseB(line, &f)
		case 'C':
			if !p.taskDone {
				err = p.parseC(lines[i:], &f)
			}
		case 'D':
			err = p.parseD(line, &f)
		case 'E':
			err = p.parseE(line, &f)
		case 'F':
			err = p.parseF(line, &f)
		case 'G':
			err = p.parseG(line, &f)
		case 'H':
			err = p.parseH(line, &f)
		case 'I':
			err = p.parseI(line)
		case 'J':
			err = p.parseJ(line)
		case 'K':
			err = p.parseK(line, &f)
		case 'L':
			err = p.parseL(line, &f)
		default:
			err = fmt.Errorf("invalid record :: %v", line)
		}
		if err != nil {
			return f, err
		}
	}

	return f, nil
}

type field struct {
	start int64
	end   int64
	tlc   string
}

// IGCParser gives functionality to parse IGC flight files.
type IGCParser struct {
	IFields  []field
	JFields  []field
	taskDone bool
	numSat   int
}

func (p *IGCParser) parseA(line string, f *Flight) error {
	if len(line) < 7 {
		return fmt.Errorf("line too short :: %v", line)
	}
	f.Header.Manufacturer = line[1:4]
	f.Header.UniqueID = line[4:7]
	f.Header.AdditionalData = line[7:]
	return nil
}

func (p *IGCParser) parseB(line string, f *Flight) error {
	if len(line) < 37 {
		return fmt.Errorf("line too short :: %v", line)
	}
	pt := NewPoint()
	var err error
	pt.Time, err = time.Parse(TimeFormat, line[1:7])
	if err != nil {
		return err
	}
	pt.Latitude = spatial.DMD2Decimal(line[7:15])
	pt.Longitude = spatial.DMD2Decimal(line[15:24])
	if line[24] == 'A' || line[24] == 'V' {
		pt.FixValidity = line[24]
	} else {
		return fmt.Errorf("invalid fix validity :: %v", line[24])
	}
	pt.PressureAltitude, err = strconv.ParseInt(line[25:30], 10, 64)
	if err != nil {
		return err
	}
	pt.GNSSAltitude, err = strconv.ParseInt(line[30:35], 10, 64)
	if err != nil {
		return err
	}
	for _, f := range p.IFields {
		pt.IData[f.tlc] = line[f.start-1 : f.end]
	}
	pt.NumSatellites = p.numSat
	f.Points = append(f.Points, pt)
	return nil
}

func (p *IGCParser) parseC(lines []string, f *Flight) error {
	line := lines[0]
	if len(line) < 25 {
		return fmt.Errorf("wrong line size :: %v", line)
	}
	var err error
	var nTP int
	if nTP, err = strconv.Atoi(line[23:25]); err != nil {
		return fmt.Errorf("invalid number of turnpoints :: %v", line)
	}
	if len(lines) < 5+nTP {
		return fmt.Errorf("invalid number of C record lines :: %v", lines)
	}
	if f.Task.DeclarationDate, err = time.Parse(DateFormat+TimeFormat, lines[0][1:13]); err != nil {
		f.Task.DeclarationDate = time.Time{}
	}
	if f.Task.FlightDate, err = time.Parse(DateFormat, lines[0][13:19]); err != nil {
		f.Task.FlightDate = time.Time{}
	}
	if f.Task.Number, err = strconv.Atoi(line[19:23]); err != nil {
		return err
	}
	f.Task.Description = line[25:]
	if f.Task.Takeoff, err = p.taskPoint(lines[1]); err != nil {
		return err
	}
	if f.Task.Start, err = p.taskPoint(lines[2]); err != nil {
		return err
	}
	for i := 0; i < nTP; i++ {
		var tp Point
		if tp, err = p.taskPoint(lines[3+i]); err != nil {
			return err
		}
		f.Task.Turnpoints = append(f.Task.Turnpoints, tp)
	}
	if f.Task.Finish, err = p.taskPoint(lines[3+nTP]); err != nil {
		return err
	}
	if f.Task.Landing, err = p.taskPoint(lines[4+nTP]); err != nil {
		return err
	}
	p.taskDone = true
	return nil
}

func (p *IGCParser) taskPoint(line string) (Point, error) {
	if len(line) < 18 {
		return Point{}, fmt.Errorf("line too short :: %v", line)
	}
	return Point{
		Latitude:    spatial.DMD2Decimal(line[1:9]),
		Longitude:   spatial.DMD2Decimal(line[9:18]),
		Description: line[18:],
	}, nil
}

func (p *IGCParser) parseD(line string, f *Flight) error {
	if len(line) < 6 {
		return fmt.Errorf("line too short :: %v", line)
	}
	if line[1] == '2' {
		f.DGPSStationID = line[2:6]
	}
	return nil
}

func (p *IGCParser) parseE(line string, f *Flight) error {
	if len(line) < 10 {
		return fmt.Errorf("line too short :: %v", line)
	}
	t, err := time.Parse(TimeFormat, line[1:7])
	if err != nil {
		return err
	}
	if f.Events[t] == nil {
		f.Events[t] = make(map[string]string)
	}
	f.Events[t][line[7:10]] = line[10:]
	return nil
}

func (p *IGCParser) parseF(line string, f *Flight) error {
	if len(line) < 7 {
		return fmt.Errorf("line too short :: %v", line)
	}
	t, err := time.Parse(TimeFormat, line[1:7])
	if err != nil {
		return err
	}
	if f.Satellites[t] == nil {
		f.Satellites[t] = []int{}
	}
	for i := 7; i < len(line)-1; i = i + 2 {
		var n int
		if n, err = strconv.Atoi(line[i : i+2]); err != nil {
			return err
		}
		f.Satellites[t] = append(f.Satellites[t], n)
	}
	p.numSat = len(f.Satellites[t])
	return nil
}

func (p *IGCParser) parseG(line string, f *Flight) error {
	f.Signature = f.Signature + line[1:]
	return nil
}

func (p *IGCParser) parseH(line string, f *Flight) error {
	var err error
	if len(line) < 5 {
		return fmt.Errorf("line too short :: %v", line)
	}

	switch line[2:5] {
	case "DTE":
		if len(line) < 11 {
			return fmt.Errorf("line too short :: %v", line)
		}
		f.Header.Date, err = time.Parse(DateFormat, line[5:11])
	case "FXA":
		if len(line) < 8 {
			return fmt.Errorf("line too short :: %v", line)
		}
		f.Header.FixAccuracy, err = strconv.ParseInt(line[5:8], 10, 64)
	case "PLT":
		f.Header.Pilot = stripUpTo(line[5:], ":")
	case "CM2":
		f.Header.Crew = stripUpTo(line[5:], ":")
	case "GTY":
		f.Header.GliderType = stripUpTo(line[5:], ":")
	case "GID":
		f.Header.GliderID = stripUpTo(line[5:], ":")
	case "DTM":
		if len(line) < 8 {
			return fmt.Errorf("line too short :: %v", line)
		}
		f.Header.GPSDatum = stripUpTo(line[5:], ":")
	case "RFW":
		f.Header.FirmwareVersion = stripUpTo(line[5:], ":")
	case "RHW":
		f.Header.HardwareVersion = stripUpTo(line[5:], ":")
	case "FTY":
		f.Header.FlightRecorder = stripUpTo(line[5:], ":")
	case "GPS":
		f.Header.GPS = line[5:]
	case "PRS":
		f.Header.PressureSensor = stripUpTo(line[5:], ":")
	case "CID":
		f.Header.CompetitionID = stripUpTo(line[5:], ":")
	case "CCL":
		f.Header.CompetitionClass = stripUpTo(line[5:], ":")
	case "TZN":
		z, err := strconv.ParseFloat(stripUpTo(line[5:], ":"), 64)
		if err != nil {
			return err
		}
		f.Header.Timezone = *time.FixedZone("", int(z*3600))
	default:
		err = fmt.Errorf("unknown record :: %v", line)
	}

	return err
}

func (p *IGCParser) parseI(line string) error {
	if len(line) < 3 {
		return fmt.Errorf("line too short :: %v", line)
	}
	n, err := strconv.ParseInt(line[1:3], 10, 0)
	if err != nil {
		return fmt.Errorf("invalid number of I fields :: %v", line)
	}
	if len(line) != int(n*7+3) {
		return fmt.Errorf("wrong line size :: %v", line)
	}
	for i := 0; i < int(n); i++ {
		s := i*7 + 3
		start, _ := strconv.ParseInt(line[s:s+2], 10, 0)
		end, _ := strconv.ParseInt(line[s+2:s+4], 10, 0)
		tlc := line[s+4 : s+7]
		p.IFields = append(p.IFields, field{start: start, end: end, tlc: tlc})
	}
	return nil
}

func (p *IGCParser) parseJ(line string) error {
	if len(line) < 3 {
		return fmt.Errorf("line too short :: %v", line)
	}
	n, err := strconv.ParseInt(line[1:3], 10, 0)
	if err != nil {
		return fmt.Errorf("invalid number of J fields :: %v", line)
	}
	if len(line) != int(n*7+3) {
		return fmt.Errorf("wrong line size :: %v", line)
	}
	for i := 0; i < int(n); i++ {
		s := i*7 + 3
		start, _ := strconv.ParseInt(line[s:s+2], 10, 0)
		end, _ := strconv.ParseInt(line[s+2:s+4], 10, 0)
		tlc := line[s+4 : s+7]
		p.JFields = append(p.JFields, field{start: start, end: end, tlc: tlc})
	}
	return nil
}

func (p *IGCParser) parseK(line string, f *Flight) error {
	if len(line) < 7 {
		return fmt.Errorf("line too short :: %v", line)
	}
	t, err := time.Parse(TimeFormat, line[1:7])
	if err != nil {
		return err
	}
	fields := make(map[string]string)
	for _, f := range p.JFields {
		fields[f.tlc] = line[f.start-1 : f.end]
	}
	f.K[t] = fields
	return nil
}

func (p *IGCParser) parseL(line string, f *Flight) error {
	if len(line) < 4 {
		return fmt.Errorf("line too short :: %v", line)
	}
	f.Logbook = append(f.Logbook, LogEntry{Type: line[1:4], Text: line[4:]})
	return nil
}

func stripUpTo(s string, sep string) string {
	i := strings.Index(s, sep)
	if i == -1 {
		return s
	}
	return s[i+1:]
}
