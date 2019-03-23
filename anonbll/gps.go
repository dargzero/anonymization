package anonbll

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

const ddRegex = "([-+]?[0-9]*\\.?[0-9]+)°?[,\\s]+([-+]?[0-9]*\\.?[0-9]+)°?"
const dmsRegex = "([NS])\\s*([1-9]\\d*)°\\s*(?:([1-9]\\d*)'?)?\\s*(?:([1-9]\\d*)\"?)?\\s*,\\s*([EW])\\s*([1-9]\\d*)°\\s*(?:([1-9]\\d*)'?)?\\s*(?:([1-9]\\d*)\"?)?\\s*"

func ParseCoordinate(coordinate string) (lat, lon float64, err error) {
	format := discoverFormat(coordinate)
	if format == "DD" {
		return readDD(coordinate)
	} else if format == "DMS" {
		return readDMS(coordinate)
	} else {
		return 0, 0, errors.New("unrecognized format")
	}
}

func discoverFormat(coordinates string) string {
	if matches(dmsRegex, coordinates) {
		return "DMS"
	}
	if matches(ddRegex, coordinates) {
		return "DD"
	}
	return "ERR"
}

func matches(re, str string) bool {
	m, _ := regexp.MatchString(re, str)
	return m
}

func readDMS(coords string) (lat float64, lon float64, err error) {
	re := regexp.MustCompile(dmsRegex)
	groups := re.FindStringSubmatch(coords)
	lat, err = parseDmsCoordinateFromSegments(groups[1], groups[2], groups[3], groups[4])
	if err != nil {
		return
	}
	lon, err = parseDmsCoordinateFromSegments(groups[5], groups[6], groups[7], groups[8])
	if err != nil {
		return
	}
	return
}

func parseDmsCoordinateFromSegments(seg1, seg2, seg3, seg4 string) (float64, error) {
	lat1, err := parseDmsSegment(seg2, false)
	if err != nil {
		return 0, err
	}
	lat2, err := parseDmsSegment(seg3, true)
	if err != nil {
		return 0, err
	}
	lat3, err := parseDmsSegment(seg4, true)
	if err != nil {
		return 0, err
	}
	lat := lat1 + lat2/60.0 + lat3/360.0
	hemisphere := strings.ToLower(seg1)
	if hemisphere == "s" || hemisphere == "w" {
		lat = -lat
	}
	return lat, nil
}

func parseDmsSegment(seg string, optional bool) (val float64, err error) {
	val = 0.0
	if !optional || optional && seg != "" {
		val, err = strconv.ParseFloat(seg, 64)
	}
	return
}

func readDD(coords string) (lat, lon float64, err error) {
	re := regexp.MustCompile(ddRegex)
	groups := re.FindStringSubmatch(coords)
	lat, err = strconv.ParseFloat(groups[1], 64)
	if err != nil {
		return
	}
	lon, err = strconv.ParseFloat(groups[2], 64)
	return
}
